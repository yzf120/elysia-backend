package dao

import (
	"github.com/yzf120/elysia-backend/model/subject"
)

// SubjectDAO 科目数据访问对象
type SubjectDAO interface {
	CreateSubject(subject *subject.Subject) error
	GetSubjectById(subjectId string) (*subject.Subject, error)
	GetSubjectByCode(subjectCode string) (*subject.Subject, error)
	UpdateSubject(subjectId string, updates map[string]interface{}) error
	DeleteSubject(subjectId string) error
	ListSubjects(whereClause string, args []interface{}, limit, offset int32) ([]*subject.Subject, error)
	CountSubjects(whereClause string, args []interface{}) (int32, error)
}

type subjectDAOImpl struct{}

// NewSubjectDAO 创建科目DAO
func NewSubjectDAO() SubjectDAO {
	return &subjectDAOImpl{}
}

// CreateSubject 创建科目
func (d *subjectDAOImpl) CreateSubject(subject *subject.Subject) error {
	db := DB
	return db.Create(subject).Error
}

// GetSubjectById 根据科目ID查询科目
func (d *subjectDAOImpl) GetSubjectById(subjectId string) (*subject.Subject, error) {
	db := DB
	var subj subject.Subject
	err := db.Where("subject_id = ?", subjectId).First(&subj).Error
	if err != nil {
		return nil, err
	}
	return &subj, nil
}

// GetSubjectByCode 根据科目代码查询科目
func (d *subjectDAOImpl) GetSubjectByCode(subjectCode string) (*subject.Subject, error) {
	db := DB
	var subj subject.Subject
	err := db.Where("subject_code = ?", subjectCode).First(&subj).Error
	if err != nil {
		return nil, err
	}
	return &subj, nil
}

// UpdateSubject 更新科目信息
func (d *subjectDAOImpl) UpdateSubject(subjectId string, updates map[string]interface{}) error {
	db := DB
	return db.Model(&subject.Subject{}).Where("subject_id = ?", subjectId).Updates(updates).Error
}

// DeleteSubject 删除科目
func (d *subjectDAOImpl) DeleteSubject(subjectId string) error {
	db := DB
	return db.Where("subject_id = ?", subjectId).Delete(&subject.Subject{}).Error
}

// ListSubjects 查询科目列表
func (d *subjectDAOImpl) ListSubjects(whereClause string, args []interface{}, limit, offset int32) ([]*subject.Subject, error) {
	db := DB
	var subjects []*subject.Subject
	query := db.Model(&subject.Subject{})

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Limit(int(limit)).Offset(int(offset)).Find(&subjects).Error
	return subjects, err
}

// CountSubjects 统计科目数量
func (d *subjectDAOImpl) CountSubjects(whereClause string, args []interface{}) (int32, error) {
	db := DB
	var count int64
	query := db.Model(&subject.Subject{})

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Count(&count).Error
	return int32(count), err
}
