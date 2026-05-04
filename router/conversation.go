package router

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/authen"
	agentpb "github.com/yzf120/elysia-backend/proto/agent"
	"github.com/yzf120/elysia-backend/rpc"
	"github.com/yzf120/elysia-backend/service"
	agent_session "github.com/yzf120/elysia-session/proto/agent_session"
	conversationpb "github.com/yzf120/elysia-session/proto/conversation"
)

// 内容审核服务（在 router.Init() 中初始化）
var contentModerationService *service.ContentModerationService

// AI回复兜底文案（当AI回复违规时返回给用户）
const aiViolationFallbackReply = "抱歉，AI生成的回复内容不适合展示。请尝试换一种方式提问，或联系管理员。"

// 用户被封禁时的提示
const userBannedMessage = "您因多次发送违规内容，AI对话功能已被禁止使用。如有疑问请联系管理员。"

// AIChatRequest AI对话请求（来自前端）
type AIChatRequest struct {
	// 会话ID（可选，首轮对话为空，后续对话传入）
	SessionID string `json:"session_id,omitempty"`
	// 题目ID（编程界面开启的对话时传入，普通对话不传或传0）
	ProblemID int64 `json:"problem_id,omitempty"`
	// 问题类型标识，如 "algorithm_problem" 表示算法题
	QuestionType string `json:"question_type"`
	// 题目信息（作为上下文传给AI）
	ProblemInfo *ProblemContext `json:"problem_info,omitempty"`
	// 对话历史（包含本次用户消息）
	Messages []ChatMessage `json:"messages"`
	// 模型ID（可选，默认使用豆包）
	ModelID string `json:"model_id,omitempty"`
	// 是否开启深度思考模式
	EnableThinking bool `json:"enable_thinking,omitempty"`
	// 用户当前IDE中的代码（作为上下文传给AI）
	UserCode string `json:"user_code,omitempty"`
	// 用户当前选择的编程语言
	UserCodeLang string `json:"user_code_lang,omitempty"`
	// 运行记录加入对话：判题结果（accepted / partial_pass）
	JudgeResult string `json:"judge_result,omitempty"`
	// 运行记录加入对话：未通过的测试用例（JSON 字符串）
	FailedCases string `json:"failed_cases,omitempty"`
}

// ChatMessage 单条对话消息
type ChatMessage struct {
	Role    string `json:"role"`    // "user" 或 "assistant"
	Content string `json:"content"` // 消息内容
}

// ProblemContext 题目上下文信息
type ProblemContext struct {
	ID           int64    `json:"id"`
	Title        string   `json:"title"`
	Difficulty   string   `json:"difficulty"`
	Description  string   `json:"description"`
	InputFormat  string   `json:"input_format"`
	OutputFormat string   `json:"output_format"`
	Tags         []string `json:"tags,omitempty"`
	TimeLimit    int      `json:"time_limit"`
	MemoryLimit  int      `json:"memory_limit"`
}

func registerConversation(router *mux.Router) {
	// 学生AI答疑接口（SSE流式输出）
	router.HandleFunc("/student/ai/chat", studentAIChatHandler).Methods("POST", "OPTIONS")
	// 查询支持的模型列表
	router.HandleFunc("/student/ai/models", studentAIModelsHandler).Methods("GET", "OPTIONS")
	// 查询用户会话列表
	router.HandleFunc("/student/ai/sessions", studentAISessionsHandler).Methods("GET", "OPTIONS")
	// 查询某会话的消息列表
	router.HandleFunc("/student/ai/sessions/{sessionId}/messages", studentAISessionMessagesHandler).Methods("GET", "OPTIONS")
}

// studentAIModelsHandler 查询支持的模型列表（自动过滤管理员禁用的模型）
// GET /student/ai/models
func studentAIModelsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()

	// 通过 chat-agent RPC 查询模型列表
	resp, err := rpc.GetAgentClient().GetProxy().ListModels(ctx, &agentpb.AgentListModelsRequest{})
	if err != nil {
		log.Printf("[conversation] 查询模型列表失败: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 500, "message": "查询模型列表失败"},
			"data":  nil,
		})
		return
	}

	// 从数据库查询禁用的模型ID列表，过滤掉禁用的模型
	if aiModelConfigService != nil {
		disabledIds, err := aiModelConfigService.ListDisabledModelIds()
		if err != nil {
			log.Printf("[conversation] 查询禁用模型列表失败: %v", err)
			// 查询失败不阻断，返回全量模型
		} else if len(disabledIds) > 0 {
			disabledSet := make(map[string]bool, len(disabledIds))
			for _, id := range disabledIds {
				disabledSet[id] = true
			}
			// 过滤掉禁用的模型
			filteredModels := make([]agentpb.AgentModelInfo, 0, len(resp.Models))
			for _, m := range resp.Models {
				if !disabledSet[m.ModelID] {
					filteredModels = append(filteredModels, m)
				}
			}
			resp.Models = filteredModels
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{"code": 0},
		"data":  resp,
	})
}

// studentAIChatHandler 学生AI答疑处理器（SSE流式输出）
// POST /student/ai/chat
func studentAIChatHandler(w http.ResponseWriter, r *http.Request) {
	// 处理 OPTIONS 预检请求
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	reqCtx := r.Context()

	// 验证学生身份
	studentId, ok := authen.GetRoleIDFromContext(reqCtx)
	if !ok || studentId == "" {
		http.Error(w, "未授权，请先登录", http.StatusUnauthorized)
		return
	}

	// 解析请求体
	request := &AIChatRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		http.Error(w, "请求参数错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(request.Messages) == 0 {
		http.Error(w, "messages 不能为空", http.StatusBadRequest)
		return
	}

	// 确定用户角色
	userRole := "student"
	if strings.HasPrefix(studentId, "tea_") {
		userRole = "teacher"
	}

	// ===== 1. 检查用户是否因违规被封禁 =====
	banned, violationCount, err := contentModerationService.IsUserBanned(studentId)
	if err != nil {
		log.Printf("[conversation] 检查用户封禁状态失败，学生: %s, err: %v", studentId, err)
		// 查询失败时不阻断，继续处理
	} else if banned {
		log.Printf("[conversation] 用户已被封禁，学生: %s, 违规次数: %d", studentId, violationCount)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 4031, "message": userBannedMessage},
			"data":  map[string]interface{}{"violation_count": violationCount},
		})
		return
	}

	// ===== 2. 同步审核用户消息（Q） =====
	userMsg := ""
	for i := len(request.Messages) - 1; i >= 0; i-- {
		if request.Messages[i].Role == "user" {
			userMsg = request.Messages[i].Content
			break
		}
	}

	if userMsg != "" {
		moderateCtx, moderateCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer moderateCancel()

		// 文本审核
		pureText := stripImageContent(userMsg)
		if strings.TrimSpace(pureText) != "" {
			textResult, err := contentModerationService.ModerateText(moderateCtx, studentId, userRole, request.SessionID, "user", pureText)
			if err != nil {
				log.Printf("[conversation] 用户消息文本审核异常，学生: %s, err: %v", studentId, err)
			} else if !textResult.Passed {
				// 用户消息违规，返回警示，不调用AI
				log.Printf("[conversation] 用户消息违规被拦截，学生: %s, label: %s", studentId, textResult.Label)
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": map[string]interface{}{"code": 4032, "message": textResult.Message},
					"data":  map[string]interface{}{"blocked": true, "label": textResult.Label},
				})
				return
			}
		}

		// 图片审核
		imageURLs := extractImageURLs(userMsg)
		for i, imgURL := range imageURLs {
			imgReq := buildImageModerationRequest(imgURL, studentId, request.SessionID)
			if imgReq == nil {
				continue
			}
			imgResult, err := contentModerationService.ModerateImage(moderateCtx, studentId, userRole, request.SessionID, "user", imgReq)
			if err != nil {
				log.Printf("[conversation] 用户消息图片#%d审核异常，学生: %s, err: %v", i+1, studentId, err)
			} else if !imgResult.Passed {
				log.Printf("[conversation] 用户消息图片#%d违规被拦截，学生: %s, label: %s", i+1, studentId, imgResult.Label)
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": map[string]interface{}{"code": 4032, "message": "您发送的图片包含违规内容，已被系统拦截。请规范用语，多次违规将被禁止使用AI对话功能。"},
					"data":  map[string]interface{}{"blocked": true, "label": imgResult.Label},
				})
				return
			}
		}
	}

	// 设置 SSE 响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no") // 禁用 nginx 缓冲

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "不支持流式响应", http.StatusInternalServerError)
		return
	}

	// 选择模型
	modelID := request.ModelID
	if modelID == "" {
		modelID = "doubao-seed-1-6-lite-251015" // 默认豆包模型
	}

	// 构建系统提示词
	systemPrompt := buildSystemPrompt(request.QuestionType, request.ProblemInfo, request.UserCode, request.UserCodeLang, request.JudgeResult, request.FailedCases)

	// 调试日志：打印系统提示词构建的关键参数
	hasProblemInfo := request.ProblemInfo != nil
	problemTitle := ""
	if hasProblemInfo {
		problemTitle = request.ProblemInfo.Title
	}
	log.Printf("[conversation][DEBUG] 系统提示词构建参数: questionType=%s, hasProblemInfo=%v, problemTitle=%s, userCodeLen=%d, judgeResult=%s, failedCasesLen=%d, systemPromptLen=%d",
		request.QuestionType, hasProblemInfo, problemTitle, len(request.UserCode), request.JudgeResult, len(request.FailedCases), len(systemPrompt))
	if len(systemPrompt) > 500 {
		log.Printf("[conversation][DEBUG] 系统提示词前500字符: %s", systemPrompt[:500])
	} else {
		log.Printf("[conversation][DEBUG] 系统提示词全文: %s", systemPrompt)
	}

	// 构建发送给 chat-agent 的消息列表
	agentMessages := make([]agentpb.AgentChatMessage, 0, len(request.Messages))
	for _, msg := range request.Messages {
		agentMessages = append(agentMessages, agentpb.AgentChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 使用独立的 context 用于 RPC 流式调用，避免 HTTP 框架超时导致 context 被取消
	// 同时监听客户端断开，主动取消 RPC 调用
	rpcCtx, rpcCancel := context.WithCancel(context.Background())
	defer rpcCancel()

	// 监听客户端断开（reqCtx 取消时同步取消 rpcCtx）
	go func() {
		select {
		case <-reqCtx.Done():
			log.Printf("[conversation] 客户端断开连接，取消 RPC 调用，学生: %s", studentId)
			rpcCancel()
		case <-rpcCtx.Done():
		}
	}()

	// 调用 chat-agent 的流式 RPC
	agentReq := &agentpb.AgentStreamChatRequest{
		ModelID:      modelID,
		Messages:     agentMessages,
		SystemPrompt: systemPrompt,
	}

	// 构建 ExtraParams（透传上下文信息给 Agent）
	agentReq.ExtraParams = map[string]string{
		"user_id":    studentId,
		"user_role":  userRole,
		"session_id": request.SessionID,
	}
	if request.EnableThinking {
		agentReq.ExtraParams["enable_thinking"] = "true"
		log.Printf("[conversation] 深度思考模式已开启，模型: %s，学生: %s", modelID, studentId)
	}
	if request.ProblemID > 0 {
		agentReq.ExtraParams["problem_id"] = fmt.Sprintf("%d", request.ProblemID)
	}
	if request.ProblemInfo != nil {
		if pBytes, err := json.Marshal(request.ProblemInfo); err == nil {
			agentReq.ExtraParams["problem_info"] = string(pBytes)
		}
	}
	if request.UserCode != "" {
		agentReq.ExtraParams["student_code"] = request.UserCode
		agentReq.ExtraParams["language"] = request.UserCodeLang
	}
	// 运行记录加入对话：透传判题结果和未通过用例
	if request.JudgeResult != "" {
		agentReq.ExtraParams["judge_result"] = request.JudgeResult
		log.Printf("[conversation] 运行记录加入对话，judge_result: %s，学生: %s", request.JudgeResult, studentId)
	}
	if request.FailedCases != "" {
		agentReq.ExtraParams["failed_cases"] = request.FailedCases
	}

	agentStream, err := rpc.GetAgentClient().GetProxy().StreamChat(rpcCtx, agentReq)
	if err != nil {
		log.Printf("[conversation] 调用 chat-agent StreamChat 失败: %v", err)
		// 通过 SSE 发送错误事件
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
		flusher.Flush()
		return
	}

	log.Printf("[conversation] 开始接收 chat-agent 流式响应，学生: %s", studentId)

	// 用于拼接完整的 AI 回复内容
	var aiReplyBuilder strings.Builder
	streamErr := false

	// 逐个接收 chat-agent 的流式响应，通过 SSE 发送给前端
	for {
		// 检查客户端是否已断开（rpcCtx 被取消）
		select {
		case <-rpcCtx.Done():
			log.Printf("[conversation] 客户端已断开，停止接收流式响应，学生: %s", studentId)
			streamErr = true
			goto streamEnd
		default:
		}

		chunk, err := agentStream.Recv()
		if err == io.EOF {
			// 流结束，发送结束事件
			fmt.Fprintf(w, "event: done\ndata: [DONE]\n\n")
			flusher.Flush()
			break
		}
		if err != nil {
			// 区分客户端主动断开和真正的错误
			if rpcCtx.Err() != nil {
				log.Printf("[conversation] 客户端断开导致 RPC 流关闭，学生: %s", studentId)
			} else {
				log.Printf("[conversation] 接收 chat-agent 响应失败: %v", err)
				fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
				flusher.Flush()
			}
			streamErr = true
			break
		}

		// 累积 AI 回复内容
		if chunk.Content != "" {
			aiReplyBuilder.WriteString(chunk.Content)
		}

		// 过滤无意义的空 chunk（content 为空且不是结束标记），避免向前端发送大量无用数据
		if chunk.Content == "" && !chunk.IsEnd && chunk.FinishReason == "" &&
			chunk.PromptTokens == 0 && chunk.CompletionTokens == 0 && chunk.TotalTokens == 0 {
			continue
		}

		// 将 chunk 序列化为 JSON 并通过 SSE 发送
		chunkData, jsonErr := json.Marshal(map[string]interface{}{
			"content":           chunk.Content,
			"is_end":            chunk.IsEnd,
			"finish_reason":     chunk.FinishReason,
			"prompt_tokens":     chunk.PromptTokens,
			"completion_tokens": chunk.CompletionTokens,
			"total_tokens":      chunk.TotalTokens,
		})
		if jsonErr != nil {
			log.Printf("[conversation] 序列化 chunk 失败: %v", jsonErr)
			continue
		}

		// 写入 SSE 数据，检查写入错误（客户端断开时写入会失败）
		_, writeErr := fmt.Fprintf(w, "data: %s\n\n", string(chunkData))
		if writeErr != nil {
			log.Printf("[conversation] SSE 写入失败（客户端可能已断开），学生: %s, err: %v", studentId, writeErr)
			streamErr = true
			break
		}
		flusher.Flush()

		// 如果是最后一个 chunk，退出
		if chunk.IsEnd {
			fmt.Fprintf(w, "event: done\ndata: [DONE]\n\n")
			flusher.Flush()
			break
		}
	}
streamEnd:

	log.Printf("[conversation] SSE 流式响应完成，学生: %s", studentId)

	// ===== 3. 异步审核 AI 回复 + 存储会话记录 =====
	// 审核和存储都在 goroutine 中执行，不阻塞 HTTP 连接关闭
	aiReply := aiReplyBuilder.String()
	if streamErr || strings.TrimSpace(aiReply) == "" {
		aiReply = "抱歉，AI 助教暂时无法回答，请稍后再试。"
		log.Printf("[conversation] AI 回复异常，使用兜底回复，学生: %s", studentId)
	}

	go func() {
		// 审核 AI 回复内容
		if strings.TrimSpace(aiReply) != "" {
			aiModerateCtx, aiModerateCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer aiModerateCancel()

			aiPureText := stripImageContent(aiReply)
			if strings.TrimSpace(aiPureText) != "" {
				aiTextResult, err := contentModerationService.ModerateText(aiModerateCtx, studentId, userRole, request.SessionID, "ai", aiPureText)
				if err != nil {
					log.Printf("[conversation] AI回复文本审核异常，学生: %s, err: %v", studentId, err)
				} else if !aiTextResult.Passed {
					// AI回复违规，存储时使用兜底文案
					log.Printf("[conversation] AI回复违规，使用兜底文案，学生: %s, label: %s", studentId, aiTextResult.Label)
					aiReply = aiViolationFallbackReply
				}
			}
		}

		// 存储会话和对话记录到 session 服务

		sessionSvcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		sessionID := request.SessionID
		// 计算本轮 Q/A 的消息序号
		// 已有消息数 = (len(request.Messages) - 1) 条历史 + 本次用户消息
		// 本次用户消息 seq = 已有历史消息数 + 1（奇数）
		// 本次 AI 回复 seq = 已有历史消息数 + 2（偶数）
		historyCount := int32(len(request.Messages) - 1) // 不含本次用户消息
		userSeq := historyCount + 1
		aiSeq := historyCount + 2

		if sessionID == "" {
			// 首轮对话：创建新会话
			sessionTitle := userMsg
			if len([]rune(sessionTitle)) > 30 {
				runes := []rune(sessionTitle)
				sessionTitle = string(runes[:30]) + "..."
			}
			createRsp, err := rpc.GetSessionClient().CreateSession(sessionSvcCtx, &agent_session.CreateSessionRequest{
				UserId:       studentId,
				SessionTitle: sessionTitle,
				ProblemId:    request.ProblemID,
			})
			if err != nil {
				log.Printf("[conversation] 创建会话失败，学生: %s, err: %v", studentId, err)
				return
			}
			if createRsp.Code != 200 {
				log.Printf("[conversation] 创建会话返回错误，学生: %s, code: %d, msg: %s", studentId, createRsp.Code, createRsp.Message)
				return
			}
			sessionID = createRsp.SessionId
			log.Printf("[conversation] 创建会话成功，学生: %s, sessionID: %s", studentId, sessionID)
		}

		// 写入用户消息（奇数 seq）
		_, err := rpc.GetSessionClient().CreateConversation(sessionSvcCtx, &conversationpb.CreateConversationRequest{
			SessionId:   sessionID,
			UserId:      studentId,
			ModelId:     modelID,
			MessageType: conversationpb.MessageType_MESSAGE_TYPE_TEXT,
			SenderType:  conversationpb.SenderType_SENDER_TYPE_USER,
			Content:     userMsg,
			MessageSeq:  userSeq,
		})
		if err != nil {
			log.Printf("[conversation] 写入用户消息失败，sessionID: %s, err: %v", sessionID, err)
			return
		}

		// 写入 AI 回复（偶数 seq）
		_, err = rpc.GetSessionClient().CreateConversation(sessionSvcCtx, &conversationpb.CreateConversationRequest{
			SessionId:   sessionID,
			UserId:      studentId,
			ModelId:     modelID,
			MessageType: conversationpb.MessageType_MESSAGE_TYPE_TEXT,
			SenderType:  conversationpb.SenderType_SENDER_TYPE_AGENT,
			Content:     aiReply,
			MessageSeq:  aiSeq,
		})
		if err != nil {
			log.Printf("[conversation] 写入 AI 回复失败，sessionID: %s, err: %v", sessionID, err)
			return
		}

		log.Printf("[conversation] 会话记录存储完成，sessionID: %s, userSeq: %d, aiSeq: %d", sessionID, userSeq, aiSeq)
	}()
}

// studentAISessionsHandler 获取用户AI会话列表
// GET /student/ai/sessions?page=1&page_size=20&problem_id=0
func studentAISessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()
	studentId, ok := authen.GetRoleIDFromContext(ctx)
	if !ok || studentId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 401, "message": "未授权"},
			"data":  nil,
		})
		return
	}

	// 解析分页参数
	page := int32(1)
	pageSize := int32(20)
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := fmt.Sscanf(p, "%d", &page); v == 0 || err != nil {
			page = 1
		}
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if v, err := fmt.Sscanf(ps, "%d", &pageSize); v == 0 || err != nil {
			pageSize = 20
		}
	}

	svcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.GetSessionClient().ListSessionsByUser(svcCtx, &agent_session.ListSessionsByUserRequest{
		UserId:   studentId,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		log.Printf("[conversation] 获取会话列表失败，学生: %s, err: %v", studentId, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 500, "message": "获取会话列表失败"},
			"data":  nil,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{"code": 0},
		"data":  resp,
	})
}

// studentAISessionMessagesHandler 获取某会话的消息列表
// GET /student/ai/sessions/{sessionId}/messages?page=1&page_size=100
func studentAISessionMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()
	studentId, ok := authen.GetRoleIDFromContext(ctx)
	if !ok || studentId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 401, "message": "未授权"},
			"data":  nil,
		})
		return
	}

	vars := mux.Vars(r)
	sessionId := vars["sessionId"]
	if sessionId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 400, "message": "sessionId 不能为空"},
			"data":  nil,
		})
		return
	}

	page := int32(1)
	pageSize := int32(100)
	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	svcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := rpc.GetSessionClient().ListConversations(svcCtx, &conversationpb.ListConversationsRequest{
		SessionId: sessionId,
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		log.Printf("[conversation] 获取会话消息失败，sessionId: %s, err: %v", sessionId, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 500, "message": "获取会话消息失败"},
			"data":  nil,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{"code": 0},
		"data":  resp,
	})
}

// buildImageModerationRequest 根据图片URL构建图片审核请求
func buildImageModerationRequest(imgURL, userId, sessionId string) *rpc.ImageModerationRequest {
	imgReq := &rpc.ImageModerationRequest{
		BizType: "default",
		UserId:  userId,
		DataId:  sessionId,
	}

	if strings.HasPrefix(imgURL, "http") {
		imgReq.FileUrl = imgURL
	} else if strings.HasPrefix(imgURL, "data:image/") {
		parts := strings.SplitN(imgURL, ",", 2)
		if len(parts) == 2 {
			decoded, err := base64Decode(parts[1])
			if err != nil {
				log.Printf("[ContentModeration] 图片base64解码失败，userId: %s, err: %v", userId, err)
				return nil
			}
			imgReq.FileContent = decoded
		}
	} else {
		return nil
	}

	if imgReq.FileUrl == "" && len(imgReq.FileContent) == 0 {
		return nil
	}

	return imgReq
}

// extractImageURLs 从消息内容中提取图片URL
// 支持格式：[IMAGE:data:image/...;base64,...] 和 [IMAGE:https://...]
func extractImageURLs(content string) []string {
	const imagePrefix = "[IMAGE:"
	const imageSuffix = "]"

	var urls []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, imagePrefix) && strings.HasSuffix(line, imageSuffix) {
			url := line[len(imagePrefix) : len(line)-len(imageSuffix)]
			if url != "" {
				urls = append(urls, url)
			}
		}
	}
	return urls
}

// stripImageContent 从消息内容中剥离图片数据，只保留纯文本
// 将 [IMAGE:...] 格式的行移除，返回纯文本内容
func stripImageContent(content string) string {
	const imagePrefix = "[IMAGE:"
	const imageSuffix = "]"

	var textLines []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// 跳过图片行
		if strings.HasPrefix(trimmed, imagePrefix) && strings.HasSuffix(trimmed, imageSuffix) {
			continue
		}
		textLines = append(textLines, line)
	}
	return strings.Join(textLines, "\n")
}

// base64Decode 解码 base64 字符串为原始字节
func base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// buildSystemPrompt 根据问题类型和题目信息构建系统提示词
func buildSystemPrompt(questionType string, problemInfo *ProblemContext, userCode string, userCodeLang string, judgeResult string, failedCases string) string {
	basePrompt := "你是一位专业的编程助教，擅长帮助学生理解算法和编程问题。请用清晰、易懂的方式回答学生的问题，可以给出思路提示，但不要直接给出完整答案，鼓励学生自己思考。回复格式要求：用自然段落组织内容，可以用加粗强调关键词、用列表梳理步骤、用代码块展示代码。"

	if questionType == "algorithm_problem" && problemInfo != nil {
		// 明确告知 AI 它已经拥有完整的题目上下文
		problemDesc := "\n\n【重要】你已经拥有学生当前正在做的题目的完整信息，不需要再向学生询问题目内容。请直接基于以下题目信息回答学生的问题。\n"
		problemDesc += "\n【当前题目信息】\n"
		problemDesc += "题目名称：" + problemInfo.Title + "\n"
		if problemInfo.Difficulty != "" {
			diffMap := map[string]string{"easy": "简单", "medium": "中等", "hard": "困难"}
			diff := diffMap[problemInfo.Difficulty]
			if diff == "" {
				diff = problemInfo.Difficulty
			}
			problemDesc += "难度：" + diff + "\n"
		}
		if problemInfo.Description != "" {
			problemDesc += "题目描述：" + problemInfo.Description + "\n"
		}
		if problemInfo.InputFormat != "" {
			problemDesc += "输入格式：" + problemInfo.InputFormat + "\n"
		}
		if problemInfo.OutputFormat != "" {
			problemDesc += "输出格式：" + problemInfo.OutputFormat + "\n"
		}
		if len(problemInfo.Tags) > 0 {
			tags := ""
			for i, t := range problemInfo.Tags {
				if i > 0 {
					tags += "、"
				}
				tags += t
			}
			problemDesc += "相关标签：" + tags + "\n"
		}
		problemDesc += "\n请结合以上题目信息，帮助学生理解题意、分析思路，但不要直接给出完整代码解答。"

		// 加入运行记录信息（判题结果和未通过用例）
		if judgeResult != "" {
			judgeLabel := judgeResult
			switch judgeResult {
			case "accepted":
				judgeLabel = "全部通过"
			case "partial_pass":
				judgeLabel = "部分通过"
			case "wrong_answer":
				judgeLabel = "答案错误"
			case "compile_error":
				judgeLabel = "编译错误"
			case "runtime_error":
				judgeLabel = "运行时错误"
			case "time_limit_exceeded":
				judgeLabel = "超时"
			}
			problemDesc += "\n\n【运行记录 - 判题结果：" + judgeLabel + "】\n"
			problemDesc += "学生已将运行记录加入对话，请直接结合以下判题信息分析问题，不要再向学生索要运行结果。\n"
			if failedCases != "" {
				problemDesc += "未通过的测试用例详情：\n```json\n" + failedCases + "\n```\n"
				problemDesc += "请根据以上未通过的测试用例，分析学生代码的问题所在，引导学生思考如何修复。\n"
			}
		}

		// 加入用户当前IDE中的代码作为上下文
		if userCode != "" {
			langLabel := userCodeLang
			if langLabel == "" {
				langLabel = "未知语言"
			}
			problemDesc += "\n\n【学生当前代码（" + langLabel + "）】\n你已经拥有学生的代码，不需要再向学生索要代码。\n```" + userCodeLang + "\n" + userCode + "\n```\n请严格遵守以下规则：\n1. 当学生的问题涉及代码分析（如'为什么没通过'、'哪里有问题'、'帮我看看'等），直接结合上面的代码给出针对性分析。\n2. 若学生问的是题目思路、算法原理等与代码无关的问题，不要主动提及学生的代码。\n3. 绝对不要说'请提供你的代码'或'请贴出你的代码'之类的话，因为你已经有了。"
		}

		return basePrompt + problemDesc
	}

	return basePrompt
}
