package dao

import (
	"strings"

	"github.com/yzf120/elysia-backend/model/platform"
)

// PlatformContentDAO 平台内容 DAO
type PlatformContentDAO struct{}

func NewPlatformContentDAO() *PlatformContentDAO {
	return &PlatformContentDAO{}
}

func (d *PlatformContentDAO) CreateSystemAnnouncement(item *platform.SystemAnnouncement) error {
	return DB.Create(item).Error
}

func (d *PlatformContentDAO) UpdateSystemAnnouncement(announcementId string, updates map[string]interface{}) error {
	return DB.Model(&platform.SystemAnnouncement{}).Where("announcement_id = ?", announcementId).Updates(updates).Error
}

func (d *PlatformContentDAO) GetSystemAnnouncementByAnnouncementId(announcementId string) (*platform.SystemAnnouncement, error) {
	var item platform.SystemAnnouncement
	if err := DB.Where("announcement_id = ?", announcementId).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *PlatformContentDAO) DeleteSystemAnnouncement(announcementId string) error {
	return DB.Where("announcement_id = ?", announcementId).Delete(&platform.SystemAnnouncement{}).Error
}

func (d *PlatformContentDAO) ListSystemAnnouncements(keyword string, status *int32, publishedOnly bool, page, pageSize int) ([]*platform.SystemAnnouncement, int64, error) {
	query := DB.Model(&platform.SystemAnnouncement{})
	if publishedOnly {
		query = query.Where("status = ?", 1)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("(title LIKE ? OR content LIKE ?)", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []*platform.SystemAnnouncement
	offset := (page - 1) * pageSize
	if err := query.Order("publish_time DESC").Order("create_time DESC").Limit(pageSize).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (d *PlatformContentDAO) CreateBookshelfItem(item *platform.BookshelfItem) error {
	return DB.Create(item).Error
}

func (d *PlatformContentDAO) UpdateBookshelfItem(itemId string, updates map[string]interface{}) error {
	return DB.Model(&platform.BookshelfItem{}).Where("item_id = ?", itemId).Updates(updates).Error
}

func (d *PlatformContentDAO) GetBookshelfItemByItemId(itemId string) (*platform.BookshelfItem, error) {
	var item platform.BookshelfItem
	if err := DB.Where("item_id = ?", itemId).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *PlatformContentDAO) DeleteBookshelfItem(itemId string) error {
	return DB.Where("item_id = ?", itemId).Delete(&platform.BookshelfItem{}).Error
}

func (d *PlatformContentDAO) ListBookshelfItems(keyword string, status *int32, publishedOnly bool, page, pageSize int) ([]*platform.BookshelfItem, int64, error) {
	query := DB.Model(&platform.BookshelfItem{})
	if publishedOnly {
		query = query.Where("status = ?", 1)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("(title LIKE ? OR description LIKE ?)", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []*platform.BookshelfItem
	offset := (page - 1) * pageSize
	if err := query.Order("sort_order DESC").Order("publish_time DESC").Order("create_time DESC").Limit(pageSize).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
