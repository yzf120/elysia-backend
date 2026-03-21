package content_security

import (
	"time"
)

// ViolationRecord 内容安全违规记录
type ViolationRecord struct {
	Id         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId     string    `gorm:"column:user_id;type:varchar(64);index;not null" json:"user_id"`         // 用户ID（学生ID或教师ID）
	UserRole   string    `gorm:"column:user_role;type:varchar(16);not null" json:"user_role"`           // 用户角色：student / teacher
	SessionId  string    `gorm:"column:session_id;type:varchar(128)" json:"session_id"`                 // 会话ID
	SenderType string    `gorm:"column:sender_type;type:varchar(16);not null" json:"sender_type"`       // 发送方：user / ai
	Content    string    `gorm:"column:content;type:text" json:"content"`                               // 违规内容（截取前500字符）
	Suggestion string    `gorm:"column:suggestion;type:varchar(16);not null" json:"suggestion"`         // 审核建议：Block / Review
	Label      string    `gorm:"column:label;type:varchar(32);not null" json:"label"`                   // 违规标签：Porn / Abuse / Ad 等
	Score      int32     `gorm:"column:score;type:int;not null;default:0" json:"score"`                 // 违规置信度 0-100
	RequestId  string    `gorm:"column:request_id;type:varchar(128)" json:"request_id"`                 // 腾讯云审核请求ID
	ContentType string   `gorm:"column:content_type;type:varchar(16);not null;default:'text'" json:"content_type"` // 内容类型：text / image
	CreateTime time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
}

// TableName 指定表名
func (ViolationRecord) TableName() string {
	return "ai_violation_records"
}
