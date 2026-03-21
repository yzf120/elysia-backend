package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yzf120/elysia-backend/authen"
	"github.com/yzf120/elysia-backend/rpc"
	agent_session "github.com/yzf120/elysia-session/proto/agent_session"
)

// registerFavorite 注册收藏相关路由
func registerFavorite(router *mux.Router) {
	// 收藏会话
	router.HandleFunc("/ai/favorite", favoriteSessionHandler).Methods("POST", "OPTIONS")
	// 取消收藏会话
	router.HandleFunc("/ai/unfavorite", unfavoriteSessionHandler).Methods("POST", "OPTIONS")
	// 查询收藏列表
	router.HandleFunc("/ai/favorites", listFavoritesHandler).Methods("GET", "OPTIONS")
	// 检查会话是否已收藏
	router.HandleFunc("/ai/favorite/check", checkFavoriteHandler).Methods("GET", "OPTIONS")
}

// FavoriteRequest 收藏/取消收藏请求
type FavoriteRequest struct {
	SessionID string `json:"session_id"`
}

// favoriteSessionHandler 收藏会话
func favoriteSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()
	_, ok := authen.GetRoleIDFromContext(ctx)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 401, "message": "未授权"},
		})
		return
	}

	var req FavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 400, "message": "请求参数错误"},
		})
		return
	}

	if req.SessionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 400, "message": "session_id 不能为空"},
		})
		return
	}

	// 通过 tRPC 调用 elysia-session 的收藏接口
	sessionClient := rpc.GetSessionClient()
	resp, err := sessionClient.FavoriteSession(ctx, &agent_session.FavoriteSessionRequest{
		SessionId: req.SessionID,
	})
	if err != nil {
		log.Printf("[favorite] 收藏会话失败, sessionID: %s, err: %v", req.SessionID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 500, "message": "收藏失败"},
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{"code": resp.Code, "message": resp.Message},
	})
}

// unfavoriteSessionHandler 取消收藏会话
func unfavoriteSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()
	_, ok := authen.GetRoleIDFromContext(ctx)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 401, "message": "未授权"},
		})
		return
	}

	var req FavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 400, "message": "请求参数错误"},
		})
		return
	}

	if req.SessionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 400, "message": "session_id 不能为空"},
		})
		return
	}

	sessionClient := rpc.GetSessionClient()
	resp, err := sessionClient.UnfavoriteSession(ctx, &agent_session.UnfavoriteSessionRequest{
		SessionId: req.SessionID,
	})
	if err != nil {
		log.Printf("[favorite] 取消收藏失败, sessionID: %s, err: %v", req.SessionID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 500, "message": "取消收藏失败"},
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{"code": resp.Code, "message": resp.Message},
	})
}

// listFavoritesHandler 查询收藏列表
func listFavoritesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()
	userID, ok := authen.GetRoleIDFromContext(ctx)
	if !ok || userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 401, "message": "未授权"},
		})
		return
	}

	// 解析分页参数
	page := int32(1)
	pageSize := int32(20)
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = int32(v)
		}
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = int32(v)
		}
	}

	sessionClient := rpc.GetSessionClient()
	resp, err := sessionClient.ListFavoriteSessions(ctx, &agent_session.ListFavoriteSessionsRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		log.Printf("[favorite] 查询收藏列表失败, userID: %s, err: %v", userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 500, "message": "查询收藏列表失败"},
		})
		return
	}

	// 转换为前端需要的格式
	items := make([]map[string]interface{}, 0)
	if resp.Sessions != nil {
		for _, s := range resp.Sessions {
			items = append(items, map[string]interface{}{
				"session_id":    s.SessionId,
				"session_title": s.SessionTitle,
				"problem_id":    s.ProblemId,
				"is_favorited":  s.IsFavorited,
				"create_time":   s.CreateTime,
				"update_time":   s.UpdateTime,
			})
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{"code": 0},
		"data": map[string]interface{}{
			"favorites": items,
			"total":     resp.Total,
			"page":      resp.Page,
			"page_size": resp.PageSize,
		},
	})
}

// checkFavoriteHandler 检查会话是否已收藏
func checkFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()
	_, ok := authen.GetRoleIDFromContext(ctx)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 401, "message": "未授权"},
		})
		return
	}

	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 400, "message": "session_id 不能为空"},
		})
		return
	}

	sessionClient := rpc.GetSessionClient()
	resp, err := sessionClient.CheckFavorite(ctx, &agent_session.CheckFavoriteRequest{
		SessionId: sessionID,
	})
	if err != nil {
		log.Printf("[favorite] 检查收藏状态失败, sessionID: %s, err: %v", sessionID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": 500, "message": "检查收藏状态失败"},
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{"code": 0},
		"data": map[string]interface{}{
			"is_favorited": resp.IsFavorited,
		},
	})
}