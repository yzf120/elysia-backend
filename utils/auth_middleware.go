package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// AuthMiddleware 认证中间件，只对需要认证的路由生效
func AuthMiddleware(jwtService *JWTService, publicRoutes []string) func(http.Handler) http.Handler {
	// 将白名单路由转换为 map 以便快速查找
	publicRouteMap := make(map[string]bool)
	for _, route := range publicRoutes {
		publicRouteMap[route] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 检查当前路由是否在白名单中（不需要认证）
			if publicRouteMap[r.URL.Path] {
				// 白名单路由，直接放行
				next.ServeHTTP(w, r)
				return
			}

			// 获取 token 从 Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// 对于需要认证的 API 返回 401，对于 Web 页面重定向到登录页
				if isAPIRequest(r) {
					http.Error(w, "Unauthorized: Missing authorization token", http.StatusUnauthorized)
				} else {
					http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				}
				return
			}

			// 检查 Bearer token 格式
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				if isAPIRequest(r) {
					http.Error(w, "Unauthorized: Invalid authorization header format", http.StatusUnauthorized)
				} else {
					http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				}
				return
			}

			tokenString := parts[1]

			// 验证 token
			userID, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				if isAPIRequest(r) {
					http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
				} else {
					http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				}
				return
			}

			// 将用户 ID 添加到请求上下文中
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// isAPIRequest 判断是否是 API 请求（根据 Content-Type 或路径判断）
func isAPIRequest(r *http.Request) bool {
	// 如果请求头包含 application/json，则认为是 API 请求
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "application/x-www-form-urlencoded") {
		return true
	}

	// 如果路径以 /api/ 开头，也认为是 API 请求
	if strings.HasPrefix(r.URL.Path, "/api/") {
		return true
	}

	// 其他情况视为 Web 请求
	return false
}

// GetUserIDFromContext 从上下文中获取用户 ID
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userID").(string)
	return userID, ok
}

// RequireAuth 检查用户是否已认证的辅助函数（直接使用，不通过中间件）
func RequireAuth(jwtService *JWTService, w http.ResponseWriter, r *http.Request) (string, bool) {
	// 获取 token 从 Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized: Missing authorization token", http.StatusUnauthorized)
		return "", false
	}

	// 检查 Bearer token 格式
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Unauthorized: Invalid authorization header format", http.StatusUnauthorized)
		return "", false
	}

	tokenString := parts[1]

	// 验证 token
	userID, err := jwtService.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
		return "", false
	}

	return userID, true
}

// AddLogoutHandler 添加登出路由处理
func AddLogoutHandler(router *mux.Router, jwtService *JWTService) {
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
