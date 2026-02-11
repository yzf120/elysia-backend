package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/model/subject"
	"gorm.io/gorm"
)

// TeacherSubjectService 教师-科目关联服务接口
type TeacherSubjectService interface {
	AssignSubjectToTeacher(teacherId, subjectId string, startDate time.Time, remark string) error
	RemoveSubjectFromTeacher(teacherId, subjectId string) error
	UpdateTeacherSubject(id int64, updates map[string]interface{}) error
	GetTeacherSubjects(teacherId string) ([]*subject.Subject, error)
	GetSubjectTeachers(subjectId string) ([]string, error)
	ListTeacherSubjectRelations(teacherId, subjectId string, status int32, page, pageSize int32) ([]*subject.TeacherSubject, int32, error)
	StopTeachingSubject(teacherId, subjectId string, endDate time.Time) error
	ResumeTeachingSubject(teacherId, subjectId string) error
}

type teacherSubjectServiceImpl struct {
	teacherSubjectDAO dao.TeacherSubjectDAO
	teacherDAO        dao.TeacherDAO
	subjectDAO        dao.SubjectDAO
}

// NewTeacherSubjectService 创建教师-科目关联服务
func NewTeacherSubjectService() TeacherSubjectService {
	return &teacherSubjectServiceImpl{
		teacherSubjectDAO: dao.NewTeacherSubjectDAO(),
		teacherDAO:        dao.NewTeacherDAO(),
		subjectDAO:        dao.NewSubjectDAO(),
	}
}

// AssignSubjectToTeacher 为教师分配科目
func (s *teacherSubjectServiceImpl) AssignSubjectToTeacher(teacherId, subjectId string, startDate time.Time, remark string) error {
	// 检查教师是否存在
	_, err := s.teacherDAO.GetTeacherById(teacherId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("教师不存在")
		}
		return fmt.Errorf("查询教师失败: %w", err)
	}

	// 检查科目是否存在
	_, err = s.subjectDAO.GetSubjectById(subjectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("科目不存在")
		}
		return fmt.Errorf("查询科目失败: %w", err)
	}

	// 检查是否已经分配
	existingRelation, err := s.teacherSubjectDAO.GetTeacherSubject(teacherId, subjectId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("检查关联关系失败: %w", err)
	}
	if existingRelation != nil {
		return errors.New("该教师已被分配此科目")
	}

	// 创建关联
	teacherSubject := &subject.TeacherSubject{
		TeacherId:  teacherId,
		SubjectId:  subjectId,
		StartDate:  startDate,
		Status:     1,
		Remark:     remark,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := s.teacherSubjectDAO.CreateTeacherSubject(teacherSubject); err != nil {
		return fmt.Errorf("分配科目失败: %w", err)
	}

	return nil
}

// RemoveSubjectFromTeacher 移除教师的科目
func (s *teacherSubjectServiceImpl) RemoveSubjectFromTeacher(teacherId, subjectId string) error {
	// 检查关联是否存在
	_, err := s.teacherSubjectDAO.GetTeacherSubject(teacherId, subjectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("关联关系不存在")
		}
		return fmt.Errorf("查询关联关系失败: %w", err)
	}

	if err := s.teacherSubjectDAO.DeleteByTeacherAndSubject(teacherId, subjectId); err != nil {
		return fmt.Errorf("移除科目失败: %w", err)
	}

	return nil
}

// UpdateTeacherSubject 更新教师-科目关联信息
func (s *teacherSubjectServiceImpl) UpdateTeacherSubject(id int64, updates map[string]interface{}) error {
	// 检查关联是否存在
	_, err := s.teacherSubjectDAO.GetTeacherSubjectById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("关联关系不存在")
		}
		return fmt.Errorf("查询关联关系失败: %w", err)
	}

	if err := s.teacherSubjectDAO.UpdateTeacherSubject(id, updates); err != nil {
		return fmt.Errorf("更新关联信息失败: %w", err)
	}

	return nil
}

// GetTeacherSubjects 获取教师教授的所有科目
func (s *teacherSubjectServiceImpl) GetTeacherSubjects(teacherId string) ([]*subject.Subject, error) {
	subjects, err := s.teacherSubjectDAO.GetSubjectsByTeacherId(teacherId)
	if err != nil {
		return nil, fmt.Errorf("获取教师科目失败: %w", err)
	}
	return subjects, nil
}

// GetSubjectTeachers 获取教授某科目的所有教师
func (s *teacherSubjectServiceImpl) GetSubjectTeachers(subjectId string) ([]string, error) {
	teacherIds, err := s.teacherSubjectDAO.GetTeachersBySubjectId(subjectId)
	if err != nil {
		return nil, fmt.Errorf("获取科目教师失败: %w", err)
	}
	return teacherIds, nil
}

// ListTeacherSubjectRelations 查询教师-科目关联列表
func (s *teacherSubjectServiceImpl) ListTeacherSubjectRelations(teacherId, subjectId string, status int32, page, pageSize int32) ([]*subject.TeacherSubject, int32, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var whereClause string
	var args []interface{}

	if teacherId != "" && subjectId != "" && status >= 0 {
		whereClause = "teacher_id = ? AND subject_id = ? AND status = ?"
		args = []interface{}{teacherId, subjectId, status}
	} else if teacherId != "" && subjectId != "" {
		whereClause = "teacher_id = ? AND subject_id = ?"
		args = []interface{}{teacherId, subjectId}
	} else if teacherId != "" && status >= 0 {
		whereClause = "teacher_id = ? AND status = ?"
		args = []interface{}{teacherId, status}
	} else if subjectId != "" && status >= 0 {
		whereClause = "subject_id = ? AND status = ?"
		args = []interface{}{subjectId, status}
	} else if teacherId != "" {
		whereClause = "teacher_id = ?"
		args = []interface{}{teacherId}
	} else if subjectId != "" {
		whereClause = "subject_id = ?"
		args = []interface{}{subjectId}
	} else if status >= 0 {
		whereClause = "status = ?"
		args = []interface{}{status}
	}

	relations, err := s.teacherSubjectDAO.ListTeacherSubjects(whereClause, args, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("查询关联列表失败: %w", err)
	}

	total, err := s.teacherSubjectDAO.CountTeacherSubjects(whereClause, args)
	if err != nil {
		return nil, 0, fmt.Errorf("统计关联数量失败: %w", err)
	}

	return relations, total, nil
}

// StopTeachingSubject 停止教授某科目
func (s *teacherSubjectServiceImpl) StopTeachingSubject(teacherId, subjectId string, endDate time.Time) error {
	relation, err := s.teacherSubjectDAO.GetTeacherSubject(teacherId, subjectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("关联关系不存在")
		}
		return fmt.Errorf("查询关联关系失败: %w", err)
	}

	updates := map[string]interface{}{
		"status":   0,
		"end_date": endDate,
	}

	if err := s.teacherSubjectDAO.UpdateTeacherSubject(relation.Id, updates); err != nil {
		return fmt.Errorf("停止教授科目失败: %w", err)
	}

	return nil
}

// ResumeTeachingSubject 恢复教授某科目
func (s *teacherSubjectServiceImpl) ResumeTeachingSubject(teacherId, subjectId string) error {
	relation, err := s.teacherSubjectDAO.GetTeacherSubject(teacherId, subjectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("关联关系不存在")
		}
		return fmt.Errorf("查询关联关系失败: %w", err)
	}

	updates := map[string]interface{}{
		"status": 1,
	}

	if err := s.teacherSubjectDAO.UpdateTeacherSubject(relation.Id, updates); err != nil {
		return fmt.Errorf("恢复教授科目失败: %w", err)
	}

	return nil
}
