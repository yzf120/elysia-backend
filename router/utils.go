package router

import (
	"encoding/json"
	"github.com/yzf120/elysia-backend/service"
	"net/http"

	"github.com/yzf120/elysia-backend/errs"
)

var smsService *service.SMSService

// setResponseHeaders 设置响应头
func setResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

// writeErrorResponse 写入错误响应
func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	errResp := &errs.BaseResponse{
		Data:  nil,
		Error: errs.NewError(statusCode, message),
	}
	respBytes, _ := json.Marshal(errResp)
	w.WriteHeader(statusCode)
	w.Write(respBytes)
}

// writeSuccessResponse 写入成功响应
func writeSuccessResponse(w http.ResponseWriter, data interface{}) {
	resp := &errs.BaseResponse{
		Data: data,
	}
	respBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}
