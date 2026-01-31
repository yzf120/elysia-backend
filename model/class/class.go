package class

import (
	"time"
)

// Class 班级数据模型
type Class struct {
	Id              int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ClassId         string    `gorm:"column:class_id;type:varchar(64);uniqueIndex;not null" json:"class_id"`
	ClassName       string    `gorm:"column:class_name;type:varchar(128);not null" json:"class_name"`
	ClassCode       string    `gorm:"column:class_code;type:varchar(32);uniqueIndex;not null" json:"class_code"`
	TeacherId       string    `gorm:"column:teacher_id;type:varchar(64);not null" json:"teacher_id"`
	Subject         string    `gorm:"column:subject;type:varchar(128)" json:"subject"`
	Semester        string    `gorm:"column:semester;type:varchar(32)" json:"semester"`
	MaxStudents     int32     `gorm:"column:max_students;type:int;not null;default:100" json:"max_students"`
	CurrentStudents int32     `gorm:"column:current_students;type:int;not null;default:0" json:"current_students"`
	Description     string    `gorm:"column:description;type:text" json:"description"`
	Announcement    string    `gorm:"column:announcement;type:text" json:"announcement"`
	QrCodeUrl       string    `gorm:"column:qr_code_url;type:varchar(512)" json:"qr_code_url"`
	Status          int32     `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`
	CreateTime      time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (Class) TableName() string {
	return "class"
}

// ClassMember 班级成员关联数据模型
type ClassMember struct {
	Id         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ClassId    string    `gorm:"column:class_id;type:varchar(64);not null" json:"class_id"`
	StudentId  string    `gorm:"column:student_id;type:varchar(64);not null" json:"student_id"`
	JoinTime   time.Time `gorm:"column:join_time;type:datetime;autoCreateTime" json:"join_time"`
	Status     int32     `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`
	Remark     string    `gorm:"column:remark;type:varchar(512)" json:"remark"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (ClassMember) TableName() string {
	return "class_member"
}