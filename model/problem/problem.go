package problem

import (
	"time"
)

// Problem 题目数据模型
type Problem struct {
	Id                  int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title               string    `gorm:"column:title;type:varchar(255);not null" json:"title"`
	TitleSlug           string    `gorm:"column:title_slug;type:varchar(255);uniqueIndex;not null" json:"title_slug"`
	Difficulty          string    `gorm:"column:difficulty;type:enum('简单','中等','困难');not null;default:'简单'" json:"difficulty"`
	Tags                string    `gorm:"column:tags;type:varchar(500);not null" json:"tags"`
	Description         string    `gorm:"column:description;type:text;not null" json:"description"`
	Explanation         string    `gorm:"column:explanation;type:text" json:"explanation"`
	Hint                string    `gorm:"column:hint;type:text" json:"hint"`
	Constraints         string    `gorm:"column:constraints;type:text" json:"constraints"`
	AdvancedRequirement string    `gorm:"column:advanced_requirement;type:text" json:"advanced_requirement"`
	TestCases           string    `gorm:"column:test_cases;type:json;not null" json:"test_cases"`
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime;autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Problem) TableName() string {
	return "problem"
}
