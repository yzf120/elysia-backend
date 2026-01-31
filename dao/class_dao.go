package dao

import (
	"github.com/yzf120/elysia-backend/client"
	"github.com/yzf120/elysia-backend/model/class"
)

// ClassDAO 班级数据访问对象
type ClassDAO interface {
	CreateClass(class *class.Class) error
	GetClassById(classId string) (*class.Class, error)
	GetClassByCode(classCode string) (*class.Class, error)
	UpdateClass(classId string, updates map[string]interface{}) error
	DeleteClass(classId string) error
	ListClasses(whereClause string, args []interface{}, limit, offset int32) ([]*class.Class, error)
	CountClasses(whereClause string, args []interface{}) (int32, error)
	ListClassesByTeacherId(teacherId string, limit, offset int32) ([]*class.Class, error)
}

type classDAOImpl struct{}

// NewClassDAO 创建班级DAO
func NewClassDAO() ClassDAO {
	return &classDAOImpl{}
}

// CreateClass 创建班级
func (d *classDAOImpl) CreateClass(class *class.Class) error {
	db := client.GetMySQLClient().GormDB
	return db.Create(class).Error
}

// GetClassById 根据班级ID查询班级
func (d *classDAOImpl) GetClassById(classId string) (*class.Class, error) {
	db := client.GetMySQLClient().GormDB
	var class class.Class
	err := db.Where("class_id = ?", classId).First(&class).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

// GetClassByCode 根据班级验证码查询班级
func (d *classDAOImpl) GetClassByCode(classCode string) (*class.Class, error) {
	db := client.GetMySQLClient().GormDB
	var class class.Class
	err := db.Where("class_code = ?", classCode).First(&class).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

// UpdateClass 更新班级信息
func (d *classDAOImpl) UpdateClass(classId string, updates map[string]interface{}) error {
	db := client.GetMySQLClient().GormDB
	return db.Model(&class.Class{}).Where("class_id = ?", classId).Updates(updates).Error
}

// DeleteClass 删除班级
func (d *classDAOImpl) DeleteClass(classId string) error {
	db := client.GetMySQLClient().GormDB
	return db.Where("class_id = ?", classId).Delete(&class.Class{}).Error
}

// ListClasses 查询班级列表
func (d *classDAOImpl) ListClasses(whereClause string, args []interface{}, limit, offset int32) ([]*class.Class, error) {
	db := client.GetMySQLClient().GormDB
	var classes []*class.Class
	query := db.Model(&class.Class{})

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Limit(int(limit)).Offset(int(offset)).Find(&classes).Error
	return classes, err
}

// CountClasses 统计班级数量
func (d *classDAOImpl) CountClasses(whereClause string, args []interface{}) (int32, error) {
	db := client.GetMySQLClient().GormDB
	var count int64
	query := db.Model(&class.Class{})

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Count(&count).Error
	return int32(count), err
}

// ListClassesByTeacherId 根据教师ID查询班级列表
func (d *classDAOImpl) ListClassesByTeacherId(teacherId string, limit, offset int32) ([]*class.Class, error) {
	db := client.GetMySQLClient().GormDB
	var classes []*class.Class
	err := db.Where("teacher_id = ?", teacherId).Limit(int(limit)).Offset(int(offset)).Find(&classes).Error
	return classes, err
}

// ClassMemberDAO 班级成员数据访问对象
type ClassMemberDAO interface {
	AddMember(member *class.ClassMember) error
	RemoveMember(classId, studentId string) error
	GetMember(classId, studentId string) (*class.ClassMember, error)
	ListMembersByClassId(classId string, limit, offset int32) ([]*class.ClassMember, error)
	ListClassesByStudentId(studentId string, limit, offset int32) ([]*class.ClassMember, error)
	CountMembersByClassId(classId string) (int32, error)
	UpdateMemberStatus(classId, studentId string, status int32) error
}

type classMemberDAOImpl struct{}

// NewClassMemberDAO 创建班级成员DAO
func NewClassMemberDAO() ClassMemberDAO {
	return &classMemberDAOImpl{}
}

// AddMember 添加班级成员
func (d *classMemberDAOImpl) AddMember(member *class.ClassMember) error {
	db := client.GetMySQLClient().GormDB
	return db.Create(member).Error
}

// RemoveMember 移除班级成员
func (d *classMemberDAOImpl) RemoveMember(classId, studentId string) error {
	db := client.GetMySQLClient().GormDB
	return db.Where("class_id = ? AND student_id = ?", classId, studentId).Delete(&class.ClassMember{}).Error
}

// GetMember 查询班级成员
func (d *classMemberDAOImpl) GetMember(classId, studentId string) (*class.ClassMember, error) {
	db := client.GetMySQLClient().GormDB
	var member class.ClassMember
	err := db.Where("class_id = ? AND student_id = ?", classId, studentId).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// ListMembersByClassId 根据班级ID查询成员列表
func (d *classMemberDAOImpl) ListMembersByClassId(classId string, limit, offset int32) ([]*class.ClassMember, error) {
	db := client.GetMySQLClient().GormDB
	var members []*class.ClassMember
	err := db.Where("class_id = ? AND status = 1", classId).Limit(int(limit)).Offset(int(offset)).Find(&members).Error
	return members, err
}

// ListClassesByStudentId 根据学生ID查询班级列表
func (d *classMemberDAOImpl) ListClassesByStudentId(studentId string, limit, offset int32) ([]*class.ClassMember, error) {
	db := client.GetMySQLClient().GormDB
	var members []*class.ClassMember
	err := db.Where("student_id = ? AND status = 1", studentId).Limit(int(limit)).Offset(int(offset)).Find(&members).Error
	return members, err
}

// CountMembersByClassId 统计班级成员数量
func (d *classMemberDAOImpl) CountMembersByClassId(classId string) (int32, error) {
	db := client.GetMySQLClient().GormDB
	var count int64
	err := db.Model(&class.ClassMember{}).Where("class_id = ? AND status = 1", classId).Count(&count).Error
	return int32(count), err
}

// UpdateMemberStatus 更新成员状态
func (d *classMemberDAOImpl) UpdateMemberStatus(classId, studentId string, status int32) error {
	db := client.GetMySQLClient().GormDB
	return db.Model(&class.ClassMember{}).Where("class_id = ? AND student_id = ?", classId, studentId).Update("status", status).Error
}
