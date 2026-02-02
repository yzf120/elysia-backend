package auth

// LoginWithSMSRequest 手机号+验证码登录请求
type LoginWithSMSRequest struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}

// LoginWithSMSResponse 手机号+验证码登录响应
type LoginWithSMSResponse struct {
	Code     int32       `json:"code"`
	Message  string      `json:"message"`
	UserInfo interface{} `json:"user_info"`
	Token    string      `json:"token"`
}
