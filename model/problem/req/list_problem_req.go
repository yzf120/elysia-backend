package req

// ListProblemsRequest 题库列表搜索请求
type ListProblemsRequest struct {
	Keyword    string `json:"keyword"`     // 搜索关键词（题目标题模糊匹配）
	Difficulty string `json:"difficulty"`  // 难度筛选：简单/中等/困难，空表示不筛选
	Page       int    `json:"page"`        // 页码，从1开始
	PageSize   int    `json:"page_size"`   // 每页数量，默认20
}
