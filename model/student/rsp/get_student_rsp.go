package rsp

// GetStudentResponse 获取学生信息响应
type GetStudentResponse struct {
	Code    int32        `json:"code"`    // 响应码 0-成功 其他-失败
	Message string       `json:"message"` // 响应消息
	Student *StudentInfo `json:"student"` // 学生信息
}