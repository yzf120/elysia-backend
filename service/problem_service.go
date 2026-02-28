package service

import (
	"context"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/problem"
)

// ProblemService 题目服务
type ProblemService struct {
	problemDAO dao.ProblemDAO
}

// NewProblemService 创建题目服务
func NewProblemService() *ProblemService {
	return &ProblemService{
		problemDAO: dao.NewProblemDAO(),
	}
}

// CreateProblem 创建题目
func (s *ProblemService) CreateProblem(ctx context.Context, p *problem.Problem) (*problem.Problem, error) {
	if p.Title == "" || p.TitleSlug == "" || p.Description == "" || p.TestCases == "" {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "必填参数不能为空")
	}
	if err := s.problemDAO.CreateProblem(p); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "创建题目失败: "+err.Error())
	}
	return p, nil
}

// GetProblemById 根据题目ID查询题目
func (s *ProblemService) GetProblemById(id int64) (*problem.Problem, error) {
	p, err := s.problemDAO.GetProblemById(id)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "查询题目失败: "+err.Error())
	}
	if p == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "题目不存在")
	}
	return p, nil
}

// UpdateProblem 更新题目
func (s *ProblemService) UpdateProblem(id int64, updates map[string]interface{}) (*problem.Problem, error) {
	existing, err := s.problemDAO.GetProblemById(id)
	if err != nil || existing == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "题目不存在")
	}
	if err := s.problemDAO.UpdateProblem(id, updates); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "更新题目失败: "+err.Error())
	}
	updated, err := s.problemDAO.GetProblemById(id)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "查询题目失败: "+err.Error())
	}
	return updated, nil
}

// DeleteProblem 删除题目
func (s *ProblemService) DeleteProblem(id int64) error {
	existing, err := s.problemDAO.GetProblemById(id)
	if err != nil || existing == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "题目不存在")
	}
	if err := s.problemDAO.DeleteProblem(id); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "删除题目失败: "+err.Error())
	}
	return nil
}
