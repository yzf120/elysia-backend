package rsp

// UpdateStudentResponse 更新学生信息响应
type UpdateStudentResponse struct {
	Code    int32  `json:"code"`    // 响应码 0-成功 其他-失败
	Message string `json:"message"` // 响应消息
}