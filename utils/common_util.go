package utils

import (
	"crypto/rand"
	"fmt"
	"time"
)

// GenerateAdminId 生成管理员ID
func GenerateAdminId() string {
	// 使用时间戳+随机数生成唯一的管理员ID
	timestamp := time.Now().UnixNano() / 1000000 // 毫秒级时间戳
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)

	// 格式：admin_时间戳_随机数
	return fmt.Sprintf("admin_%d_%x", timestamp, randomBytes)
}

// GenerateRandomPassword 生成随机密码
func GenerateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}

// ValidateAdminUsername 验证管理员用户名格式
func ValidateAdminUsername(username string) error {
	if len(username) < 3 || len(username) > 128 {
		return fmt.Errorf("用户名长度必须在3-128个字符之间")
	}

	// 检查是否包含非法字符
	for _, char := range username {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-' || char == '.') {
			return fmt.Errorf("用户名只能包含字母、数字、下划线、连字符和点号")
		}
	}

	return nil
}

// ValidateAdminPassword 验证管理员密码强度
func ValidateAdminPassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("密码长度不能少于6位")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char == '!' || char == '@' || char == '#' || char == '$' ||
			char == '%' || char == '^' || char == '&' || char == '*' ||
			char == '(' || char == ')' || char == '-' || char == '_' ||
			char == '+' || char == '=' || char == '[' || char == ']' ||
			char == '{' || char == '}' || char == '|' || char == ';' ||
			char == ':' || char == ',' || char == '.' || char == '<' ||
			char == '>' || char == '/' || char == '?' || char == '~':
			hasSpecial = true
		}
	}

	// 建议密码包含大小写字母、数字和特殊字符
	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return fmt.Errorf("密码应包含大小写字母、数字和特殊字符")
	}

	return nil
}

// MaskPassword 密码掩码显示
func MaskPassword(password string) string {
	if len(password) == 0 {
		return ""
	}
	return "******"
}

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	// 简单的邮箱格式验证
	if len(email) < 5 || len(email) > 128 {
		return false
	}

	var hasAt, hasDot bool
	for i, char := range email {
		if char == '@' {
			hasAt = true
			if i == 0 || i == len(email)-1 {
				return false
			}
		}
		if char == '.' && hasAt {
			hasDot = true
		}
	}

	return hasAt && hasDot
}

// GenerateDefaultAdminPassword 生成默认管理员密码
func GenerateDefaultAdminPassword() string {
	return "Admin@123" // 默认密码，建议首次登录后修改
}
