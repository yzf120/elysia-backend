package req

// CreateProblemRequest 创建题目请求
type CreateProblemRequest struct {
	Title               string `json:"title"`
	TitleSlug           string `json:"title_slug"`
	Difficulty          string `json:"difficulty"`
	Tags                string `json:"tags"`
	Description         string `json:"description"`
	Explanation         string `json:"explanation"`
	Hint                string `json:"hint"`
	Constraints         string `json:"constraints"`
	AdvancedRequirement string `json:"advanced_requirement"`
	TestCases           string `json:"test_cases"`
}
