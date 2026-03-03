package dao

import (
	"github.com/yzf120/elysia-backend/model/subject"
)

// SemesterDAO 学期数据访问对象
type SemesterDAO interface {
	GetSemesterById(semesterId string) (*subject.Semester, error)
	GetSemesterByName(semesterName string) (*subject.Semester, error)
	ListSemesters(status int32) ([]*subject.Semester, error)
}

type semesterDAOImpl struct{}

// NewSemesterDAO 创建学期DAO
func NewSemesterDAO() SemesterDAO {
	return &semesterDAOImpl{}
}

// GetSemesterById 根据学期ID查询学期信息
func (d *semesterDAOImpl) GetSemesterById(semesterId string) (*subject.Semester, error) {
	db := DB
	var semester subject.Semester
	err := db.Where("semester_id = ? AND status = 1", semesterId).First(&semester).Error
	if err != nil {
		return nil, err
	}
	return &semester, nil
}

// GetSemesterByName 根据学期名称查询学期信息（如：2026春）
func (d *semesterDAOImpl) GetSemesterByName(semesterName string) (*subject.Semester, error) {
	db := DB
	var semester subject.Semester
	err := db.Where("semester_name = ? AND status = 1", semesterName).First(&semester).Error
	if err != nil {
		return nil, err
	}
	return &semester, nil
}

// ListSemesters 查询学期列表
func (d *semesterDAOImpl) ListSemesters(status int32) ([]*subject.Semester, error) {
	db := DB
	var semesters []*subject.Semester
	query := db.Model(&subject.Semester{})
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	err := query.Order("year DESC, term DESC").Find(&semesters).Error
	return semesters, err
}
