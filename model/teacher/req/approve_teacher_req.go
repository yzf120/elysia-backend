package req

// ApproveTeacherRequest 审批教师请求
type ApproveTeacherRequest struct {
	Approved bool   `json:"approved"` // 是否通过
	Remark   string `json:"remark"`   // 审批意见
}
