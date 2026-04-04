package class

import "time"

// 资料类型常量
const (
	MaterialTypePDF   = "pdf"
	MaterialTypeWord  = "word"
	MaterialTypeText  = "text"
	MaterialTypeVideo = "video"
)

// SectionMaterial 小节学习资料数据模型
type SectionMaterial struct {
	Id           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MaterialId   string    `gorm:"column:material_id;type:varchar(64);uniqueIndex;not null" json:"material_id"`
	SectionId    string    `gorm:"column:section_id;type:varchar(64);not null;index:idx_section_id" json:"section_id"`
	ChapterId    string    `gorm:"column:chapter_id;type:varchar(64);not null;index:idx_chapter_id" json:"chapter_id"`
	ClassId      string    `gorm:"column:class_id;type:varchar(64);not null;index:idx_class_id" json:"class_id"`
	TeacherId    string    `gorm:"column:teacher_id;type:varchar(64);not null" json:"teacher_id"`
	Title        string    `gorm:"column:title;type:varchar(256);not null" json:"title"`
	Description  string    `gorm:"column:description;type:text" json:"description"`
	MaterialType string    `gorm:"column:material_type;type:varchar(32);not null" json:"material_type"`
	FileName     string    `gorm:"column:file_name;type:varchar(256);not null;default:''" json:"file_name"`
	FilePath     string    `gorm:"column:file_path;type:varchar(512);not null;default:''" json:"file_path"`
	FileSize     int64     `gorm:"column:file_size;type:bigint;not null;default:0" json:"file_size"`
	MimeType     string    `gorm:"column:mime_type;type:varchar(128);not null;default:''" json:"mime_type"`
	SortOrder    int32     `gorm:"column:sort_order;type:int;not null;default:0" json:"sort_order"`
	Status       int32     `gorm:"column:status;type:tinyint;not null;default:1" json:"status"`
	CreateTime   time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

// TableName 指定表名
func (SectionMaterial) TableName() string {
	return "section_material"
}
