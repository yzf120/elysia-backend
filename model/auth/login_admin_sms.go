package auth

// AdminLoginWithSMSRequest 管理员手机号+验证码登录请求
type AdminLoginWithSMSRequest struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}

// AdminLoginWithSMSResponse 管理员手机号+验证码登录响应
type AdminLoginWithSMSResponse struct {
	Code     int32       `json:"code"`
	Message  string      `json:"message"`
	UserInfo interface{} `json:"user_info"`
	Token    string      `json:"token"`
}
