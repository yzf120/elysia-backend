package req

// CreateStudentRequest 创建学生信息请求
type CreateStudentRequest struct {
	UserId           string   `json:"user_id"`            // 用户ID（必填）
	Major            string   `json:"major"`              // 专业（必填）
	Grade            string   `json:"grade"`              // 年级（必填）
	ProgrammingLevel string   `json:"programming_level"`  // 编程基础（必填）
	Interests        []string `json:"interests"`          // 兴趣爱好（可选）
	LearningTags     []string `json:"learning_tags"`      // 学习标签（可选）
}