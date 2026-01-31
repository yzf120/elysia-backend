package rsp

// LoginTeacherResponse 教师登录响应
type LoginTeacherResponse struct {
	Code    int32        `json:"code"`    // 响应码 0-成功 其他-失败
	Message string       `json:"message"` // 响应消息
	Teacher *TeacherInfo `json:"teacher"` // 教师信息
	User    *UserInfo    `json:"user"`    // 用户信息
	Token   string       `json:"token"`   // 登录令牌
}