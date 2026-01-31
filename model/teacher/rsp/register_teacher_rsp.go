package rsp

// RegisterTeacherResponse 教师注册响应
type RegisterTeacherResponse struct {
	Code      int32  `json:"code"`       // 响应码 0-成功 其他-失败
	Message   string `json:"message"`    // 响应消息
	TeacherId string `json:"teacher_id"` // 教师ID
}