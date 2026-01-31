package rsp

// GetTeacherResponse 获取教师信息响应
type GetTeacherResponse struct {
	Code    int32        `json:"code"`    // 响应码 0-成功 其他-失败
	Message string       `json:"message"` // 响应消息
	Teacher *TeacherInfo `json:"teacher"` // 教师信息
}