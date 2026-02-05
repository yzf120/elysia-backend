package student

import (
	"time"
)

// Student 学生数据模型
type Student struct {
	Id               int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	StudentId        string    `gorm:"column:student_id;type:varchar(64);uniqueIndex;not null" json:"student_id"`
	PhoneNumber      string    `gorm:"column:phone_number;type:varchar(20);uniqueIndex;not null" json:"phone_number"`
	Password         string    `gorm:"column:password;type:varchar(255)" json:"-"`
	StudentName      string    `gorm:"column:student_name;type:varchar(100)" json:"student_name"`
	StudentNumber    string    `gorm:"column:student_number;type:varchar(32)" json:"student_number"`
	Email            string    `gorm:"column:email;type:varchar(100)" json:"email"`
	Gender           int32     `gorm:"column:gender;type:tinyint;default:0" json:"gender"`
	ImageURL         string    `gorm:"column:image_url;type:varchar(255)" json:"image_url"`
	Major            string    `gorm:"column:major;type:varchar(128)" json:"major"`
	Grade            string    `gorm:"column:grade;type:varchar(32)" json:"grade"`
	ProgrammingLevel string    `gorm:"column:programming_level;type:varchar(32)" json:"programming_level"`
	Interests        string    `gorm:"column:interests;type:text" json:"interests"`
	LearningTags     string    `gorm:"column:learning_tags;type:text" json:"learning_tags"`
	LearningProgress string    `gorm:"column:learning_progress;type:text" json:"learning_progress"`
	Status           int32     `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime       time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (Student) TableName() string {
	return "students"
}