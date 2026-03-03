package service_impl

import (
	"context"

	"github.com/yzf120/elysia-backend/consts"
	"github.com/yzf120/elysia-backend/errs"
	codeReq "github.com/yzf120/elysia-backend/model/code/req"
	codeRsp "github.com/yzf120/elysia-backend/model/code/rsp"
	"github.com/yzf120/elysia-backend/service"
)

// CodeRunServiceImpl 代码运行服务实现（只做出入参处理）
type CodeRunServiceImpl struct {
	codeRunService *service.CodeRunService
}

// NewCodeRunServiceImpl 创建代码运行服务实现
func NewCodeRunServiceImpl() *CodeRunServiceImpl {
	return &CodeRunServiceImpl{
		codeRunService: service.NewCodeRunService(),
	}
}

// SubmitCodeRun 提交代码运行任务
func (s *CodeRunServiceImpl) SubmitCodeRun(ctx context.Context, studentId string, request *codeReq.CodeRunRequest) (*codeRsp.CodeRunResponse, error) {
	record, err := s.codeRunService.SubmitCodeRun(ctx, studentId, request.ProblemId, request.Language, request.Code, request.RunType, request.TestInput)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &codeRsp.CodeRunResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}
	return &codeRsp.CodeRunResponse{
		Code:    consts.SuccessCode,
		Message: "代码已提交，正在评测中",
		RunId:   record.Id,
	}, nil
}

// GetCodeRunResult 查询代码运行结果
func (s *CodeRunServiceImpl) GetCodeRunResult(ctx context.Context, request *codeReq.GetCodeRunResultRequest) (*codeRsp.CodeRunResultResponse, error) {
	record, err := s.codeRunService.GetCodeRunResult(request.RunId)
	if err != nil {
		code, msg := errs.ParseCommonError(err.Error())
		return &codeRsp.CodeRunResultResponse{
			Code:    int32(code),
			Message: msg,
		}, nil
	}
	return &codeRsp.CodeRunResultResponse{
		Code:    consts.SuccessCode,
		Message: consts.MessageQuerySuccess,
		Result: &codeRsp.CodeRunResult{
			RunId:      record.Id,
			Status:     record.Status,
			Output:     record.Output,
			ErrorMsg:   record.ErrorMsg,
			TimeCost:   record.TimeCost,
			MemoryUsed: record.MemoryUsed,
			RunType:    record.RunType,
			Language:   record.Language,
			CreatedAt:  record.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
