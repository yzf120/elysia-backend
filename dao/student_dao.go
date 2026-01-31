package dao

import (
	"github.com/yzf120/elysia-backend/client"
	"github.com/yzf120/elysia-backend/model"
)

// StudentDAO 学生数据访问对象
type StudentDAO interface {
	CreateStudent(student *model.Student) error
	GetStudentById(studentId string) (*model.Student, error)
	GetStudentByUserId(userId string) (*model.Student, error)
	UpdateStudent(studentId string, updates map[string]interface{}) error
	DeleteStudent(studentId string) error
	ListStudents(whereClause string, args []interface{}, limit, offset int32) ([]*model.Student, error)
	CountStudents(whereClause string, args []interface{}) (int32, error)
}

type studentDAOImpl struct{}

// NewStudentDAO 创建学生DAO
func NewStudentDAO() StudentDAO {
	return &studentDAOImpl{}
}

// CreateStudent 创建学生
func (d *studentDAOImpl) CreateStudent(student *model.Student) error {
	db := client.GetMySQLClient().GormDB
	return db.Create(student).Error
}

// GetStudentById 根据学生ID查询学生
func (d *studentDAOImpl) GetStudentById(studentId string) (*model.Student, error) {
	db := client.GetMySQLClient().GormDB
	var student model.Student
	err := db.Where("student_id = ?", studentId).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// GetStudentByUserId 根据用户ID查询学生
func (d *studentDAOImpl) GetStudentByUserId(userId string) (*model.Student, error) {
	db := client.GetMySQLClient().GormDB
	var student model.Student
	err := db.Where("user_id = ?", userId).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// UpdateStudent 更新学生信息
func (d *studentDAOImpl) UpdateStudent(studentId string, updates map[string]interface{}) error {
	db := client.GetMySQLClient().GormDB
	return db.Model(&model.Student{}).Where("student_id = ?", studentId).Updates(updates).Error
}

// DeleteStudent 删除学生
func (d *studentDAOImpl) DeleteStudent(studentId string) error {
	db := client.GetMySQLClient().GormDB
	return db.Where("student_id = ?", studentId).Delete(&model.Student{}).Error
}

// ListStudents 查询学生列表
func (d *studentDAOImpl) ListStudents(whereClause string, args []interface{}, limit, offset int32) ([]*model.Student, error) {
	db := client.GetMySQLClient().GormDB
	var students []*model.Student
	query := db.Model(&model.Student{})

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Limit(int(limit)).Offset(int(offset)).Find(&students).Error
	return students, err
}

// CountStudents 统计学生数量
func (d *studentDAOImpl) CountStudents(whereClause string, args []interface{}) (int32, error) {
	db := client.GetMySQLClient().GormDB
	var count int64
	query := db.Model(&model.Student{})

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Count(&count).Error
	return int32(count), err
}
