package profile

import "time"

// StudentProfile 学生做题画像（每次做题后自动更新）
type StudentProfile struct {
	Id                    int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	StudentId             string     `gorm:"column:student_id;type:varchar(64);uniqueIndex;not null" json:"student_id"`
	TotalSubmissions      int        `gorm:"column:total_submissions;not null;default:0" json:"total_submissions"`
	AcceptedCount         int        `gorm:"column:accepted_count;not null;default:0" json:"accepted_count"`
	WrongAnswerCount      int        `gorm:"column:wrong_answer_count;not null;default:0" json:"wrong_answer_count"`
	CompileErrorCount     int        `gorm:"column:compile_error_count;not null;default:0" json:"compile_error_count"`
	RuntimeErrorCount     int        `gorm:"column:runtime_error_count;not null;default:0" json:"runtime_error_count"`
	TLECount              int        `gorm:"column:tle_count;not null;default:0" json:"tle_count"`
	MLECount              int        `gorm:"column:mle_count;not null;default:0" json:"mle_count"`
	AcceptRate            float64    `gorm:"column:accept_rate;type:decimal(5,4);not null;default:0" json:"accept_rate"`
	SolvedProblemCount    int        `gorm:"column:solved_problem_count;not null;default:0" json:"solved_problem_count"`
	AttemptedProblemCount int        `gorm:"column:attempted_problem_count;not null;default:0" json:"attempted_problem_count"`
	PreferredLanguage     string     `gorm:"column:preferred_language;type:varchar(32)" json:"preferred_language"`
	LanguageStats         string     `gorm:"column:language_stats;type:text" json:"language_stats"`
	AvgTimeCost           int64      `gorm:"column:avg_time_cost;default:0" json:"avg_time_cost"`
	CommonErrors          string     `gorm:"column:common_errors;type:text" json:"common_errors"`
	DifficultyLevel       string     `gorm:"column:difficulty_level;type:varchar(32);default:beginner" json:"difficulty_level"`
	LastSubmitTime        *time.Time `gorm:"column:last_submit_time" json:"last_submit_time"`
	CreatedAt             time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (StudentProfile) TableName() string {
	return "student_profile"
}
