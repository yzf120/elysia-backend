package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/student"
)

// StudentService 学生服务
type StudentService struct {
	studentDAO dao.StudentDAO
	userDAO    dao.UserDAO
}

// NewStudentService 创建学生服务
func NewStudentService() *StudentService {
	return &StudentService{
		studentDAO: dao.NewStudentDAO(),
		userDAO:    dao.NewUserDAO(),
	}
}

// CreateStudent 创建学生信息（注册时补充信息）
func (s *StudentService) CreateStudent(userId, major, grade, programmingLevel string, interests, learningTags []string) (*student.Student, error) {
	// 检查用户是否存在
	user, err := s.userDAO.GetUserById(userId)
	if err != nil || user == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "用户不存在")
	}

	// 检查是否已经创建学生信息
	existingStudent, _ := s.studentDAO.GetStudentByUserId(userId)
	if existingStudent != nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "学生信息已存在")
	}

	// 生成学生ID
	studentId := fmt.Sprintf("stu_%d", time.Now().UnixNano())

	// 转换interests和learningTags为JSON
	interestsJSON, _ := json.Marshal(interests)
	learningTagsJSON, _ := json.Marshal(learningTags)

	// 构建学生模型
	student := &student.Student{
		StudentId:        studentId,
		UserId:           userId,
		Major:            major,
		Grade:            grade,
		ProgrammingLevel: programmingLevel,
		Interests:        string(interestsJSON),
		LearningTags:     string(learningTagsJSON),
		Status:           1,
	}

	// 创建学生
	if err := s.studentDAO.CreateStudent(student); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "创建学生信息失败: "+err.Error())
	}

	return student, nil
}

// GetStudentByUserId 根据用户ID获取学生信息
func (s *StudentService) GetStudentByUserId(userId string) (*student.Student, error) {
	student, err := s.studentDAO.GetStudentByUserId(userId)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "查询学生信息失败: "+err.Error())
	}
	if student == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "学生信息不存在")
	}
	return student, nil
}

// UpdateStudent 更新学生信息
func (s *StudentService) UpdateStudent(studentId string, updates map[string]interface{}) (*student.Student, error) {
	// 检查学生是否存在
	existingStudent, err := s.studentDAO.GetStudentById(studentId)
	if err != nil || existingStudent == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "学生信息不存在")
	}

	// 执行更新
	if err := s.studentDAO.UpdateStudent(studentId, updates); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "更新学生信息失败: "+err.Error())
	}

	// 查询更新后的学生信息
	updatedStudent, err := s.studentDAO.GetStudentById(studentId)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "查询学生信息失败: "+err.Error())
	}

	return updatedStudent, nil
}

// UpdateLearningProgress 更新学习进度
func (s *StudentService) UpdateLearningProgress(studentId string, progress map[string]interface{}) error {
	// 转换为JSON
	progressJSON, err := json.Marshal(progress)
	if err != nil {
		return errs.NewCommonError(errs.ErrBadRequest, "学习进度格式错误")
	}

	updates := map[string]interface{}{
		"learning_progress": string(progressJSON),
	}

	if err := s.studentDAO.UpdateStudent(studentId, updates); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "更新学习进度失败: "+err.Error())
	}

	return nil
}

// ListStudents 查询学生列表
func (s *StudentService) ListStudents(page, pageSize int32, filters map[string]interface{}) ([]*student.Student, int32, error) {
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

	if major, ok := filters["major"].(string); ok && major != "" {
		whereClause += " AND major = ?"
		args = append(args, major)
	}
	if grade, ok := filters["grade"].(string); ok && grade != "" {
		whereClause += " AND grade = ?"
		args = append(args, grade)
	}
	if programmingLevel, ok := filters["programming_level"].(string); ok && programmingLevel != "" {
		whereClause += " AND programming_level = ?"
		args = append(args, programmingLevel)
	}

	// 查询总数
	total, err := s.studentDAO.CountStudents(whereClause, args)
	if err != nil {
		return nil, 0, errs.NewCommonError(errs.ErrInternal, "统计学生数量失败: "+err.Error())
	}

	// 查询列表
	offset := (page - 1) * pageSize
	students, err := s.studentDAO.ListStudents(whereClause, args, pageSize, offset)
	if err != nil {
		return nil, 0, errs.NewCommonError(errs.ErrInternal, "查询学生列表失败: "+err.Error())
	}

	return students, total, nil
}
