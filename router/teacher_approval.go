package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/authen"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/teacher/req"
	"github.com/yzf120/elysia-backend/service_impl"
)

var teacherApprovalService = service_impl.NewTeacherApprovalServiceImpl()

// RegisterTeacherApprovalRoutes 注册教师审批单路由
func RegisterTeacherApprovalRoutes(protectedRouter *mux.Router) {
	// 创建管理员专用子路由，应用管理员认证中间件
	adminRouter := protectedRouter.PathPrefix("/admin").Subrouter()
	adminRouter.Use(authen.AdminAuthMiddleware)

	// 获取审批单详情（根据审批单ID）- 需要认证
	protectedRouter.HandleFunc("/teacher/approvals/{approval_id}", getApprovalByIdHandler).Methods("GET")

	// 获取教师的审批单（根据教师ID）- 需要认证
	protectedRouter.HandleFunc("/teacher/{teacher_id}/approval", getApprovalByTeacherIdHandler).Methods("GET")

	// 查询审批单列表（管理员）- 使用管理员路由
	adminRouter.HandleFunc("/teacher/approvals", listApprovalsHandler).Methods("GET")

	// 审批教师（管理员）- 使用管理员路由
	adminRouter.HandleFunc("/teacher/approvals/{approval_id}/approve", approveTeacherHandler).Methods("POST")

	// 删除审批单（管理员）- 使用管理员路由
	adminRouter.HandleFunc("/teacher/approvals/{approval_id}", deleteApprovalHandler).Methods("DELETE")
}

// getApprovalByIdHandler 获取审批单详情处理器
func getApprovalByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 获取路径参数
	vars := mux.Vars(r)
	approvalId := vars["approval_id"]

	// 调用服务
	resp, err := teacherApprovalService.GetApprovalById(ctx, approvalId)
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

// getApprovalByTeacherIdHandler 获取教师的审批单处理器
func getApprovalByTeacherIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 获取路径参数
	vars := mux.Vars(r)
	teacherId := vars["teacher_id"]

	// 调用服务
	resp, err := teacherApprovalService.GetApprovalByTeacherId(ctx, teacherId)
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

// listApprovalsHandler 查询审批单列表处理器
func listApprovalsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 从上下文获取管理员ID（由 AdminAuthMiddleware 设置）
	adminId, ok := authen.GetAdminIDFromContext(ctx)
	if !ok || adminId == "" {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusUnauthorized, "未授权：需要管理员权限"),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(respBytes)
		return
	}

	// 解析请求体
	request := &req.ListApprovalsRequest{}
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
	resp, err := teacherApprovalService.ListApprovals(ctx, request)
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

// approveTeacherHandler 审批教师处理器
func approveTeacherHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 从上下文获取管理员ID（由 AdminAuthMiddleware 设置）
	adminId, ok := authen.GetAdminIDFromContext(ctx)
	if !ok || adminId == "" {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusUnauthorized, "未授权：需要管理员权限"),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(respBytes)
		return
	}

	// 获取路径参数
	vars := mux.Vars(r)
	approvalId := vars["approval_id"]

	// 解析请求体
	request := &req.ApproveTeacherRequest{}
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
	resp, err := teacherApprovalService.ApproveTeacher(ctx, approvalId, adminId, request)
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

// deleteApprovalHandler 删除审批单处理器
func deleteApprovalHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 从上下文获取管理员ID（由 AdminAuthMiddleware 设置）
	adminId, ok := authen.GetAdminIDFromContext(ctx)
	if !ok || adminId == "" {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusUnauthorized, "未授权：需要管理员权限"),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(respBytes)
		return
	}

	// 获取路径参数
	vars := mux.Vars(r)
	approvalId := vars["approval_id"]

	// 调用服务
	resp, err := teacherApprovalService.DeleteApproval(ctx, approvalId, adminId)
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
