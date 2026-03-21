package intent

import (
	"time"
)

// IntentPromptTemplate 意图提示词模板数据模型
type IntentPromptTemplate struct {
	Id              int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IntentCode      string    `gorm:"column:intent_code;type:varchar(30);not null;index" json:"intent_code"`       // 关联意图编码
	TemplateType    string    `gorm:"column:template_type;type:varchar(30);not null;index" json:"template_type"`   // 模板类型
	TemplateName    string    `gorm:"column:template_name;type:varchar(100);not null" json:"template_name"`        // 模板名称
	TemplateContent string    `gorm:"column:template_content;type:text;not null" json:"template_content"`          // 模板内容
	IsActive        int       `gorm:"column:is_active;type:tinyint;not null;default:1;index" json:"is_active"`     // 是否启用
	Version         int       `gorm:"column:version;type:int;not null;default:1" json:"version"`                   // 版本号
	CreateTime      time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (IntentPromptTemplate) TableName() string {
	return "oj_intent_prompt_template"
}
