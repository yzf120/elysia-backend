package router

import (
	"encoding/json"
	"net/http"
	"strconv"

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
