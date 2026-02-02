package authen

import (
	"context"
	"net/http"

	"github.com/yzf120/elysia-backend/consts"
)

// AdminAuthMiddleware 管理员认证中间件
func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从 context 获取用户信息（由 Authen 中间件设置）
		userInfo, ok := GetUserInfoFromContext(r.Context())
		if !ok || userInfo.UserID == "" {
			http.Error(w, "未授权：需要登录", http.StatusUnauthorized)
			return
		}

		// 检查用户类型是否为管理员
		if userInfo.UserType != consts.RoleAdmin {
			http.Error(w, "未授权：需要管理员权限", http.StatusForbidden)
			return
		}

		// 检查是否有管理员 ID（RoleID）
		if userInfo.RoleID == "" {
			http.Error(w, "未授权：管理员信息不存在", http.StatusUnauthorized)
			return
		}

		// 将管理员 ID 添加到请求上下文中（保持向后兼容）
		ctx := context.WithValue(r.Context(), "admin_id", userInfo.RoleID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetAdminIDFromContext 从上下文中获取管理员 ID
func GetAdminIDFromContext(ctx context.Context) (string, bool) {
	// 优先从新的方式获取
	if roleID, ok := GetRoleIDFromContext(ctx); ok {
		return roleID, ok
	}
	// 兼容旧的方式
	adminID, ok := ctx.Value("admin_id").(string)
	return adminID, ok
}
