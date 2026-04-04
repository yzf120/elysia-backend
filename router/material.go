package router

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/service"
)

var materialService *service.MaterialService

// registerMaterial 注册学习资料相关路由
func registerMaterial(publicRouter *mux.Router, protectedRouter *mux.Router) {
	// 教师操作：上传/删除学习资料
	protectedRouter.HandleFunc("/teacher/material/upload", uploadMaterialHandler).Methods("POST")
	protectedRouter.HandleFunc("/teacher/material/delete", deleteMaterialHandler).Methods("POST")

	// 查询：师生共用
	protectedRouter.HandleFunc("/material/list", listMaterialsHandler).Methods("POST")

	// 文件访问：师生共用
	protectedRouter.HandleFunc("/material/file/{material_id}/view", viewMaterialFileHandler).Methods("GET")
	protectedRouter.HandleFunc("/material/file/{material_id}/download", downloadMaterialFileHandler).Methods("GET")
}

// uploadMaterialHandler 教师上传学习资料（multipart/form-data）
func uploadMaterialHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	// 解析 multipart form（最大 200MB）
	if err := r.ParseMultipartForm(200 << 20); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "表单解析失败: "+err.Error())
		return
	}

	teacherId := strings.TrimSpace(r.FormValue("teacher_id"))
	sectionId := strings.TrimSpace(r.FormValue("section_id"))
	title := strings.TrimSpace(r.FormValue("title"))
	description := strings.TrimSpace(r.FormValue("description"))
	materialType := strings.TrimSpace(r.FormValue("material_type"))

	// 获取上传文件（text类型可以没有文件）
	file, fileHeader, err := r.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		writeErrorResponse(w, http.StatusBadRequest, "读取文件失败: "+err.Error())
		return
	}
	if file != nil {
		_ = file.Close()
	}

	material, svcErr := materialService.CreateMaterial(teacherId, sectionId, title, description, materialType, fileHeader)
	if svcErr != nil {
		writeErrorResponse(w, http.StatusBadRequest, svcErr.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message":  "上传成功",
		"material": material,
	})
}

// deleteMaterialHandler 教师删除学习资料
func deleteMaterialHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	var reqBody struct {
		TeacherId  string `json:"teacher_id"`
		MaterialId string `json:"material_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := materialService.DeleteMaterial(reqBody.TeacherId, reqBody.MaterialId); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "删除成功",
	})
}

// listMaterialsHandler 查询小节下的学习资料列表（师生共用）
func listMaterialsHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	var reqBody struct {
		SectionId string `json:"section_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	materials, err := materialService.GetMaterialsBySection(reqBody.SectionId)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 构造返回数据（隐藏 file_path 对学生不可见）
	type MaterialInfo struct {
		MaterialId   string `json:"material_id"`
		SectionId    string `json:"section_id"`
		Title        string `json:"title"`
		Description  string `json:"description"`
		MaterialType string `json:"material_type"`
		FileName     string `json:"file_name"`
		FileSize     int64  `json:"file_size"`
		MimeType     string `json:"mime_type"`
		SortOrder    int32  `json:"sort_order"`
		CreateTime   string `json:"create_time"`
	}

	list := make([]MaterialInfo, 0, len(materials))
	for _, m := range materials {
		list = append(list, MaterialInfo{
			MaterialId:   m.MaterialId,
			SectionId:    m.SectionId,
			Title:        m.Title,
			Description:  m.Description,
			MaterialType: m.MaterialType,
			FileName:     m.FileName,
			FileSize:     m.FileSize,
			MimeType:     m.MimeType,
			SortOrder:    m.SortOrder,
			CreateTime:   m.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message":   consts.MessageQuerySuccess,
		"materials": list,
	})
}

// viewMaterialFileHandler 查看学习资料文件（inline）
func viewMaterialFileHandler(w http.ResponseWriter, r *http.Request) {
	materialId := mux.Vars(r)["material_id"]
	material, filePath, err := materialService.GetMaterialFile(materialId)
	if err != nil {
		setResponseHeaders(w)
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	if strings.TrimSpace(material.MimeType) != "" {
		w.Header().Set("Content-Type", material.MimeType)
	}
	w.Header().Set("Content-Disposition", "inline; filename*=UTF-8''"+url.PathEscape(material.FileName))
	http.ServeFile(w, r, filePath)
}

// downloadMaterialFileHandler 下载学习资料文件（attachment）
// 注意：视频类型不支持下载
func downloadMaterialFileHandler(w http.ResponseWriter, r *http.Request) {
	materialId := mux.Vars(r)["material_id"]
	material, filePath, err := materialService.GetMaterialFile(materialId)
	if err != nil {
		setResponseHeaders(w)
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	// 视频类型不允许下载
	if material.MaterialType == "video" {
		setResponseHeaders(w)
		writeErrorResponse(w, http.StatusForbidden, "视频资料仅支持在线查看，不支持下载")
		return
	}

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+url.PathEscape(material.FileName))
	http.ServeFile(w, r, filePath)
}
