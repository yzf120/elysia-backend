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
	adminUserService *service_impl.AdminUserServiceImpl
	adminAuthService *service.AdminAuthService
)

// registerAdmin 注册管理员相关路由
func registerAdmin(publicRouter *mux.Router, protectedRouter *mux.Router) {
	// 管理员端注册登录接口
	adminAuthRouter := publicRouter.PathPrefix("/admin/auth").Subrouter()
	adminAuthRouter.HandleFunc("/send-code-register", adminSendCodeHandler).Methods("POST")
	adminAuthRouter.HandleFunc("/send-code-login", adminSendCodeHandler).Methods("POST")
	adminAuthRouter.HandleFunc("/register-sms", adminRegisterWithSMSHandler).Methods("POST")
	adminAuthRouter.HandleFunc("/login-sms", adminLoginWithSMSHandler).Methods("POST")
	adminAuthRouter.HandleFunc("/login-password", adminLoginWithPasswordHandler).Methods("POST")

	// 以下接口需要认证（受保护路由）
	protectedRouter.HandleFunc("/admin/create", createAdminUserHandler).Methods("POST")
	protectedRouter.HandleFunc("/admin/get", getAdminUserHandler).Methods("GET")
	protectedRouter.HandleFunc("/admin/list", listAdminUsersHandler).Methods("POST")
	protectedRouter.HandleFunc("/admin/update-password", updateAdminUserPasswordHandler).Methods("POST")
	protectedRouter.HandleFunc("/admin/update-status", updateAdminUserStatusHandler).Methods("POST")

	// 管理员个人信息相关接口（需要认证）
	protectedRouter.HandleFunc("/admin/profile", getAdminProfileHandler).Methods("GET")
	protectedRouter.HandleFunc("/admin/profile/password", updateAdminPasswordHandler).Methods("POST")
	protectedRouter.HandleFunc("/admin/profile/email", updateAdminEmailHandler).Methods("POST")
	protectedRouter.HandleFunc("/admin/auth/logout", adminLogoutHandler).Methods("POST")
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

// ==================== 管理员个人信息处理器函数 ====================

// getAdminProfileHandler 获取当前管理员个人信息
func getAdminProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	// 从上下文获取管理员ID
	adminId, ok := ctx.Value("admin_id").(string)
	if !ok || adminId == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "未授权")
		return
	}

	request := &adminPb.GetAdminUserRequest{
		AdminId: adminId,
	}

	// 调用服务
	resp, err := adminUserService.GetAdminUser(ctx, request)
	if err != nil || resp.Code != consts.SuccessCode {
		msg := "查询失败"
		if resp != nil {
			msg = resp.Message
		}
		writeErrorResponse(w, http.StatusInternalServerError, msg)
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message":   "查询成功",
		"user_info": resp.AdminUser,
	})
}

// updateAdminPasswordHandler 更新当前管理员密码
func updateAdminPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	// 从上下文获取管理员ID
	adminId, ok := ctx.Value("admin_id").(string)
	if !ok || adminId == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "未授权")
		return
	}

	// 解析请求体
	var reqBody struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	request := &adminPb.UpdateAdminUserPasswordRequest{
		AdminId:     adminId,
		OldPassword: reqBody.OldPassword,
		NewPassword: reqBody.NewPassword,
	}

	// 调用服务
	resp, err := adminUserService.UpdateAdminUserPassword(ctx, request)
	if err != nil || resp.Code != consts.SuccessCode {
		msg := "修改失败"
		if resp != nil {
			msg = resp.Message
		}
		writeErrorResponse(w, http.StatusBadRequest, msg)
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "密码修改成功，请重新登录",
	})
}

// updateAdminEmailHandler 更新当前管理员邮箱
func updateAdminEmailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	// 从上下文获取管理员ID
	adminId, ok := ctx.Value("admin_id").(string)
	if !ok || adminId == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "未授权")
		return
	}

	// 解析请求体
	var reqBody struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// TODO: 验证邮箱验证码
	// 此处简化处理，直接调用更新邮箱服务
	if err := adminAuthService.UpdateAdminEmail(ctx, adminId, reqBody.Email); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "邮箱绑定成功",
	})
}

// adminLogoutHandler 管理员退出登录
func adminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	// 当前简单实现，客户端自行清除token
	writeSuccessResponse(w, map[string]interface{}{
		"message": "退出成功",
	})
}
