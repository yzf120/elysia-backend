package req

// ListApprovalsRequest 查询审批单列表请求
type ListApprovalsRequest struct {
	Page           int32  `json:"page"`            // 页码
	PageSize       int32  `json:"page_size"`       // 每页数量
	ApprovalStatus *int32 `json:"approval_status"` // 审批状态（0-待审批 1-审批通过 2-审批驳回）
	Department     string `json:"department"`      // 部门筛选
	TeacherName    string `json:"teacher_name"`    // 教师姓名模糊搜索
}
