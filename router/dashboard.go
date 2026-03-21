package router

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	llmpb "github.com/yzf120/elysia-llm-tool/proto/llm"
	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/rpc"
)

// RegisterDashboardRoutes 注册Dashboard监控路由（管理员接口）
func RegisterDashboardRoutes(protectedRouter *mux.Router) {
	// 查询模型推理用量
	protectedRouter.HandleFunc("/admin/dashboard/usage", getInferenceUsageHandler).Methods("GET")
	// 获取Dashboard统计概览
	protectedRouter.HandleFunc("/admin/dashboard/stats", getDashboardStatsHandler).Methods("GET")
}

// getInferenceUsageHandler 查询模型推理用量
// GET /admin/dashboard/usage?interval=Day&start_time=2025-01-01&end_time=2025-01-07&endpoint=ep-xxx
func getInferenceUsageHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	interval := r.URL.Query().Get("interval")
	if interval == "" {
		interval = "Day"
	}
	startTime := r.URL.Query().Get("start_time")
	endTime := r.URL.Query().Get("end_time")

	// 默认查询近7天
	if startTime == "" || endTime == "" {
		now := time.Now()
		endTime = now.Format("2006-01-02")
		startTime = now.AddDate(0, 0, -7).Format("2006-01-02")
	}

	// 构建 RPC 请求
	req := &llmpb.GetInferenceUsageRequest{
		QueryInterval: interval,
		StartTime:     startTime,
		EndTime:       endTime,
	}

	// 可选：按接入点过滤
	endpoint := r.URL.Query().Get("endpoint")
	if endpoint != "" {
		req.Filters = append(req.Filters, &llmpb.UsageFilter{
			Key:    "ModelEndpoint",
			Values: []string{endpoint},
		})
	}

	// 可选：按模型名称过滤
	modelName := r.URL.Query().Get("model_name")
	if modelName != "" {
		req.Filters = append(req.Filters, &llmpb.UsageFilter{
			Key:       "ModelName",
			ValueLike: modelName,
		})
	}

	// 调用 llm-tool RPC
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := rpc.GetLLMToolClient().GetProxy().GetInferenceUsage(ctx, req)
	if err != nil {
		log.Printf("[Dashboard] 查询推理用量失败: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "查询推理用量失败: "+err.Error())
		return
	}

	// 将 proto 响应转换为前端友好的 JSON 格式
	result := convertUsageResponse(resp)
	writeSuccessResponse(w, result)
}

// getDashboardStatsHandler 获取Dashboard统计概览
// GET /admin/dashboard/stats
func getDashboardStatsHandler(w http.ResponseWriter, r *http.Request) {
	setResponseHeaders(w)

	stats := make(map[string]interface{})

	// 统计学生总数
	studentDAO := dao.NewStudentDAO()
	studentCount, err := studentDAO.CountStudents("", nil)
	if err != nil {
		log.Printf("[Dashboard] 统计学生数失败: %v", err)
		studentCount = 0
	}

	// 统计教师总数
	teacherDAO := dao.NewTeacherDAO()
	teacherCount, err := teacherDAO.CountTeachers("", nil)
	if err != nil {
		log.Printf("[Dashboard] 统计教师数失败: %v", err)
		teacherCount = 0
	}

	// 统计待审核数（approval_status = 0 表示待审核）
	teacherApprovalDAO := dao.NewTeacherApprovalDAO()
	pendingCount, err := teacherApprovalDAO.CountApprovals("approval_status = ?", []interface{}{0})
	if err != nil {
		log.Printf("[Dashboard] 统计待审核数失败: %v", err)
		pendingCount = 0
	}

	stats["total_students"] = studentCount
	stats["total_teachers"] = teacherCount
	stats["total_users"] = studentCount + teacherCount
	stats["pending_audits"] = pendingCount

	writeSuccessResponse(w, stats)
}

// usageResult 前端友好的用量数据结构
type usageResult struct {
	// 按天/小时聚合的数据点，用于图表展示
	DataPoints []usageDataPoint `json:"data_points"`
	// 汇总数据
	Summary usageSummary `json:"summary"`
}

type usageDataPoint struct {
	Date         string `json:"date"`
	Hour         string `json:"hour,omitempty"`
	InputTokens  int64  `json:"input_tokens"`
	OutputTokens int64  `json:"output_tokens"`
	TotalTokens  int64  `json:"total_tokens"`
	ReqCount     int64  `json:"req_count"`
}

type usageSummary struct {
	TotalInputTokens  int64 `json:"total_input_tokens"`
	TotalOutputTokens int64 `json:"total_output_tokens"`
	TotalTokens       int64 `json:"total_tokens"`
	TotalRequests     int64 `json:"total_requests"`
}

// convertUsageResponse 将 proto 响应转换为前端友好格式
func convertUsageResponse(resp *llmpb.GetInferenceUsageResponse) *usageResult {
	result := &usageResult{
		DataPoints: make([]usageDataPoint, 0),
	}

	if resp == nil || len(resp.Data) == 0 {
		return result
	}

	// 建立字段名到索引的映射
	fieldIndex := make(map[string]int)
	for i, f := range resp.Fields {
		fieldIndex[f.Name] = i
	}

	for _, row := range resp.Data {
		dp := usageDataPoint{}

		if idx, ok := fieldIndex["Day"]; ok && idx < len(row.Values) {
			dp.Date = row.Values[idx]
		}
		if idx, ok := fieldIndex["Hour"]; ok && idx < len(row.Values) {
			dp.Hour = row.Values[idx]
		}
		if idx, ok := fieldIndex["InputTokens"]; ok && idx < len(row.Values) {
			dp.InputTokens = parseInt64(row.Values[idx])
		}
		if idx, ok := fieldIndex["OutputTokens"]; ok && idx < len(row.Values) {
			dp.OutputTokens = parseInt64(row.Values[idx])
		}
		if idx, ok := fieldIndex["TotalTokens"]; ok && idx < len(row.Values) {
			dp.TotalTokens = parseInt64(row.Values[idx])
		}
		if idx, ok := fieldIndex["ReqCnt"]; ok && idx < len(row.Values) {
			dp.ReqCount = parseInt64(row.Values[idx])
		}

		result.DataPoints = append(result.DataPoints, dp)

		// 累加汇总
		result.Summary.TotalInputTokens += dp.InputTokens
		result.Summary.TotalOutputTokens += dp.OutputTokens
		result.Summary.TotalTokens += dp.TotalTokens
		result.Summary.TotalRequests += dp.ReqCount
	}

	return result
}

// parseInt64 安全地将字符串转为 int64
func parseInt64(s string) int64 {
	var n int64
	json.Unmarshal([]byte(s), &n)
	return n
}
