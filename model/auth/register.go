package auth

// RegisterWithSMSRequest 手机号+验证码注册请求
type RegisterWithSMSRequest struct {
	PhoneNumber   string `json:"phone_number"`
	Code          string `json:"code"`
	StudentNumber string `json:"student_number"`
	Password      string `json:"password"`
}

// RegisterWithSMSResponse 手机号+验证码注册响应
type RegisterWithSMSResponse struct {
	Code     int32       `json:"code"`
	Message  string      `json:"message"`
	UserInfo interface{} `json:"user_info"`
	Token    string      `json:"token"`
}
