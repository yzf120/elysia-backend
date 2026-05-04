package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/service"
)

// 模型配置服务（在 router.Init() 中初始化）
var aiModelConfigService *service.AIModelConfigService

// RegisterAIModelConfigRoutes 注册AI模型配置管理路由（管理员）
func RegisterAIModelConfigRoutes(protectedRouter *mux.Router) {
	// 获取所有模型配置列表（管理员）
	protectedRouter.HandleFunc("/admin/ai-models", adminListAIModelsHandler).Methods("GET", "OPTIONS")
	// 切换模型启用/禁用状态（管理员）
	protectedRouter.HandleFunc("/admin/ai-models/toggle", adminToggleAIModelHandler).Methods("POST", "OPTIONS")
}

// adminListAIModelsHandler 获取所有AI模型配置列表
// GET /admin/ai-models
func adminListAIModelsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	setResponseHeaders(w)

	models, err := aiModelConfigService.ListAllModels()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "查询模型配置失败: "+err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"models": models,
	})
}

// adminToggleAIModelHandler 切换模型启用/禁用状态
// POST /admin/ai-models/toggle
// Body: {"model_id": "doubao-seed-2-0-lite-260215", "enabled": true}
func adminToggleAIModelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	setResponseHeaders(w)

	var reqBody struct {
		ModelId string `json:"model_id"`
		Enabled bool   `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if reqBody.ModelId == "" {
		writeErrorResponse(w, http.StatusBadRequest, "model_id 不能为空")
		return
	}

	if err := aiModelConfigService.ToggleModelStatus(reqBody.ModelId, reqBody.Enabled); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	action := "禁用"
	if reqBody.Enabled {
		action = "启用"
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message":  "模型已" + action,
		"model_id": reqBody.ModelId,
		"enabled":  reqBody.Enabled,
	})
}
