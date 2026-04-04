package service

import (
	"encoding/json"
	"log"
	"sort"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/model/profile"
)

// StudentProfileService 学生画像服务
type StudentProfileService struct {
	profileDAO dao.StudentProfileDAO
}

// NewStudentProfileService 创建学生画像服务
func NewStudentProfileService() *StudentProfileService {
	return &StudentProfileService{
		profileDAO: dao.NewStudentProfileDAO(),
	}
}

// statusCountResult 状态统计结果（用于 SQL 聚合查询）
type statusCountResult struct {
	Status string `gorm:"column:status"`
	Count  int    `gorm:"column:cnt"`
}

// langCountResult 语言统计结果
type langCountResult struct {
	Language string `gorm:"column:language"`
	Count    int    `gorm:"column:cnt"`
}

// UpdateProfileAfterSubmit 做题提交后更新学生画像
// 该方法会从 code_run 表中重新聚合统计数据，确保画像数据的准确性
func (s *StudentProfileService) UpdateProfileAfterSubmit(studentId string) {
	if studentId == "" {
		return
	}

	db := dao.DB

	// 1. 统计各状态的提交次数（仅 submit 类型）
	var statusCounts []statusCountResult
	err := db.Table("code_run").
		Select("status, COUNT(*) as cnt").
		Where("student_id = ? AND run_type = 'submit'", studentId).
		Group("status").
		Scan(&statusCounts).Error
	if err != nil {
		log.Printf("[StudentProfile] 统计提交状态失败: %v", err)
		return
	}

	// 解析各状态计数
	var totalSubmissions, acceptedCount, wrongAnswerCount, compileErrorCount, runtimeErrorCount, tleCount, mleCount int
	for _, sc := range statusCounts {
		totalSubmissions += sc.Count
		switch sc.Status {
		case "accepted":
			acceptedCount = sc.Count
		case "wrong_answer":
			wrongAnswerCount = sc.Count
		case "compile_error":
			compileErrorCount = sc.Count
		case "runtime_error":
			runtimeErrorCount = sc.Count
		case "time_limit_exceeded":
			tleCount = sc.Count
		case "memory_limit_exceeded":
			mleCount = sc.Count
		}
	}

	// 2. 计算通过率
	var acceptRate float64
	if totalSubmissions > 0 {
		acceptRate = float64(acceptedCount) / float64(totalSubmissions)
	}

	// 3. 统计已解决题目数（至少一次 accepted 的不同题目）
	var solvedCount int64
	err = db.Table("code_run").
		Where("student_id = ? AND run_type = 'submit' AND status = 'accepted'", studentId).
		Distinct("problem_id").
		Count(&solvedCount).Error
	if err != nil {
		log.Printf("[StudentProfile] 统计已解决题目数失败: %v", err)
		solvedCount = 0
	}

	// 4. 统计尝试过的题目数
	var attemptedCount int64
	err = db.Table("code_run").
		Where("student_id = ? AND run_type = 'submit'", studentId).
		Distinct("problem_id").
		Count(&attemptedCount).Error
	if err != nil {
		log.Printf("[StudentProfile] 统计尝试题目数失败: %v", err)
		attemptedCount = 0
	}

	// 5. 统计各编程语言使用次数
	var langCounts []langCountResult
	err = db.Table("code_run").
		Select("language, COUNT(*) as cnt").
		Where("student_id = ? AND run_type = 'submit'", studentId).
		Group("language").
		Order("cnt DESC").
		Scan(&langCounts).Error

	preferredLang := ""
	langStatsMap := make(map[string]int)
	if err == nil && len(langCounts) > 0 {
		preferredLang = langCounts[0].Language
		for _, lc := range langCounts {
			if lc.Language != "" {
				langStatsMap[lc.Language] = lc.Count
			}
		}
	}
	langStatsJSON, _ := json.Marshal(langStatsMap)

	// 6. 统计平均执行时间（仅 accepted 的记录）
	var avgTimeCost float64
	err = db.Table("code_run").
		Select("COALESCE(AVG(time_cost), 0)").
		Where("student_id = ? AND run_type = 'submit' AND status = 'accepted' AND time_cost > 0", studentId).
		Scan(&avgTimeCost).Error
	if err != nil {
		avgTimeCost = 0
	}

	// 7. 统计常见错误类型（按频率排序，取 Top 3）
	commonErrors := s.calcCommonErrors(wrongAnswerCount, compileErrorCount, runtimeErrorCount, tleCount, mleCount)
	commonErrorsJSON, _ := json.Marshal(commonErrors)

	// 8. 推断能力等级
	difficultyLevel := s.inferDifficultyLevel(totalSubmissions, acceptRate, int(solvedCount))

	// 9. 获取最近提交时间
	var lastSubmitTime *time.Time
	var lastRecord struct {
		CreatedAt time.Time `gorm:"column:created_at"`
	}
	err = db.Table("code_run").
		Select("created_at").
		Where("student_id = ? AND run_type = 'submit'", studentId).
		Order("created_at DESC").
		Limit(1).
		Scan(&lastRecord).Error
	if err == nil && !lastRecord.CreatedAt.IsZero() {
		lastSubmitTime = &lastRecord.CreatedAt
	}

	// 10. 构建画像并写入数据库
	p := &profile.StudentProfile{
		StudentId:             studentId,
		TotalSubmissions:      totalSubmissions,
		AcceptedCount:         acceptedCount,
		WrongAnswerCount:      wrongAnswerCount,
		CompileErrorCount:     compileErrorCount,
		RuntimeErrorCount:     runtimeErrorCount,
		TLECount:              tleCount,
		MLECount:              mleCount,
		AcceptRate:            acceptRate,
		SolvedProblemCount:    int(solvedCount),
		AttemptedProblemCount: int(attemptedCount),
		PreferredLanguage:     preferredLang,
		LanguageStats:         string(langStatsJSON),
		AvgTimeCost:           int64(avgTimeCost),
		CommonErrors:          string(commonErrorsJSON),
		DifficultyLevel:       difficultyLevel,
		LastSubmitTime:        lastSubmitTime,
	}

	if err := s.profileDAO.UpsertProfile(p); err != nil {
		log.Printf("[StudentProfile] 更新画像失败 (student_id=%s): %v", studentId, err)
		return
	}

	log.Printf("[StudentProfile] 画像更新成功: student_id=%s, submissions=%d, accept_rate=%.2f%%, solved=%d, level=%s",
		studentId, totalSubmissions, acceptRate*100, solvedCount, difficultyLevel)
}

// calcCommonErrors 计算常见错误类型（按频率排序，取 Top 3）
func (s *StudentProfileService) calcCommonErrors(wrongAnswer, compileError, runtimeError, tle, mle int) []string {
	type errorStat struct {
		Name  string
		Count int
	}

	errors := []errorStat{
		{"wrong_answer", wrongAnswer},
		{"compile_error", compileError},
		{"runtime_error", runtimeError},
		{"time_limit_exceeded", tle},
		{"memory_limit_exceeded", mle},
	}

	// 按频率降序排序
	sort.Slice(errors, func(i, j int) bool {
		return errors[i].Count > errors[j].Count
	})

	// 取 Top 3（过滤掉 count=0 的）
	var result []string
	for _, e := range errors {
		if e.Count > 0 && len(result) < 3 {
			result = append(result, e.Name)
		}
	}
	return result
}

// inferDifficultyLevel 根据做题数据推断能力等级
func (s *StudentProfileService) inferDifficultyLevel(totalSubmissions int, acceptRate float64, solvedCount int) string {
	// 提交次数太少，无法判断，默认 beginner
	if totalSubmissions < 3 {
		return "beginner"
	}

	// 综合评分（满分 100）
	score := 0.0

	// 通过率权重 40%（通过率越高分越高）
	score += acceptRate * 40

	// 解题数量权重 30%（解题越多分越高，上限 30 题满分）
	solvedScore := float64(solvedCount) / 30.0
	if solvedScore > 1.0 {
		solvedScore = 1.0
	}
	score += solvedScore * 30

	// 提交活跃度权重 30%（提交越多分越高，上限 100 次满分）
	submitScore := float64(totalSubmissions) / 100.0
	if submitScore > 1.0 {
		submitScore = 1.0
	}
	score += submitScore * 30

	// 根据综合评分判断等级
	switch {
	case score >= 60:
		return "advanced"
	case score >= 30:
		return "intermediate"
	default:
		return "beginner"
	}
}

// GetProfile 获取学生画像
func (s *StudentProfileService) GetProfile(studentId string) (*profile.StudentProfile, error) {
	return s.profileDAO.GetProfileByStudentId(studentId)
}
