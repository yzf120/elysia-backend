package dao

import (
	"github.com/yzf120/elysia-backend/client"
	"github.com/yzf120/elysia-backend/model"
)

// TeacherDAO 教师数据访问对象
type TeacherDAO interface {
	CreateTeacher(teacher *model.Teacher) error
	GetTeacherById(teacherId string) (*model.Teacher, error)
	GetTeacherByUserId(userId string) (*model.Teacher, error)
	GetTeacherByEmployeeNumber(employeeNumber string) (*model.Teacher, error)
	GetTeacherBySchoolEmail(schoolEmail string) (*model.Teacher, error)
	UpdateTeacher(teacherId string, updates map[string]interface{}) error
	DeleteTeacher(teacherId string) error
	ListTeachers(whereClause string, args []interface{}, limit, offset int32) ([]*model.Teacher, error)
	CountTeachers(whereClause string, args []interface{}) (int32, error)
}

type teacherDAOImpl struct{}

// NewTeacherDAO 创建教师DAO
func NewTeacherDAO() TeacherDAO {
	return &teacherDAOImpl{}
}

// CreateTeacher 创建教师
func (d *teacherDAOImpl) CreateTeacher(teacher *model.Teacher) error {
	db := client.GetMySQLClient().GormDB
	return db.Create(teacher).Error
}

// GetTeacherById 根据教师ID查询教师
func (d *teacherDAOImpl) GetTeacherById(teacherId string) (*model.Teacher, error) {
	db := client.GetMySQLClient().GormDB
	var teacher model.Teacher
	err := db.Where("teacher_id = ?", teacherId).First(&teacher).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// GetTeacherByUserId 根据用户ID查询教师
func (d *teacherDAOImpl) GetTeacherByUserId(userId string) (*model.Teacher, error) {
	db := client.GetMySQLClient().GormDB
	var teacher model.Teacher
	err := db.Where("user_id = ?", userId).First(&teacher).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// GetTeacherByEmployeeNumber 根据工号查询教师
func (d *teacherDAOImpl) GetTeacherByEmployeeNumber(employeeNumber string) (*model.Teacher, error) {
	db := client.GetMySQLClient().GormDB
	var teacher model.Teacher
	err := db.Where("employee_number = ?", employeeNumber).First(&teacher).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// GetTeacherBySchoolEmail 根据学校邮箱查询教师
func (d *teacherDAOImpl) GetTeacherBySchoolEmail(schoolEmail string) (*model.Teacher, error) {
	db := client.GetMySQLClient().GormDB
	var teacher model.Teacher
	err := db.Where("school_email = ?", schoolEmail).First(&teacher).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// UpdateTeacher 更新教师信息
func (d *teacherDAOImpl) UpdateTeacher(teacherId string, updates map[string]interface{}) error {
	db := client.GetMySQLClient().GormDB
	return db.Model(&model.Teacher{}).Where("teacher_id = ?", teacherId).Updates(updates).Error
}

// DeleteTeacher 删除教师
func (d *teacherDAOImpl) DeleteTeacher(teacherId string) error {
	db := client.GetMySQLClient().GormDB
	return db.Where("teacher_id = ?", teacherId).Delete(&model.Teacher{}).Error
}

// ListTeachers 查询教师列表
func (d *teacherDAOImpl) ListTeachers(whereClause string, args []interface{}, limit, offset int32) ([]*model.Teacher, error) {
	db := client.GetMySQLClient().GormDB
	var teachers []*model.Teacher
	query := db.Model(&model.Teacher{})

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Limit(int(limit)).Offset(int(offset)).Find(&teachers).Error
	return teachers, err
}

// CountTeachers 统计教师数量
func (d *teacherDAOImpl) CountTeachers(whereClause string, args []interface{}) (int32, error) {
	db := client.GetMySQLClient().GormDB
	var count int64
	query := db.Model(&model.Teacher{})

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Count(&count).Error
	return int32(count), err
}
