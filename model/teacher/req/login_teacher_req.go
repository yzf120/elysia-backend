package req

// LoginTeacherRequest 教师登录请求
type LoginTeacherRequest struct {
	PhoneNumber string `json:"phone_number"` // 手机号（必填）
	Password    string `json:"password"`     // 密码（必填）
}