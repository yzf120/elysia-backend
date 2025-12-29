package router

import (
	"net/http"

	"github.com/yzf120/elysia-backend/errs"
)

// respondWithError 返回错误响应
func respondWithError(w http.ResponseWriter, httpCode int, errCode int, message string) {
	resp := errs.NewCommonErrRspV2(errCode, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(resp.Serialize())
}

// respondWithSuccess 返回成功响应
func respondWithSuccess(w http.ResponseWriter, data interface{}) {
	resp := errs.GetCommonSuccessResp(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp.Serialize())
}

// respondWithJSON 返回JSON响应（用于直接返回proto响应）
func respondWithJSON(w http.ResponseWriter, httpCode int, data interface{}) {
	resp := errs.NewCommonRspV2(0, "success", data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(resp.Serialize())
}
