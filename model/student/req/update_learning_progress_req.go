package req

// UpdateLearningProgressRequest 更新学习进度请求
type UpdateLearningProgressRequest struct {
	StudentId string                 `json:"student_id"` // 学生ID（必填）
	Progress  map[string]interface{} `json:"progress"`   // 学习进度（必填）
}