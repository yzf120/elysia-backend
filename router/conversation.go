package router

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/authen"
	agentpb "github.com/yzf120/elysia-backend/proto/agent"
	"github.com/yzf120/elysia-backend/rpc"
)

// AIChatRequest AI对话请求（来自前端）
type AIChatRequest struct {
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
			break
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
