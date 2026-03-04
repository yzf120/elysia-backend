package dao

import (
	"github.com/yzf120/elysia-backend/model/problem"
)

// ProblemDAO 题目数据访问对象
type ProblemDAO interface {
	CreateProblem(p *problem.Problem) error
	GetProblemById(id int64) (*problem.Problem, error)
	UpdateProblem(id int64, updates map[string]interface{}) error
	DeleteProblem(id int64) error
	ListProblems(keyword, difficulty string, page, pageSize int) ([]*problem.Problem, int64, error)
}

type problemDAOImpl struct{}

// NewProblemDAO 创建题目DAO
func NewProblemDAO() ProblemDAO {
	return &problemDAOImpl{}
}

// CreateProblem 创建题目
func (d *problemDAOImpl) CreateProblem(p *problem.Problem) error {
	return DB.Create(p).Error
}

// GetProblemById 根据题目ID查询题目
func (d *problemDAOImpl) GetProblemById(id int64) (*problem.Problem, error) {
	var p problem.Problem
	err := DB.Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdateProblem 更新题目信息
func (d *problemDAOImpl) UpdateProblem(id int64, updates map[string]interface{}) error {
	return DB.Model(&problem.Problem{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteProblem 删除题目
func (d *problemDAOImpl) DeleteProblem(id int64) error {
	return DB.Where("id = ?", id).Delete(&problem.Problem{}).Error
}

// ListProblems 分页查询题库列表，支持关键词和难度筛选
func (d *problemDAOImpl) ListProblems(keyword, difficulty string, page, pageSize int) ([]*problem.Problem, int64, error) {
	var problems []*problem.Problem
	var total int64

	query := DB.Model(&problem.Problem{})
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Select("id, title, title_slug, difficulty, tags").
		Order("id ASC").
		Offset(offset).Limit(pageSize).
		Find(&problems).Error; err != nil {
		return nil, 0, err
	}
	return problems, total, nil
}
