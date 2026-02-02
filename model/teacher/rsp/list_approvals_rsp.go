package rsp

import "github.com/yzf120/elysia-backend/model/teacher"

// ListApprovalsResponse 查询审批单列表响应
type ListApprovalsResponse struct {
	Code      int32                      `json:"code"`
	Message   string                     `json:"message"`
	Approvals []*teacher.TeacherApproval `json:"approvals,omitempty"`
	Total     int32                      `json:"total"`
	Page      int32                      `json:"page"`
	PageSize  int32                      `json:"page_size"`
}
