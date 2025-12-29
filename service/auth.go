package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/model"
	authPb "github.com/yzf120/elysia-backend/proto/auth"
	userPb "github.com/yzf120/elysia-backend/proto/user"
	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务
type AuthService struct {
	authDAO     dao.AuthDAO
	userService *UserService
}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{
		authDAO:     dao.NewAuthDAO(),
		userService: NewUserService(),
	}
}

// Register 用户注册
func (s *AuthService) Register(req *authPb.RegisterRequest) (*authPb.UserInfo, string, error) {
	// 参数校验
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, "", err
	}

	// 调用用户创建接口
	createUserReq := &userPb.CreateUserRequest{
		UserName:        req.UserName,
		Password:        req.Password,
		PhoneNumber:     req.PhoneNumber,
		Email:           req.Email,
		Gender:          req.Gender,
		ChineseName:     req.ChineseName,
		WxMiniAppOpenId: req.WxMiniAppOpenId,
		ImageUrl:        req.ImageUrl,
		UserType:        req.UserType,
		RegisterSource:  "phone", // 默认注册来源为手机号
	}

	// 创建用户
	user, err := s.userService.CreateUser(createUserReq)
	if err != nil {
		return nil, "", err
	}

	// 生成登录令牌
	token, err := s.generateToken(user.UserId)
	if err != nil {
		return nil, "", fmt.Errorf("生成令牌失败: %v", err)
	}

	// 转换为认证用户信息
	userInfo := s.convertUserToAuthUserInfo(user)

	return userInfo, token, nil
}

// Login 用户登录
func (s *AuthService) Login(req *authPb.LoginRequest) (*authPb.UserInfo, string, error) {
	// 参数校验
	if err := s.validateLoginRequest(req); err != nil {
		return nil, "", err
	}

	// 根据手机号查询用户（包含密码）
	userModel, hashedPassword, err := s.authDAO.GetUserByPhoneNumberWithPassword(req.PhoneNumber)
	if err != nil {
		return nil, "", fmt.Errorf("查询用户失败: %v", err)
	}

	if userModel == nil {
		return nil, "", fmt.Errorf("用户不存在")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return nil, "", fmt.Errorf("密码错误")
	}

	// 检查用户状态
	if userModel.Status != 2 {
		return nil, "", fmt.Errorf("用户状态异常，无法登录")
	}

	// 生成登录令牌
	token, err := s.generateToken(userModel.UserID)
	if err != nil {
		return nil, "", fmt.Errorf("生成令牌失败: %v", err)
	}

	// 转换为认证用户信息
	userInfo := s.convertModelToAuthUserInfo(userModel)

	return userInfo, token, nil
}

// validateRegisterRequest 校验注册请求
func (s *AuthService) validateRegisterRequest(req *authPb.RegisterRequest) error {
	if req.UserName == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if req.Password == "" {
		return fmt.Errorf("密码不能为空")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("密码长度不能少于6位")
	}
	if req.PhoneNumber == "" {
		return fmt.Errorf("手机号不能为空")
	}
	// 可以添加手机号格式校验
	if len(req.PhoneNumber) != 11 {
		return fmt.Errorf("手机号格式不正确")
	}
	return nil
}

// validateLoginRequest 校验登录请求
func (s *AuthService) validateLoginRequest(req *authPb.LoginRequest) error {
	if req.PhoneNumber == "" {
		return fmt.Errorf("手机号不能为空")
	}
	if req.Password == "" {
		return fmt.Errorf("密码不能为空")
	}
	return nil
}

// generateToken 生成登录令牌（简单实现，实际应使用JWT）
func (s *AuthService) generateToken(userID string) (string, error) {
	// 生成32字节随机数
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// 转换为十六进制字符串
	token := hex.EncodeToString(b)
	// 实际应该将token存储到Redis等缓存中，并设置过期时间
	return token, nil
}

// convertUserToAuthUserInfo 将用户proto转换为认证用户信息
func (s *AuthService) convertUserToAuthUserInfo(user *userPb.User) *authPb.UserInfo {
	if user == nil {
		return nil
	}
	return &authPb.UserInfo{
		UserId:          user.UserId,
		UserName:        user.UserName,
		Email:           user.Email,
		Gender:          user.Gender,
		PhoneNumber:     user.PhoneNumber,
		WxMiniAppOpenId: user.WxMiniAppOpenId,
		ChineseName:     user.ChineseName,
		Status:          user.Status,
		CreateTime:      user.CreateTime,
		ImageUrl:        user.ImageUrl,
		RegisterSource:  user.RegisterSource,
		UserType:        user.UserType,
	}
}

// convertModelToAuthUserInfo 将数据模型转换为认证用户信息
func (s *AuthService) convertModelToAuthUserInfo(model *model.User) *authPb.UserInfo {
	if model == nil {
		return nil
	}
	return &authPb.UserInfo{
		UserId:          model.UserID,
		UserName:        model.UserName,
		Email:           model.Email,
		Gender:          model.Gender,
		PhoneNumber:     model.PhoneNumber,
		WxMiniAppOpenId: model.WxMiniAppOpenID,
		ChineseName:     model.ChineseName,
		Status:          model.Status,
		CreateTime:      model.CreateTime.Format("2006-01-02 15:04:05"),
		ImageUrl:        model.ImageURL,
		RegisterSource:  model.RegisterSource,
		UserType:        model.UserType,
	}
}
