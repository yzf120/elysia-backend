package dao

import (
	"github.com/yzf120/elysia-backend/model/student"
)

// StudentDAO 学生数据访问对象
type StudentDAO interface {
	CreateStudent(student *student.Student) error
	GetStudentById(studentId string) (*student.Student, error)
	GetStudentByPhoneNumber(phoneNumber string) (*student.Student, error)
	GetStudentByStudentNumber(studentNumber string) (*student.Student, error)
	UpdateStudent(studentId string, updates map[string]interface{}) error
	DeleteStudent(studentId string) error
	ListStudents(whereClause string, args []interface{}, limit, offset int32) ([]*student.Student, error)
	CountStudents(whereClause string, args []interface{}) (int32, error)
}

type studentDAOImpl struct{}

// NewStudentDAO 创建学生DAO
func NewStudentDAO() StudentDAO {
	return &studentDAOImpl{}
}

// CreateStudent 创建学生
func (d *studentDAOImpl) CreateStudent(student *student.Student) error {
	return DB.Create(student).Error
}

// GetStudentById 根据学生ID查询学生
func (d *studentDAOImpl) GetStudentById(studentId string) (*student.Student, error) {
	var student student.Student
	err := DB.Where("student_id = ?", studentId).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// GetStudentByPhoneNumber 根据手机号查询学生
func (d *studentDAOImpl) GetStudentByPhoneNumber(phoneNumber string) (*student.Student, error) {
	var student student.Student
	err := DB.Where("phone_number = ?", phoneNumber).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// GetStudentByStudentNumber 根据学号查询学生
func (d *studentDAOImpl) GetStudentByStudentNumber(studentNumber string) (*student.Student, error) {
	var student student.Student
	err := DB.Where("student_number = ?", studentNumber).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// UpdateStudent 更新学生信息
func (d *studentDAOImpl) UpdateStudent(studentId string, updates map[string]interface{}) error {
	return DB.Model(&student.Student{}).Where("student_id = ?", studentId).Updates(updates).Error
}

// DeleteStudent 删除学生
func (d *studentDAOImpl) DeleteStudent(studentId string) error {
	return DB.Where("student_id = ?", studentId).Delete(&student.Student{}).Error
}

// ListStudents 查询学生列表
func (d *studentDAOImpl) ListStudents(whereClause string, args []interface{}, limit, offset int32) ([]*student.Student, error) {
	var students []*student.Student
	query := DB.Model(&student.Student{})

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Limit(int(limit)).Offset(int(offset)).Find(&students).Error
	return students, err
}

// CountStudents 统计学生数量
func (d *studentDAOImpl) CountStudents(whereClause string, args []interface{}) (int32, error) {
	var count int64
	query := DB.Model(&student.Student{})

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Count(&count).Error
	return int32(count), err
}
