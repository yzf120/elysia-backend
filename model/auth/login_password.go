package auth

// LoginWithPasswordRequest 学号+密码登录请求
type LoginWithPasswordRequest struct {
	StudentNumber string `json:"student_number"`
	Password      string `json:"password"`
}

// LoginWithPasswordResponse 手机号+密码登录响应
type LoginWithPasswordResponse struct {
	Code     int32       `json:"code"`
	Message  string      `json:"message"`
	UserInfo interface{} `json:"user_info"`
	Token    string      `json:"token"`
}
