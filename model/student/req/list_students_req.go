package req

// ListStudentsRequest 查询学生列表请求
type ListStudentsRequest struct {
	Page             int32  `json:"page"`               // 页码（从1开始）
	PageSize         int32  `json:"page_size"`          // 每页数量
	Major            string `json:"major"`              // 专业筛选（可选）
	Grade            string `json:"grade"`              // 年级筛选（可选）
	ProgrammingLevel string `json:"programming_level"`  // 编程基础筛选（可选）
}