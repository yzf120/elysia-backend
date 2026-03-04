package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/service_impl"
)

var chapterService *service_impl.ChapterServiceImpl

// registerChapter 注册章节相关路由
func registerChapter(publicRouter *mux.Router, protectedRouter *mux.Router) {
	// 教师操作：增删改（仅教师）
	protectedRouter.HandleFunc("/teacher/chapter/create", createChapterHandler).Methods("POST")
	protectedRouter.HandleFunc("/teacher/chapter/update", updateChapterHandler).Methods("POST")
	protectedRouter.HandleFunc("/teacher/chapter/delete", deleteChapterHandler).Methods("POST")
	protectedRouter.HandleFunc("/teacher/chapter/reorder", reorderChaptersHandler).Methods("POST")

	// 小节操作（仅教师）
	protectedRouter.HandleFunc("/teacher/section/create", createSectionHandler).Methods("POST")
	protectedRouter.HandleFunc("/teacher/section/update", updateSectionHandler).Methods("POST")
	protectedRouter.HandleFunc("/teacher/section/delete", deleteSectionHandler).Methods("POST")
	protectedRouter.HandleFunc("/teacher/section/reorder", reorderSectionsHandler).Methods("POST")

	// 查询：师生共用
	protectedRouter.HandleFunc("/class/chapters", getClassChaptersHandler).Methods("POST")
}

// createChapterHandler 创建章节
func createChapterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)
	req := &service_impl.CreateChapterRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := chapterService.CreateChapter(ctx, req)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.Code != 0 {
		errResp := &errs.BaseResponse{Data: nil, Error: errs.NewError(int(resp.Code), resp.Message)}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
		return
	}
	writeSuccessResponse(w, resp)
}

// updateChapterHandler 更新章节
func updateChapterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)
	req := &service_impl.UpdateChapterRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := chapterService.UpdateChapter(ctx, req)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.Code != 0 {
		errResp := &errs.BaseResponse{Data: nil, Error: errs.NewError(int(resp.Code), resp.Message)}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
		return
	}
	writeSuccessResponse(w, resp)
}

// deleteChapterHandler 删除章节
func deleteChapterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)
	req := &service_impl.DeleteChapterRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := chapterService.DeleteChapter(ctx, req)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.Code != 0 {
		errResp := &errs.BaseResponse{Data: nil, Error: errs.NewError(int(resp.Code), resp.Message)}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
		return
	}
	writeSuccessResponse(w, resp)
}

// reorderChaptersHandler 调整章节排序
func reorderChaptersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)
	req := &service_impl.ReorderChaptersRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := chapterService.ReorderChapters(ctx, req)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.Code != 0 {
		errResp := &errs.BaseResponse{Data: nil, Error: errs.NewError(int(resp.Code), resp.Message)}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
		return
	}
	writeSuccessResponse(w, resp)
}

// createSectionHandler 创建小节
func createSectionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)
	req := &service_impl.CreateSectionRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := chapterService.CreateSection(ctx, req)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.Code != 0 {
		errResp := &errs.BaseResponse{Data: nil, Error: errs.NewError(int(resp.Code), resp.Message)}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
		return
	}
	writeSuccessResponse(w, resp)
}

// updateSectionHandler 更新小节
func updateSectionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)
	req := &service_impl.UpdateSectionRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := chapterService.UpdateSection(ctx, req)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.Code != 0 {
		errResp := &errs.BaseResponse{Data: nil, Error: errs.NewError(int(resp.Code), resp.Message)}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
		return
	}
	writeSuccessResponse(w, resp)
}

// deleteSectionHandler 删除小节
func deleteSectionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)
	req := &service_impl.DeleteSectionRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := chapterService.DeleteSection(ctx, req)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.Code != 0 {
		errResp := &errs.BaseResponse{Data: nil, Error: errs.NewError(int(resp.Code), resp.Message)}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
		return
	}
	writeSuccessResponse(w, resp)
}

// reorderSectionsHandler 调整小节排序
func reorderSectionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)
	req := &service_impl.ReorderSectionsRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := chapterService.ReorderSections(ctx, req)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.Code != 0 {
		errResp := &errs.BaseResponse{Data: nil, Error: errs.NewError(int(resp.Code), resp.Message)}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
		return
	}
	writeSuccessResponse(w, resp)
}

// getClassChaptersHandler 查询班级章节列表（师生共用）
func getClassChaptersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)
	req := &service_impl.GetClassChaptersRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := chapterService.GetClassChapters(ctx, req)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.Code != 0 {
		errResp := &errs.BaseResponse{Data: nil, Error: errs.NewError(int(resp.Code), resp.Message)}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
		return
	}
	writeSuccessResponse(w, resp)
}
