package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	adminModel "github.com/yzf120/elysia-backend/model/admin"
	"github.com/yzf120/elysia-backend/model/platform"
	"gorm.io/gorm"
)

const (
	platformBookshelfUploadDir  = "uploads/platform-bookshelf"
	maxBookshelfAttachmentSize  = 50 * 1024 * 1024
	systemAnnouncementPublished = 1
	platformBookshelfPublished  = 1
)

type PlatformContentService struct {
	platformDAO *dao.PlatformContentDAO
	adminDAO    dao.AdminUserDAO
}

type ListQuery struct {
	Page     int
	PageSize int
	Keyword  string
	Status   *int32
}

type SaveSystemAnnouncementInput struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Priority  string `json:"priority"`
	Published bool   `json:"published"`
}

type SaveBookshelfItemInput struct {
	Title           string
	Description     string
	ContentType     string
	ExternalURL     string
	Published       bool
	SortOrder       int32
	ClearAttachment bool
}

type SystemAnnouncementDTO struct {
	AnnouncementID string     `json:"announcement_id"`
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	Priority       string     `json:"priority"`
	Published      bool       `json:"published"`
	PublisherName  string     `json:"publisher_name"`
	ViewCount      int64      `json:"view_count"`
	PublishTime    *time.Time `json:"publish_time,omitempty"`
	CreateTime     time.Time  `json:"create_time"`
	UpdateTime     time.Time  `json:"update_time"`
}

type BookshelfItemDTO struct {
	ItemID           string     `json:"item_id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	ContentType      string     `json:"content_type"`
	ContentTypeLabel string     `json:"content_type_label"`
	ExternalURL      string     `json:"external_url"`
	AttachmentName   string     `json:"attachment_name"`
	AttachmentSize   int64      `json:"attachment_size"`
	HasAttachment    bool       `json:"has_attachment"`
	Published        bool       `json:"published"`
	SortOrder        int32      `json:"sort_order"`
	PublishTime      *time.Time `json:"publish_time,omitempty"`
	CreateTime       time.Time  `json:"create_time"`
	UpdateTime       time.Time  `json:"update_time"`
}

func NewPlatformContentService() *PlatformContentService {
	return &PlatformContentService{
		platformDAO: dao.NewPlatformContentDAO(),
		adminDAO:    dao.NewAdminUserDAO(),
	}
}

func (s *PlatformContentService) ListAdminSystemAnnouncements(query ListQuery) ([]*SystemAnnouncementDTO, int64, error) {
	page, pageSize := normalizePage(query.Page, query.PageSize)
	items, total, err := s.platformDAO.ListSystemAnnouncements(query.Keyword, query.Status, false, page, pageSize)
	if err != nil {
		return nil, 0, errs.NewCommonError(http.StatusInternalServerError, "查询系统公告失败: "+err.Error())
	}
	return mapAnnouncementList(items), total, nil
}

func (s *PlatformContentService) ListUserSystemAnnouncements(query ListQuery) ([]*SystemAnnouncementDTO, int64, error) {
	page, pageSize := normalizePage(query.Page, query.PageSize)
	items, total, err := s.platformDAO.ListSystemAnnouncements(query.Keyword, nil, true, page, pageSize)
	if err != nil {
		return nil, 0, errs.NewCommonError(http.StatusInternalServerError, "查询系统公告失败: "+err.Error())
	}
	return mapAnnouncementList(items), total, nil
}

func (s *PlatformContentService) CreateSystemAnnouncement(adminId string, input SaveSystemAnnouncementInput) (*SystemAnnouncementDTO, error) {
	adminUser, err := s.ensureAdmin(adminId)
	if err != nil {
		return nil, err
	}
	priorityCode, err := parsePriority(input.Priority)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.Title) == "" {
		return nil, errs.NewCommonError(http.StatusBadRequest, "公告标题不能为空")
	}
	if strings.TrimSpace(input.Content) == "" {
		return nil, errs.NewCommonError(http.StatusBadRequest, "公告内容不能为空")
	}

	item := &platform.SystemAnnouncement{
		AnnouncementId:   fmt.Sprintf("ann_%d", time.Now().UnixNano()),
		Title:            strings.TrimSpace(input.Title),
		Content:          strings.TrimSpace(input.Content),
		Priority:         priorityCode,
		Status:           publishedToStatus(input.Published),
		PublisherAdminId: adminUser.AdminId,
		PublisherName:    adminUser.RealName,
	}
	if input.Published {
		now := time.Now()
		item.PublishTime = &now
	}
	if err := s.platformDAO.CreateSystemAnnouncement(item); err != nil {
		return nil, errs.NewCommonError(http.StatusInternalServerError, "创建系统公告失败: "+err.Error())
	}
	return mapAnnouncement(item), nil
}

func (s *PlatformContentService) UpdateSystemAnnouncement(adminId, announcementId string, input SaveSystemAnnouncementInput) (*SystemAnnouncementDTO, error) {
	adminUser, err := s.ensureAdmin(adminId)
	if err != nil {
		return nil, err
	}
	priorityCode, err := parsePriority(input.Priority)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(announcementId) == "" {
		return nil, errs.NewCommonError(http.StatusBadRequest, "公告ID不能为空")
	}
	if strings.TrimSpace(input.Title) == "" {
		return nil, errs.NewCommonError(http.StatusBadRequest, "公告标题不能为空")
	}
	if strings.TrimSpace(input.Content) == "" {
		return nil, errs.NewCommonError(http.StatusBadRequest, "公告内容不能为空")
	}

	existing, err := s.platformDAO.GetSystemAnnouncementByAnnouncementId(announcementId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewCommonError(http.StatusNotFound, "系统公告不存在")
		}
		return nil, errs.NewCommonError(http.StatusInternalServerError, "查询系统公告失败: "+err.Error())
	}

	updates := map[string]interface{}{
		"title":              strings.TrimSpace(input.Title),
		"content":            strings.TrimSpace(input.Content),
		"priority":           priorityCode,
		"status":             publishedToStatus(input.Published),
		"publisher_admin_id": adminUser.AdminId,
		"publisher_name":     adminUser.RealName,
	}
	if input.Published {
		now := time.Now()
		updates["publish_time"] = &now
	} else {
		updates["publish_time"] = nil
	}
	if err := s.platformDAO.UpdateSystemAnnouncement(announcementId, updates); err != nil {
		return nil, errs.NewCommonError(http.StatusInternalServerError, "更新系统公告失败: "+err.Error())
	}

	existing.Title = strings.TrimSpace(input.Title)
	existing.Content = strings.TrimSpace(input.Content)
	existing.Priority = priorityCode
	existing.Status = publishedToStatus(input.Published)
	existing.PublisherAdminId = adminUser.AdminId
	existing.PublisherName = adminUser.RealName
	if input.Published {
		now := time.Now()
		existing.PublishTime = &now
	} else {
		existing.PublishTime = nil
	}
	return mapAnnouncement(existing), nil
}

func (s *PlatformContentService) DeleteSystemAnnouncement(adminId, announcementId string) error {
	if _, err := s.ensureAdmin(adminId); err != nil {
		return err
	}
	if strings.TrimSpace(announcementId) == "" {
		return errs.NewCommonError(http.StatusBadRequest, "公告ID不能为空")
	}
	if _, err := s.platformDAO.GetSystemAnnouncementByAnnouncementId(announcementId); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errs.NewCommonError(http.StatusNotFound, "系统公告不存在")
		}
		return errs.NewCommonError(http.StatusInternalServerError, "查询系统公告失败: "+err.Error())
	}
	if err := s.platformDAO.DeleteSystemAnnouncement(announcementId); err != nil {
		return errs.NewCommonError(http.StatusInternalServerError, "删除系统公告失败: "+err.Error())
	}
	return nil
}

func (s *PlatformContentService) ListAdminBookshelfItems(query ListQuery) ([]*BookshelfItemDTO, int64, error) {
	page, pageSize := normalizePage(query.Page, query.PageSize)
	items, total, err := s.platformDAO.ListBookshelfItems(query.Keyword, query.Status, false, page, pageSize)
	if err != nil {
		return nil, 0, errs.NewCommonError(http.StatusInternalServerError, "查询平台书架失败: "+err.Error())
	}
	return mapBookshelfList(items), total, nil
}

func (s *PlatformContentService) ListUserBookshelfItems(query ListQuery) ([]*BookshelfItemDTO, int64, error) {
	page, pageSize := normalizePage(query.Page, query.PageSize)
	items, total, err := s.platformDAO.ListBookshelfItems(query.Keyword, nil, true, page, pageSize)
	if err != nil {
		return nil, 0, errs.NewCommonError(http.StatusInternalServerError, "查询平台书架失败: "+err.Error())
	}
	return mapBookshelfList(items), total, nil
}

func (s *PlatformContentService) CreateBookshelfItem(adminId string, input SaveBookshelfItemInput, fileHeader *multipart.FileHeader) (*BookshelfItemDTO, error) {
	adminUser, err := s.ensureAdmin(adminId)
	if err != nil {
		return nil, err
	}
	if err := validateBookshelfInput(input, fileHeader, false, false); err != nil {
		return nil, err
	}

	attachmentName, storageName, mimeType, size, err := storeAttachment(fileHeader)
	if err != nil {
		return nil, err
	}

	item := &platform.BookshelfItem{
		ItemId:                fmt.Sprintf("bks_%d", time.Now().UnixNano()),
		Title:                 strings.TrimSpace(input.Title),
		Description:           strings.TrimSpace(input.Description),
		ContentType:           normalizeContentType(input.ContentType),
		ExternalURL:           normalizeURL(input.ExternalURL),
		AttachmentName:        attachmentName,
		AttachmentStorageName: storageName,
		AttachmentMimeType:    mimeType,
		AttachmentSize:        size,
		SortOrder:             input.SortOrder,
		Status:                publishedToStatus(input.Published),
		CreatorAdminId:        adminUser.AdminId,
		UpdaterAdminId:        adminUser.AdminId,
	}
	if input.Published {
		now := time.Now()
		item.PublishTime = &now
	}
	if err := s.platformDAO.CreateBookshelfItem(item); err != nil {
		removeStoredAttachment(storageName)
		return nil, errs.NewCommonError(http.StatusInternalServerError, "创建平台书架内容失败: "+err.Error())
	}
	return mapBookshelfItem(item), nil
}

func (s *PlatformContentService) UpdateBookshelfItem(adminId, itemId string, input SaveBookshelfItemInput, fileHeader *multipart.FileHeader) (*BookshelfItemDTO, error) {
	adminUser, err := s.ensureAdmin(adminId)
	if err != nil {
		return nil, err
	}
	existing, err := s.platformDAO.GetBookshelfItemByItemId(itemId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewCommonError(http.StatusNotFound, "平台书架内容不存在")
		}
		return nil, errs.NewCommonError(http.StatusInternalServerError, "查询平台书架内容失败: "+err.Error())
	}

	hadAttachment := existing.AttachmentStorageName != ""
	if err := validateBookshelfInput(input, fileHeader, true, hadAttachment && !input.ClearAttachment); err != nil {
		return nil, err
	}

	newAttachmentName, newStorageName, newMimeType, newSize, err := storeAttachment(fileHeader)
	if err != nil {
		return nil, err
	}

	keepAttachment := hadAttachment && !input.ClearAttachment && fileHeader == nil && (normalizeContentType(input.ContentType) == "attachment" || normalizeContentType(input.ContentType) == "mixed")
	updates := map[string]interface{}{
		"title":            strings.TrimSpace(input.Title),
		"description":      strings.TrimSpace(input.Description),
		"content_type":     normalizeContentType(input.ContentType),
		"external_url":     normalizeURL(input.ExternalURL),
		"sort_order":       input.SortOrder,
		"status":           publishedToStatus(input.Published),
		"updater_admin_id": adminUser.AdminId,
	}
	if input.Published {
		now := time.Now()
		updates["publish_time"] = &now
	} else {
		updates["publish_time"] = nil
	}
	if newStorageName != "" {
		updates["attachment_name"] = newAttachmentName
		updates["attachment_storage_name"] = newStorageName
		updates["attachment_mime_type"] = newMimeType
		updates["attachment_size"] = newSize
	} else if !keepAttachment {
		updates["attachment_name"] = ""
		updates["attachment_storage_name"] = ""
		updates["attachment_mime_type"] = ""
		updates["attachment_size"] = 0
	}

	if err := s.platformDAO.UpdateBookshelfItem(itemId, updates); err != nil {
		removeStoredAttachment(newStorageName)
		return nil, errs.NewCommonError(http.StatusInternalServerError, "更新平台书架内容失败: "+err.Error())
	}

	oldStorageName := existing.AttachmentStorageName
	existing.Title = strings.TrimSpace(input.Title)
	existing.Description = strings.TrimSpace(input.Description)
	existing.ContentType = normalizeContentType(input.ContentType)
	existing.ExternalURL = normalizeURL(input.ExternalURL)
	existing.SortOrder = input.SortOrder
	existing.Status = publishedToStatus(input.Published)
	existing.UpdaterAdminId = adminUser.AdminId
	if input.Published {
		now := time.Now()
		existing.PublishTime = &now
	} else {
		existing.PublishTime = nil
	}
	if newStorageName != "" {
		existing.AttachmentName = newAttachmentName
		existing.AttachmentStorageName = newStorageName
		existing.AttachmentMimeType = newMimeType
		existing.AttachmentSize = newSize
	} else if !keepAttachment {
		existing.AttachmentName = ""
		existing.AttachmentStorageName = ""
		existing.AttachmentMimeType = ""
		existing.AttachmentSize = 0
	}
	if oldStorageName != "" && oldStorageName != existing.AttachmentStorageName {
		removeStoredAttachment(oldStorageName)
	}
	return mapBookshelfItem(existing), nil
}

func (s *PlatformContentService) DeleteBookshelfItem(adminId, itemId string) error {
	if _, err := s.ensureAdmin(adminId); err != nil {
		return err
	}
	existing, err := s.platformDAO.GetBookshelfItemByItemId(itemId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errs.NewCommonError(http.StatusNotFound, "平台书架内容不存在")
		}
		return errs.NewCommonError(http.StatusInternalServerError, "查询平台书架内容失败: "+err.Error())
	}
	if err := s.platformDAO.DeleteBookshelfItem(itemId); err != nil {
		return errs.NewCommonError(http.StatusInternalServerError, "删除平台书架内容失败: "+err.Error())
	}
	removeStoredAttachment(existing.AttachmentStorageName)
	return nil
}

func (s *PlatformContentService) GetBookshelfAttachmentFile(itemId string) (*platform.BookshelfItem, string, error) {
	if strings.TrimSpace(itemId) == "" {
		return nil, "", errs.NewCommonError(http.StatusBadRequest, "条目ID不能为空")
	}
	item, err := s.platformDAO.GetBookshelfItemByItemId(itemId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", errs.NewCommonError(http.StatusNotFound, "平台书架内容不存在")
		}
		return nil, "", errs.NewCommonError(http.StatusInternalServerError, "查询平台书架内容失败: "+err.Error())
	}
	if item.AttachmentStorageName == "" {
		return nil, "", errs.NewCommonError(http.StatusNotFound, "该条目没有附件")
	}
	filePath := filepath.Join(platformBookshelfUploadDir, filepath.Base(item.AttachmentStorageName))
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil, "", errs.NewCommonError(http.StatusNotFound, "附件文件不存在")
		}
		return nil, "", errs.NewCommonError(http.StatusInternalServerError, "读取附件失败: "+err.Error())
	}
	return item, filePath, nil
}

func (s *PlatformContentService) ensureAdmin(adminId string) (*adminModel.AdminUser, error) {
	if strings.TrimSpace(adminId) == "" {
		return nil, errs.NewCommonError(http.StatusUnauthorized, "未授权：管理员信息不存在")
	}
	adminUser, err := s.adminDAO.GetAdminUserByAdminId(adminId)
	if err != nil {
		return nil, errs.NewCommonError(http.StatusInternalServerError, "查询管理员信息失败: "+err.Error())
	}
	if adminUser == nil {
		return nil, errs.NewCommonError(http.StatusForbidden, "未授权：管理员不存在")
	}
	if adminUser.Status != 1 {
		return nil, errs.NewCommonError(http.StatusForbidden, "管理员账号已禁用")
	}
	return adminUser, nil
}

func parsePriority(priority string) (int32, error) {
	switch strings.ToLower(strings.TrimSpace(priority)) {
	case "", "normal":
		return 2, nil
	case "high":
		return 3, nil
	case "low":
		return 1, nil
	default:
		return 0, errs.NewCommonError(http.StatusBadRequest, "无效的公告优先级")
	}
}

func priorityLabel(priority int32) string {
	switch priority {
	case 3:
		return "high"
	case 1:
		return "low"
	default:
		return "normal"
	}
}

func normalizeContentType(contentType string) string {
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "link":
		return "link"
	case "attachment":
		return "attachment"
	case "mixed":
		return "mixed"
	default:
		return "text"
	}
}

func contentTypeLabel(contentType string) string {
	switch normalizeContentType(contentType) {
	case "link":
		return "链接"
	case "attachment":
		return "附件"
	case "mixed":
		return "图文+链接/附件"
	default:
		return "文本"
	}
}

func normalizeURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	return raw
}

func validateBookshelfInput(input SaveBookshelfItemInput, fileHeader *multipart.FileHeader, isUpdate bool, hasExistingAttachment bool) error {
	if strings.TrimSpace(input.Title) == "" {
		return errs.NewCommonError(http.StatusBadRequest, "书架标题不能为空")
	}
	contentType := normalizeContentType(input.ContentType)
	if (contentType == "text" || contentType == "mixed") && strings.TrimSpace(input.Description) == "" {
		return errs.NewCommonError(http.StatusBadRequest, "请输入文本内容或说明")
	}
	if contentType == "link" || contentType == "mixed" {
		if normalizeURL(input.ExternalURL) == "" {
			return errs.NewCommonError(http.StatusBadRequest, "请输入有效链接")
		}
		if _, err := url.ParseRequestURI(normalizeURL(input.ExternalURL)); err != nil {
			return errs.NewCommonError(http.StatusBadRequest, "请输入有效链接")
		}
	}
	if contentType == "attachment" || contentType == "mixed" {
		if fileHeader == nil && !hasExistingAttachment {
			return errs.NewCommonError(http.StatusBadRequest, "请上传附件")
		}
	}
	if fileHeader != nil && fileHeader.Size > maxBookshelfAttachmentSize {
		return errs.NewCommonError(http.StatusBadRequest, "附件大小不能超过 50MB")
	}
	if isUpdate && contentType == "text" && input.ClearAttachment && strings.TrimSpace(input.Description) == "" {
		return errs.NewCommonError(http.StatusBadRequest, "纯文本内容不能为空")
	}
	return nil
}

func storeAttachment(fileHeader *multipart.FileHeader) (string, string, string, int64, error) {
	if fileHeader == nil {
		return "", "", "", 0, nil
	}
	if fileHeader.Size > maxBookshelfAttachmentSize {
		return "", "", "", 0, errs.NewCommonError(http.StatusBadRequest, "附件大小不能超过 50MB")
	}
	originalName := strings.TrimSpace(filepath.Base(fileHeader.Filename))
	if originalName == "" || originalName == "." {
		return "", "", "", 0, errs.NewCommonError(http.StatusBadRequest, "附件文件名无效")
	}
	if err := os.MkdirAll(platformBookshelfUploadDir, 0755); err != nil {
		return "", "", "", 0, errs.NewCommonError(http.StatusInternalServerError, "创建附件目录失败: "+err.Error())
	}
	storageName := fmt.Sprintf("%d%s", time.Now().UnixNano(), strings.ToLower(filepath.Ext(originalName)))
	fullPath := filepath.Join(platformBookshelfUploadDir, storageName)

	src, err := fileHeader.Open()
	if err != nil {
		return "", "", "", 0, errs.NewCommonError(http.StatusInternalServerError, "读取上传附件失败: "+err.Error())
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", "", "", 0, errs.NewCommonError(http.StatusInternalServerError, "创建附件文件失败: "+err.Error())
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		removeStoredAttachment(storageName)
		return "", "", "", 0, errs.NewCommonError(http.StatusInternalServerError, "保存附件失败: "+err.Error())
	}
	return originalName, storageName, fileHeader.Header.Get("Content-Type"), fileHeader.Size, nil
}

func removeStoredAttachment(storageName string) {
	storageName = strings.TrimSpace(filepath.Base(storageName))
	if storageName == "" || storageName == "." {
		return
	}
	_ = os.Remove(filepath.Join(platformBookshelfUploadDir, storageName))
}

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func publishedToStatus(published bool) int32 {
	if published {
		return 1
	}
	return 0
}

func mapAnnouncement(item *platform.SystemAnnouncement) *SystemAnnouncementDTO {
	if item == nil {
		return nil
	}
	return &SystemAnnouncementDTO{
		AnnouncementID: item.AnnouncementId,
		Title:          item.Title,
		Content:        item.Content,
		Priority:       priorityLabel(item.Priority),
		Published:      item.Status == systemAnnouncementPublished,
		PublisherName:  item.PublisherName,
		ViewCount:      item.ViewCount,
		PublishTime:    item.PublishTime,
		CreateTime:     item.CreateTime,
		UpdateTime:     item.UpdateTime,
	}
}

func mapAnnouncementList(items []*platform.SystemAnnouncement) []*SystemAnnouncementDTO {
	result := make([]*SystemAnnouncementDTO, 0, len(items))
	for _, item := range items {
		result = append(result, mapAnnouncement(item))
	}
	return result
}

func mapBookshelfItem(item *platform.BookshelfItem) *BookshelfItemDTO {
	if item == nil {
		return nil
	}
	return &BookshelfItemDTO{
		ItemID:           item.ItemId,
		Title:            item.Title,
		Description:      item.Description,
		ContentType:      item.ContentType,
		ContentTypeLabel: contentTypeLabel(item.ContentType),
		ExternalURL:      item.ExternalURL,
		AttachmentName:   item.AttachmentName,
		AttachmentSize:   item.AttachmentSize,
		HasAttachment:    item.AttachmentStorageName != "",
		Published:        item.Status == platformBookshelfPublished,
		SortOrder:        item.SortOrder,
		PublishTime:      item.PublishTime,
		CreateTime:       item.CreateTime,
		UpdateTime:       item.UpdateTime,
	}
}

func mapBookshelfList(items []*platform.BookshelfItem) []*BookshelfItemDTO {
	result := make([]*BookshelfItemDTO, 0, len(items))
	for _, item := range items {
		result = append(result, mapBookshelfItem(item))
	}
	return result
}
