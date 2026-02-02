package req

// GetTeacherRequest 获取教师信息请求
type GetTeacherRequest struct {
	TeacherId string `json:"teacher_id"` // 教师ID（必填）
}
