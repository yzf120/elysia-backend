package auth

// AdminLoginWithPasswordRequest 管理员手机号+密码登录请求
type AdminLoginWithPasswordRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

// AdminLoginWithPasswordResponse 管理员手机号+密码登录响应
type AdminLoginWithPasswordResponse struct {
	Code     int32       `json:"code"`
	Message  string      `json:"message"`
	UserInfo interface{} `json:"user_info"`
	Token    string      `json:"token"`
}
