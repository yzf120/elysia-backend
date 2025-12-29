package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model"
	authPb "github.com/yzf120/elysia-backend/proto/auth"
	userPb "github.com/yzf120/elysia-backend/proto/user"
	"github.com/yzf120/elysia-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务
type AuthService struct {
	authDAO                 dao.AuthDAO
	userService             *UserService
	verificationCodeService *utils.VerificationCodeService
	smsClient               *utils.TencentSMSClient
}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{
		authDAO:                 dao.NewAuthDAO(),
		userService:             NewUserService(),
		verificationCodeService: utils.NewVerificationCodeService(),
		smsClient:               utils.NewTencentSMSClient(),
	}
}

// 验证码校验
func (s *AuthService) VerifyCode(ctx context.Context, phoneNumber, code, codeType string) error {
	if phoneNumber == "" || code == "" || codeType == "" {
		return errs.NewCommonError(errs.ErrBadRequest, errs.ErrorMessages[errs.ErrBadRequest])
	}
	err := s.verificationCodeService.VerifyCode(phoneNumber, code, codeType)
	if err != nil {
		return errs.NewCommonError(errs.ErrSmsCodeInCollect, errs.ErrorMessages[errs.ErrSmsCodeInCollect])
	}
	return nil
}

// Register 用户注册
func (s *AuthService) Register(ctx context.Context, req *authPb.RegisterRequest) (*authPb.UserInfo, string, error) {
	// 参数校验
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, err.Error())
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
		return nil, "", errs.NewCommonError(errs.ErrInternal, "创建用户失败: "+err.Error())
	}

	// 生成登录令牌
	token, err := s.generateToken(user.UserId)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "生成令牌失败: "+err.Error())
	}

	// 转换为认证用户信息
	userInfo := s.convertUserToAuthUserInfo(user)

	return userInfo, token, nil
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *authPb.LoginRequest) (*authPb.UserInfo, string, error) {
	// 参数校验
	if err := s.validateLoginRequest(req); err != nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, err.Error())
	}

	// 根据手机号查询用户（包含密码）
	userModel, hashedPassword, err := s.authDAO.GetUserByPhoneNumberWithPassword(req.PhoneNumber)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "查询用户失败: "+err.Error())
	}

	if userModel == nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "用户不存在")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "密码错误")
	}

	// 检查用户状态
	if userModel.Status != 2 {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "用户状态异常，无法登录")
	}

	// 生成登录令牌
	token, err := s.generateToken(userModel.UserId)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "生成令牌失败: "+err.Error())
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
		UserId:          model.UserId,
		UserName:        model.UserName,
		Email:           model.Email,
		Gender:          model.Gender,
		PhoneNumber:     model.PhoneNumber,
		WxMiniAppOpenId: model.WxMiniAppOpenId,
		ChineseName:     model.ChineseName,
		Status:          model.Status,
		CreateTime:      model.CreateTime.Format("2006-01-02 15:04:05"),
		ImageUrl:        model.ImageURL,
		RegisterSource:  model.RegisterSource,
		UserType:        model.UserType,
	}
}

// SendVerificationCode 发送验证码
func (s *AuthService) SendVerificationCode(ctx context.Context, phoneNumber, codeType string) error {
	// 参数校验
	if phoneNumber == "" {
		return errs.NewCommonError(errs.ErrBadRequest, "手机号不能为空")
	}
	if len(phoneNumber) != 11 {
		return errs.NewCommonError(errs.ErrBadRequest, "手机号格式不正确")
	}
	if codeType != "register" && codeType != "login" {
		return errs.NewCommonError(errs.ErrBadRequest, "验证码类型不正确")
	}

	// 检查发送频率（60秒内只能发送一次）
	canSend, waitTime, err := s.verificationCodeService.CheckSendFrequency(phoneNumber, 60*time.Second)
	if err != nil {
		return errs.NewCommonError(errs.ErrInternal, "检查发送频率失败: "+err.Error())
	}
	if !canSend {
		return errs.NewCommonError(errs.ErrBadRequest, fmt.Sprintf("发送过于频繁，请%d秒后再试", int(waitTime.Seconds())))
	}

	// 如果是注册，检查手机号是否已存在
	if codeType == "register" {
		existingUser, _ := s.userService.userDAO.GetUserByPhoneNumber(phoneNumber)
		if existingUser != nil {
			return errs.NewCommonError(errs.ErrBadRequest, "该手机号已注册")
		}
	}

	// 如果是登录，检查用户是否存在
	if codeType == "login" {
		existingUser, _ := s.userService.userDAO.GetUserByPhoneNumber(phoneNumber)
		if existingUser == nil {
			return errs.NewCommonError(errs.ErrBadRequest, "该手机号未注册")
		}
	}

	// 生成验证码
	code := utils.GenerateVerificationCode()

	// 保存验证码到Redis（5分钟有效期）
	if err := s.verificationCodeService.SaveVerificationCode(phoneNumber, code, codeType, 5*time.Minute); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "保存验证码失败: "+err.Error())
	}

	// 获取短信模板ID
	templateId := os.Getenv("TENCENT_SMS_TEMPLATE_ID")
	if templateId == "" {
		return errs.NewCommonError(errs.ErrInternal, "短信模板ID未配置")
	}

	// 发送短信
	if err := s.smsClient.SendVerificationCode(phoneNumber, code, templateId); err != nil {
		// 发送失败，删除Redis中的验证码
		s.verificationCodeService.DeleteVerificationCode(phoneNumber, codeType)
		return errs.NewCommonError(errs.ErrInternal, "发送短信失败: "+err.Error())
	}

	return nil
}

// RegisterWithSMS 手机号+验证码注册
func (s *AuthService) RegisterWithSMS(ctx context.Context, phoneNumber, code string) (*authPb.UserInfo, string, error) {
	// 验证验证码
	if err := s.verificationCodeService.VerifyCode(phoneNumber, code, "register"); err != nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "验证码验证失败: "+err.Error())
	}

	// 生成默认用户名（手机号）
	userName := "user_" + phoneNumber

	// 创建用户（不需要密码）
	createUserReq := &userPb.CreateUserRequest{
		UserName:       userName,
		PhoneNumber:    phoneNumber,
		RegisterSource: "sms",
	}

	user, err := s.userService.CreateUser(createUserReq)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "创建用户失败: "+err.Error())
	}

	// 生成登录令牌
	token, err := s.generateToken(user.UserId)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "生成令牌失败: "+err.Error())
	}

	// 转换为认证用户信息
	userInfo := s.convertUserToAuthUserInfo(user)

	return userInfo, token, nil
}

// LoginWithSMS 手机号+验证码登录
func (s *AuthService) LoginWithSMS(phoneNumber, code string) (*authPb.UserInfo, string, error) {
	// 验证验证码
	if err := s.verificationCodeService.VerifyCode(phoneNumber, code, "login"); err != nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "验证码验证失败: "+err.Error())
	}

	// 查询用户
	userModel, err := s.userService.userDAO.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "查询用户失败: "+err.Error())
	}

	if userModel == nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "用户不存在")
	}

	// 检查用户状态
	if userModel.Status != 2 {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "用户状态异常，无法登录")
	}

	// 生成登录令牌
	token, err := s.generateToken(userModel.UserId)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "生成令牌失败: "+err.Error())
	}

	// 转换为认证用户信息
	userInfo := s.convertModelToAuthUserInfo(userModel)

	return userInfo, token, nil
}

// LoginWithPassword 手机号+密码登录
func (s *AuthService) LoginWithPassword(ctx context.Context, phoneNumber, password string) (*authPb.UserInfo, string, error) {
	// 参数校验
	if phoneNumber == "" {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "手机号不能为空")
	}
	if password == "" {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "密码不能为空")
	}

	// 根据手机号查询用户（包含密码）
	userModel, hashedPassword, err := s.authDAO.GetUserByPhoneNumberWithPassword(phoneNumber)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "查询用户失败: "+err.Error())
	}

	if userModel == nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "用户不存在")
	}

	// 如果用户没有设置密码
	if hashedPassword == "" {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "该账号未设置密码，请使用验证码登录")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "密码错误")
	}

	// 检查用户状态
	if userModel.Status != 2 {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "用户状态异常，无法登录")
	}

	// 生成登录令牌
	token, err := s.generateToken(userModel.UserId)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "生成令牌失败: "+err.Error())
	}

	// 转换为认证用户信息
	userInfo := s.convertModelToAuthUserInfo(userModel)

	return userInfo, token, nil
}
