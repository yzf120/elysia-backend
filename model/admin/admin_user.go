package admin

import (
	"time"
)

// AdminUser 管理员用户数据模型
type AdminUser struct {
	Id                 int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AdminId            string    `gorm:"column:admin_id;type:varchar(64);uniqueIndex;not null" json:"admin_id"`
	Username           string    `gorm:"column:username;type:varchar(128);uniqueIndex;not null" json:"username"`
	Password           string    `gorm:"column:password;type:varchar(512);not null" json:"-"`
	RealName           string    `gorm:"column:real_name;type:varchar(128);not null" json:"real_name"`
	Email              string    `gorm:"column:email;type:varchar(128);uniqueIndex;not null" json:"email"`
	Role               string    `gorm:"column:role;type:varchar(32);not null;default:'admin'" json:"role"`
	Status             int32     `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`
	LastLoginTime      time.Time `gorm:"column:last_login_time;type:datetime" json:"last_login_time"`
	LastLoginIp        string    `gorm:"column:last_login_ip;type:varchar(64)" json:"last_login_ip"`
	LoginFailCount     int32     `gorm:"column:login_fail_count;type:int;not null;default:0" json:"login_fail_count"`
	PasswordUpdateTime time.Time `gorm:"column:password_update_time;type:datetime;default:CURRENT_TIMESTAMP" json:"password_update_time"`
	CreateTime         time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime         time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
	Remark             string    `gorm:"column:remark;type:varchar(512)" json:"remark"`
}

// TableName 指定表名
func (AdminUser) TableName() string {
	return "admin_user"
}

// AdminOperationLog 管理员操作日志数据模型
type AdminOperationLog struct {
	Id              int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AdminId         string    `gorm:"column:admin_id;type:varchar(64);not null" json:"admin_id"`
	OperationType   string    `gorm:"column:operation_type;type:varchar(32);not null" json:"operation_type"`
	OperationDetail string    `gorm:"column:operation_detail;type:text" json:"operation_detail"`
	IpAddress       string    `gorm:"column:ip_address;type:varchar(64)" json:"ip_address"`
	UserAgent       string    `gorm:"column:user_agent;type:varchar(512)" json:"user_agent"`
	OperationTime   time.Time `gorm:"column:operation_time;type:datetime;autoCreateTime" json:"operation_time"`
}

// TableName 指定表名
func (AdminOperationLog) TableName() string {
	return "admin_operation_log"
}
