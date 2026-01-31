package req

// UpdateTeacherRequest 更新教师信息请求
type UpdateTeacherRequest struct {
	TeacherId        string   `json:"teacher_id"`        // 教师ID（必填）
	TeachingSubjects []string `json:"teaching_subjects"` // 授课科目（可选）
	TeachingYears    int32    `json:"teaching_years"`    // 教龄（可选）
	Department       string   `json:"department"`       // 所属院系（可选）
	Title            string   `json:"title"`            // 职称（可选）
}