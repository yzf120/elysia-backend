package rsp

// UserInfo 用户基本信息
type UserInfo struct {
	UserId      string `json:"user_id"`      // 用户ID
	UserName    string `json:"user_name"`    // 用户名
	Email       string `json:"email"`        // 邮箱
	PhoneNumber string `json:"phone_number"` // 手机号
	ChineseName string `json:"chinese_name"` // 中文姓名
	UserType    string `json:"user_type"`    // 用户类型
	Status      int32  `json:"status"`       // 状态
}