package subject

import (
	"time"
)

// Subject 科目数据模型
type Subject struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SubjectId   string    `gorm:"column:subject_id;type:varchar(64);uniqueIndex;not null" json:"subject_id"`
	SubjectName string    `gorm:"column:subject_name;type:varchar(128);not null" json:"subject_name"`
	SubjectCode string    `gorm:"column:subject_code;type:varchar(32);uniqueIndex;not null" json:"subject_code"`
	Category    string    `gorm:"column:category;type:varchar(64)" json:"category"`            // 科目分类：理科、文科、艺术等
	Description string    `gorm:"column:description;type:text" json:"description"`             // 科目描述
	Credits     int32     `gorm:"column:credits;type:int;default:0" json:"credits"`            // 学分
	Status      int32     `gorm:"column:status;type:tinyint;not null;default:1" json:"status"` // 状态：1-启用，0-禁用
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (Subject) TableName() string {
	return "subjects"
}

// TeacherSubject 教师-科目关联数据模型
type TeacherSubject struct {
	Id         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TeacherId  string    `gorm:"column:teacher_id;type:varchar(64);not null;index:idx_teacher_subject" json:"teacher_id"`
	SubjectId  string    `gorm:"column:subject_id;type:varchar(64);not null;index:idx_teacher_subject" json:"subject_id"`
	StartDate  time.Time `gorm:"column:start_date;type:date" json:"start_date"`               // 开始教授日期
	EndDate    time.Time `gorm:"column:end_date;type:date" json:"end_date"`                   // 结束教授日期（可为空）
	Status     int32     `gorm:"column:status;type:tinyint;not null;default:1" json:"status"` // 状态：1-在教，0-已停止
	Remark     string    `gorm:"column:remark;type:varchar(512)" json:"remark"`               // 备注
	CreateTime time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (TeacherSubject) TableName() string {
	return "teacher_subjects"
}
