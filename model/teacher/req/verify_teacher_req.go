package req

// VerifyTeacherRequest 审核教师请求
type VerifyTeacherRequest struct {
	TeacherId  string `json:"teacher_id"`  // 教师ID（必填）
	VerifierId string `json:"verifier_id"` // 审核人ID（必填）
	Approved   bool   `json:"approved"`    // 是否通过（必填）
	Remark     string `json:"remark"`      // 审核备注（可选）
}