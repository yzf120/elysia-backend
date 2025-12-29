package model

import (
	"time"
)

// User 用户数据模型
type User struct {
	Id              int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId          string    `gorm:"column:user_id;type:varchar(64);uniqueIndex;not null" json:"user_id"`
	UserName        string    `gorm:"column:user_name;type:varchar(100);not null" json:"user_name"`
	Password        string    `gorm:"column:password;type:varchar(255);not null" json:"-"`
	Email           string    `gorm:"column:email;type:varchar(100)" json:"email"`
	Gender          int32     `gorm:"column:gender;type:tinyint;default:0" json:"gender"`
	PhoneNumber     string    `gorm:"column:phone_number;type:varchar(20);uniqueIndex" json:"phone_number"`
	WxMiniAppOpenId string    `gorm:"column:wx_mini_app_open_id;type:varchar(100);uniqueIndex" json:"wx_mini_app_open_id"`
	ChineseName     string    `gorm:"column:chinese_name;type:varchar(100)" json:"chinese_name"`
	RiskLevel       int32     `gorm:"column:risk_level;type:tinyint;default:0" json:"risk_level"`
	Status          int32     `gorm:"column:status;type:tinyint;default:1" json:"status"`
	CreateTime      time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
	WhiteList       int32     `gorm:"column:white_list;type:tinyint;default:0" json:"white_list"`
	ImageURL        string    `gorm:"column:image_url;type:varchar(255)" json:"image_url"`
	RegisterSource  string    `gorm:"column:register_source;type:varchar(50)" json:"register_source"`
	UserType        string    `gorm:"column:user_type;type:varchar(50)" json:"user_type"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}
