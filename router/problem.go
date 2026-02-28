package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/errs"
	problemReq "github.com/yzf120/elysia-backend/model/problem/req"
	"github.com/yzf120/elysia-backend/service_impl"
)

var (
	problemService *service_impl.ProblemServiceImpl
)

// registerProblem 注册题目相关路由
func registerProblem(publicRouter *mux.Router, protectedRouter *mux.Router) {
	// 增删改：仅教师可操作（受保护路由，教师路由前缀）
	protectedRouter.HandleFunc("/teacher/problem/create", createProblemHandler).Methods("POST")
	protectedRouter.HandleFunc("/teacher/problem/update", updateProblemHandler).Methods("POST")
	protectedRouter.HandleFunc("/teacher/problem/delete", deleteProblemHandler).Methods("POST")

	// 查询：学生和教师均可调用（受保护路由，通用路由前缀）
	protectedRouter.HandleFunc("/problem/get", getProblemHandler).Methods("GET")
}

// createProblemHandler 创建题目处理器
func createProblemHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	request := &problemReq.CreateProblemRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	resp, err := problemService.CreateProblem(ctx, request)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// updateProblemHandler 更新题目处理器
func updateProblemHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	request := &problemReq.UpdateProblemRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	resp, err := problemService.UpdateProblem(ctx, request)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// deleteProblemHandler 删除题目处理器
func deleteProblemHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	request := &problemReq.DeleteProblemRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	resp, err := problemService.DeleteProblem(ctx, request)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// getProblemHandler 查询题目处理器（学生和教师均可调用）
func getProblemHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, "参数id无效"),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	request := &problemReq.GetProblemRequest{Id: id}
	resp, err := problemService.GetProblem(ctx, request)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}
