package auth

// VerifyCodeRequest 验证验证码请求
type VerifyCodeRequest struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
	CodeType    string `json:"code_type"` // register 或 login
}

// VerifyCodeResponse 验证验证码响应
type VerifyCodeResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
