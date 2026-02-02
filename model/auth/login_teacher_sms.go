package auth

// TeacherLoginWithSMSRequest 教师手机号+验证码登录请求
type TeacherLoginWithSMSRequest struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}

// TeacherLoginWithSMSResponse 教师手机号+验证码登录响应
type TeacherLoginWithSMSResponse struct {
	Code     int32       `json:"code"`
	Message  string      `json:"message"`
	UserInfo interface{} `json:"user_info"`
	Token    string      `json:"token"`
}
