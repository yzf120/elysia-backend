package errs

// Error 通用错误信息
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

// BaseResponse ...
type BaseResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error Error       `json:"error,omitempty"`
}

// BaseSuccessError 缺省内部错误控制
var BaseSuccessError = Error{
	Code:    0,
	Message: "success",
}

// BaseInnerServerError 缺省内部错误控制
var BaseInnerServerError = Error{
	Code:    500,
	Message: "服务开小差拉，可以稍后再试下",
}

// BaseBadRequestError ...
var BaseBadRequestError = Error{
	Code:    400,
	Message: "请求参数有误",
}

func NewError(code int, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}
