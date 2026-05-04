package agent

// AgentChatMessage 对话消息
type AgentChatMessage struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"` // 消息内容
}

// AgentStreamChatRequest 流式对话请求
type AgentStreamChatRequest struct {
	ModelID      string             `json:"model_id"`                // 模型ID，如 "doubao-seed-2-0-lite-260215"
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

// ==================== 知识库管理 RPC 消息类型 ====================

// StoreKnowledgeRequest 存储知识条目请求
type StoreKnowledgeRequest struct {
	ID         string `json:"id"`          // 文档ID（为空则自动生成）
	Content    string `json:"content"`     // 文档内容
	SourceType string `json:"source_type"` // 来源类型：knowledge_base / problem_bank / error_pattern
	SourceID   string `json:"source_id"`   // 关联的知识点/题目 ID
	Tags       string `json:"tags"`        // 标签（逗号分隔）
	FileName   string `json:"file_name"`   // 原始文件名
	FileSize   int64  `json:"file_size"`   // 文件大小（字节）
	FileType   string `json:"file_type"`   // 文件类型
}

// StoreKnowledgeResponse 存储知识条目响应
type StoreKnowledgeResponse struct {
	Success          bool   `json:"success"`
	ID               string `json:"id"`                // 存储后的文档ID
	Message          string `json:"message"`           // 提示信息
	EstimatedSeconds int32  `json:"estimated_seconds"` // 预估索引构建时间（秒）
}

// DeleteKnowledgeRequest 删除知识条目请求
type DeleteKnowledgeRequest struct {
	ID string `json:"id"` // 文档ID
}

// DeleteKnowledgeResponse 删除知识条目响应
type DeleteKnowledgeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Async   bool   `json:"async"` // 是否为异步处理（分块和索引清理在后台进行）
}

// ListKnowledgeRequest 列出知识条目请求
type ListKnowledgeRequest struct {
	Page     int32 `json:"page"`      // 页码（从1开始）
	PageSize int32 `json:"page_size"` // 每页数量
}

// ListKnowledgeResponse 列出知识条目响应
type ListKnowledgeResponse struct {
	Total int32              `json:"total"` // 总数
	Items []KnowledgeDocItem `json:"items"` // 文档列表
}

// KnowledgeDocItem 知识文档条目
type KnowledgeDocItem struct {
	ID         string `json:"id"`
	Content    string `json:"content"`
	SourceType string `json:"source_type"`
	SourceID   string `json:"source_id"`
	Tags       string `json:"tags"`
	FileName   string `json:"file_name"`
	FileSize   int64  `json:"file_size"`
	FileType   string `json:"file_type"`
	Status     int    `json:"status"`
	CreateTime string `json:"create_time"`
}

// SearchKnowledgeRequest 检索知识请求
type SearchKnowledgeRequest struct {
	Query      string   `json:"query"`       // 检索查询
	Keywords   []string `json:"keywords"`    // 关键词
	TopK       int32    `json:"top_k"`       // 返回数量
	SourceType string   `json:"source_type"` // 限定来源类型
}

// SearchKnowledgeResponse 检索知识响应
type SearchKnowledgeResponse struct {
	Items []KnowledgeSearchResult `json:"items"`
}

// KnowledgeSearchResult 知识检索结果
type KnowledgeSearchResult struct {
	ID         string  `json:"id"`
	Content    string  `json:"content"`
	SourceType string  `json:"source_type"`
	SourceID   string  `json:"source_id"`
	Tags       string  `json:"tags"`
	Score      float64 `json:"score"` // 相关度评分
}
