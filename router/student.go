package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/student/req"
	"github.com/yzf120/elysia-backend/service_impl"
)

var studentService = service_impl.NewStudentServiceImpl()

// registerStudent 注册学生相关路由
func registerStudent(router *mux.Router) {
	router.HandleFunc("/api/student/create", createStudentHandler).Methods("POST")
	router.HandleFunc("/api/student/get", getStudentHandler).Methods("GET")
	router.HandleFunc("/api/student/update", updateStudentHandler).Methods("POST")
	router.HandleFunc("/api/student/list", listStudentsHandler).Methods("POST")
	router.HandleFunc("/api/student/update-progress", updateLearningProgressHandler).Methods("POST")
}

// createStudentHandler 创建学生信息处理器
func createStudentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析请求体
	request := &req.CreateStudentRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	// 调用服务
	resp, err := studentService.CreateStudent(ctx, request)
	if err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	// 序列化并返回响应
	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// getStudentHandler 获取学生信息处理器
func getStudentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 从查询参数获取用户ID
	userId := r.URL.Query().Get("user_id")

	request := &req.GetStudentRequest{
		UserId: userId,
	}

	// 调用服务
	resp, err := studentService.GetStudent(ctx, request)
	if err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	// 序列化并返回响应
	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// updateStudentHandler 更新学生信息处理器
func updateStudentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析请求体
	request := &req.UpdateStudentRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	// 调用服务
	resp, err := studentService.UpdateStudent(ctx, request)
	if err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	// 序列化并返回响应
	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// listStudentsHandler 查询学生列表处理器
func listStudentsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析请求体
	request := &req.ListStudentsRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	// 调用服务
	resp, err := studentService.ListStudents(ctx, request)
	if err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	// 序列化并返回响应
	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// updateLearningProgressHandler 更新学习进度处理器
func updateLearningProgressHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析请求体
	request := &req.UpdateLearningProgressRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusBadRequest, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	// 调用服务
	resp, err := studentService.UpdateLearningProgress(ctx, request)
	if err != nil {
		// 构建错误响应
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(int(resp.Code))
		w.Write(respBytes)
		return
	}

	// 序列化并返回响应
	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}
