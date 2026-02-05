package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/yzf120/elysia-backend/client"
	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/middleware"
	"github.com/yzf120/elysia-backend/router"
	"log"
	"trpc.group/trpc-go/trpc-go"
	thttp "trpc.group/trpc-go/trpc-go/http"
)

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到.env文件，使用系统环境变量")
	}

	// 初始化数据库
	err = dao.InitDB()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer dao.CloseDB()

	// 初始化Redis
	err = client.InitRedisClient()
	if err != nil {
		log.Fatalf("Redis初始化失败: %v", err)
	}
	defer client.GetRedisClient().Close()

	r := mux.NewRouter()
	
	// 初始化路由器（必须在RegisterRouter之前调用，以初始化Service实例）
	router.Init()
	
	router.RegisterRouter(r)

	// 创建带 CORS 的 handler（包装整个路由器）
	corsHandler := middleware.CORS(r)

	// 创建trpc服务器
	s := trpc.NewServer()

	// 注册http服务（使用带 CORS 的 handler）
	thttp.RegisterNoProtocolServiceMux(s.Service("trpc.elysia.backend.http"), corsHandler)

	// 启动服务器
	if err := s.Serve(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
