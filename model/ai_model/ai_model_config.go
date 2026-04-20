package ai_model

import (
	"time"
)

// AIModelConfig AI模型配置数据模型
type AIModelConfig struct {
	Id          int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ModelId     string    `gorm:"column:model_id;type:varchar(100);uniqueIndex;not null" json:"model_id"` // 模型ID（如 doubao-seed-1-6-lite-251015）
	ModelName   string    `gorm:"column:model_name;type:varchar(100);not null" json:"model_name"`         // 模型显示名称
	Provider    string    `gorm:"column:provider;type:varchar(50);not null" json:"provider"`              // 提供商：doubao, qwen
	Description string    `gorm:"column:description;type:varchar(500)" json:"description"`                // 模型描述
	IsEnabled   int       `gorm:"column:is_enabled;type:tinyint;not null;default:1" json:"is_enabled"`    // 是否启用：1-启用 0-禁用
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (AIModelConfig) TableName() string {
	return "ai_model_config"
}
