package dao

import (
	"github.com/yzf120/elysia-backend/model/teacher"
	"gorm.io/gorm"
)

// TeacherDAO 教师数据访问对象
type TeacherDAO interface {
	CreateTeacher(teacher *teacher.Teacher) error
	GetTeacherById(teacherId string) (*teacher.Teacher, error)
	GetTeacherByPhoneNumber(phoneNumber string) (*teacher.Teacher, error)
	GetTeacherByEmployeeNumber(employeeNumber string) (*teacher.Teacher, error)
	GetTeacherBySchoolEmail(schoolEmail string) (*teacher.Teacher, error)
	UpdateTeacher(teacherId string, updates map[string]interface{}) error
	DeleteTeacher(teacherId string) error
	ListTeachers(whereClause string, args []interface{}, limit, offset int32) ([]*teacher.Teacher, error)
	ListTeachersAll(whereClause string, args []interface{}) ([]*teacher.Teacher, error)
	CountTeachers(whereClause string, args []interface{}) (int32, error)
	BatchUpdateTeacherStatus(teacherIds []string, status int32) error
}

type teacherDAOImpl struct{}

// NewTeacherDAO 创建教师DAO
func NewTeacherDAO() TeacherDAO {
	return &teacherDAOImpl{}
}

// CreateTeacher 创建教师
func (d *teacherDAOImpl) CreateTeacher(teacher *teacher.Teacher) error {
	db := DB
	return db.Create(teacher).Error
}

// GetTeacherById 根据教师ID查询教师
func (d *teacherDAOImpl) GetTeacherById(teacherId string) (*teacher.Teacher, error) {
	db := DB
	var teacher teacher.Teacher
	err := db.Where("teacher_id = ?", teacherId).First(&teacher).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// GetTeacherByPhoneNumber 根据手机号查询教师
func (d *teacherDAOImpl) GetTeacherByPhoneNumber(phoneNumber string) (*teacher.Teacher, error) {
	db := DB
	var teacher teacher.Teacher
	err := db.Where("phone_number = ?", phoneNumber).First(&teacher).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// GetTeacherByEmployeeNumber 根据工号查询教师
func (d *teacherDAOImpl) GetTeacherByEmployeeNumber(employeeNumber string) (*teacher.Teacher, error) {
	db := DB
	var teacher teacher.Teacher
	err := db.Where("employee_number = ?", employeeNumber).First(&teacher).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// GetTeacherBySchoolEmail 根据学校邮箱查询教师
func (d *teacherDAOImpl) GetTeacherBySchoolEmail(schoolEmail string) (*teacher.Teacher, error) {
	db := DB
	var teacher teacher.Teacher
	err := db.Where("school_email = ?", schoolEmail).First(&teacher).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// UpdateTeacher 更新教师信息
func (d *teacherDAOImpl) UpdateTeacher(teacherId string, updates map[string]interface{}) error {
	db := DB
	return db.Model(&teacher.Teacher{}).Where("teacher_id = ?", teacherId).Updates(updates).Error
}

// DeleteTeacher 删除教师
func (d *teacherDAOImpl) DeleteTeacher(teacherId string) error {
	db := DB
	return db.Where("teacher_id = ?", teacherId).Delete(&teacher.Teacher{}).Error
}

func (d *teacherDAOImpl) buildTeacherListQuery(whereClause string, args []interface{}) *gorm.DB {
	query := DB.Model(&teacher.Teacher{})
	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}
	return query.Order("create_time DESC")
}

// ListTeachers 查询教师列表
func (d *teacherDAOImpl) ListTeachers(whereClause string, args []interface{}, limit, offset int32) ([]*teacher.Teacher, error) {
	var teachers []*teacher.Teacher
	err := d.buildTeacherListQuery(whereClause, args).Limit(int(limit)).Offset(int(offset)).Find(&teachers).Error
	return teachers, err
}

func (d *teacherDAOImpl) ListTeachersAll(whereClause string, args []interface{}) ([]*teacher.Teacher, error) {
	var teachers []*teacher.Teacher
	err := d.buildTeacherListQuery(whereClause, args).Find(&teachers).Error
	return teachers, err
}

// CountTeachers 统计教师数量
func (d *teacherDAOImpl) CountTeachers(whereClause string, args []interface{}) (int32, error) {
	db := DB
	var count int64
	query := db.Model(&teacher.Teacher{})
	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}
	err := query.Count(&count).Error
	return int32(count), err
}

func (d *teacherDAOImpl) BatchUpdateTeacherStatus(teacherIds []string, status int32) error {
	if len(teacherIds) == 0 {
		return nil
	}
	return DB.Model(&teacher.Teacher{}).Where("teacher_id IN ?", teacherIds).Updates(map[string]interface{}{"status": status}).Error
}
