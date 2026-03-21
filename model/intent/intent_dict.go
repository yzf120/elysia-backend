package intent

import (
	"time"
)

// IntentDict 意图字典数据模型
type IntentDict struct {
	Id              int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IntentLevel1    string    `gorm:"column:intent_level1;type:varchar(50);not null" json:"intent_level1"`         // 一级意图
	IntentLevel2    string    `gorm:"column:intent_level2;type:varchar(100);not null" json:"intent_level2"`        // 二级子意图
	IntentCode      string    `gorm:"column:intent_code;type:varchar(30);uniqueIndex;not null" json:"intent_code"` // 意图编码
	Description     string    `gorm:"column:description;type:varchar(500)" json:"description"`                     // 意图描述
	MatchKeywords   string    `gorm:"column:match_keywords;type:text" json:"match_keywords"`                       // 匹配关键词
	ExampleQueries  string    `gorm:"column:example_queries;type:text" json:"example_queries"`                     // 示例问题（JSON数组）
	RewriteTemplate string    `gorm:"column:rewrite_template;type:text" json:"rewrite_template"`                   // 改写模板
	AgentRoute      string    `gorm:"column:agent_route;type:varchar(50);not null" json:"agent_route"`             // 路由Agent
	Priority        int       `gorm:"column:priority;type:int;not null;default:0" json:"priority"`                 // 优先级
	IsValid         int       `gorm:"column:is_valid;type:tinyint;not null;default:1" json:"is_valid"`             // 是否有效
	CreateTime      time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (IntentDict) TableName() string {
	return "oj_intent_dict"
}
