package service_impl

import (
	"context"

	pb "github.com/yzf120/elysia-backend/proto/user"
	"github.com/yzf120/elysia-backend/service"
)

// UserServiceImpl 用户服务实现（只做出入参处理）
type UserServiceImpl struct {
	userService *service.UserService
}

// NewUserServiceImpl 创建用户服务实现
func NewUserServiceImpl() *UserServiceImpl {
	return &UserServiceImpl{
		userService: service.NewUserService(),
	}
}

// CreateUser 创建用户
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// 调用service层处理业务逻辑
	user, err := s.userService.CreateUser(req)
	if err != nil {
		return &pb.CreateUserResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &pb.CreateUserResponse{
		Code:    0,
		Message: "创建用户成功",
		User:    user,
	}, nil
}

// GetUser 获取用户信息
func (s *UserServiceImpl) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// 调用service层处理业务逻辑
	user, err := s.userService.GetUser(req)
	if err != nil {
		code := int32(500)
		if err.Error() == "用户不存在" {
			code = 404
		}
		return &pb.GetUserResponse{
			Code:    code,
			Message: err.Error(),
		}, nil
	}

	return &pb.GetUserResponse{
		Code:    0,
		Message: "查询成功",
		User:    user,
	}, nil
}

// UpdateUser 更新用户信息
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// 调用service层处理业务逻辑
	user, err := s.userService.UpdateUser(req)
	if err != nil {
		code := int32(500)
		if err.Error() == "用户不存在" {
			code = 404
		}
		return &pb.UpdateUserResponse{
			Code:    code,
			Message: err.Error(),
		}, nil
	}

	return &pb.UpdateUserResponse{
		Code:    0,
		Message: "更新成功",
		User:    user,
	}, nil
}

// DeleteUser 删除用户
func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// 调用service层处理业务逻辑
	err := s.userService.DeleteUser(req)
	if err != nil {
		code := int32(500)
		if err.Error() == "用户不存在" {
			code = 404
		}
		return &pb.DeleteUserResponse{
			Code:    code,
			Message: err.Error(),
		}, nil
	}

	return &pb.DeleteUserResponse{
		Code:    0,
		Message: "删除成功",
	}, nil
}

// ListUsers 查询用户列表
func (s *UserServiceImpl) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// 调用service层处理业务逻辑
	users, total, err := s.userService.ListUsers(req)
	if err != nil {
		return &pb.ListUsersResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &pb.ListUsersResponse{
		Code:     0,
		Message:  "查询成功",
		Users:    users,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
