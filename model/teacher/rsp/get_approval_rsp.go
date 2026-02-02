package rsp

import "github.com/yzf120/elysia-backend/model/teacher"

// GetApprovalResponse 获取审批单响应
type GetApprovalResponse struct {
	Code     int32                    `json:"code"`
	Message  string                   `json:"message"`
	Approval *teacher.TeacherApproval `json:"approval,omitempty"`
}
