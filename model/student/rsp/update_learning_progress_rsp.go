package rsp

// UpdateLearningProgressResponse 更新学习进度响应
type UpdateLearningProgressResponse struct {
	Code    int32  `json:"code"`    // 响应码 0-成功 其他-失败
	Message string `json:"message"` // 响应消息
}