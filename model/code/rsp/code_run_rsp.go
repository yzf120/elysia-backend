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

// CodeCheckResponse 代码语法检查响应
type CodeCheckResponse struct {
	Code        int32         `json:"code"`
	Message     string        `json:"message"`
	HasError    bool          `json:"has_error"`    // 是否有编译错误
	Diagnostics []*Diagnostic `json:"diagnostics"`  // 错误诊断列表
}

// Diagnostic 单条编译错误诊断信息
type Diagnostic struct {
	Line      int    `json:"line"`       // 错误行号（从1开始）
	Column    int    `json:"column"`     // 错误列号（从1开始）
	EndLine   int    `json:"end_line"`   // 错误结束行号
	EndColumn int    `json:"end_column"` // 错误结束列号
	Severity  string `json:"severity"`   // 错误级别：error / warning
	Message   string `json:"message"`    // 错误信息
}