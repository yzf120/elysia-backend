package req

// GetStudentRequest 获取学生信息请求
type GetStudentRequest struct {
	StudentId string `json:"student_id"` // 学生ID（必填）
}
