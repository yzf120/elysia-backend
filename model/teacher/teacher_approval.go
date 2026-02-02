package teacher

import (
	"time"
)

// TeacherApproval 教师审批单数据模型
type TeacherApproval struct {
	Id               int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ApprovalId       string    `gorm:"column:approval_id;type:varchar(64);uniqueIndex;not null" json:"approval_id"`
	TeacherId        string    `gorm:"column:teacher_id;type:varchar(64);uniqueIndex;not null" json:"teacher_id"`
	EmployeeNumber   string    `gorm:"column:employee_number;type:varchar(32);not null" json:"employee_number"`
	SchoolEmail      string    `gorm:"column:school_email;type:varchar(128);not null" json:"school_email"`
	TeacherName      string    `gorm:"column:teacher_name;type:varchar(64);not null" json:"teacher_name"`
	Phone            string    `gorm:"column:phone;type:varchar(16);not null" json:"phone"`
	Department       string    `gorm:"column:department;type:varchar(128)" json:"department"`
	Title            string    `gorm:"column:title;type:varchar(64)" json:"title"`
	TeachingSubjects string    `gorm:"column:teaching_subjects;type:text" json:"teaching_subjects"`
	TeachingYears    int32     `gorm:"column:teaching_years;type:int;not null;default:0" json:"teaching_years"`
	ApplyRemark      string    `gorm:"column:apply_remark;type:varchar(512)" json:"apply_remark"`
	ApprovalStatus   int32     `gorm:"column:approval_status;type:tinyint;not null;default:0" json:"approval_status"`
	ApproverId       string    `gorm:"column:approver_id;type:varchar(64)" json:"approver_id"`
	ApproverName     string    `gorm:"column:approver_name;type:varchar(64)" json:"approver_name"`
	ApprovalRemark   string    `gorm:"column:approval_remark;type:varchar(512)" json:"approval_remark"`
	ApprovalTime     time.Time `gorm:"column:approval_time;type:datetime" json:"approval_time"`
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime       time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (TeacherApproval) TableName() string {
	return "teacher_approval"
}
