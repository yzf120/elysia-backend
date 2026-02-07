package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/student"
	"github.com/yzf120/elysia-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

// StudentAuthService 学生认证服务
type StudentAuthService struct {
	studentDAO              dao.StudentDAO
	verificationCodeService *utils.VerificationCodeService
	jwtService              *utils.JWTService
}

// NewStudentAuthService 创建学生认证服务
func NewStudentAuthService() *StudentAuthService {
	return &StudentAuthService{
		studentDAO:              dao.NewStudentDAO(),
		verificationCodeService: utils.NewVerificationCodeService(),
		jwtService:              utils.NewJWTService(),
	}
}

// RegisterWithSMS 学生手机号+验证码注册（需要学号和密码）
func (s *StudentAuthService) RegisterWithSMS(ctx context.Context, phoneNumber, code, studentNumber, password string) (*student.Student, error) {
	// 参数校验
	if studentNumber == "" || password == "" {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "学号和密码不能为空")
	}

	// 验证验证码
	if err := s.verificationCodeService.VerifyCode(phoneNumber, code, "student_register"); err != nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "验证码验证失败: "+err.Error())
	}

	// 检查手机号是否已注册
	existingStudent, _ := s.studentDAO.GetStudentByPhoneNumber(phoneNumber)
	if existingStudent != nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "该手机号已注册")
	}

	// 检查学号是否已存在
	existingStudentByNumber, _ := s.studentDAO.GetStudentByStudentNumber(studentNumber)
	if existingStudentByNumber != nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "该学号已被注册")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "密码加密失败")
	}

	// 创建学生记录
	studentId := fmt.Sprintf("stu_%d", time.Now().UnixNano())
	newStudent := &student.Student{
		StudentId:        studentId,
		StudentNumber:    studentNumber,
		PhoneNumber:      phoneNumber,
		Password:         string(hashedPassword),
		StudentName:      fmt.Sprintf("学生_%s", phoneNumber[len(phoneNumber)-4:]),
		Major:            "未设置",
		Grade:            "未设置",
		ProgrammingLevel: "初级",
		Interests:        "[]",
		LearningTags:     "[]",
		Status:           1, // 正常状态
	}

	if err := s.studentDAO.CreateStudent(newStudent); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "创建学生记录失败: "+err.Error())
	}

	return newStudent, nil
}

// LoginWithSMS 学生手机号+验证码登录
func (s *StudentAuthService) LoginWithSMS(ctx context.Context, phoneNumber, code string) (*student.Student, string, error) {
	// 验证验证码
	if err := s.verificationCodeService.VerifyCode(phoneNumber, code, "student_login"); err != nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "验证码验证失败: "+err.Error())
	}

	// 查询学生
	studentModel, err := s.studentDAO.GetStudentByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "查询学生失败: "+err.Error())
	}

	if studentModel == nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "学生不存在")
	}

	// 检查学生状态
	if studentModel.Status != 1 {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "账号已被禁用")
	}

	// 生成登录令牌
	token, err := s.jwtService.GenerateToken(studentModel.StudentId)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "生成令牌失败: "+err.Error())
	}

	return studentModel, token, nil
}

// LoginWithPassword 学生学号+密码登录
func (s *StudentAuthService) LoginWithPassword(ctx context.Context, studentNumber, password string) (*student.Student, string, error) {
	// 参数校验
	if studentNumber == "" {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "学号不能为空")
	}
	if password == "" {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "密码不能为空")
	}

	// 查询学生
	studentModel, err := s.studentDAO.GetStudentByStudentNumber(studentNumber)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "查询学生失败: "+err.Error())
	}

	if studentModel == nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "学生不存在")
	}

	// 如果学生没有设置密码
	if studentModel.Password == "" {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "该账号未设置密码，请使用验证码登录")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(studentModel.Password), []byte(password)); err != nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "密码错误")
	}

	// 检查学生状态
	if studentModel.Status != 1 {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "账号已被禁用")
	}

	// 生成登录令牌
	token, err := s.jwtService.GenerateToken(studentModel.StudentId)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "生成令牌失败: "+err.Error())
	}

	return studentModel, token, nil
}
