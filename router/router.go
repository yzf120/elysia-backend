package router

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/utils"
	"net/http"
	"strings"
)

func Init() {}

func RegisterRouter(router *mux.Router) {
	// 创建 JWT 服务实例
	jwtService := utils.NewJWTService()

	// 定义不需要认证的公开路由（如注册、登录等）
	publicRoutes := []string{
		"/api/auth/send-code",
		"/api/auth/verify-code",
		"/api/auth/register-sms",
		"/api/auth/login-sms",
		"/api/auth/login-password",
		"/api/auth/logout",
		"/api/user/create",      // 用户注册
		"/api/teacher/register", // 教师注册
		"/api/teacher/login",    // 教师登录
		"/api/admin/login",      // 管理员登录
	}

	// 创建子路由器：公开路由（无需认证）
	publicRouter := router.PathPrefix("/api").Subrouter()

	// 创建子路由器：受保护路由（需要认证）
	protectedRouter := router.PathPrefix("/api").Subrouter()

	// 为受保护路由添加认证中间件
	protectedRouter.Use(utils.AuthMiddleware(jwtService, publicRoutes))

	// 注册路由
	registerApiRouters(publicRouter, protectedRouter)
}

func registerApiRouters(publicRouter *mux.Router, protectedRouter *mux.Router) {
	// 注册认证相关接口（公开接口）
	registerAuth(publicRouter)
	// 添加登出接口（虽然是认证相关，但需要认证后才能登出，所以注册到受保护路由）
	registerLogout(protectedRouter)

	// 以下接口需要认证（注册到受保护路由）
	// 注册会话相关
	registerConversation(protectedRouter)
	// 注册会话相关
	registerSession(protectedRouter)
	// 注册智能体相关
	registerAgent(protectedRouter)
	// 注册用户相关接口（除注册外都需要认证）
	registerUserRoutes(publicRouter, protectedRouter)
	// 注册学生相关接口（需要认证）
	registerStudent(protectedRouter)
	// 注册教师相关接口（注册和登录为公开接口，其他需要认证）
	registerTeacher(publicRouter, protectedRouter)
	// 注册管理员相关接口（登录为公开接口，其他需要认证）
	registerAdmin(publicRouter, protectedRouter)
}

// registerUserRoutes 注册用户相关路由，将注册接口放在公开路由，其他放在受保护路由
func registerUserRoutes(publicRouter *mux.Router, protectedRouter *mux.Router) {
	// 用户创建接口（注册） - 公开路由
	publicRouter.HandleFunc("/api/user/create", createUserHandler).Methods("POST")

	// 其他用户接口需要认证 - 受保护路由
	protectedRouter.HandleFunc("/api/user/get", getUserHandler).Methods("GET")
	protectedRouter.HandleFunc("/api/user/update", updateUserHandler).Methods("POST")
	protectedRouter.HandleFunc("/api/user/delete", deleteUserHandler).Methods("POST")
	protectedRouter.HandleFunc("/api/user/list", listUsersHandler).Methods("POST")
}

// registerLogout 注册登出接口（需要认证）
func registerLogout(router *mux.Router) {
	http.HandleFunc("/api/auth/logout", func(w http.ResponseWriter, r *http.Request) {
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
	})
}
