package agent

// AgentChatMessage 对话消息
type AgentChatMessage struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"` // 消息内容
}

// AgentStreamChatRequest 流式对话请求
type AgentStreamChatRequest struct {
	ModelID      string             `json:"model_id"`                // 模型ID，如 "doubao-seed-1-6-lite-251015"
	Messages     []AgentChatMessage `json:"messages"`                // 对话历史（含本次用户消息）
	SystemPrompt string             `json:"system_prompt,omitempty"` // 系统提示词（可选）
	ExtraParams  map[string]string  `json:"extra_params,omitempty"`  // 额外参数
}

// AgentStreamChatResponse 流式对话响应（每个chunk）
type AgentStreamChatResponse struct {
	Content          string `json:"content"`       // 本次增量内容
	IsEnd            bool   `json:"is_end"`        // 是否是最后一个chunk
	FinishReason     string `json:"finish_reason"` // 结束原因：stop, length, content_filter
	PromptTokens     int32  `json:"prompt_tokens"`
	CompletionTokens int32  `json:"completion_tokens"`
	TotalTokens      int32  `json:"total_tokens"`
}

// AgentListModelsRequest 查询模型列表请求
type AgentListModelsRequest struct {
	Provider string `json:"provider,omitempty"` // 可选：doubao, qwen，为空则返回全部
}

// AgentModelInfo 模型信息
type AgentModelInfo struct {
	ModelID       string `json:"model_id"`       // 模型ID
	ModelName     string `json:"model_name"`     // 模型名称
	Provider      string `json:"provider"`       // 提供商
	Description   string `json:"description"`    // 描述
	SupportStream bool   `json:"support_stream"` // 是否支持流式
	SupportVision bool   `json:"support_vision"` // 是否支持视觉/多模态
}

// AgentListModelsResponse 查询模型列表响应
type AgentListModelsResponse struct {
	Models []AgentModelInfo `json:"models"`
}
