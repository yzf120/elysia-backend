package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/teacher"
	"github.com/yzf120/elysia-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

// TeacherService 教师服务
type TeacherService struct {
	teacherDAO              dao.TeacherDAO
	approvalDAO             dao.TeacherApprovalDAO
	verificationCodeService *utils.VerificationCodeService
	jwtService              *utils.JWTService
}

// NewTeacherService 创建教师服务
func NewTeacherService() *TeacherService {
	return &TeacherService{
		teacherDAO:              dao.NewTeacherDAO(),
		approvalDAO:             dao.NewTeacherApprovalDAO(),
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
	existingTeacher, _ = s.teacherDAO.GetTeacherByPhoneNumber(phoneNumber)
	if existingTeacher != nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "手机号已被注册")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "密码加密失败")
	}

	// 创建教师信息
	teacherId := fmt.Sprintf("tea_%d", time.Now().UnixNano())
	teachingSubjectsJSON, _ := json.Marshal(teachingSubjects)

	t := &teacher.Teacher{
		TeacherId:          teacherId,
		PhoneNumber:        phoneNumber,
		Password:           string(hashedPassword),
		TeacherName:        realName,
		EmployeeNumber:     employeeNumber,
		SchoolEmail:        schoolEmail,
		TeachingSubjects:   string(teachingSubjectsJSON),
		Department:         department,
		VerificationStatus: 0, // 待审核
		Status:             0, // 未激活
	}

	if err := s.teacherDAO.CreateTeacher(t); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "创建教师信息失败: "+err.Error())
	}

	// 创建审批单
	approvalId := fmt.Sprintf("APV%d", time.Now().UnixNano())
	approval := &teacher.TeacherApproval{
		ApprovalId:       approvalId,
		TeacherId:        teacherId,
		EmployeeNumber:   employeeNumber,
		SchoolEmail:      schoolEmail,
		TeacherName:      realName,
		Phone:            phoneNumber,
		Department:       department,
		TeachingSubjects: string(teachingSubjectsJSON),
		ApprovalStatus:   0, // 待审批
	}

	if err := s.approvalDAO.CreateApproval(approval); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "创建审批单失败: "+err.Error())
	}

	return t, nil
}

// VerifyTeacher 审核教师（管理员操作）
func (s *TeacherService) VerifyTeacher(teacherId, verifierId string, approved bool, remark string) error {
	// 查询教师信息
	t, err := s.teacherDAO.GetTeacherById(teacherId)
	if err != nil || t == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "教师信息不存在")
	}

	// 检查审核状态
	if t.VerificationStatus != 0 {
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

// GetTeacherById 根据教师ID获取教师信息
func (s *TeacherService) GetTeacherById(teacherId string) (*teacher.Teacher, error) {
	t, err := s.teacherDAO.GetTeacherById(teacherId)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "查询教师信息失败: "+err.Error())
	}
	if t == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "教师信息不存在")
	}
	return t, nil
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
