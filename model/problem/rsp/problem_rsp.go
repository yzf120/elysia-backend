package rsp

// CreateProblemResponse 创建题目响应
type CreateProblemResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Id      int64  `json:"id"`
}

// UpdateProblemResponse 更新题目响应
type UpdateProblemResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// DeleteProblemResponse 删除题目响应
type DeleteProblemResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// GetProblemResponse 查询题目响应
type GetProblemResponse struct {
	Code    int32        `json:"code"`
	Message string       `json:"message"`
	Problem *ProblemInfo `json:"problem"`
}

// ProblemBriefInfo 题目简要信息（用于列表展示）
type ProblemBriefInfo struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	TitleSlug  string `json:"title_slug"`
	Difficulty string `json:"difficulty"`
	Tags       string `json:"tags"`
}

// ListProblemsResponse 题库列表响应
type ListProblemsResponse struct {
	Code     int32               `json:"code"`
	Message  string              `json:"message"`
	Total    int64               `json:"total"`
	Problems []*ProblemBriefInfo `json:"problems"`
}
