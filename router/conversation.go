package router

import (
	"context"
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
	agent_session "github.com/yzf120/elysia-session/proto/agent_session"
	conversationpb "github.com/yzf120/elysia-session/proto/conversation"
)

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

// studentAIModelsHandler 查询支持的模型列表
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
	systemPrompt := buildSystemPrompt(request.QuestionType, request.ProblemInfo, request.UserCode, request.UserCodeLang)

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

	// 透传深度思考参数
	if request.EnableThinking {
		agentReq.ExtraParams = map[string]string{
			"enable_thinking": "true",
		}
		log.Printf("[conversation] 深度思考模式已开启，模型: %s，学生: %s", modelID, studentId)
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
		chunk, err := agentStream.Recv()
		if err == io.EOF {
			// 流结束，发送结束事件
			fmt.Fprintf(w, "event: done\ndata: [DONE]\n\n")
			flusher.Flush()
			break
		}
		if err != nil {
			log.Printf("[conversation] 接收 chat-agent 响应失败: %v", err)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
			flusher.Flush()
			streamErr = true
			break
		}

		// 累积 AI 回复内容
		if chunk.Content != "" {
			aiReplyBuilder.WriteString(chunk.Content)
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

		fmt.Fprintf(w, "data: %s\n\n", string(chunkData))
		flusher.Flush()

		// 如果是最后一个 chunk，退出
		if chunk.IsEnd {
			fmt.Fprintf(w, "event: done\ndata: [DONE]\n\n")
			flusher.Flush()
			break
		}
	}

	log.Printf("[conversation] SSE 流式响应完成，学生: %s", studentId)

	// ===== 异步存储会话和对话记录到 session 服务 =====
	go func() {
		// 确定 AI 最终回复内容（若流异常则使用兜底回复）
		aiReply := aiReplyBuilder.String()
		if streamErr || strings.TrimSpace(aiReply) == "" {
			aiReply = "抱歉，AI 助教暂时无法回答，请稍后再试。"
			log.Printf("[conversation] AI 回复异常，使用兜底回复，学生: %s", studentId)
		}

		// 取本次用户消息（messages 最后一条 role=user 的消息）
		userMsg := ""
		for i := len(request.Messages) - 1; i >= 0; i-- {
			if request.Messages[i].Role == "user" {
				userMsg = request.Messages[i].Content
				break
			}
		}

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

// buildSystemPrompt 根据问题类型和题目信息构建系统提示词
func buildSystemPrompt(questionType string, problemInfo *ProblemContext, userCode string, userCodeLang string) string {
	basePrompt := "你是一位专业的编程助教，擅长帮助学生理解算法和编程问题。请用清晰、易懂的方式回答学生的问题，可以给出思路提示，但不要直接给出完整答案，鼓励学生自己思考。"

	if questionType == "algorithm_problem" && problemInfo != nil {
		problemDesc := "\n\n【当前题目信息】\n"
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

		// 加入用户当前IDE中的代码作为上下文
		if userCode != "" {
			langLabel := userCodeLang
			if langLabel == "" {
				langLabel = "未知语言"
			}
			problemDesc += "\n\n【学生当前代码（" + langLabel + "）】\n```" + userCodeLang + "\n" + userCode + "\n```\n注意：学生已主动选择将以上代码作为参考上下文。请严格遵守以下规则：\n1. 只有当学生的问题明确涉及代码（如'我的代码哪里有问题'、'帮我看看代码'、'为什么我的代码报错'等）时，才结合代码给出针对性回答。\n2. 若学生问的是题目思路、算法原理、时间复杂度等与代码无关的问题，绝对不要主动提及、引用或分析学生的代码。\n3. 判断标准：学生的问题中是否出现'我的代码'、'这段代码'、'代码'等明确指向代码的词语。"
		}

		return basePrompt + problemDesc
	}

	return basePrompt
}
