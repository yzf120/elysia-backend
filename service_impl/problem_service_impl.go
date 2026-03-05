package service_impl

import (
	"context"

	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/problem"
	"github.com/yzf120/elysia-backend/model/problem/req"
	"github.com/yzf120/elysia-backend/model/problem/rsp"
	"github.com/yzf120/elysia-backend/service"
)

// ProblemServiceImpl 题目服务实现（只做出入参处理）
type ProblemServiceImpl struct {
	problemService *service.ProblemService
}

// NewProblemServiceImpl 创建题目服务实现
func NewProblemServiceImpl() *ProblemServiceImpl {
	return &ProblemServiceImpl{
		problemService: service.NewProblemService(),
	}
}

// CreateProblem 创建题目
func (s *ProblemServiceImpl) CreateProblem(ctx context.Context, request *req.CreateProblemRequest) (*rsp.CreateProblemResponse, error) {
	p := &problem.Problem{
		Title:               request.Title,
		TitleSlug:           request.TitleSlug,
		Difficulty:          request.Difficulty,
		Tags:                request.Tags,
		Description:         request.Description,
		Explanation:         request.Explanation,
		Hint:                request.Hint,
		Constraints:         request.Constraints,
		AdvancedRequirement: request.AdvancedRequirement,
		TestCases:           request.TestCases,
		Showcase:            request.Showcase,
		TimeLimit:           request.TimeLimit,
		MemoryLimit:         request.MemoryLimit,
	}
	created, err := s.problemService.CreateProblem(ctx, p)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.CreateProblemResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}
	return &rsp.CreateProblemResponse{
		Code:    consts.SuccessCode,
		Message: consts.MessageCreateProblemSuccess,
		Id:      created.Id,
	}, nil
}

// GetProblem 查询题目
func (s *ProblemServiceImpl) GetProblem(ctx context.Context, request *req.GetProblemRequest) (*rsp.GetProblemResponse, error) {
	p, err := s.problemService.GetProblemById(request.Id)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.GetProblemResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}
	return &rsp.GetProblemResponse{
		Code:    consts.SuccessCode,
		Message: consts.MessageGetProblemSuccess,
		Problem: &rsp.ProblemInfo{
			Id:                  p.Id,
			Title:               p.Title,
			TitleSlug:           p.TitleSlug,
			Difficulty:          p.Difficulty,
			Tags:                p.Tags,
			Description:         p.Description,
			Explanation:         p.Explanation,
			Hint:                p.Hint,
			Constraints:         p.Constraints,
			AdvancedRequirement: p.AdvancedRequirement,
			TestCases:           p.TestCases,
			Showcase:            p.Showcase,
			TimeLimit:           p.TimeLimit,
			MemoryLimit:         p.MemoryLimit,
			CreatedAt:           p.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:           p.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

// UpdateProblem 更新题目
func (s *ProblemServiceImpl) UpdateProblem(ctx context.Context, request *req.UpdateProblemRequest) (*rsp.UpdateProblemResponse, error) {
	updates := make(map[string]interface{})
	if request.Title != "" {
		updates["title"] = request.Title
	}
	if request.TitleSlug != "" {
		updates["title_slug"] = request.TitleSlug
	}
	if request.Difficulty != "" {
		updates["difficulty"] = request.Difficulty
	}
	if request.Tags != "" {
		updates["tags"] = request.Tags
	}
	if request.Description != "" {
		updates["description"] = request.Description
	}
	if request.Explanation != "" {
		updates["explanation"] = request.Explanation
	}
	if request.Hint != "" {
		updates["hint"] = request.Hint
	}
	if request.Constraints != "" {
		updates["constraints"] = request.Constraints
	}
	if request.AdvancedRequirement != "" {
		updates["advanced_requirement"] = request.AdvancedRequirement
	}
	if request.TestCases != "" {
		updates["test_cases"] = request.TestCases
	}
	if request.Showcase != "" {
		updates["showcase"] = request.Showcase
	}
	if request.TimeLimit > 0 {
		updates["time_limit"] = request.TimeLimit
	}
	if request.MemoryLimit > 0 {
		updates["memory_limit"] = request.MemoryLimit
	}

	_, err := s.problemService.UpdateProblem(request.Id, updates)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.UpdateProblemResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}
	return &rsp.UpdateProblemResponse{
		Code:    consts.SuccessCode,
		Message: consts.MessageUpdateProblemSuccess,
	}, nil
}

// DeleteProblem 删除题目
func (s *ProblemServiceImpl) DeleteProblem(ctx context.Context, request *req.DeleteProblemRequest) (*rsp.DeleteProblemResponse, error) {
	err := s.problemService.DeleteProblem(request.Id)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.DeleteProblemResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}
	return &rsp.DeleteProblemResponse{
		Code:    consts.SuccessCode,
		Message: consts.MessageDeleteProblemSuccess,
	}, nil
}

// ListProblems 查询题库列表
func (s *ProblemServiceImpl) ListProblems(ctx context.Context, request *req.ListProblemsRequest) (*rsp.ListProblemsResponse, error) {
	problems, total, err := s.problemService.ListProblems(request.Keyword, request.Difficulty, request.Page, request.PageSize)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &rsp.ListProblemsResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}
	briefs := make([]*rsp.ProblemBriefInfo, 0, len(problems))
	for _, p := range problems {
		briefs = append(briefs, &rsp.ProblemBriefInfo{
			Id:         p.Id,
			Title:      p.Title,
			TitleSlug:  p.TitleSlug,
			Difficulty: p.Difficulty,
			Tags:       p.Tags,
		})
	}
	return &rsp.ListProblemsResponse{
		Code:     consts.SuccessCode,
		Message:  consts.MessageListProblemsSuccess,
		Total:    total,
		Problems: briefs,
	}, nil
}
