package dao

import (
	"github.com/yzf120/elysia-backend/model/subject"
)

// TeacherSubjectDAO 教师-科目关联数据访问对象
type TeacherSubjectDAO interface {
	CreateTeacherSubject(teacherSubject *subject.TeacherSubject) error
	GetTeacherSubjectById(id int64) (*subject.TeacherSubject, error)
	GetTeacherSubject(teacherId, subjectId string) (*subject.TeacherSubject, error)
	UpdateTeacherSubject(id int64, updates map[string]interface{}) error
	DeleteTeacherSubject(id int64) error
	DeleteByTeacherAndSubject(teacherId, subjectId string) error
	ListTeacherSubjects(whereClause string, args []interface{}, limit, offset int32) ([]*subject.TeacherSubject, error)
	CountTeacherSubjects(whereClause string, args []interface{}) (int32, error)
	GetSubjectsByTeacherId(teacherId string) ([]*subject.Subject, error)
	GetTeachersBySubjectId(subjectId string) ([]string, error)
}

type teacherSubjectDAOImpl struct{}

// NewTeacherSubjectDAO 创建教师-科目关联DAO
func NewTeacherSubjectDAO() TeacherSubjectDAO {
	return &teacherSubjectDAOImpl{}
}

// CreateTeacherSubject 创建教师-科目关联
func (d *teacherSubjectDAOImpl) CreateTeacherSubject(teacherSubject *subject.TeacherSubject) error {
	db := DB
	return db.Create(teacherSubject).Error
}

// GetTeacherSubjectById 根据ID查询教师-科目关联
func (d *teacherSubjectDAOImpl) GetTeacherSubjectById(id int64) (*subject.TeacherSubject, error) {
	db := DB
	var ts subject.TeacherSubject
	err := db.Where("id = ?", id).First(&ts).Error
	if err != nil {
		return nil, err
	}
	return &ts, nil
}

// GetTeacherSubject 根据教师ID和科目ID查询关联
func (d *teacherSubjectDAOImpl) GetTeacherSubject(teacherId, subjectId string) (*subject.TeacherSubject, error) {
	db := DB
	var ts subject.TeacherSubject
	err := db.Where("teacher_id = ? AND subject_id = ?", teacherId, subjectId).First(&ts).Error
	if err != nil {
		return nil, err
	}
	return &ts, nil
}

// UpdateTeacherSubject 更新教师-科目关联信息
func (d *teacherSubjectDAOImpl) UpdateTeacherSubject(id int64, updates map[string]interface{}) error {
	db := DB
	return db.Model(&subject.TeacherSubject{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteTeacherSubject 删除教师-科目关联
func (d *teacherSubjectDAOImpl) DeleteTeacherSubject(id int64) error {
	db := DB
	return db.Where("id = ?", id).Delete(&subject.TeacherSubject{}).Error
}

// DeleteByTeacherAndSubject 根据教师ID和科目ID删除关联
func (d *teacherSubjectDAOImpl) DeleteByTeacherAndSubject(teacherId, subjectId string) error {
	db := DB
	return db.Where("teacher_id = ? AND subject_id = ?", teacherId, subjectId).Delete(&subject.TeacherSubject{}).Error
}

// ListTeacherSubjects 查询教师-科目关联列表
func (d *teacherSubjectDAOImpl) ListTeacherSubjects(whereClause string, args []interface{}, limit, offset int32) ([]*subject.TeacherSubject, error) {
	db := DB
	var teacherSubjects []*subject.TeacherSubject
	query := db.Model(&subject.TeacherSubject{})

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Limit(int(limit)).Offset(int(offset)).Find(&teacherSubjects).Error
	return teacherSubjects, err
}

// CountTeacherSubjects 统计教师-科目关联数量
func (d *teacherSubjectDAOImpl) CountTeacherSubjects(whereClause string, args []interface{}) (int32, error) {
	db := DB
	var count int64
	query := db.Model(&subject.TeacherSubject{})

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Count(&count).Error
	return int32(count), err
}

// GetSubjectsByTeacherId 根据教师ID获取其教授的所有科目
func (d *teacherSubjectDAOImpl) GetSubjectsByTeacherId(teacherId string) ([]*subject.Subject, error) {
	db := DB
	var subjects []*subject.Subject
	err := db.Table("subjects").
		Joins("INNER JOIN teacher_subjects ON subjects.subject_id = teacher_subjects.subject_id").
		Where("teacher_subjects.teacher_id = ? AND teacher_subjects.status = 1", teacherId).
		Find(&subjects).Error
	return subjects, err
}

// GetTeachersBySubjectId 根据科目ID获取教授该科目的所有教师ID
func (d *teacherSubjectDAOImpl) GetTeachersBySubjectId(subjectId string) ([]string, error) {
	db := DB
	var teacherIds []string
	err := db.Table("teacher_subjects").
		Select("teacher_id").
		Where("subject_id = ? AND status = 1", subjectId).
		Pluck("teacher_id", &teacherIds).Error
	return teacherIds, err
}
