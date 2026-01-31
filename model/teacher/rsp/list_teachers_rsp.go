package rsp

// ListTeachersResponse 查询教师列表响应
type ListTeachersResponse struct {
	Code     int32          `json:"code"`      // 响应码 0-成功 其他-失败
	Message  string         `json:"message"`   // 响应消息
	Teachers []*TeacherInfo `json:"teachers"`  // 教师列表
	Total    int32          `json:"total"`     // 总数
	Page     int32          `json:"page"`      // 当前页码
	PageSize int32          `json:"page_size"` // 每页数量
}