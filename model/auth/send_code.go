package auth

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	PhoneNumber string `json:"phone_number"`
	CodeType    string `json:"code_type"` // register 或 login
}

// SendCodeResponse 发送验证码响应
type SendCodeResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
