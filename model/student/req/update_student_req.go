package req

// UpdateStudentRequest 更新学生信息请求
type UpdateStudentRequest struct {
	StudentId        string   `json:"student_id"`         // 学生ID（必填）
	Name             string   `json:"name"`               // 姓名（可选）
	Email            string   `json:"email"`              // 邮箱（可选）
	Major            string   `json:"major"`              // 专业（可选）
	Grade            string   `json:"grade"`              // 年级（可选）
	ProgrammingLevel string   `json:"programming_level"`  // 编程基础（可选）
	Interests        []string `json:"interests"`          // 兴趣爱好（可选）
	LearningTags     []string `json:"learning_tags"`      // 学习标签（可选）
}