package rsp

// ListStudentsResponse 查询学生列表响应
type ListStudentsResponse struct {
	Code     int32          `json:"code"`      // 响应码 0-成功 其他-失败
	Message  string         `json:"message"`   // 响应消息
	Students []*StudentInfo `json:"students"`  // 学生列表
	Total    int32          `json:"total"`     // 总数
	Page     int32          `json:"page"`      // 当前页码
	PageSize int32          `json:"page_size"` // 每页数量
}