package platform

import "time"

// SystemAnnouncement 平台系统公告
// status: 0-草稿 1-已发布
// priority: 1-low 2-normal 3-high
type SystemAnnouncement struct {
	Id               int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AnnouncementId   string     `gorm:"column:announcement_id;type:varchar(64);uniqueIndex;not null" json:"announcement_id"`
	Title            string     `gorm:"column:title;type:varchar(200);not null" json:"title"`
	Content          string     `gorm:"column:content;type:text;not null" json:"content"`
	Priority         int32      `gorm:"column:priority;type:tinyint;not null;default:2" json:"priority"`
	Status           int32      `gorm:"column:status;type:tinyint;not null;default:0" json:"status"`
	PublisherAdminId string     `gorm:"column:publisher_admin_id;type:varchar(64);not null" json:"publisher_admin_id"`
	PublisherName    string     `gorm:"column:publisher_name;type:varchar(128);not null" json:"publisher_name"`
	ViewCount        int64      `gorm:"column:view_count;type:bigint;not null;default:0" json:"view_count"`
	PublishTime      *time.Time `gorm:"column:publish_time;type:datetime" json:"publish_time"`
	CreateTime       time.Time  `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime       time.Time  `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

func (SystemAnnouncement) TableName() string {
	return "system_announcements"
}

// BookshelfItem 平台书架条目
// status: 0-草稿 1-已发布
// content_type: text/link/attachment/mixed
type BookshelfItem struct {
	Id                    int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ItemId                string     `gorm:"column:item_id;type:varchar(64);uniqueIndex;not null" json:"item_id"`
	Title                 string     `gorm:"column:title;type:varchar(200);not null" json:"title"`
	Description           string     `gorm:"column:description;type:text" json:"description"`
	ContentType           string     `gorm:"column:content_type;type:varchar(32);not null" json:"content_type"`
	ExternalURL           string     `gorm:"column:external_url;type:varchar(1024)" json:"external_url"`
	AttachmentName        string     `gorm:"column:attachment_name;type:varchar(255)" json:"attachment_name"`
	AttachmentStorageName string     `gorm:"column:attachment_storage_name;type:varchar(255)" json:"attachment_storage_name"`
	AttachmentMimeType    string     `gorm:"column:attachment_mime_type;type:varchar(128)" json:"attachment_mime_type"`
	AttachmentSize        int64      `gorm:"column:attachment_size;type:bigint;not null;default:0" json:"attachment_size"`
	SortOrder             int32      `gorm:"column:sort_order;type:int;not null;default:0" json:"sort_order"`
	Status                int32      `gorm:"column:status;type:tinyint;not null;default:0" json:"status"`
	CreatorAdminId        string     `gorm:"column:creator_admin_id;type:varchar(64);not null" json:"creator_admin_id"`
	UpdaterAdminId        string     `gorm:"column:updater_admin_id;type:varchar(64);not null" json:"updater_admin_id"`
	PublishTime           *time.Time `gorm:"column:publish_time;type:datetime" json:"publish_time"`
	CreateTime            time.Time  `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime            time.Time  `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

func (BookshelfItem) TableName() string {
	return "platform_bookshelf_items"
}
