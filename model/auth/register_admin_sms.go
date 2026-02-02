package auth

// AdminRegisterWithSMSRequest 管理员手机号+验证码注册请求
type AdminRegisterWithSMSRequest struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}

// AdminRegisterWithSMSResponse 管理员手机号+验证码注册响应
type AdminRegisterWithSMSResponse struct {
	Code     int32       `json:"code"`
	Message  string      `json:"message"`
	UserInfo interface{} `json:"user_info"`
	Token    string      `json:"token"`
}
