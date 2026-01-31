package req

// RegisterTeacherRequest 教师注册请求
type RegisterTeacherRequest struct {
	PhoneNumber      string   `json:"phone_number"`      // 手机号（必填）
	Password         string   `json:"password"`          // 密码（必填）
	EmployeeNumber   string   `json:"employee_number"`   // 工号（必填）
	SchoolEmail      string   `json:"school_email"`      // 学校邮箱（必填）
	RealName         string   `json:"real_name"`         // 真实姓名（必填）
	Department       string   `json:"department"`        // 所属院系（必填）
	TeachingSubjects []string `json:"teaching_subjects"` // 授课科目（可选）
}