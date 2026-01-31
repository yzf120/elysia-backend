package rsp

// TeacherInfo 教师信息
type TeacherInfo struct {
	TeacherId          string   `json:"teacher_id"`          // 教师ID
	UserId             string   `json:"user_id"`             // 用户ID
	EmployeeNumber     string   `json:"employee_number"`     // 工号
	SchoolEmail        string   `json:"school_email"`        // 学校邮箱
	TeachingSubjects   []string `json:"teaching_subjects"`   // 授课科目
	TeachingYears      int32    `json:"teaching_years"`      // 教龄
	Department         string   `json:"department"`          // 所属院系
	Title              string   `json:"title"`               // 职称
	VerificationStatus int32    `json:"verification_status"` // 认证状态
	Status             int32    `json:"status"`              // 状态
	CreateTime         string   `json:"create_time"`         // 创建时间
	UpdateTime         string   `json:"update_time"`         // 更新时间
}