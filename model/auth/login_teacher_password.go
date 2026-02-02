package auth

// TeacherLoginWithPasswordRequest 教师工号+密码登录请求
type TeacherLoginWithPasswordRequest struct {
	EmployeeNumber string `json:"employee_number"`
	Password       string `json:"password"`
}

// TeacherLoginWithPasswordResponse 教师工号+密码登录响应
type TeacherLoginWithPasswordResponse struct {
	Code     int32       `json:"code"`
	Message  string      `json:"message"`
	UserInfo interface{} `json:"user_info"`
	Token    string      `json:"token"`
}
