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
