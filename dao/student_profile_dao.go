package dao

import (
	"github.com/yzf120/elysia-backend/model/profile"
	"gorm.io/gorm"
)

// StudentProfileDAO 学生画像数据访问对象
type StudentProfileDAO interface {
	// GetProfileByStudentId 根据学生ID查询画像
	GetProfileByStudentId(studentId string) (*profile.StudentProfile, error)
	// UpsertProfile 创建或更新画像（基于 student_id 唯一键）
	UpsertProfile(p *profile.StudentProfile) error
	// UpdateProfile 更新画像字段
	UpdateProfile(studentId string, updates map[string]interface{}) error
}

type studentProfileDAOImpl struct{}

// NewStudentProfileDAO 创建学生画像DAO
func NewStudentProfileDAO() StudentProfileDAO {
	return &studentProfileDAOImpl{}
}

// GetProfileByStudentId 根据学生ID查询画像
func (d *studentProfileDAOImpl) GetProfileByStudentId(studentId string) (*profile.StudentProfile, error) {
	var p profile.StudentProfile
	err := DB.Where("student_id = ?", studentId).First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// UpsertProfile 创建或更新画像
func (d *studentProfileDAOImpl) UpsertProfile(p *profile.StudentProfile) error {
	// 先查询是否存在
	existing, err := d.GetProfileByStudentId(p.StudentId)
	if err != nil {
		return err
	}
	if existing == nil {
		// 不存在则创建
		return DB.Create(p).Error
	}
	// 存在则更新（排除 id 和 student_id）
	return DB.Model(&profile.StudentProfile{}).Where("student_id = ?", p.StudentId).Updates(map[string]interface{}{
		"total_submissions":       p.TotalSubmissions,
		"accepted_count":          p.AcceptedCount,
		"wrong_answer_count":      p.WrongAnswerCount,
		"compile_error_count":     p.CompileErrorCount,
		"runtime_error_count":     p.RuntimeErrorCount,
		"tle_count":               p.TLECount,
		"mle_count":               p.MLECount,
		"accept_rate":             p.AcceptRate,
		"solved_problem_count":    p.SolvedProblemCount,
		"attempted_problem_count": p.AttemptedProblemCount,
		"preferred_language":      p.PreferredLanguage,
		"language_stats":          p.LanguageStats,
		"avg_time_cost":           p.AvgTimeCost,
		"common_errors":           p.CommonErrors,
		"difficulty_level":        p.DifficultyLevel,
		"last_submit_time":        p.LastSubmitTime,
	}).Error
}

// UpdateProfile 更新画像字段
func (d *studentProfileDAOImpl) UpdateProfile(studentId string, updates map[string]interface{}) error {
	return DB.Model(&profile.StudentProfile{}).Where("student_id = ?", studentId).Updates(updates).Error
}
