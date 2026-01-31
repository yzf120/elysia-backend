package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model"
	"github.com/yzf120/elysia-backend/model/teacher"
	"github.com/yzf120/elysia-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

// TeacherService 教师服务
type TeacherService struct {
	teacherDAO              dao.TeacherDAO
	userDAO                 dao.UserDAO
	verificationCodeService *utils.VerificationCodeService
	jwtService              *utils.JWTService
}

// NewTeacherService 创建教师服务
func NewTeacherService() *TeacherService {
	return &TeacherService{
		teacherDAO:              dao.NewTeacherDAO(),
		userDAO:                 dao.NewUserDAO(),
		verificationCodeService: utils.NewVerificationCodeService(),
		jwtService:              utils.NewJWTService(),
	}
}

// RegisterTeacher 教师注册（工号+学校邮箱双重验证）
func (s *TeacherService) RegisterTeacher(ctx context.Context, phoneNumber, password, employeeNumber, schoolEmail, realName, department string, teachingSubjects []string) (*teacher.Teacher, error) {
	// 参数校验
	if phoneNumber == "" || password == "" || employeeNumber == "" || schoolEmail == "" {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "必填参数不能为空")
	}

	// 检查工号是否已存在
	existingTeacher, _ := s.teacherDAO.GetTeacherByEmployeeNumber(employeeNumber)
	if existingTeacher != nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "工号已被注册")
	}

	// 检查学校邮箱是否已存在
	existingTeacher, _ = s.teacherDAO.GetTeacherBySchoolEmail(schoolEmail)
	if existingTeacher != nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "学校邮箱已被注册")
	}

	// 检查手机号是否已存在
	existingUser, _ := s.userDAO.GetUserByPhoneNumber(phoneNumber)
	if existingUser != nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "手机号已被注册")
	}

	// 创建用户账号
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "密码加密失败")
	}

	userId := fmt.Sprintf("user_%d", time.Now().UnixNano())
	user := &model.User{
		UserId:         userId,
		UserName:       realName,
		Password:       string(hashedPassword),
		Email:          schoolEmail,
		PhoneNumber:    phoneNumber,
		ChineseName:    realName,
		RegisterSource: "teacher",
		UserType:       "teacher",
		Status:         2, // 可用状态
	}

	if err := s.userDAO.CreateUser(user); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "创建用户失败: "+err.Error())
	}

	// 创建教师信息
	teacherId := fmt.Sprintf("tea_%d", time.Now().UnixNano())
	teachingSubjectsJSON, _ := json.Marshal(teachingSubjects)

	teacher := &teacher.Teacher{
		TeacherId:          teacherId,
		UserId:             userId,
		EmployeeNumber:     employeeNumber,
		SchoolEmail:        schoolEmail,
		TeachingSubjects:   string(teachingSubjectsJSON),
		Department:         department,
		VerificationStatus: 0, // 待审核
		Status:             0, // 未激活
	}

	if err := s.teacherDAO.CreateTeacher(teacher); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "创建教师信息失败: "+err.Error())
	}

	return teacher, nil
}

// LoginTeacher 教师登录
func (s *TeacherService) LoginTeacher(ctx context.Context, phoneNumber, password string) (*teacher.Teacher, *model.User, string, error) {
	// 参数校验
	if phoneNumber == "" || password == "" {
		return nil, nil, "", errs.NewCommonError(errs.ErrBadRequest, "手机号和密码不能为空")
	}

	// 查询用户
	user, err := s.userDAO.GetUserByPhoneNumber(phoneNumber)
	if err != nil || user == nil {
		return nil, nil, "", errs.NewCommonError(errs.ErrBadRequest, "用户不存在")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, "", errs.NewCommonError(errs.ErrBadRequest, "密码错误")
	}

	// 检查用户类型
	if user.UserType != "teacher" {
		return nil, nil, "", errs.NewCommonError(errs.ErrBadRequest, "该账号不是教师账号")
	}

	// 查询教师信息
	teacher, err := s.teacherDAO.GetTeacherByUserId(user.UserId)
	if err != nil || teacher == nil {
		return nil, nil, "", errs.NewCommonError(errs.ErrBadRequest, "教师信息不存在")
	}

	// 检查教师状态
	if teacher.Status != 1 {
		return nil, nil, "", errs.NewCommonError(errs.ErrBadRequest, "教师账号未激活或已被禁用")
	}

	// 生成登录令牌
	token, err := s.jwtService.GenerateToken(user.UserId)
	if err != nil {
		return nil, nil, "", errs.NewCommonError(errs.ErrInternal, "生成令牌失败: "+err.Error())
	}

	return teacher, user, token, nil
}

// VerifyTeacher 审核教师（管理员操作）
func (s *TeacherService) VerifyTeacher(teacherId, verifierId string, approved bool, remark string) error {
	// 查询教师信息
	teacher, err := s.teacherDAO.GetTeacherById(teacherId)
	if err != nil || teacher == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "教师信息不存在")
	}

	// 检查审核状态
	if teacher.VerificationStatus != 0 {
		return errs.NewCommonError(errs.ErrBadRequest, "该教师已审核")
	}

	// 更新审核状态
	updates := map[string]interface{}{
		"verification_time":   time.Now(),
		"verifier_id":         verifierId,
		"verification_remark": remark,
	}

	if approved {
		updates["verification_status"] = 1 // 已通过
		updates["status"] = 1              // 激活账号
	} else {
		updates["verification_status"] = 2 // 已驳回
	}

	if err := s.teacherDAO.UpdateTeacher(teacherId, updates); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "更新教师信息失败: "+err.Error())
	}

	return nil
}

// GetTeacherByUserId 根据用户ID获取教师信息
func (s *TeacherService) GetTeacherByUserId(userId string) (*teacher.Teacher, error) {
	teacher, err := s.teacherDAO.GetTeacherByUserId(userId)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "查询教师信息失败: "+err.Error())
	}
	if teacher == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "教师信息不存在")
	}
	return teacher, nil
}

// UpdateTeacher 更新教师信息
func (s *TeacherService) UpdateTeacher(teacherId string, updates map[string]interface{}) (*teacher.Teacher, error) {
	// 检查教师是否存在
	existingTeacher, err := s.teacherDAO.GetTeacherById(teacherId)
	if err != nil || existingTeacher == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "教师信息不存在")
	}

	// 执行更新
	if err := s.teacherDAO.UpdateTeacher(teacherId, updates); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "更新教师信息失败: "+err.Error())
	}

	// 查询更新后的教师信息
	updatedTeacher, err := s.teacherDAO.GetTeacherById(teacherId)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "查询教师信息失败: "+err.Error())
	}

	return updatedTeacher, nil
}

// ListTeachers 查询教师列表
func (s *TeacherService) ListTeachers(page, pageSize int32, filters map[string]interface{}) ([]*teacher.Teacher, int32, error) {
	// 参数校验和默认值设置
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 构建查询条件
	whereClause := "1=1"
	var args []interface{}

	if department, ok := filters["department"].(string); ok && department != "" {
		whereClause += " AND department = ?"
		args = append(args, department)
	}
	if verificationStatus, ok := filters["verification_status"].(int32); ok {
		whereClause += " AND verification_status = ?"
		args = append(args, verificationStatus)
	}
	if status, ok := filters["status"].(int32); ok {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	// 查询总数
	total, err := s.teacherDAO.CountTeachers(whereClause, args)
	if err != nil {
		return nil, 0, errs.NewCommonError(errs.ErrInternal, "统计教师数量失败: "+err.Error())
	}

	// 查询列表
	offset := (page - 1) * pageSize
	teachers, err := s.teacherDAO.ListTeachers(whereClause, args, pageSize, offset)
	if err != nil {
		return nil, 0, errs.NewCommonError(errs.ErrInternal, "查询教师列表失败: "+err.Error())
	}

	return teachers, total, nil
}
