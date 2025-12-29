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

// VerifyCodeRequest 验证验证码请求
type VerifyCodeRequest struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
	CodeType    string `json:"code_type"` // register 或 login
}

type VerifyCodeResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// RegisterWithSMSRequest 手机号+验证码注册请求
type RegisterWithSMSRequest struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}

type RegisterWithSMSResponse struct {
	Code     int32       `json:"code"`
	Message  string      `json:"message"`
	UserInfo interface{} `json:"user_info"`
	Token    string      `json:"token"`
}

// LoginWithSMSRequest 手机号+验证码登录请求
type LoginWithSMSRequest struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}

type LoginWithSMSResponse struct {
	Code     int32       `json:"code"`
	Message  string      `json:"message"`
	UserInfo interface{} `json:"user_info"`
	Token    string      `json:"token"`
}

// LoginWithPasswordRequest 手机号+密码登录请求
type LoginWithPasswordRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginWithPasswordResponse struct {
	Code     int32       `json:"code"`
	Message  string      `json:"message"`
	UserInfo interface{} `json:"user_info"`
	Token    string      `json:"token"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	Code     int32       `json:"code"`
	Message  string      `json:"message"`
	UserInfo interface{} `json:"user_info"`
	Token    string      `json:"token"`
}
