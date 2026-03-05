package rsp

// CodeRunResponse 提交代码运行任务的响应
type CodeRunResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	RunId   int64  `json:"run_id"` // 运行记录ID，用于轮询查询结果
}

// CodeRunResultResponse 查询代码运行结果的响应
type CodeRunResultResponse struct {
	Code    int32          `json:"code"`
	Message string         `json:"message"`
	Result  *CodeRunResult `json:"result,omitempty"`
}

// CodeRunResult 代码运行结果详情
type CodeRunResult struct {
	RunId      int64  `json:"run_id"`
	Status     string `json:"status"`      // pending/running/accepted/wrong_answer/time_limit_exceeded/memory_limit_exceeded/compile_error/runtime_error
	Output     string `json:"output"`      // 实际输出
	ErrorMsg   string `json:"error_msg"`   // 错误信息
	TimeCost   int64  `json:"time_cost"`   // 执行时间（毫秒）
	MemoryUsed int64  `json:"memory_used"` // 内存使用（KB）
	RunType    string `json:"run_type"`    // test/submit
	Language   string `json:"language"`
	Code       string `json:"code"` // 提交的代码
	CreatedAt  string `json:"created_at"`
}

// ListCodeRunRecordsResponse 查询运行记录列表的响应
type ListCodeRunRecordsResponse struct {
	Code    int32            `json:"code"`
	Message string           `json:"message"`
	Records []*CodeRunResult `json:"records"`
}
