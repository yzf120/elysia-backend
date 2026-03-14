package router

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/authen"
	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/service"
)

var platformContentService *service.PlatformContentService

func RegisterPlatformContentRoutes(protectedRouter *mux.Router) {
	adminRouter := protectedRouter.PathPrefix("/admin").Subrouter()
	adminRouter.Use(authen.AdminAuthMiddleware)

	adminRouter.HandleFunc("/system-announcements", createSystemAnnouncementHandler).Methods("POST")
	adminRouter.HandleFunc("/system-announcements", listAdminSystemAnnouncementsHandler).Methods("GET")
	adminRouter.HandleFunc("/system-announcements/{announcement_id}", updateSystemAnnouncementHandler).Methods("PUT")
	adminRouter.HandleFunc("/system-announcements/{announcement_id}", deleteSystemAnnouncementHandler).Methods("DELETE")

	adminRouter.HandleFunc("/platform-bookshelf", createBookshelfItemHandler).Methods("POST")
	adminRouter.HandleFunc("/platform-bookshelf", listAdminBookshelfItemsHandler).Methods("GET")
	adminRouter.HandleFunc("/platform-bookshelf/{item_id}", updateBookshelfItemHandler).Methods("PUT")
	adminRouter.HandleFunc("/platform-bookshelf/{item_id}", deleteBookshelfItemHandler).Methods("DELETE")

	protectedRouter.HandleFunc("/system-announcements", listUserSystemAnnouncementsHandler).Methods("GET")
	protectedRouter.HandleFunc("/platform-bookshelf", listUserBookshelfItemsHandler).Methods("GET")
	protectedRouter.HandleFunc("/platform-bookshelf/files/{item_id}/view", viewBookshelfAttachmentHandler).Methods("GET")
	protectedRouter.HandleFunc("/platform-bookshelf/files/{item_id}/download", downloadBookshelfAttachmentHandler).Methods("GET")
}

func createSystemAnnouncementHandler(w http.ResponseWriter, r *http.Request) {
	adminId, ok := authen.GetAdminIDFromContext(r.Context())
	if !ok || adminId == "" {
		writeError(w, http.StatusUnauthorized, "未授权：需要管理员权限")
		return
	}

	var req service.SaveSystemAnnouncementInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求参数有误")
		return
	}

	announcement, err := platformContentService.CreateSystemAnnouncement(adminId, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code":         consts.SuccessCode,
		"message":      "创建成功",
		"announcement": announcement,
	})
}

func listAdminSystemAnnouncementsHandler(w http.ResponseWriter, r *http.Request) {
	query := buildListQuery(r, true)
	announcements, total, err := platformContentService.ListAdminSystemAnnouncements(query)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code":          consts.SuccessCode,
		"message":       consts.MessageQuerySuccess,
		"announcements": announcements,
		"total":         total,
		"page":          query.Page,
		"page_size":     query.PageSize,
	})
}

func updateSystemAnnouncementHandler(w http.ResponseWriter, r *http.Request) {
	adminId, ok := authen.GetAdminIDFromContext(r.Context())
	if !ok || adminId == "" {
		writeError(w, http.StatusUnauthorized, "未授权：需要管理员权限")
		return
	}

	var req service.SaveSystemAnnouncementInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求参数有误")
		return
	}

	announcementId := mux.Vars(r)["announcement_id"]
	announcement, err := platformContentService.UpdateSystemAnnouncement(adminId, announcementId, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code":         consts.SuccessCode,
		"message":      consts.MessageUpdateSuccess,
		"announcement": announcement,
	})
}

func deleteSystemAnnouncementHandler(w http.ResponseWriter, r *http.Request) {
	adminId, ok := authen.GetAdminIDFromContext(r.Context())
	if !ok || adminId == "" {
		writeError(w, http.StatusUnauthorized, "未授权：需要管理员权限")
		return
	}

	announcementId := mux.Vars(r)["announcement_id"]
	if err := platformContentService.DeleteSystemAnnouncement(adminId, announcementId); err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code":    consts.SuccessCode,
		"message": consts.MessageDeleteSuccess,
	})
}

func listUserSystemAnnouncementsHandler(w http.ResponseWriter, r *http.Request) {
	query := buildListQuery(r, false)
	announcements, total, err := platformContentService.ListUserSystemAnnouncements(query)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code":          consts.SuccessCode,
		"message":       consts.MessageQuerySuccess,
		"announcements": announcements,
		"total":         total,
		"page":          query.Page,
		"page_size":     query.PageSize,
	})
}

func createBookshelfItemHandler(w http.ResponseWriter, r *http.Request) {
	adminId, ok := authen.GetAdminIDFromContext(r.Context())
	if !ok || adminId == "" {
		writeError(w, http.StatusUnauthorized, "未授权：需要管理员权限")
		return
	}

	input, fileHeader, err := parseBookshelfForm(r)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	item, err := platformContentService.CreateBookshelfItem(adminId, input, fileHeader)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code":           consts.SuccessCode,
		"message":        "创建成功",
		"bookshelf_item": item,
	})
}

func listAdminBookshelfItemsHandler(w http.ResponseWriter, r *http.Request) {
	query := buildListQuery(r, true)
	items, total, err := platformContentService.ListAdminBookshelfItems(query)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code":            consts.SuccessCode,
		"message":         consts.MessageQuerySuccess,
		"bookshelf_items": items,
		"total":           total,
		"page":            query.Page,
		"page_size":       query.PageSize,
	})
}

func updateBookshelfItemHandler(w http.ResponseWriter, r *http.Request) {
	adminId, ok := authen.GetAdminIDFromContext(r.Context())
	if !ok || adminId == "" {
		writeError(w, http.StatusUnauthorized, "未授权：需要管理员权限")
		return
	}

	input, fileHeader, err := parseBookshelfForm(r)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	itemId := mux.Vars(r)["item_id"]
	item, err := platformContentService.UpdateBookshelfItem(adminId, itemId, input, fileHeader)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code":           consts.SuccessCode,
		"message":        consts.MessageUpdateSuccess,
		"bookshelf_item": item,
	})
}

func deleteBookshelfItemHandler(w http.ResponseWriter, r *http.Request) {
	adminId, ok := authen.GetAdminIDFromContext(r.Context())
	if !ok || adminId == "" {
		writeError(w, http.StatusUnauthorized, "未授权：需要管理员权限")
		return
	}

	itemId := mux.Vars(r)["item_id"]
	if err := platformContentService.DeleteBookshelfItem(adminId, itemId); err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code":    consts.SuccessCode,
		"message": consts.MessageDeleteSuccess,
	})
}

func listUserBookshelfItemsHandler(w http.ResponseWriter, r *http.Request) {
	query := buildListQuery(r, false)
	items, total, err := platformContentService.ListUserBookshelfItems(query)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code":            consts.SuccessCode,
		"message":         consts.MessageQuerySuccess,
		"bookshelf_items": items,
		"total":           total,
		"page":            query.Page,
		"page_size":       query.PageSize,
	})
}

func viewBookshelfAttachmentHandler(w http.ResponseWriter, r *http.Request) {
	itemId := mux.Vars(r)["item_id"]
	item, filePath, err := platformContentService.GetBookshelfAttachmentFile(itemId)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	setFileHeaders(w)
	if strings.TrimSpace(item.AttachmentMimeType) != "" {
		w.Header().Set("Content-Type", item.AttachmentMimeType)
	}
	w.Header().Set("Content-Disposition", "inline; filename*=UTF-8''"+url.PathEscape(item.AttachmentName))
	http.ServeFile(w, r, filePath)
}

func downloadBookshelfAttachmentHandler(w http.ResponseWriter, r *http.Request) {
	itemId := mux.Vars(r)["item_id"]
	item, filePath, err := platformContentService.GetBookshelfAttachmentFile(itemId)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	setFileHeaders(w)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+url.PathEscape(item.AttachmentName))
	http.ServeFile(w, r, filePath)
}

func buildListQuery(r *http.Request, allowStatus bool) service.ListQuery {
	page := parseIntWithDefault(r.URL.Query().Get("page"), 1)
	pageSize := parseIntWithDefault(r.URL.Query().Get("page_size"), 10)
	query := service.ListQuery{
		Page:     page,
		PageSize: pageSize,
		Keyword:  strings.TrimSpace(r.URL.Query().Get("keyword")),
	}
	if allowStatus {
		if statusStr := strings.TrimSpace(r.URL.Query().Get("status")); statusStr != "" {
			statusVal := int32(parseIntWithDefault(statusStr, -1))
			if statusVal >= 0 {
				query.Status = &statusVal
			}
		}
	}
	return query
}

func parseBookshelfForm(r *http.Request) (service.SaveBookshelfItemInput, *multipart.FileHeader, error) {
	if err := r.ParseMultipartForm(60 << 20); err != nil {
		return service.SaveBookshelfItemInput{}, nil, errs.NewCommonError(http.StatusBadRequest, "表单解析失败")
	}
	input := service.SaveBookshelfItemInput{
		Title:           strings.TrimSpace(r.FormValue("title")),
		Description:     strings.TrimSpace(r.FormValue("description")),
		ContentType:     strings.TrimSpace(r.FormValue("content_type")),
		ExternalURL:     strings.TrimSpace(r.FormValue("external_url")),
		Published:       parseBoolWithDefault(r.FormValue("published")),
		SortOrder:       int32(parseIntWithDefault(r.FormValue("sort_order"), 0)),
		ClearAttachment: parseBoolWithDefault(r.FormValue("clear_attachment")),
	}

	file, fileHeader, err := r.FormFile("attachment")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return input, nil, nil
		}
		return service.SaveBookshelfItemInput{}, nil, errs.NewCommonError(http.StatusBadRequest, "读取附件失败")
	}
	_ = file.Close()
	return input, fileHeader, nil
}

func writeServiceError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	message := "服务内部错误"
	var commonErr *errs.CommonError
	if errors.As(err, &commonErr) {
		message = commonErr.Message
		if code, convErr := strconv.Atoi(commonErr.Code); convErr == nil && code >= 100 && code <= 599 {
			statusCode = code
		}
	} else if err != nil {
		message = err.Error()
	}
	writeError(w, statusCode, message)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]interface{}{
		"code":    statusCode,
		"message": message,
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	setJSONHeaders(w)
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func setJSONHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
}

func setFileHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
}

func parseIntWithDefault(raw string, defaultValue int) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return defaultValue
	}
	return value
}

func parseBoolWithDefault(raw string) bool {
	value, err := strconv.ParseBool(strings.TrimSpace(raw))
	if err != nil {
		return false
	}
	return value
}
