package teacher

import (
	"time"
)

// Teacher 教师数据模型
type Teacher struct {
	Id                 int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TeacherId          string    `gorm:"column:teacher_id;type:varchar(64);uniqueIndex;not null" json:"teacher_id"`
	PhoneNumber        string    `gorm:"column:phone_number;type:varchar(20);uniqueIndex;not null" json:"phone_number"`
	Password           string    `gorm:"column:password;type:varchar(255)" json:"-"`
	TeacherName        string    `gorm:"column:teacher_name;type:varchar(100)" json:"teacher_name"`
	EmployeeNumber     string    `gorm:"column:employee_number;type:varchar(32);uniqueIndex;not null" json:"employee_number"`
	SchoolEmail        string    `gorm:"column:school_email;type:varchar(128);uniqueIndex;not null" json:"school_email"`
	Gender             int32     `gorm:"column:gender;type:tinyint;default:0" json:"gender"`
	ImageURL           string    `gorm:"column:image_url;type:varchar(255)" json:"image_url"`
	TeachingSubjects   string    `gorm:"column:teaching_subjects;type:text" json:"teaching_subjects"`
	TeachingYears      int32     `gorm:"column:teaching_years;type:int;not null;default:0" json:"teaching_years"`
	Department         string    `gorm:"column:department;type:varchar(128)" json:"department"`
	Title              string    `gorm:"column:title;type:varchar(64)" json:"title"`
	VerificationStatus int32     `gorm:"column:verification_status;type:tinyint;not null;default:0" json:"verification_status"`
	VerificationTime   time.Time `gorm:"column:verification_time;type:datetime" json:"verification_time"`
	VerifierId         string    `gorm:"column:verifier_id;type:varchar(64)" json:"verifier_id"`
	VerificationRemark string    `gorm:"column:verification_remark;type:varchar(512)" json:"verification_remark"`
	Status             int32     `gorm:"column:status;type:tinyint;not null;default:0" json:"status"`
	CreateTime         time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime         time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (Teacher) TableName() string {
	return "teacher"
}