package req

// CodeRunRequest 代码运行/测试请求
type CodeRunRequest struct {
	ProblemId int64  `json:"problem_id"` // 题目ID
	Language  string `json:"language"`   // 语言：python/java/go/cpp/c
	Code      string `json:"code"`       // 用户代码
	RunType   string `json:"run_type"`   // test（测试样例）或 submit（提交）
	TestInput string `json:"test_input"` // 测试模式下直接传入的样例输入（run_type=test 时使用）
}

// GetCodeRunResultRequest 查询代码运行结果请求
type GetCodeRunResultRequest struct {
	RunId int64 `json:"run_id"` // 运行记录ID
}
