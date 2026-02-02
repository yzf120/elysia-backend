package router

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/auth"
	adminPb "github.com/yzf120/elysia-backend/proto/admin"
	"github.com/yzf120/elysia-backend/service"
	"github.com/yzf120/elysia-backend/service_impl"
)

var (
	adminUserService = service_impl.NewAdminUserServiceImpl()
	adminAuthService = service.NewAdminAuthService()
)

// registerAdmin 注册管理员相关路由
func registerAdmin(publicRouter *mux.Router, protectedRouter *mux.Router) {
	// 管理员端注册登录接口
	adminAuthRouter := publicRouter.PathPrefix("/api/admin/auth").Subrouter()
	adminAuthRouter.HandleFunc("/send-code-register", adminSendCodeHandler).Methods("POST")
	adminAuthRouter.HandleFunc("/send-code-login", adminSendCodeHandler).Methods("POST")
	adminAuthRouter.HandleFunc("/register-sms", adminRegisterWithSMSHandler).Methods("POST")
	adminAuthRouter.HandleFunc("/login-sms", adminLoginWithSMSHandler).Methods("POST")
	adminAuthRouter.HandleFunc("/login-password", adminLoginWithPasswordHandler).Methods("POST")

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
			Error: errs.NewError(http.StatusInternalServerError, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBytes)
		return
	}

	// 序列化并返回响应
	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

// ==================== 管理员认证处理器函数 ====================

// adminSendCodeHandler 管理员端发送验证码
func adminSendCodeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	req := &auth.SendCodeRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	codeType := consts.Register
	if strings.Contains(r.URL.Path, consts.Login) {
		codeType = consts.Login
	}

	if err := smsService.SendVerificationCode(ctx, req.PhoneNumber, consts.RoleAdmin, codeType); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "验证码发送成功",
	})
}

// adminRegisterWithSMSHandler 管理员注册
func adminRegisterWithSMSHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	req := &auth.AdminRegisterWithSMSRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	admin, token, err := adminAuthService.RegisterWithSMS(ctx, req.PhoneNumber, req.Code)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message":   "注册成功，默认密码为 Admin@123，请及时修改",
		"user_info": admin,
		"token":     token,
	})
}

// adminLoginWithSMSHandler 管理员验证码登录
func adminLoginWithSMSHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	req := &auth.AdminLoginWithSMSRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	admin, token, err := adminAuthService.LoginWithSMS(ctx, req.PhoneNumber, req.Code)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message":   "登录成功",
		"user_info": admin,
		"token":     token,
	})
}

// adminLoginWithPasswordHandler 管理员手机号密码登录
func adminLoginWithPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	req := &auth.AdminLoginWithPasswordRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	admin, token, err := adminAuthService.LoginWithPassword(ctx, req.PhoneNumber, req.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message":   "登录成功",
		"user_info": admin,
		"token":     token,
	})
}
