package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/authen"
	"github.com/yzf120/elysia-backend/service"
	"github.com/yzf120/elysia-backend/service_impl"
	"github.com/yzf120/elysia-backend/utils"
)

func Init() {
	// 初始化所有Service实例
	// 注意：此函数必须在 dao.InitDB() 之后调用
	adminUserService = service_impl.NewAdminUserServiceImpl()
	adminAuthService = service.NewAdminAuthService()
	studentService = service_impl.NewStudentServiceImpl()
	studentAuthService = service.NewStudentAuthService()
	teacherService = service_impl.NewTeacherServiceImpl()
	teacherAuthService = service.NewTeacherAuthService()
	smsService = service.NewSMSService()
	problemService = service_impl.NewProblemServiceImpl()
}

func RegisterRouter(router *mux.Router) {
	// 创建子路由器：公开路由（无需认证）
	publicRouter := router.PathPrefix("/api").Subrouter()

	// 创建子路由器：受保护路由（需要认证）
	protectedRouter := router.PathPrefix("/api").Subrouter()

	// 为受保护路由添加统一的身份认证中间件
	protectedRouter.Use(authen.Authen)

	// 注册路由
	registerApiRouters(publicRouter, protectedRouter)
}

func registerApiRouters(publicRouter *mux.Router, protectedRouter *mux.Router) {
	// 管理员相关接口（包含认证接口）
	registerAdmin(publicRouter, protectedRouter)
	// 教师相关接口（包含认证接口）
	registerTeacher(publicRouter, protectedRouter)
	// 学生相关接口（包含认证接口）
	registerStudent(publicRouter, protectedRouter)

	// 登出接口（虽然是认证相关，但需要认证后才能登出，所以注册到受保护路由）
	registerLogout(protectedRouter)

	// 对话相关
	registerConversation(protectedRouter)
	// 会话相关
	registerSession(protectedRouter)

	// 智能体相关
	registerAgent(protectedRouter)

	// 教师审批单相关接口（需要认证）
	RegisterTeacherApprovalRoutes(protectedRouter)

	// 题目相关接口（增删改仅教师，查询学生和教师均可）
	registerProblem(publicRouter, protectedRouter)

}

// registerLogout 注册登出接口（需要认证）
func registerLogout(router *mux.Router) {
	// 学生登出接口
	router.HandleFunc("/student/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 获取 token 从 Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing authorization token", http.StatusUnauthorized)
			return
		}

		// 检查 Bearer token 格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized: Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// 创建 JWT 服务实例
		jwtService := utils.NewJWTService()

		// 先验证 token 以获取用户 ID
		userID, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
			return
		}

		// 使 token 失效（登出）
		if err := jwtService.InvalidateToken(userID, tokenString); err != nil {
			http.Error(w, fmt.Sprintf("Logout failed: %v", err), http.StatusInternalServerError)
			return
		}

		// 返回成功响应
		resp := map[string]interface{}{
			"success": true,
			"message": "成功登出",
		}
		json.NewEncoder(w).Encode(resp)
	}).Methods("POST")

	// 教师登出接口
	router.HandleFunc("/teacher/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 获取 token 从 Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing authorization token", http.StatusUnauthorized)
			return
		}

		// 检查 Bearer token 格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized: Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// 创建 JWT 服务实例
		jwtService := utils.NewJWTService()

		// 先验证 token 以获取用户 ID
		userID, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
			return
		}

		// 使 token 失效（登出）
		if err := jwtService.InvalidateToken(userID, tokenString); err != nil {
			http.Error(w, fmt.Sprintf("Logout failed: %v", err), http.StatusInternalServerError)
			return
		}

		// 返回成功响应
		resp := map[string]interface{}{
			"success": true,
			"message": "成功登出",
		}
		json.NewEncoder(w).Encode(resp)
	}).Methods("POST")

	// 管理员登出接口
	router.HandleFunc("/admin/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 获取 token 从 Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing authorization token", http.StatusUnauthorized)
			return
		}

		// 检查 Bearer token 格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized: Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// 创建 JWT 服务实例
		jwtService := utils.NewJWTService()

		// 先验证 token 以获取用户 ID
		userID, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
			return
		}

		// 使 token 失效（登出）
		if err := jwtService.InvalidateToken(userID, tokenString); err != nil {
			http.Error(w, fmt.Sprintf("Logout failed: %v", err), http.StatusInternalServerError)
			return
		}

		// 返回成功响应
		resp := map[string]interface{}{
			"success": true,
			"message": "成功登出",
		}
		json.NewEncoder(w).Encode(resp)
	}).Methods("POST")
}
