package test

import (
	"bytes"
	"encoding/json"
	"github.com/yzf120/elysia-backend/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	pb "github.com/yzf120/elysia-backend/proto/user"
)

// TestCreateUserHandler 测试创建用户接口
func TestCreateUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "成功创建用户",
			requestBody: pb.CreateUserRequest{
				PhoneNumber:     "13800138000",
				Nickname:        "测试用户",
				Avatar:          "https://example.com/avatar.jpg",
				Gender:          1,
				WxMiniAppOpenId: "wx_open_id_123",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp pb.CreateUserResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Errorf("解析响应失败: %v", err)
				}
				if resp.Code != 0 {
					t.Errorf("期望 code=0, 实际 code=%d", resp.Code)
				}
			},
		},
		{
			name:           "请求体格式错误",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusBadRequest {
					t.Errorf("期望状态码 %d, 实际 %d", http.StatusBadRequest, rec.Code)
				}
			},
		},
		{
			name: "缺少必填字段",
			requestBody: pb.CreateUserRequest{
				Nickname: "测试用户",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				// 根据业务逻辑验证响应
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 准备请求体
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("序列化请求体失败: %v", err)
				}
			}

			// 创建请求
			req := httptest.NewRequest("POST", "/api/user/create", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// 创建响应记录器
			rec := httptest.NewRecorder()

			// 调用处理器
			router.createUserHandler(rec, req)

			// 检查响应
			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}

// TestGetUserHandler 测试获取用户接口
func TestGetUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "通过用户ID查询",
			queryParams: map[string]string{
				"user_id": "user_123",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp pb.GetUserResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Errorf("解析响应失败: %v", err)
				}
			},
		},
		{
			name: "通过手机号查询",
			queryParams: map[string]string{
				"phone_number": "13800138000",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp pb.GetUserResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Errorf("解析响应失败: %v", err)
				}
			},
		},
		{
			name: "通过微信OpenID查询",
			queryParams: map[string]string{
				"wx_mini_app_open_id": "wx_open_id_123",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp pb.GetUserResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Errorf("解析响应失败: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 构建查询参数
			url := "/api/user/get?"
			for key, value := range tt.queryParams {
				url += key + "=" + value + "&"
			}

			// 创建请求
			req := httptest.NewRequest("GET", url, nil)

			// 创建响应记录器
			rec := httptest.NewRecorder()

			// 调用处理器
			router.getUserHandler(rec, req)

			// 检查响应
			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}

// TestUpdateUserHandler 测试更新用户接口
func TestUpdateUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "成功更新用户",
			requestBody: pb.UpdateUserRequest{
				UserId:   "user_123",
				Nickname: "新昵称",
				Avatar:   "https://example.com/new_avatar.jpg",
				Gender:   2,
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp pb.UpdateUserResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Errorf("解析响应失败: %v", err)
				}
			},
		},
		{
			name:           "请求体格式错误",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusBadRequest {
					t.Errorf("期望状态码 %d, 实际 %d", http.StatusBadRequest, rec.Code)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 准备请求体
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("序列化请求体失败: %v", err)
				}
			}

			// 创建请求
			req := httptest.NewRequest("POST", "/api/user/update", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// 创建响应记录器
			rec := httptest.NewRecorder()

			// 调用处理器
			router.updateUserHandler(rec, req)

			// 检查响应
			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}

// TestDeleteUserHandler 测试删除用户接口
func TestDeleteUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "成功删除用户",
			requestBody: pb.DeleteUserRequest{
				UserId: "user_123",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp pb.DeleteUserResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Errorf("解析响应失败: %v", err)
				}
			},
		},
		{
			name:           "请求体格式错误",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusBadRequest {
					t.Errorf("期望状态码 %d, 实际 %d", http.StatusBadRequest, rec.Code)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 准备请求体
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("序列化请求体失败: %v", err)
				}
			}

			// 创建请求
			req := httptest.NewRequest("POST", "/api/user/delete", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// 创建响应记录器
			rec := httptest.NewRecorder()

			// 调用处理器
			router.deleteUserHandler(rec, req)

			// 检查响应
			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}

// TestListUsersHandler 测试查询用户列表接口
func TestListUsersHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "查询所有用户",
			requestBody: pb.ListUsersRequest{
				Page:     1,
				PageSize: 10,
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp pb.ListUsersResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Errorf("解析响应失败: %v", err)
				}
			},
		},
		{
			name: "按性别筛选",
			requestBody: pb.ListUsersRequest{
				Page:     1,
				PageSize: 10,
				Gender:   1,
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp pb.ListUsersResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Errorf("解析响应失败: %v", err)
				}
			},
		},
		{
			name:           "请求体格式错误",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusBadRequest {
					t.Errorf("期望状态码 %d, 实际 %d", http.StatusBadRequest, rec.Code)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 准备请求体
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("序列化请求体失败: %v", err)
				}
			}

			// 创建请求
			req := httptest.NewRequest("POST", "/api/user/list", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// 创建响应记录器
			rec := httptest.NewRecorder()

			// 调用处理器
			router.listUsersHandler(rec, req)

			// 检查响应
			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}

// TestRegisterUser 测试路由注册
func TestRegisterUser(t *testing.T) {
	router := mux.NewRouter()
	router.registerUser(router)

	// 测试路由是否正确注册
	routes := []struct {
		path   string
		method string
	}{
		{"/api/user/create", "POST"},
		{"/api/user/get", "GET"},
		{"/api/user/update", "POST"},
		{"/api/user/delete", "POST"},
		{"/api/user/list", "POST"},
	}

	for _, route := range routes {
		req := httptest.NewRequest(route.method, route.path, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 如果路由不存在，会返回 404
		if rec.Code == http.StatusNotFound {
			t.Errorf("路由 %s %s 未正确注册", route.method, route.path)
		}
	}
}
