package service_impl

import (
	"context"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/auth"
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

func (s *AuthServiceImpl) VerifyCode(ctx context.Context, req *auth.VerifyCodeRequest) (*auth.VerifyCodeResponse, error) {
	err := s.authService.VerifyCode(ctx, req.PhoneNumber, req.Code, req.CodeType)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &auth.VerifyCodeResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &auth.VerifyCodeResponse{
		Code:    0,
		Message: "验证码验证成功",
	}, nil
}

// Register 用户注册
func (s *AuthServiceImpl) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 调用service层处理业务逻辑
	userInfo, token, err := s.authService.Register(ctx, req)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &pb.RegisterResponse{
			Code:    int32(code),
			Message: msg,
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
	userInfo, token, err := s.authService.Login(ctx, req)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &pb.LoginResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &pb.LoginResponse{
		Code:    0,
		Message: "登录成功",
		User:    userInfo,
		Token:   token,
	}, nil
}

func (s *AuthServiceImpl) SendVerificationCode(ctx context.Context, req *auth.SendCodeRequest) (*auth.SendCodeResponse, error) {
	// 调用service层处理业务逻辑
	err := s.authService.SendVerificationCode(ctx, req.PhoneNumber, req.CodeType)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &auth.SendCodeResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &auth.SendCodeResponse{
		Code:    0,
		Message: "验证码发送成功",
	}, nil
}

func (s *AuthServiceImpl) RegisterWithSMS(ctx context.Context, req *auth.RegisterWithSMSRequest) (*auth.RegisterWithSMSResponse, error) {
	// 调用service层处理业务逻辑
	userInfo, token, err := s.authService.RegisterWithSMS(ctx, req.PhoneNumber, req.Code)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &auth.RegisterWithSMSResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &auth.RegisterWithSMSResponse{
		Code:     0,
		Message:  "注册成功",
		UserInfo: userInfo,
		Token:    token,
	}, nil
}

// LoginWithSMS 手机号+验证码登录
func (s *AuthServiceImpl) LoginWithSMS(ctx context.Context, req *auth.LoginWithSMSRequest) (*auth.LoginWithSMSResponse, error) {
	// 调用service层处理业务逻辑
	userInfo, token, err := s.authService.LoginWithSMS(req.PhoneNumber, req.Code)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &auth.LoginWithSMSResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &auth.LoginWithSMSResponse{
		Code:     0,
		Message:  "登录成功",
		UserInfo: userInfo,
		Token:    token,
	}, nil
}

// LoginWithPasswordHandler 手机号+密码登录处理器
func (s *AuthServiceImpl) LoginWithPassword(ctx context.Context, req *auth.LoginWithPasswordRequest) (*auth.LoginWithPasswordResponse, error) {
	// 调用service层处理业务逻辑
	userInfo, token, err := s.authService.LoginWithPassword(ctx, req.PhoneNumber, req.Password)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &auth.LoginWithPasswordResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}

	return &auth.LoginWithPasswordResponse{
		Code:     0,
		Message:  "登录成功",
		UserInfo: userInfo,
		Token:    token,
	}, nil
}
