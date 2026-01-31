package req

// GetTeacherRequest 获取教师信息请求
type GetTeacherRequest struct {
	UserId string `json:"user_id"` // 用户ID（必填）
}