package router

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/authen"
	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/service"
)

var adminUserManagementService *service.AdminUserManagementService

func RegisterAdminUserManagementRoutes(protectedRouter *mux.Router) {
	adminRouter := protectedRouter.PathPrefix("/admin/users").Subrouter()
	adminRouter.Use(authen.AdminAuthMiddleware)

	adminRouter.HandleFunc("/students", listAdminStudentsHandler).Methods("GET")
	adminRouter.HandleFunc("/students/status", batchUpdateStudentStatusHandler).Methods("POST")
	adminRouter.HandleFunc("/students/export", exportStudentsHandler).Methods("GET")

	adminRouter.HandleFunc("/teachers", listAdminTeachersHandler).Methods("GET")
	adminRouter.HandleFunc("/teachers/status", batchUpdateTeacherStatusHandler).Methods("POST")
	adminRouter.HandleFunc("/teachers/export", exportTeachersHandler).Methods("GET")
}

func listAdminStudentsHandler(w http.ResponseWriter, r *http.Request) {
	input := service.AdminStudentListInput{
		Page:     parseIntWithDefault(r.URL.Query().Get("page"), 1),
		PageSize: parseIntWithDefault(r.URL.Query().Get("page_size"), 10),
		Keyword:  strings.TrimSpace(r.URL.Query().Get("keyword")),
		Major:    strings.TrimSpace(r.URL.Query().Get("major")),
		Grade:    strings.TrimSpace(r.URL.Query().Get("grade")),
		Status:   strings.TrimSpace(r.URL.Query().Get("status")),
	}
	students, total, err := adminUserManagementService.ListStudents(input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code":      consts.SuccessCode,
		"message":   consts.MessageQuerySuccess,
		"students":  students,
		"total":     total,
		"page":      input.Page,
		"page_size": input.PageSize,
	})
}

func listAdminTeachersHandler(w http.ResponseWriter, r *http.Request) {
	input := service.AdminTeacherListInput{
		Page:               parseIntWithDefault(r.URL.Query().Get("page"), 1),
		PageSize:           parseIntWithDefault(r.URL.Query().Get("page_size"), 10),
		Keyword:            strings.TrimSpace(r.URL.Query().Get("keyword")),
		Department:         strings.TrimSpace(r.URL.Query().Get("department")),
		VerificationStatus: strings.TrimSpace(r.URL.Query().Get("verification_status")),
		Status:             strings.TrimSpace(r.URL.Query().Get("status")),
	}
	teachers, total, err := adminUserManagementService.ListTeachers(input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code":      consts.SuccessCode,
		"message":   consts.MessageQuerySuccess,
		"teachers":  teachers,
		"total":     total,
		"page":      input.Page,
		"page_size": input.PageSize,
	})
}

func batchUpdateStudentStatusHandler(w http.ResponseWriter, r *http.Request) {
	var input service.AdminBatchStudentStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "请求参数有误")
		return
	}
	updatedCount, err := adminUserManagementService.BatchUpdateStudentsStatus(input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code":          consts.SuccessCode,
		"message":       consts.MessageUpdateSuccess,
		"updated_count": updatedCount,
	})
}

func batchUpdateTeacherStatusHandler(w http.ResponseWriter, r *http.Request) {
	var input service.AdminBatchTeacherStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "请求参数有误")
		return
	}
	updatedCount, err := adminUserManagementService.BatchUpdateTeachersStatus(input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"code":          consts.SuccessCode,
		"message":       consts.MessageUpdateSuccess,
		"updated_count": updatedCount,
	})
}

func exportStudentsHandler(w http.ResponseWriter, r *http.Request) {
	input := service.AdminStudentListInput{
		Keyword: strings.TrimSpace(r.URL.Query().Get("keyword")),
		Major:   strings.TrimSpace(r.URL.Query().Get("major")),
		Grade:   strings.TrimSpace(r.URL.Query().Get("grade")),
		Status:  strings.TrimSpace(r.URL.Query().Get("status")),
	}
	content, fileName, err := adminUserManagementService.ExportStudents(input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	setFileHeaders(w)
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+url.PathEscape(fileName))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

func exportTeachersHandler(w http.ResponseWriter, r *http.Request) {
	input := service.AdminTeacherListInput{
		Keyword:            strings.TrimSpace(r.URL.Query().Get("keyword")),
		Department:         strings.TrimSpace(r.URL.Query().Get("department")),
		VerificationStatus: strings.TrimSpace(r.URL.Query().Get("verification_status")),
		Status:             strings.TrimSpace(r.URL.Query().Get("status")),
	}
	content, fileName, err := adminUserManagementService.ExportTeachers(input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	setFileHeaders(w)
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+url.PathEscape(fileName))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}
