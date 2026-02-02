package auth

// LoginWithUsernameRequest 用户名+密码登录请求（管理员专用）
type LoginWithUsernameRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginWithUsernameResponse 用户名+密码登录响应
type LoginWithUsernameResponse struct {
	Code     int32       `json:"code"`
	Message  string      `json:"message"`
	UserInfo interface{} `json:"user_info"`
	Token    string      `json:"token"`
}
