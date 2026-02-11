package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/model/subject"
	"gorm.io/gorm"
)

// SubjectService 科目服务接口
type SubjectService interface {
	CreateSubject(subjectName, subjectCode, category, description string, credits int32) (*subject.Subject, error)
	GetSubjectById(subjectId string) (*subject.Subject, error)
	GetSubjectByCode(subjectCode string) (*subject.Subject, error)
	UpdateSubject(subjectId string, updates map[string]interface{}) error
	DeleteSubject(subjectId string) error
	ListSubjects(category string, status int32, page, pageSize int32) ([]*subject.Subject, int32, error)
	EnableSubject(subjectId string) error
	DisableSubject(subjectId string) error
}

type subjectServiceImpl struct {
	subjectDAO dao.SubjectDAO
}

// NewSubjectService 创建科目服务
func NewSubjectService() SubjectService {
	return &subjectServiceImpl{
		subjectDAO: dao.NewSubjectDAO(),
	}
}

// CreateSubject 创建科目
func (s *subjectServiceImpl) CreateSubject(subjectName, subjectCode, category, description string, credits int32) (*subject.Subject, error) {
	// 检查科目代码是否已存在
	existingSubject, err := s.subjectDAO.GetSubjectByCode(subjectCode)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("检查科目代码失败: %w", err)
	}
	if existingSubject != nil {
		return nil, errors.New("科目代码已存在")
	}

	// 创建科目
	newSubject := &subject.Subject{
		SubjectId:   fmt.Sprintf("subj_%d", time.Now().UnixNano()),
		SubjectName: subjectName,
		SubjectCode: subjectCode,
		Category:    category,
		Description: description,
		Credits:     credits,
		Status:      1,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}

	if err := s.subjectDAO.CreateSubject(newSubject); err != nil {
		return nil, fmt.Errorf("创建科目失败: %w", err)
	}

	return newSubject, nil
}

// GetSubjectById 根据科目ID获取科目信息
func (s *subjectServiceImpl) GetSubjectById(subjectId string) (*subject.Subject, error) {
	subj, err := s.subjectDAO.GetSubjectById(subjectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("科目不存在")
		}
		return nil, fmt.Errorf("获取科目信息失败: %w", err)
	}
	return subj, nil
}

// GetSubjectByCode 根据科目代码获取科目信息
func (s *subjectServiceImpl) GetSubjectByCode(subjectCode string) (*subject.Subject, error) {
	subj, err := s.subjectDAO.GetSubjectByCode(subjectCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("科目不存在")
		}
		return nil, fmt.Errorf("获取科目信息失败: %w", err)
	}
	return subj, nil
}

// UpdateSubject 更新科目信息
func (s *subjectServiceImpl) UpdateSubject(subjectId string, updates map[string]interface{}) error {
	// 检查科目是否存在
	_, err := s.GetSubjectById(subjectId)
	if err != nil {
		return err
	}

	// 如果更新科目代码，检查是否重复
	if newCode, ok := updates["subject_code"].(string); ok {
		existingSubject, err := s.subjectDAO.GetSubjectByCode(newCode)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("检查科目代码失败: %w", err)
		}
		if existingSubject != nil && existingSubject.SubjectId != subjectId {
			return errors.New("科目代码已存在")
		}
	}

	if err := s.subjectDAO.UpdateSubject(subjectId, updates); err != nil {
		return fmt.Errorf("更新科目信息失败: %w", err)
	}

	return nil
}

// DeleteSubject 删除科目
func (s *subjectServiceImpl) DeleteSubject(subjectId string) error {
	// 检查科目是否存在
	_, err := s.GetSubjectById(subjectId)
	if err != nil {
		return err
	}

	if err := s.subjectDAO.DeleteSubject(subjectId); err != nil {
		return fmt.Errorf("删除科目失败: %w", err)
	}

	return nil
}

// ListSubjects 查询科目列表
func (s *subjectServiceImpl) ListSubjects(category string, status int32, page, pageSize int32) ([]*subject.Subject, int32, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var whereClause string
	var args []interface{}

	if category != "" && status >= 0 {
		whereClause = "category = ? AND status = ?"
		args = []interface{}{category, status}
	} else if category != "" {
		whereClause = "category = ?"
		args = []interface{}{category}
	} else if status >= 0 {
		whereClause = "status = ?"
		args = []interface{}{status}
	}

	subjects, err := s.subjectDAO.ListSubjects(whereClause, args, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("查询科目列表失败: %w", err)
	}

	total, err := s.subjectDAO.CountSubjects(whereClause, args)
	if err != nil {
		return nil, 0, fmt.Errorf("统计科目数量失败: %w", err)
	}

	return subjects, total, nil
}

// EnableSubject 启用科目
func (s *subjectServiceImpl) EnableSubject(subjectId string) error {
	return s.UpdateSubject(subjectId, map[string]interface{}{"status": 1})
}

// DisableSubject 禁用科目
func (s *subjectServiceImpl) DisableSubject(subjectId string) error {
	return s.UpdateSubject(subjectId, map[string]interface{}{"status": 0})
}
