package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/utils"
)

// SMSService 短信服务
type SMSService struct {
	studentDAO              dao.StudentDAO
	teacherDAO              dao.TeacherDAO
	adminUserDAO            dao.AdminUserDAO
	verificationCodeService *utils.VerificationCodeService
	smsClient               *utils.TencentSMSClient
}

// NewSMSService 创建短信服务
func NewSMSService() *SMSService {
	return &SMSService{
		studentDAO:              dao.NewStudentDAO(),
		teacherDAO:              dao.NewTeacherDAO(),
		adminUserDAO:            dao.NewAdminUserDAO(),
		verificationCodeService: utils.NewVerificationCodeService(),
		smsClient:               utils.NewTencentSMSClient(),
	}
}

// SendVerificationCode 发送验证码
// userType: student, teacher, admin
// codeType: register, login
func (s *SMSService) SendVerificationCode(ctx context.Context, phoneNumber, userType, codeType string) error {
	// 参数校验
	if phoneNumber == "" {
		return errs.NewCommonError(errs.ErrBadRequest, "手机号不能为空")
	}
	if len(phoneNumber) != 11 {
		return errs.NewCommonError(errs.ErrBadRequest, "手机号格式不正确")
	}
	if userType != consts.RoleStudent && userType != consts.RoleTeacher && userType != consts.RoleAdmin {
		return errs.NewCommonError(errs.ErrBadRequest, "用户类型不正确")
	}
	if codeType != consts.Register && codeType != consts.Login {
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

	// 根据用户类型和验证码类型检查用户是否存在
	codeKey := fmt.Sprintf("%s_%s", userType, codeType)

	if codeType == consts.Register {
		// 注册时检查手机号是否已存在
		var exists bool
		switch userType {
		case consts.RoleStudent:
			student, _ := s.studentDAO.GetStudentByPhoneNumber(phoneNumber)
			exists = student != nil
		case consts.RoleTeacher:
			teacher, _ := s.teacherDAO.GetTeacherByPhoneNumber(phoneNumber)
			exists = teacher != nil
		case consts.RoleAdmin:
			admin, _ := s.adminUserDAO.GetAdminUserByPhoneNumber(phoneNumber)
			exists = admin != nil
		}

		if exists {
			return errs.NewCommonError(errs.ErrBadRequest, "该手机号已注册")
		}
	} else if codeType == consts.Login {
		// 登录时检查用户是否存在
		var exists bool
		switch userType {
		case consts.RoleStudent:
			student, _ := s.studentDAO.GetStudentByPhoneNumber(phoneNumber)
			exists = student != nil
		case consts.RoleTeacher:
			teacher, _ := s.teacherDAO.GetTeacherByPhoneNumber(phoneNumber)
			exists = teacher != nil
		case consts.RoleAdmin:
			admin, _ := s.adminUserDAO.GetAdminUserByPhoneNumber(phoneNumber)
			exists = admin != nil
		}

		if !exists {
			return errs.NewCommonError(errs.ErrBadRequest, "该手机号未注册")
		}
	}

	// 生成验证码
	code := utils.GenerateVerificationCode()

	// 保存验证码到Redis（5分钟有效期）
	if err := s.verificationCodeService.SaveVerificationCode(phoneNumber, code, codeKey, 5*time.Minute); err != nil {
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
		s.verificationCodeService.DeleteVerificationCode(phoneNumber, codeKey)
		return errs.NewCommonError(errs.ErrInternal, "发送短信失败: "+err.Error())
	}

	return nil
}
