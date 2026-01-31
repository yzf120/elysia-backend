package req

// GetStudentRequest 获取学生信息请求
type GetStudentRequest struct {
	UserId string `json:"user_id"` // 用户ID（必填）
}