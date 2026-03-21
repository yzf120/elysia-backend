package intent

import (
	"time"
)

// UserIntentRecord 用户意图记录数据模型
type UserIntentRecord struct {
	Id               int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId           string    `gorm:"column:user_id;type:varchar(64);not null;index" json:"user_id"`              // 用户ID
	SessionId        string    `gorm:"column:session_id;type:varchar(128);index" json:"session_id"`                // 会话ID
	QuestionId       string    `gorm:"column:question_id;type:varchar(64)" json:"question_id"`                     // 关联题目ID
	OriginalRequest  string    `gorm:"column:original_request;type:text;not null" json:"original_request"`         // 原始请求
	IntentCode       string    `gorm:"column:intent_code;type:varchar(30);not null;index" json:"intent_code"`      // 意图编码
	IntentLevel1     string    `gorm:"column:intent_level1;type:varchar(50);not null;index" json:"intent_level1"`  // 一级意图
	RewrittenRequest string    `gorm:"column:rewritten_request;type:text" json:"rewritten_request"`                // 改写后请求
	IntentConfidence float64   `gorm:"column:intent_confidence;type:decimal(5,2)" json:"intent_confidence"`        // 置信度
	ResponseTimeMs   int       `gorm:"column:response_time_ms;type:int" json:"response_time_ms"`                   // 识别耗时
	RecognizeStatus  int       `gorm:"column:recognize_status;type:tinyint;not null;default:1" json:"recognize_status"` // 识别状态
	CreateTime       time.Time `gorm:"column:create_time;type:datetime;autoCreateTime;index" json:"create_time"`
}

// TableName 指定表名
func (UserIntentRecord) TableName() string {
	return "oj_user_intent_record"
}
