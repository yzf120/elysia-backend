package rsp

// CreateStudentResponse 创建学生信息响应
type CreateStudentResponse struct {
	Code      int32  `json:"code"`       // 响应码 0-成功 其他-失败
	Message   string `json:"message"`    // 响应消息
	StudentId string `json:"student_id"` // 学生ID
}