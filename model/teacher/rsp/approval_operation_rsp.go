package rsp

// ApproveTeacherResponse 审批教师响应
type ApproveTeacherResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// DeleteApprovalResponse 删除审批单响应
type DeleteApprovalResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
