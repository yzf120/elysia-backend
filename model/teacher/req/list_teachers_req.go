package req

// ListTeachersRequest 查询教师列表请求
type ListTeachersRequest struct {
	Page               int32  `json:"page"`                // 页码（从1开始）
	PageSize           int32  `json:"page_size"`           // 每页数量
	Department         string `json:"department"`          // 院系筛选（可选）
	VerificationStatus int32  `json:"verification_status"` // 认证状态筛选（可选，-1表示不筛选）
	Status             int32  `json:"status"`              // 状态筛选（可选，-1表示不筛选）
}