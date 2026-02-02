package service_impl

import (
	"context"

	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/errs"
	adminPb "github.com/yzf120/elysia-backend/proto/admin"
	"github.com/yzf120/elysia-backend/service"
)

// AdminUserServiceImpl 管理员用户服务实现
type AdminUserServiceImpl struct {
	adminUserService *service.AdminUserService
}

// NewAdminUserServiceImpl 创建管理员用户服务实现
func NewAdminUserServiceImpl() *AdminUserServiceImpl {
	return &AdminUserServiceImpl{
		adminUserService: service.NewAdminUserService(),
	}
}

// CreateAdminUser 创建管理员用户
func (s *AdminUserServiceImpl) CreateAdminUser(ctx context.Context, req *adminPb.CreateAdminUserRequest) (*adminPb.CreateAdminUserResponse, error) {
	// 调用service层处理业务逻辑
	adminUserInfo, err := s.adminUserService.CreateAdminUser(req)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &adminPb.CreateAdminUserResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &adminPb.CreateAdminUserResponse{
		Code:      consts.SuccessCode,
		Message:   consts.MessageCreateUserSuccess,
		AdminUser: adminUserInfo,
	}, nil
}

// UpdateAdminUserPassword 更新管理员用户密码
func (s *AdminUserServiceImpl) UpdateAdminUserPassword(ctx context.Context, req *adminPb.UpdateAdminUserPasswordRequest) (*adminPb.UpdateAdminUserPasswordResponse, error) {
	// 调用service层处理密码更新逻辑
	err := s.adminUserService.UpdateAdminUserPassword(req.AdminId, req.OldPassword, req.NewPassword)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &adminPb.UpdateAdminUserPasswordResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &adminPb.UpdateAdminUserPasswordResponse{
		Code:    consts.SuccessCode,
		Message: consts.MessageUpdateSuccess,
	}, nil
}

// GetAdminUser 查询管理员用户信息
func (s *AdminUserServiceImpl) GetAdminUser(ctx context.Context, req *adminPb.GetAdminUserRequest) (*adminPb.GetAdminUserResponse, error) {
	// 根据管理员ID查询
	adminUserInfo, err := s.adminUserService.GetAdminUserByAdminId(req.AdminId)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &adminPb.GetAdminUserResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &adminPb.GetAdminUserResponse{
		Code:      consts.SuccessCode,
		Message:   consts.MessageQuerySuccess,
		AdminUser: adminUserInfo,
	}, nil
}

// ListAdminUsers 查询管理员用户列表
func (s *AdminUserServiceImpl) ListAdminUsers(ctx context.Context, req *adminPb.ListAdminUsersRequest) (*adminPb.ListAdminUsersResponse, error) {
	// 参数校验
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 调用service层查询列表
	adminUserInfos, total, err := s.adminUserService.ListAdminUsers(int(req.Page), int(req.PageSize), req.Role, req.Status)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &adminPb.ListAdminUsersResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &adminPb.ListAdminUsersResponse{
		Code:       consts.SuccessCode,
		Message:    consts.MessageQuerySuccess,
		AdminUsers: adminUserInfos,
		Total:      int32(total),
	}, nil
}

// UpdateAdminUserStatus 更新管理员用户状态
func (s *AdminUserServiceImpl) UpdateAdminUserStatus(ctx context.Context, req *adminPb.UpdateAdminUserStatusRequest) (*adminPb.UpdateAdminUserStatusResponse, error) {
	// 调用service层更新状态
	err := s.adminUserService.UpdateAdminUserStatus(req.AdminId, req.Status)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &adminPb.UpdateAdminUserStatusResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &adminPb.UpdateAdminUserStatusResponse{
		Code:    consts.SuccessCode,
		Message: consts.MessageUpdateSuccess,
	}, nil
}
