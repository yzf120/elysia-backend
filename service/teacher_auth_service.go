package service

import (
	"context"
	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/teacher"
	"github.com/yzf120/elysia-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

// TeacherAuthService 教师认证服务
type TeacherAuthService struct {
	teacherDAO              dao.TeacherDAO
	teacherApprovalDAO      dao.TeacherApprovalDAO
	verificationCodeService *utils.VerificationCodeService
	smsClient               *utils.TencentSMSClient
	jwtService              *utils.JWTService
}

// NewTeacherAuthService 创建教师认证服务
func NewTeacherAuthService() *TeacherAuthService {
	return &TeacherAuthService{
		teacherDAO:              dao.NewTeacherDAO(),
		teacherApprovalDAO:      dao.NewTeacherApprovalDAO(),
		verificationCodeService: utils.NewVerificationCodeService(),
		smsClient:               utils.NewTencentSMSClient(),
		jwtService:              utils.NewJWTService(),
	}
}

// LoginWithSMS 教师手机号+验证码登录
func (s *TeacherAuthService) LoginWithSMS(ctx context.Context, phoneNumber, code string) (*teacher.Teacher, string, error) {
	// 验证验证码
	if err := s.verificationCodeService.VerifyCode(phoneNumber, code, "teacher_login"); err != nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "验证码验证失败: "+err.Error())
	}

	// 查询教师
	teacherModel, err := s.teacherDAO.GetTeacherByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "查询教师失败: "+err.Error())
	}

	if teacherModel == nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "教师不存在")
	}

	// 检查教师状态
	if teacherModel.Status == 2 {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "账号已被禁用")
	}

	// 生成登录令牌
	token, err := s.jwtService.GenerateToken(teacherModel.TeacherId)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "生成令牌失败: "+err.Error())
	}

	return teacherModel, token, nil
}

// LoginWithPassword 教师手机号+密码登录
func (s *TeacherAuthService) LoginWithPassword(ctx context.Context, phoneNumber, password string) (*teacher.Teacher, string, error) {
	// 参数校验
	if phoneNumber == "" {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "手机号不能为空")
	}
	if password == "" {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "密码不能为空")
	}

	// 查询教师
	teacherModel, err := s.teacherDAO.GetTeacherByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "查询教师失败: "+err.Error())
	}

	if teacherModel == nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "教师不存在")
	}

	// 如果教师没有设置密码
	if teacherModel.Password == "" {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "该账号未设置密码，请使用验证码登录")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(teacherModel.Password), []byte(password)); err != nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "密码错误")
	}

	// 检查教师状态
	if teacherModel.Status == 2 {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "账号已被禁用")
	}

	// 生成登录令牌
	token, err := s.jwtService.GenerateToken(teacherModel.TeacherId)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "生成令牌失败: "+err.Error())
	}

	return teacherModel, token, nil
}
