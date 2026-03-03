package code

import "time"

// CodeRun 代码运行记录
type CodeRun struct {
	Id         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProblemId  int64     `gorm:"column:problem_id;not null;index" json:"problem_id"`
	StudentId  string    `gorm:"column:student_id;type:varchar(64);not null;index" json:"student_id"`
	Language   string    `gorm:"column:language;type:varchar(32);not null" json:"language"`
	Code       string    `gorm:"column:code;type:longtext;not null" json:"code"`
	RunType    string    `gorm:"column:run_type;type:enum('test','submit');not null;default:'test'" json:"run_type"`
	Status     string    `gorm:"column:status;type:enum('pending','running','accepted','wrong_answer','time_limit_exceeded','memory_limit_exceeded','compile_error','runtime_error');not null;default:'pending'" json:"status"`
	Output     string    `gorm:"column:output;type:text" json:"output"`
	ErrorMsg   string    `gorm:"column:error_msg;type:text" json:"error_msg"`
	TimeCost   int64     `gorm:"column:time_cost" json:"time_cost"`     // 执行时间（毫秒）
	MemoryUsed int64     `gorm:"column:memory_used" json:"memory_used"` // 内存使用（KB）
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (CodeRun) TableName() string {
	return "code_run"
}
