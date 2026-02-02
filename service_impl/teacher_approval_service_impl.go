package service_impl

import (
	"context"

	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/teacher/req"
	"github.com/yzf120/elysia-backend/model/teacher/rsp"
	"github.com/yzf120/elysia-backend/service"
)

// TeacherApprovalServiceImpl 教师审批单服务实现（只做出入参处理）
type TeacherApprovalServiceImpl struct {
	approvalService *service.TeacherApprovalService
}

// NewTeacherApprovalServiceImpl 创建教师审批单服务实现
func NewTeacherApprovalServiceImpl() *TeacherApprovalServiceImpl {
	return &TeacherApprovalServiceImpl{
		approvalService: service.NewTeacherApprovalService(),
	}
}

// GetApprovalById 根据审批单ID获取审批单
func (s *TeacherApprovalServiceImpl) GetApprovalById(ctx context.Context, approvalId string) (*rsp.GetApprovalResponse, error) {
	// 调用service层处理业务逻辑
	approval, err := s.approvalService.GetApprovalById(approvalId)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.GetApprovalResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &rsp.GetApprovalResponse{
		Code:     consts.SuccessCode,
		Message:  consts.MessageQuerySuccess,
		Approval: approval,
	}, nil
}

// GetApprovalByTeacherId 根据教师ID获取审批单
func (s *TeacherApprovalServiceImpl) GetApprovalByTeacherId(ctx context.Context, teacherId string) (*rsp.GetApprovalResponse, error) {
	// 调用service层处理业务逻辑
	approval, err := s.approvalService.GetApprovalByTeacherId(teacherId)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.GetApprovalResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &rsp.GetApprovalResponse{
		Code:     consts.SuccessCode,
		Message:  consts.MessageQuerySuccess,
		Approval: approval,
	}, nil
}

// ListApprovals 查询审批单列表
func (s *TeacherApprovalServiceImpl) ListApprovals(ctx context.Context, request *req.ListApprovalsRequest) (*rsp.ListApprovalsResponse, error) {
	// 构建过滤条件
	filters := make(map[string]interface{})
	if request.ApprovalStatus != nil {
		filters["approval_status"] = *request.ApprovalStatus
	}
	if request.Department != "" {
		filters["department"] = request.Department
	}
	if request.TeacherName != "" {
		filters["teacher_name"] = request.TeacherName
	}

	// 调用service层处理业务逻辑
	approvals, total, err := s.approvalService.ListApprovals(request.Page, request.PageSize, filters)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.ListApprovalsResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &rsp.ListApprovalsResponse{
		Code:      consts.SuccessCode,
		Message:   consts.MessageQuerySuccess,
		Approvals: approvals,
		Total:     total,
		Page:      request.Page,
		PageSize:  request.PageSize,
	}, nil
}

// ApproveTeacher 审批教师
func (s *TeacherApprovalServiceImpl) ApproveTeacher(ctx context.Context, approvalId, adminId string, request *req.ApproveTeacherRequest) (*rsp.ApproveTeacherResponse, error) {
	// 调用service层处理业务逻辑
	err := s.approvalService.ApproveTeacher(approvalId, adminId, request.Approved, request.Remark)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.ApproveTeacherResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &rsp.ApproveTeacherResponse{
		Code:    consts.SuccessCode,
		Message: "审批成功",
	}, nil
}

// DeleteApproval 删除审批单
func (s *TeacherApprovalServiceImpl) DeleteApproval(ctx context.Context, approvalId, adminId string) (*rsp.DeleteApprovalResponse, error) {
	// 调用service层处理业务逻辑
	err := s.approvalService.DeleteApproval(approvalId, adminId)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.DeleteApprovalResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &rsp.DeleteApprovalResponse{
		Code:    consts.SuccessCode,
		Message: "删除成功",
	}, nil
}
