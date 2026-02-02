package req

// CreateStudentRequest 创建学生信息请求
type CreateStudentRequest struct {
	StudentId        string   `json:"student_id"`        // 学生ID（必填）
	Major            string   `json:"major"`             // 专业（必填）
	Grade            string   `json:"grade"`             // 年级（必填）
	ProgrammingLevel string   `json:"programming_level"` // 编程基础（必填）
	Interests        []string `json:"interests"`         // 兴趣爱好（可选）
	LearningTags     []string `json:"learning_tags"`     // 学习标签（可选）
}

// LoginWithPasswordRequest 学号+密码登录请求
type LoginWithPasswordRequest struct {
	StudentNumber string `json:"student_number"`
	Password      string `json:"password"`
}
