package router

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/auth"
	"github.com/yzf120/elysia-backend/model/student/req"
	"github.com/yzf120/elysia-backend/service"
	"github.com/yzf120/elysia-backend/service_impl"
)

var (
	studentService     = service_impl.NewStudentServiceImpl()
	studentAuthService = service.NewStudentAuthService()
)

// registerStudent 学生相关路由
func registerStudent(publicRouter *mux.Router, protectedRouter *mux.Router) {
	// 学生端注册登录接口
	studentAuthRouter := publicRouter.PathPrefix("/api/student/auth").Subrouter()
	studentAuthRouter.HandleFunc("/send-code-register", studentSendCodeHandler).Methods("POST")
	studentAuthRouter.HandleFunc("/send-code-login", studentSendCodeHandler).Methods("POST")
	studentAuthRouter.HandleFunc("/register-sms", studentRegisterWithSMSHandler).Methods("POST")
	studentAuthRouter.HandleFunc("/login-sms", studentLoginWithSMSHandler).Methods("POST")
	studentAuthRouter.HandleFunc("/login-password", studentLoginWithPasswordHandler).Methods("POST")

	protectedRouter.HandleFunc("/api/student/create", createStudentHandler).Methods("POST")
	protectedRouter.HandleFunc("/api/student/get", getStudentHandler).Methods("GET")
	protectedRouter.HandleFunc("/api/student/update", updateStudentHandler).Methods("POST")
	protectedRouter.HandleFunc("/api/student/list", listStudentsHandler).Methods("POST")
	protectedRouter.HandleFunc("/api/student/update-progress", updateLearningProgressHandler).Methods("POST")
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

	// 从查询参数获取学生ID
	studentId := r.URL.Query().Get("student_id")

	request := &req.GetStudentRequest{
		StudentId: studentId,
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

// ==================== 学生认证处理器函数 ====================

// studentSendCodeHandler 学生端发送验证码
func studentSendCodeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	req := &auth.SendCodeRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 从URL路径判断验证码类型
	codeType := consts.Register
	if strings.Contains(r.URL.Path, consts.Login) {
		codeType = consts.Login
	}

	if err := smsService.SendVerificationCode(ctx, req.PhoneNumber, consts.RoleStudent, codeType); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message": "验证码发送成功",
	})
}

// studentRegisterWithSMSHandler 学生注册
func studentRegisterWithSMSHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	req := &auth.RegisterWithSMSRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	student, err := studentAuthService.RegisterWithSMS(ctx, req.PhoneNumber, req.Code, req.StudentNumber, req.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message":   "注册成功",
		"user_info": student,
	})
}

// studentLoginWithSMSHandler 学生验证码登录
func studentLoginWithSMSHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	req := &auth.LoginWithSMSRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	student, token, err := studentAuthService.LoginWithSMS(ctx, req.PhoneNumber, req.Code)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message":   "登录成功",
		"user_info": student,
		"token":     token,
	})
}

// studentLoginWithPasswordHandler 学生密码登录
func studentLoginWithPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	req := &auth.LoginWithPasswordRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	student, token, err := studentAuthService.LoginWithPassword(ctx, req.StudentNumber, req.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"message":   "登录成功",
		"user_info": student,
		"token":     token,
	})
}
