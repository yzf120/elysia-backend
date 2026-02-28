package rsp

// ProblemInfo 题目信息
type ProblemInfo struct {
	Id                  int64  `json:"id"`
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
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}
