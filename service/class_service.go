package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	classModel "github.com/yzf120/elysia-backend/model/class"
)

// ClassService 班级服务
type ClassService struct {
	classDAO       dao.ClassDAO
	classMemberDAO dao.ClassMemberDAO
	teacherDAO     dao.TeacherDAO
	studentDAO     dao.StudentDAO
}

// NewClassService 创建班级服务
func NewClassService() *ClassService {
	return &ClassService{
		classDAO:       dao.NewClassDAO(),
		classMemberDAO: dao.NewClassMemberDAO(),
		teacherDAO:     dao.NewTeacherDAO(),
		studentDAO:     dao.NewStudentDAO(),
	}
}

// CreateClass 创建班级（教师操作）
func (s *ClassService) CreateClass(teacherId, className, subject, semester, description string, maxStudents int32) (*classModel.Class, error) {
	// 参数校验
	if teacherId == "" || className == "" || subject == "" || semester == "" {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "必填参数不能为空")
	}

	// 检查教师是否存在
	teacher, err := s.teacherDAO.GetTeacherById(teacherId)
	if err != nil || teacher == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "教师信息不存在")
	}

	// 检查教师状态
	if teacher.Status != 1 {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "教师账号未激活")
	}

	// 生成班级ID和验证码
	classId := fmt.Sprintf("cls_%d", time.Now().UnixNano())
	classCode := s.generateClassCode()

	// 构建班级模型
	class := &classModel.Class{
		ClassId:         classId,
		ClassName:       className,
		ClassCode:       classCode,
		TeacherId:       teacherId,
		Subject:         subject,
		Semester:        semester,
		MaxStudents:     maxStudents,
		CurrentStudents: 0,
		Description:     description,
		Status:          1, // 进行中
	}

	// 创建班级
	if err := s.classDAO.CreateClass(class); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "创建班级失败: "+err.Error())
	}

	return class, nil
}

// generateClassCode 生成班级验证码（6位数字+字母）
func (s *ClassService) generateClassCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// JoinClass 学生加入班级
func (s *ClassService) JoinClass(studentId, classCode string) error {
	// 查询班级
	class, err := s.classDAO.GetClassByCode(classCode)
	if err != nil || class == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "班级不存在或验证码错误")
	}

	// 检查班级状态
	if class.Status != 1 {
		return errs.NewCommonError(errs.ErrBadRequest, "班级已结束或已归档")
	}

	// 检查学生是否存在
	student, err := s.studentDAO.GetStudentById(studentId)
	if err != nil || student == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "学生信息不存在")
	}

	// 检查是否已加入
	existingMember, _ := s.classMemberDAO.GetMember(class.ClassId, studentId)
	if existingMember != nil && existingMember.Status == 1 {
		return errs.NewCommonError(errs.ErrBadRequest, "已加入该班级")
	}

	// 检查班级人数是否已满
	if class.CurrentStudents >= class.MaxStudents {
		return errs.NewCommonError(errs.ErrBadRequest, "班级人数已满")
	}

	// 添加班级成员
	member := &classModel.ClassMember{
		ClassId:   class.ClassId,
		StudentId: studentId,
		Status:    1,
	}

	if err := s.classMemberDAO.AddMember(member); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "加入班级失败: "+err.Error())
	}

	// 更新班级人数
	updates := map[string]interface{}{
		"current_students": class.CurrentStudents + 1,
	}
	if err := s.classDAO.UpdateClass(class.ClassId, updates); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "更新班级人数失败: "+err.Error())
	}

	return nil
}

// LeaveClass 学生退出班级
func (s *ClassService) LeaveClass(studentId, classId string) error {
	// 查询班级成员
	member, err := s.classMemberDAO.GetMember(classId, studentId)
	if err != nil || member == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "未加入该班级")
	}

	// 更新成员状态为已退出
	if err := s.classMemberDAO.UpdateMemberStatus(classId, studentId, 0); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "退出班级失败: "+err.Error())
	}

	// 更新班级人数
	class, err := s.classDAO.GetClassById(classId)
	if err == nil && class != nil {
		updates := map[string]interface{}{
			"current_students": class.CurrentStudents - 1,
		}
		s.classDAO.UpdateClass(classId, updates)
	}

	return nil
}

// RemoveStudent 教师移除学生
func (s *ClassService) RemoveStudent(teacherId, classId, studentId string) error {
	// 查询班级
	class, err := s.classDAO.GetClassById(classId)
	if err != nil || class == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "班级不存在")
	}

	// 检查是否是班级创建者
	if class.TeacherId != teacherId {
		return errs.NewCommonError(errs.ErrBadRequest, "无权限操作")
	}

	// 移除学生
	if err := s.classMemberDAO.RemoveMember(classId, studentId); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "移除学生失败: "+err.Error())
	}

	// 更新班级人数
	updates := map[string]interface{}{
		"current_students": class.CurrentStudents - 1,
	}
	if err := s.classDAO.UpdateClass(classId, updates); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "更新班级人数失败: "+err.Error())
	}

	return nil
}

// GetClassMembers 获取班级成员列表
func (s *ClassService) GetClassMembers(classId string, page, pageSize int32) ([]*classModel.ClassMember, int32, error) {
	// 参数校验和默认值设置
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 查询班级成员
	offset := (page - 1) * pageSize
	members, err := s.classMemberDAO.ListMembersByClassId(classId, pageSize, offset)
	if err != nil {
		return nil, 0, errs.NewCommonError(errs.ErrInternal, "查询班级成员失败: "+err.Error())
	}

	// 统计总数
	total, err := s.classMemberDAO.CountMembersByClassId(classId)
	if err != nil {
		return nil, 0, errs.NewCommonError(errs.ErrInternal, "统计班级成员失败: "+err.Error())
	}

	return members, total, nil
}

// GetStudentClasses 获取学生加入的班级列表
func (s *ClassService) GetStudentClasses(studentId string, page, pageSize int32) ([]*classModel.ClassMember, int32, error) {
	// 参数校验和默认值设置
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 查询学生班级
	offset := (page - 1) * pageSize
	members, err := s.classMemberDAO.ListClassesByStudentId(studentId, pageSize, offset)
	if err != nil {
		return nil, 0, errs.NewCommonError(errs.ErrInternal, "查询学生班级失败: "+err.Error())
	}

	// 这里简化处理，实际应该统计总数
	total := int32(len(members))

	return members, total, nil
}

// GetTeacherClasses 获取教师创建的班级列表
func (s *ClassService) GetTeacherClasses(teacherId string, page, pageSize int32) ([]*classModel.Class, int32, error) {
	// 参数校验和默认值设置
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 查询教师班级
	offset := (page - 1) * pageSize
	classes, err := s.classDAO.ListClassesByTeacherId(teacherId, pageSize, offset)
	if err != nil {
		return nil, 0, errs.NewCommonError(errs.ErrInternal, "查询教师班级失败: "+err.Error())
	}

	// 这里简化处理，实际应该统计总数
	total := int32(len(classes))

	return classes, total, nil
}

// UpdateClass 更新班级信息
func (s *ClassService) UpdateClass(teacherId, classId string, updates map[string]interface{}) (*classModel.Class, error) {
	// 查询班级
	class, err := s.classDAO.GetClassById(classId)
	if err != nil || class == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "班级不存在")
	}

	// 检查是否是班级创建者
	if class.TeacherId != teacherId {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "无权限操作")
	}

	// 执行更新
	if err := s.classDAO.UpdateClass(classId, updates); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "更新班级信息失败: "+err.Error())
	}

	// 查询更新后的班级信息
	updatedClass, err := s.classDAO.GetClassById(classId)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "查询班级信息失败: "+err.Error())
	}

	return updatedClass, nil
}

// GetClassByCode 根据验证码获取班级信息
func (s *ClassService) GetClassByCode(classCode string) (*classModel.Class, error) {
	class, err := s.classDAO.GetClassByCode(classCode)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "查询班级失败: "+err.Error())
	}
	if class == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "班级不存在")
	}
	return class, nil
}
