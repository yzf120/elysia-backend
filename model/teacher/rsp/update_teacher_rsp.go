package rsp

// UpdateTeacherResponse 更新教师信息响应
type UpdateTeacherResponse struct {
	Code    int32  `json:"code"`    // 响应码 0-成功 其他-失败
	Message string `json:"message"` // 响应消息
}