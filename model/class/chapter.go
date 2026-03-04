package class

import "time"

// SectionType 小节类型
const (
	SectionTypeProblem    = 1 // 算法题
	SectionTypeDiscussion = 2 // 讨论话题
)

// ClassChapter 班级章节数据模型
type ClassChapter struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ChapterId   string    `gorm:"column:chapter_id;type:varchar(64);uniqueIndex;not null" json:"chapter_id"`
	ClassId     string    `gorm:"column:class_id;type:varchar(64);not null;index:idx_class_id" json:"class_id"`
	Title       string    `gorm:"column:title;type:varchar(256);not null" json:"title"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	SortOrder   int32     `gorm:"column:sort_order;type:int;not null;default:0" json:"sort_order"`
	Status      int32     `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (ClassChapter) TableName() string {
	return "class_chapter"
}

// ClassSection 章节小节数据模型
type ClassSection struct {
	Id          int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SectionId   string `gorm:"column:section_id;type:varchar(64);uniqueIndex;not null" json:"section_id"`
	ChapterId   string `gorm:"column:chapter_id;type:varchar(64);not null;index:idx_chapter_id" json:"chapter_id"`
	ClassId     string `gorm:"column:class_id;type:varchar(64);not null;index:idx_class_id" json:"class_id"`
	Title       string `gorm:"column:title;type:varchar(256);not null" json:"title"`
	Description string `gorm:"column:description;type:text" json:"description"`
	SectionType int32  `gorm:"column:section_type;type:tinyint;not null;default:1" json:"section_type"`
	// 算法题关联字段（section_type=1 时使用，关联题库）
	ProblemId string `gorm:"column:problem_id;type:varchar(64);not null;default:''" json:"problem_id"`
	// 讨论内容字段（section_type=2 时使用）
	DiscussionTitle   string    `gorm:"column:discussion_title;type:varchar(256);not null;default:''" json:"discussion_title"`
	DiscussionContent string    `gorm:"column:discussion_content;type:text" json:"discussion_content"`
	SortOrder         int32     `gorm:"column:sort_order;type:int;not null;default:0" json:"sort_order"`
	Status            int32     `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`
	CreateTime        time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime        time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (ClassSection) TableName() string {
	return "class_section"
}
