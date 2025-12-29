package main

import (
	"github.com/yzf120/elysia-backend/client"
	"github.com/yzf120/elysia-backend/proto/helloworld"
	pb "github.com/yzf120/elysia-backend/proto/user"
	"github.com/yzf120/elysia-backend/router"
	"github.com/yzf120/elysia-backend/service_impl"
	"log"

	"github.com/joho/godotenv"
	"github.com/yzf120/elysia-backend/dao"
	"trpc.group/trpc-go/trpc-go"
)

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到.env文件，使用系统环境变量")
	}

	// 初始化MySQL客户端
	err = client.InitMySQLClient()
	if err != nil {
		log.Fatalf("MySQL客户端初始化失败: %v", err)
	}
	defer client.GetMySQLClient().Close()

	// 初始化数据库（保持兼容性）
	err = dao.InitDB()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer dao.CloseDB()

	// 创建trpc服务器
	s := trpc.NewServer()

	// 注册helloworld服务
	helloworld.RegisterGreeterService(s.Service("trpc.elysia.backend.helloworld"),
		&service_impl.HelloWorldImpl{})

	// 注册用户服务
	pb.RegisterUserServiceService(s.Service("trpc.elysia.backend.user"),
		service_impl.NewUserServiceImpl())

	router.Init()

	// 启动服务器
	if err := s.Serve(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
