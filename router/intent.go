package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/model/intent"
	"github.com/yzf120/elysia-backend/service"
)

var intentService *service.IntentService

// RegisterIntentRoutes 注册意图管理路由（管理员接口，需要认证）
func RegisterIntentRoutes(protectedRouter *mux.Router) {
	// ==================== 意图字典管理 ====================
	protectedRouter.HandleFunc("/admin/intent/dict", listIntentDictsHandler).Methods("GET")
	protectedRouter.HandleFunc("/admin/intent/dict", createIntentDictHandler).Methods("POST")
	protectedRouter.HandleFunc("/admin/intent/dict/{id}", getIntentDictHandler).Methods("GET")
	protectedRouter.HandleFunc("/admin/intent/dict/{id}", updateIntentDictHandler).Methods("PUT")
	protectedRouter.HandleFunc("/admin/intent/dict/{id}", deleteIntentDictHandler).Methods("DELETE")
	protectedRouter.HandleFunc("/admin/intent/dict/{id}/status", updateIntentDictStatusHandler).Methods("POST")

	// ==================== 意图提示词模板管理 ====================
	protectedRouter.HandleFunc("/admin/intent/prompt", listPromptTemplatesHandler).Methods("GET")
	protectedRouter.HandleFunc("/admin/intent/prompt", createPromptTemplateHandler).Methods("POST")
	protectedRouter.HandleFunc("/admin/intent/prompt/{id}", getPromptTemplateHandler).Methods("GET")
	protectedRouter.HandleFunc("/admin/intent/prompt/{id}", updatePromptTemplateHandler).Methods("PUT")
	protectedRouter.HandleFunc("/admin/intent/prompt/{id}", deletePromptTemplateHandler).Methods("DELETE")
	protectedRouter.HandleFunc("/admin/intent/prompt/{id}/status", updatePromptTemplateStatusHandler).Methods("POST")

	// ==================== 意图记录查询（只读） ====================
	protectedRouter.HandleFunc("/admin/intent/records", listIntentRecordsHandler).Methods("GET")
	protectedRouter.HandleFunc("/admin/intent/stats", getIntentStatsHandler).Methods("GET")
}

// ==================== 意图字典处理器 ====================

// listIntentDictsHandler 查询意图字典列表
// GET /admin/intent/dict?page=1&page_size=20&intent_level1=xxx&is_valid=1
func listIntentDictsHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	page := getQueryInt(r, "page", 1)
	pageSize := getQueryInt(r, "page_size", 20)
	intentLevel1 := r.URL.Query().Get("intent_level1")
	isValid := getQueryInt(r, "is_valid", -1)

	if pageSize > 100 {
		pageSize = 100
	}

	dicts, total, err := intentService.ListIntentDicts(page, pageSize, intentLevel1, isValid)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "查询意图字典失败: "+err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"list":  dicts,
		"total": total,
		"page":  page,
	})
}

// createIntentDictHandler 创建意图字典
// POST /admin/intent/dict
func createIntentDictHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	var dict intent.IntentDict
	if err := json.NewDecoder(r.Body).Decode(&dict); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 参数校验
	if dict.IntentCode == "" {
		writeErrorResponse(w, http.StatusBadRequest, "意图编码不能为空")
		return
	}
	if dict.IntentLevel1 == "" {
		writeErrorResponse(w, http.StatusBadRequest, "一级意图不能为空")
		return
	}
	if dict.IntentLevel2 == "" {
		writeErrorResponse(w, http.StatusBadRequest, "二级子意图不能为空")
		return
	}
	if dict.AgentRoute == "" {
		writeErrorResponse(w, http.StatusBadRequest, "路由Agent不能为空")
		return
	}

	if err := intentService.CreateIntentDict(&dict); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "创建成功",
		"data":    dict,
	})
}

// getIntentDictHandler 获取单个意图字典
// GET /admin/intent/dict/{id}
func getIntentDictHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	id, err := getPathInt(r, "id")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "无效的ID")
		return
	}

	dict, err := intentService.GetIntentDictById(id)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	writeSuccessResponse(w, dict)
}

// updateIntentDictHandler 更新意图字典
// PUT /admin/intent/dict/{id}
func updateIntentDictHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	id, err := getPathInt(r, "id")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "无效的ID")
		return
	}

	var dict intent.IntentDict
	if err := json.NewDecoder(r.Body).Decode(&dict); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}
	dict.Id = id

	if err := intentService.UpdateIntentDict(&dict); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "更新成功",
	})
}

// deleteIntentDictHandler 删除意图字典（软删除）
// DELETE /admin/intent/dict/{id}
func deleteIntentDictHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	id, err := getPathInt(r, "id")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := intentService.DeleteIntentDict(id); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "删除成功",
	})
}

// updateIntentDictStatusHandler 更新意图字典状态
// POST /admin/intent/dict/{id}/status
func updateIntentDictStatusHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	id, err := getPathInt(r, "id")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "无效的ID")
		return
	}

	var reqBody struct {
		IsValid int `json:"is_valid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := intentService.UpdateIntentDictStatus(id, reqBody.IsValid); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "状态更新成功",
	})
}

// ==================== 意图提示词模板处理器 ====================

// listPromptTemplatesHandler 查询模板列表
// GET /admin/intent/prompt?page=1&page_size=20&intent_code=xxx&template_type=xxx
func listPromptTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	page := getQueryInt(r, "page", 1)
	pageSize := getQueryInt(r, "page_size", 20)
	intentCode := r.URL.Query().Get("intent_code")
	templateType := r.URL.Query().Get("template_type")

	if pageSize > 100 {
		pageSize = 100
	}

	tpls, total, err := intentService.ListPromptTemplates(page, pageSize, intentCode, templateType)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "查询模板列表失败: "+err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"list":  tpls,
		"total": total,
		"page":  page,
	})
}

// createPromptTemplateHandler 创建提示词模板
// POST /admin/intent/prompt
func createPromptTemplateHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	var tpl intent.IntentPromptTemplate
	if err := json.NewDecoder(r.Body).Decode(&tpl); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 参数校验
	if tpl.IntentCode == "" {
		writeErrorResponse(w, http.StatusBadRequest, "意图编码不能为空")
		return
	}
	if tpl.TemplateType == "" {
		writeErrorResponse(w, http.StatusBadRequest, "模板类型不能为空")
		return
	}
	if tpl.TemplateName == "" {
		writeErrorResponse(w, http.StatusBadRequest, "模板名称不能为空")
		return
	}
	if tpl.TemplateContent == "" {
		writeErrorResponse(w, http.StatusBadRequest, "模板内容不能为空")
		return
	}

	if err := intentService.CreatePromptTemplate(&tpl); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "创建成功",
		"data":    tpl,
	})
}

// getPromptTemplateHandler 获取单个模板
// GET /admin/intent/prompt/{id}
func getPromptTemplateHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	id, err := getPathInt(r, "id")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "无效的ID")
		return
	}

	tpl, err := intentService.GetPromptTemplateById(id)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	writeSuccessResponse(w, tpl)
}

// updatePromptTemplateHandler 更新模板
// PUT /admin/intent/prompt/{id}
func updatePromptTemplateHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	id, err := getPathInt(r, "id")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "无效的ID")
		return
	}

	var tpl intent.IntentPromptTemplate
	if err := json.NewDecoder(r.Body).Decode(&tpl); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}
	tpl.Id = id

	if err := intentService.UpdatePromptTemplate(&tpl); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "更新成功",
	})
}

// deletePromptTemplateHandler 删除模板
// DELETE /admin/intent/prompt/{id}
func deletePromptTemplateHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	id, err := getPathInt(r, "id")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "无效的ID")
		return
	}

	if err := intentService.DeletePromptTemplate(id); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "删除成功",
	})
}

// updatePromptTemplateStatusHandler 更新模板状态
// POST /admin/intent/prompt/{id}/status
func updatePromptTemplateStatusHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	id, err := getPathInt(r, "id")
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "无效的ID")
		return
	}

	var reqBody struct {
		IsActive int `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := intentService.UpdatePromptTemplateStatus(id, reqBody.IsActive); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "状态更新成功",
	})
}

// ==================== 意图记录查询处理器 ====================

// listIntentRecordsHandler 查询意图记录列表
// GET /admin/intent/records?page=1&page_size=20&user_id=xxx&intent_code=xxx&intent_level1=xxx
func listIntentRecordsHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	page := getQueryInt(r, "page", 1)
	pageSize := getQueryInt(r, "page_size", 20)
	userId := r.URL.Query().Get("user_id")
	intentCode := r.URL.Query().Get("intent_code")
	intentLevel1 := r.URL.Query().Get("intent_level1")

	if pageSize > 100 {
		pageSize = 100
	}

	records, total, err := intentService.ListIntentRecords(page, pageSize, userId, intentCode, intentLevel1)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "查询意图记录失败: "+err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"list":  records,
		"total": total,
		"page":  page,
	})
}

// getIntentStatsHandler 获取意图统计数据
// GET /admin/intent/stats
func getIntentStatsHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	stats, err := intentService.GetIntentStats()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "获取意图统计失败: "+err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"intent_count": stats,
	})
}

// ==================== 工具函数 ====================

// getQueryInt 从URL查询参数获取整数值
func getQueryInt(r *http.Request, key string, defaultValue int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return intVal
}

// getPathInt 从URL路径参数获取整数值
func getPathInt(r *http.Request, key string) (int, error) {
	vars := mux.Vars(r)
	val := vars[key]
	return strconv.Atoi(val)
}
