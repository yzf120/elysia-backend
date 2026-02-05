[//]: # (# CORS 问题修复文档)

[//]: # ()
[//]: # (## 🎯 问题描述)

[//]: # ()
[//]: # (前端（`http://localhost:3000`）向后端（`http://localhost:8001`）发送请求时，浏览器报错：)

[//]: # ()
[//]: # (```)

[//]: # (Access to XMLHttpRequest at 'http://localhost:8001/api/admin/auth/login-password' )

[//]: # (from origin 'http://localhost:3000' has been blocked by CORS policy: )

[//]: # (Response to preflight request doesn't pass access control check: )

[//]: # (No 'Access-Control-Allow-Origin' header is present on the requested resource.)

[//]: # (```)

[//]: # ()
[//]: # (## 🔍 问题原因)

[//]: # ()
[//]: # (1. **缺少 CORS 中间件**：后端没有配置跨域资源共享（CORS）中间件)

[//]: # (2. **OPTIONS 请求未处理**：浏览器发送的 preflight 请求（OPTIONS）没有得到正确响应)

[//]: # ()
[//]: # (## ✅ 解决方案)

[//]: # ()
[//]: # (### 1. 创建 CORS 中间件)

[//]: # ()
[//]: # (创建文件：`middleware/cors.go`)

[//]: # ()
[//]: # (```go)

[//]: # (package middleware)

[//]: # ()
[//]: # (import &#40;)

[//]: # (	"net/http")

[//]: # (&#41;)

[//]: # ()
[//]: # (// CORS 跨域资源共享中间件)

[//]: # (func CORS&#40;next http.Handler&#41; http.Handler {)

[//]: # (	return http.HandlerFunc&#40;func&#40;w http.ResponseWriter, r *http.Request&#41; {)

[//]: # (		// 设置 CORS 响应头（对所有请求都设置）)

[//]: # (		origin := r.Header.Get&#40;"Origin"&#41;)

[//]: # (		if origin == "" {)

[//]: # (			origin = "*")

[//]: # (		})

[//]: # (		)
[//]: # (		w.Header&#40;&#41;.Set&#40;"Access-Control-Allow-Origin", origin&#41;)

[//]: # (		w.Header&#40;&#41;.Set&#40;"Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH"&#41;)

[//]: # (		w.Header&#40;&#41;.Set&#40;"Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin"&#41;)

[//]: # (		w.Header&#40;&#41;.Set&#40;"Access-Control-Allow-Credentials", "true"&#41;)

[//]: # (		w.Header&#40;&#41;.Set&#40;"Access-Control-Max-Age", "86400"&#41; // 24小时)

[//]: # ()
[//]: # (		// 处理 preflight 请求（OPTIONS 请求）)

[//]: # (		if r.Method == "OPTIONS" {)

[//]: # (			w.WriteHeader&#40;http.StatusOK&#41;)

[//]: # (			return)

[//]: # (		})

[//]: # ()
[//]: # (		// 继续处理其他请求)

[//]: # (		next.ServeHTTP&#40;w, r&#41;)

[//]: # (	}&#41;)

[//]: # (})

[//]: # (```)

[//]: # ()
[//]: # (### 2. 在 main.go 中应用 CORS 中间件)

[//]: # ()
[//]: # (修改 `main.go`：)

[//]: # ()
[//]: # (```go)

[//]: # (package main)

[//]: # ()
[//]: # (import &#40;)

[//]: # (	"github.com/gorilla/mux")

[//]: # (	"github.com/joho/godotenv")

[//]: # (	"github.com/yzf120/elysia-backend/client")

[//]: # (	"github.com/yzf120/elysia-backend/dao")

[//]: # (	"github.com/yzf120/elysia-backend/middleware"  // 添加这行)

[//]: # (	"github.com/yzf120/elysia-backend/router")

[//]: # (	"log")

[//]: # (	"trpc.group/trpc-go/trpc-go")

[//]: # (	thttp "trpc.group/trpc-go/trpc-go/http")

[//]: # (&#41;)

[//]: # ()
[//]: # (func main&#40;&#41; {)

[//]: # (	// ... 数据库和Redis初始化代码 ...)

[//]: # ()
[//]: # (	r := mux.NewRouter&#40;&#41;)

[//]: # (	)
[//]: # (	// 初始化路由器)

[//]: # (	router.Init&#40;&#41;)

[//]: # (	router.RegisterRouter&#40;r&#41;)

[//]: # ()
[//]: # (	// 创建带 CORS 的 handler（包装整个路由器）)

[//]: # (	corsHandler := middleware.CORS&#40;r&#41;)

[//]: # ()
[//]: # (	// 创建trpc服务器)

[//]: # (	s := trpc.NewServer&#40;&#41;)

[//]: # ()
[//]: # (	// 注册http服务（使用带 CORS 的 handler）)

[//]: # (	thttp.RegisterNoProtocolServiceMux&#40;s.Service&#40;"trpc.elysia.backend.http"&#41;, corsHandler&#41;)

[//]: # ()
[//]: # (	// 启动服务器)

[//]: # (	if err := s.Serve&#40;&#41;; err != nil {)

[//]: # (		log.Fatalf&#40;"服务器启动失败: %v", err&#41;)

[//]: # (	})

[//]: # (})

[//]: # (```)

[//]: # ()
[//]: # (## 🧪 测试验证)

[//]: # ()
[//]: # (### 1. 测试 OPTIONS 请求（Preflight）)

[//]: # ()
[//]: # (```bash)

[//]: # (curl -X OPTIONS http://localhost:8001/api/admin/auth/login-password \)

[//]: # (  -H "Origin: http://localhost:3000" \)

[//]: # (  -H "Access-Control-Request-Method: POST" \)

[//]: # (  -H "Access-Control-Request-Headers: Content-Type, Authorization" \)

[//]: # (  -i)

[//]: # (```)

[//]: # ()
[//]: # (**预期结果**：)

[//]: # (```)

[//]: # (HTTP/1.1 200 OK)

[//]: # (Access-Control-Allow-Origin: http://localhost:3000)

[//]: # (Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, PATCH)

[//]: # (Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With, Accept, Origin)

[//]: # (Access-Control-Allow-Credentials: true)

[//]: # (Access-Control-Max-Age: 86400)

[//]: # (```)

[//]: # ()
[//]: # (### 2. 测试实际 POST 请求)

[//]: # ()
[//]: # (```bash)

[//]: # (curl -X POST http://localhost:8001/api/admin/auth/login-password \)

[//]: # (  -H "Content-Type: application/json" \)

[//]: # (  -H "Origin: http://localhost:3000" \)

[//]: # (  -d '{"phone_number":"13800138000","password":"123456Admin"}' \)

[//]: # (  -i)

[//]: # (```)

[//]: # ()
[//]: # (**预期结果**：)

[//]: # (```)

[//]: # (HTTP/1.1 200 OK)

[//]: # (Access-Control-Allow-Origin: *)

[//]: # (Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, PATCH)

[//]: # (...)

[//]: # (```)

[//]: # ()
[//]: # (### 3. 在浏览器中测试)

[//]: # ()
[//]: # (1. 确保后端服务运行在 `http://localhost:8001`)

[//]: # (2. 确保前端服务运行在 `http://localhost:3000`)

[//]: # (3. 打开浏览器访问 `http://localhost:3000`)

[//]: # (4. 打开开发者工具（F12）→ Network 标签)

[//]: # (5. 尝试登录)

[//]: # (6. 检查网络请求：)

[//]: # (   - ✅ 应该看到 OPTIONS 请求返回 200)

[//]: # (   - ✅ 应该看到 POST 请求返回 200)

[//]: # (   - ✅ 不应该有 CORS 错误)

[//]: # ()
[//]: # (## 📋 修改文件清单)

[//]: # ()
[//]: # (1. ✅ **新建**：`middleware/cors.go` - CORS 中间件)

[//]: # (2. ✅ **修改**：`main.go` - 应用 CORS 中间件)

[//]: # ()
[//]: # (## 🎉 修复效果)

[//]: # ()
[//]: # (- ✅ 前端可以正常向后端发送请求)

[//]: # (- ✅ OPTIONS preflight 请求正常响应)

[//]: # (- ✅ POST/GET 等请求正常响应)

[//]: # (- ✅ 浏览器不再报 CORS 错误)

[//]: # (- ✅ 支持跨域携带凭证（cookies）)

[//]: # ()
[//]: # (## 💡 技术要点)

[//]: # ()
[//]: # (1. **CORS 中间件位置**：必须在所有路由处理之前应用)

[//]: # (2. **OPTIONS 请求处理**：必须返回 200 状态码和正确的 CORS 头)

[//]: # (3. **Origin 处理**：动态设置 `Access-Control-Allow-Origin` 以支持凭证传递)

[//]: # (4. **Handler 包装**：使用 `middleware.CORS&#40;r&#41;` 包装整个路由器)

[//]: # ()
[//]: # (## 🔐 安全建议)

[//]: # ()
[//]: # (生产环境中，建议：)

[//]: # (1. 不要使用 `Access-Control-Allow-Origin: *`)

[//]: # (2. 明确指定允许的域名列表)

[//]: # (3. 根据请求的 Origin 动态返回允许的域名)

[//]: # (4. 限制允许的 HTTP 方法和请求头)

[//]: # ()
[//]: # (## 📚 相关资源)

[//]: # ()
[//]: # (- [MDN - CORS]&#40;https://developer.mozilla.org/zh-CN/docs/Web/HTTP/CORS&#41;)

[//]: # (- [Gorilla Mux 文档]&#40;https://github.com/gorilla/mux&#41;)

[//]: # (- [Go HTTP 中间件模式]&#40;https://www.alexedwards.net/blog/making-and-using-middleware&#41;)
