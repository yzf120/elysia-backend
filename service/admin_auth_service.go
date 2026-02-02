package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/admin"
	"github.com/yzf120/elysia-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

// AdminAuthService 管理员认证服务
type AdminAuthService struct {
	adminUserDAO            dao.AdminUserDAO
	verificationCodeService *utils.VerificationCodeService
	smsClient               *utils.TencentSMSClient
	jwtService              *utils.JWTService
}

// NewAdminAuthService 创建管理员认证服务
func NewAdminAuthService() *AdminAuthService {
	return &AdminAuthService{
		adminUserDAO:            dao.NewAdminUserDAO(),
		verificationCodeService: utils.NewVerificationCodeService(),
		smsClient:               utils.NewTencentSMSClient(),
		jwtService:              utils.NewJWTService(),
	}
}

// RegisterWithSMS 管理员手机号+验证码注册
func (s *AdminAuthService) RegisterWithSMS(ctx context.Context, phoneNumber, code string) (*admin.AdminUser, string, error) {
	// 验证验证码
	if err := s.verificationCodeService.VerifyCode(phoneNumber, code, "admin_register"); err != nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "验证码验证失败: "+err.Error())
	}

	// 检查手机号是否已注册
	existingAdmin, _ := s.adminUserDAO.GetAdminUserByPhoneNumber(phoneNumber)
	if existingAdmin != nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "该手机号已注册")
	}

	// 生成默认密码（需要首次登录后修改）
	defaultPassword := "Admin@123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "密码加密失败: "+err.Error())
	}

	// 创建管理员记录
	adminId := fmt.Sprintf("adm_%d", time.Now().UnixNano())
	newAdmin := &admin.AdminUser{
		AdminId:     adminId,
		Username:    fmt.Sprintf("admin_%s", phoneNumber),
		PhoneNumber: phoneNumber,
		Password:    string(hashedPassword),
		RealName:    fmt.Sprintf("管理员_%s", phoneNumber[len(phoneNumber)-4:]),
		Email:       fmt.Sprintf("admin_%s@admin.com", phoneNumber),
		Role:        consts.RoleAdmin,
		Status:      consts.AdminStatusInactive, // 注册时初始化为未激活
	}

	if err := s.adminUserDAO.CreateAdminUser(newAdmin); err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "创建管理员记录失败: "+err.Error())
	}

	// 生成登录令牌
	token, err := s.jwtService.GenerateToken(adminId)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "生成令牌失败: "+err.Error())
	}

	return newAdmin, token, nil
}

// LoginWithSMS 管理员手机号+验证码登录
func (s *AdminAuthService) LoginWithSMS(ctx context.Context, phoneNumber, code string) (*admin.AdminUser, string, error) {
	// 验证验证码
	if err := s.verificationCodeService.VerifyCode(phoneNumber, code, "admin_login"); err != nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "验证码验证失败: "+err.Error())
	}

	// 查询管理员
	adminModel, err := s.adminUserDAO.GetAdminUserByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "查询管理员失败: "+err.Error())
	}

	if adminModel == nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "管理员不存在")
	}

	// 检查管理员状态
	if adminModel.Status != 1 {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "账号已被禁用")
	}

	// 生成登录令牌
	token, err := s.jwtService.GenerateToken(adminModel.AdminId)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "生成令牌失败: "+err.Error())
	}

	// 更新登录信息
	s.adminUserDAO.UpdateAdminUserLoginInfo(adminModel.AdminId, "", time.Now())

	return adminModel, token, nil
}

// LoginWithPassword 管理员手机号+密码登录
func (s *AdminAuthService) LoginWithPassword(ctx context.Context, phoneNumber, password string) (*admin.AdminUser, string, error) {
	// 参数校验
	if phoneNumber == "" {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "手机号不能为空")
	}
	if password == "" {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "密码不能为空")
	}

	// 查询管理员
	adminModel, err := s.adminUserDAO.GetAdminUserByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "查询管理员失败: "+err.Error())
	}

	if adminModel == nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "管理员不存在")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(adminModel.Password), []byte(password)); err != nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "密码错误")
	}

	// 检查管理员状态
	if adminModel.Status != 1 {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "账号已被禁用")
	}

	// 生成登录令牌
	token, err := s.jwtService.GenerateToken(adminModel.AdminId)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "生成令牌失败: "+err.Error())
	}

	// 更新登录信息
	s.adminUserDAO.UpdateAdminUserLoginInfo(adminModel.AdminId, "", time.Now())

	return adminModel, token, nil
}
