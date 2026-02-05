package authen

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/utils"
)

// UserContextKey 用户信息在 context 中的 key
type UserContextKey string

const (
	UserIDKey   UserContextKey = "user_id"
	UserTypeKey UserContextKey = "user_type"
	RoleIDKey   UserContextKey = "role_id" // 学生ID/教师ID/管理员ID
)

// UserInfo 用户信息结构
type UserInfo struct {
	UserID   string // 用户ID
	UserType string // 用户类型：student/teacher/admin
	RoleID   string // 角色ID（学生ID/教师ID/管理员ID）
}

// Authen 统一身份认证中间件
func Authen(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置响应头
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 获取 token 从 Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// 判断是否是 API 请求
			if isAPIRequest(r) {
				respondError(w, http.StatusUnauthorized, "未授权：缺少认证令牌")
			} else {
				// Web 页面重定向到登录页
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			}
			return
		}

		// 检查 Bearer token 格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			if isAPIRequest(r) {
				respondError(w, http.StatusUnauthorized, "未授权：无效的认证令牌格式")
			} else {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			}
			return
		}

		tokenString := parts[1]

		// 创建 JWT 服务实例
		jwtService := utils.NewJWTService()

		// 验证 token 并获取用户 ID
		userID, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			if isAPIRequest(r) {
				respondError(w, http.StatusUnauthorized, fmt.Sprintf("未授权：%v", err))
			} else {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			}
			return
		}

		// 根据用户ID前缀判断用户类型
		var userType string
		if strings.HasPrefix(userID, "stu_") {
			userType = "student"
		} else if strings.HasPrefix(userID, "tea_") {
			userType = "teacher"
		} else if strings.HasPrefix(userID, "admin_") {
			userType = "admin"
		} else {
			// 如果无法识别用户类型，返回错误
			respondError(w, http.StatusUnauthorized, "未授权：无效的用户ID格式")
			return
		}

		// 根据用户类型查询对应的角色信息
		userInfo := &UserInfo{
			UserID:   userID,
			UserType: userType,
			RoleID:   userID, // 默认RoleID就是userID本身
		}

		// 根据用户类型验证用户是否存在并获取详细信息
		switch userType {
		case "student":
			studentDAO := dao.NewStudentDAO()
			student, err := studentDAO.GetStudentById(userID)
			if err != nil || student == nil {
				respondError(w, http.StatusUnauthorized, "未授权：学生不存在")
				return
			}
			userInfo.RoleID = student.StudentId
		case "teacher":
			teacherDAO := dao.NewTeacherDAO()
			teacher, err := teacherDAO.GetTeacherById(userID)
			if err != nil || teacher == nil {
				respondError(w, http.StatusUnauthorized, "未授权：教师不存在")
				return
			}
			userInfo.RoleID = teacher.TeacherId
		case "admin":
			adminDAO := dao.NewAdminUserDAO()
			admin, err := adminDAO.GetAdminUserByAdminId(userID)
			if err != nil || admin == nil {
				respondError(w, http.StatusUnauthorized, "未授权：管理员不存在")
				return
			}
			userInfo.RoleID = admin.AdminId
		}

		// 将用户信息添加到请求上下文中
		ctx := context.WithValue(r.Context(), UserIDKey, userInfo.UserID)
		ctx = context.WithValue(ctx, UserTypeKey, userInfo.UserType)
		ctx = context.WithValue(ctx, RoleIDKey, userInfo.RoleID)

		// 继续处理请求
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// isAPIRequest 判断是否是 API 请求
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

	return false
}

// respondError 返回错误响应
func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	resp := map[string]interface{}{
		"code":    statusCode,
		"message": message,
	}
	json.NewEncoder(w).Encode(resp)
}

// GetUserIDFromContext 从上下文中获取用户 ID
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetUserTypeFromContext 从上下文中获取用户类型
func GetUserTypeFromContext(ctx context.Context) (string, bool) {
	userType, ok := ctx.Value(UserTypeKey).(string)
	return userType, ok
}

// GetRoleIDFromContext 从上下文中获取角色 ID（学生ID/教师ID/管理员ID）
func GetRoleIDFromContext(ctx context.Context) (string, bool) {
	roleID, ok := ctx.Value(RoleIDKey).(string)
	return roleID, ok
}

// GetUserInfoFromContext 从上下文中获取完整的用户信息
func GetUserInfoFromContext(ctx context.Context) (*UserInfo, bool) {
	userID, ok1 := GetUserIDFromContext(ctx)
	userType, ok2 := GetUserTypeFromContext(ctx)
	roleID, ok3 := GetRoleIDFromContext(ctx)

	if !ok1 || !ok2 {
		return nil, false
	}

	return &UserInfo{
		UserID:   userID,
		UserType: userType,
		RoleID:   roleID, // 可能为空
	}, ok3
}
