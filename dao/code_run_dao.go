package dao

import (
	"github.com/yzf120/elysia-backend/model/code"
)

// CodeRunDAO 代码运行记录数据访问对象
type CodeRunDAO interface {
	CreateCodeRun(r *code.CodeRun) error
	GetCodeRunById(id int64) (*code.CodeRun, error)
	UpdateCodeRun(id int64, updates map[string]interface{}) error
	ListCodeRunsByStudent(studentId string, problemId int64, limit int) ([]*code.CodeRun, error)
	// BatchGetAcceptedProblems 批量查询学生已完全通过（accepted）的题目ID集合
	BatchGetAcceptedProblems(studentId string, problemIds []int64) (map[int64]bool, error)
}

type codeRunDAOImpl struct{}

// NewCodeRunDAO 创建代码运行DAO
func NewCodeRunDAO() CodeRunDAO {
	return &codeRunDAOImpl{}
}

// CreateCodeRun 创建代码运行记录
func (d *codeRunDAOImpl) CreateCodeRun(r *code.CodeRun) error {
	return DB.Create(r).Error
}

// GetCodeRunById 根据ID查询运行记录
func (d *codeRunDAOImpl) GetCodeRunById(id int64) (*code.CodeRun, error) {
	var r code.CodeRun
	err := DB.Where("id = ?", id).First(&r).Error
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// UpdateCodeRun 更新运行记录
func (d *codeRunDAOImpl) UpdateCodeRun(id int64, updates map[string]interface{}) error {
	return DB.Model(&code.CodeRun{}).Where("id = ?", id).Updates(updates).Error
}

// ListCodeRunsByStudent 查询学生的运行记录列表
func (d *codeRunDAOImpl) ListCodeRunsByStudent(studentId string, problemId int64, limit int) ([]*code.CodeRun, error) {
	var records []*code.CodeRun
	query := DB.Where("student_id = ? AND problem_id = ?", studentId, problemId).
		Order("created_at DESC").
		Limit(limit)
	err := query.Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// BatchGetAcceptedProblems 批量查询学生已完全通过（accepted）的题目ID集合
// 只查 run_type='submit' 且 status='accepted' 的记录
func (d *codeRunDAOImpl) BatchGetAcceptedProblems(studentId string, problemIds []int64) (map[int64]bool, error) {
	result := make(map[int64]bool)
	if len(problemIds) == 0 {
		return result, nil
	}
	var acceptedIds []int64
	err := DB.Model(&code.CodeRun{}).
		Select("DISTINCT problem_id").
		Where("student_id = ? AND problem_id IN ? AND run_type = 'submit' AND status = 'accepted'", studentId, problemIds).
		Pluck("problem_id", &acceptedIds).Error
	if err != nil {
		return nil, err
	}
	for _, id := range acceptedIds {
		result[id] = true
	}
	return result, nil
}
