package service_impl

import (
	"context"

	pb "github.com/yzf120/elysia-backend/proto/auth"
	"github.com/yzf120/elysia-backend/service"
)

// AuthServiceImpl 认证服务实现
type AuthServiceImpl struct {
	authService *service.AuthService
}

// NewAuthServiceImpl 创建认证服务实现
func NewAuthServiceImpl() *AuthServiceImpl {
	return &AuthServiceImpl{
		authService: service.NewAuthService(),
	}
}

// Register 用户注册
func (s *AuthServiceImpl) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 调用service层处理业务逻辑
	userInfo, token, err := s.authService.Register(req)
	if err != nil {
		return &pb.RegisterResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &pb.RegisterResponse{
		Code:    0,
		Message: "注册成功",
		User:    userInfo,
		Token:   token,
	}, nil
}

// Login 用户登录
func (s *AuthServiceImpl) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 调用service层处理业务逻辑
	userInfo, token, err := s.authService.Login(req)
	if err != nil {
		return &pb.LoginResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &pb.LoginResponse{
		Code:    0,
		Message: "登录成功",
		User:    userInfo,
		Token:   token,
	}, nil
}
