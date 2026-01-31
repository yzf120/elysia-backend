package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/errs"
	adminPb "github.com/yzf120/elysia-backend/proto/admin"
	"github.com/yzf120/elysia-backend/service_impl"
)

var adminUserService = service_impl.NewAdminUserServiceImpl()

// registerAdmin 注册管理员相关路由
func registerAdmin(publicRouter *mux.Router, protectedRouter *mux.Router) {
	// 管理员登录接口（公开路由）
	publicRouter.HandleFunc("/api/admin/login", loginAdminUserHandler).Methods("POST")

	// 以下接口需要认证（受保护路由）
	protectedRouter.HandleFunc("/api/admin/create", createAdminUserHandler).Methods("POST")
	protectedRouter.HandleFunc("/api/admin/get", getAdminUserHandler).Methods("GET")
	protectedRouter.HandleFunc("/api/admin/list", listAdminUsersHandler).Methods("POST")
	protectedRouter.HandleFunc("/api/admin/update-password", updateAdminUserPasswordHandler).Methods("POST")
	protectedRouter.HandleFunc("/api/admin/update-status", updateAdminUserStatusHandler).Methods("POST")
}

// createAdminUserHandler 创建管理员用户处理器
func createAdminUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析请求体
	request := &adminPb.CreateAdminUserRequest{}
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
	resp, err := adminUserService.CreateAdminUser(ctx, request)
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

// loginAdminUserHandler 管理员登录处理器
func loginAdminUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析请求体
	request := &adminPb.LoginAdminUserRequest{}
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

	// 获取客户端IP地址
	if request.IpAddress == "" {
		request.IpAddress = r.RemoteAddr
	}

	// 调用服务
	resp, err := adminUserService.LoginAdminUser(ctx, request)
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

// getAdminUserHandler 获取管理员用户信息处理器
func getAdminUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 从查询参数获取管理员ID
	adminId := r.URL.Query().Get("admin_id")

	request := &adminPb.GetAdminUserRequest{
		AdminId: adminId,
	}

	// 调用服务
	resp, err := adminUserService.GetAdminUser(ctx, request)
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

// listAdminUsersHandler 查询管理员用户列表处理器
func listAdminUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析请求体
	request := &adminPb.ListAdminUsersRequest{}
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
	resp, err := adminUserService.ListAdminUsers(ctx, request)
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

// updateAdminUserPasswordHandler 更新管理员用户密码处理器
func updateAdminUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析请求体
	request := &adminPb.UpdateAdminUserPasswordRequest{}
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
	resp, err := adminUserService.UpdateAdminUserPassword(ctx, request)
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

// updateAdminUserStatusHandler 更新管理员用户状态处理器
func updateAdminUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 解析请求体
	request := &adminPb.UpdateAdminUserStatusRequest{}
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
	resp, err := adminUserService.UpdateAdminUserStatus(ctx, request)
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
