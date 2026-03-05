package router

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/authen"
	"github.com/yzf120/elysia-backend/errs"
	codeReq "github.com/yzf120/elysia-backend/model/code/req"
	"github.com/yzf120/elysia-backend/service_impl"
)

var (
	codeRunService *service_impl.CodeRunServiceImpl
)

// registerCodeRun 注册代码运行相关路由（学生端，需要认证）
func registerCodeRun(protectedRouter *mux.Router) {
	// 提交代码运行/测试任务
	protectedRouter.HandleFunc("/student/code/run", submitCodeRunHandler).Methods("POST")
	// 查询代码运行结果（轮询）
	protectedRouter.HandleFunc("/student/code/result", getCodeRunResultHandler).Methods("GET")
	// 查询学生某题的运行记录列表（最新10条，倒序）
	protectedRouter.HandleFunc("/student/code/records", listCodeRunRecordsHandler).Methods("GET")
	// 批量查询学生已完全通过的题目ID集合（用于课程目录打钩）
	protectedRouter.HandleFunc("/student/code/progress", getCodeProgressHandler).Methods("GET")
}

// submitCodeRunHandler 提交代码运行任务处理器
func submitCodeRunHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	// 从 JWT 中获取学生ID（RoleID 即为学生ID）
	studentId, ok := authen.GetRoleIDFromContext(ctx)
	if !ok || studentId == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "未授权，请先登录")
		return
	}

	request := &codeReq.CodeRunRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if request.ProblemId <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "problem_id 无效")
		return
	}
	if request.Code == "" {
		writeErrorResponse(w, http.StatusBadRequest, "代码不能为空")
		return
	}
	if request.Language == "" {
		writeErrorResponse(w, http.StatusBadRequest, "language 不能为空")
		return
	}
	if request.RunType == "" {
		request.RunType = "test"
	}

	resp, err := codeRunService.SubmitCodeRun(ctx, studentId, request)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusInternalServerError, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBytes)
		return
	}

	if resp.Code != 0 {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"run_id":  resp.RunId,
		"message": resp.Message,
	})
}

// getCodeRunResultHandler 查询代码运行结果处理器（前端轮询）
func getCodeRunResultHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	runIdStr := r.URL.Query().Get("run_id")
	runId, err := strconv.ParseInt(runIdStr, 10, 64)
	if err != nil || runId <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "run_id 无效")
		return
	}

	request := &codeReq.GetCodeRunResultRequest{RunId: runId}
	resp, err := codeRunService.GetCodeRunResult(ctx, request)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusInternalServerError, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBytes)
		return
	}

	if resp.Code != 0 {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	writeSuccessResponse(w, resp.Result)
}

// listCodeRunRecordsHandler 查询学生某题的运行记录列表（最新10条，倒序）
func listCodeRunRecordsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	// 从 JWT 中获取学生ID
	studentId, ok := authen.GetRoleIDFromContext(ctx)
	if !ok || studentId == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "未授权，请先登录")
		return
	}

	problemIdStr := r.URL.Query().Get("problem_id")
	problemId, err := strconv.ParseInt(problemIdStr, 10, 64)
	if err != nil || problemId <= 0 {
		writeErrorResponse(w, http.StatusBadRequest, "problem_id 无效")
		return
	}

	resp, err := codeRunService.ListCodeRunRecords(ctx, studentId, problemId)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusInternalServerError, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBytes)
		return
	}

	if resp.Code != 0 {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(int(resp.Code), resp.Message),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBytes)
		return
	}

	writeSuccessResponse(w, map[string]interface{}{
		"records": resp.Records,
	})
}

// getCodeProgressHandler 批量查询学生已完全通过的题目ID集合
// GET /student/code/progress?problem_ids=1,2,3
func getCodeProgressHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	setResponseHeaders(w)

	// 从 JWT 中获取学生ID
	studentId, ok := authen.GetRoleIDFromContext(ctx)
	if !ok || studentId == "" {
		writeErrorResponse(w, http.StatusUnauthorized, "未授权，请先登录")
		return
	}

	// 解析 problem_ids 参数（逗号分隔）
	problemIdsStr := r.URL.Query().Get("problem_ids")
	if problemIdsStr == "" {
		writeSuccessResponse(w, map[string]interface{}{
			"accepted_problem_ids": []int64{},
		})
		return
	}

	var problemIds []int64
	for _, s := range strings.Split(problemIdsStr, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil || id <= 0 {
			continue
		}
		problemIds = append(problemIds, id)
	}

	if len(problemIds) == 0 {
		writeSuccessResponse(w, map[string]interface{}{
			"accepted_problem_ids": []int64{},
		})
		return
	}

	acceptedMap, err := codeRunService.BatchGetAcceptedProblems(ctx, studentId, problemIds)
	if err != nil {
		errResp := &errs.BaseResponse{
			Data:  nil,
			Error: errs.NewError(http.StatusInternalServerError, err.Error()),
		}
		respBytes, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBytes)
		return
	}

	// 返回已通过的 problem_id 列表
	acceptedIds := make([]int64, 0, len(acceptedMap))
	for id := range acceptedMap {
		acceptedIds = append(acceptedIds, id)
	}

	writeSuccessResponse(w, map[string]interface{}{
		"accepted_problem_ids": acceptedIds,
	})
}
