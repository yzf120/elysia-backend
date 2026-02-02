package rsp

// StudentInfo 学生信息
type StudentInfo struct {
	StudentId        string   `json:"student_id"`        // 学生ID
	StudentNumber    string   `json:"student_number"`    // 学号
	Major            string   `json:"major"`             // 专业
	Grade            string   `json:"grade"`             // 年级
	ProgrammingLevel string   `json:"programming_level"` // 编程基础
	Interests        []string `json:"interests"`         // 兴趣爱好
	LearningTags     []string `json:"learning_tags"`     // 学习标签
	LearningProgress string   `json:"learning_progress"` // 学习进度
	Status           int32    `json:"status"`            // 状态
	CreateTime       string   `json:"create_time"`       // 创建时间
	UpdateTime       string   `json:"update_time"`       // 更新时间
}
