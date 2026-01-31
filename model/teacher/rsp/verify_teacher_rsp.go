package rsp

// VerifyTeacherResponse 审核教师响应
type VerifyTeacherResponse struct {
	Code    int32  `json:"code"`    // 响应码 0-成功 其他-失败
	Message string `json:"message"` // 响应消息
}