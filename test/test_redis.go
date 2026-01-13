package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到.env文件，使用系统环境变量")
	}

	// 获取 Redis 配置
	host := getRedisEnv("REDIS_HOST", "localhost")
	port := getRedisEnv("REDIS_PORT", "6379")
	password := getRedisEnv("REDIS_PASSWORD", "")
	dbStr := getRedisEnv("REDIS_DB", "0")

	fmt.Printf("尝试连接 Redis:\n")
	fmt.Printf("  Host: %s\n", host)
	fmt.Printf("  Port: %s\n", port)
	if password == "" {
		fmt.Printf("  Password: [空]\n")
	} else {
		fmt.Printf("  Password: [已设置]\n")
	}
	fmt.Printf("  DB: %s\n", dbStr)

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0
	}

	// 创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Password:     password,
		DB:           db,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx := context.Background()

	fmt.Printf("\n尝试 Ping Redis...\n")

	// 尝试连接
	start := time.Now()
	err = client.Ping(ctx).Err()
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("❌ Redis 连接失败 (耗时: %v):\n", duration)
		fmt.Printf("   错误信息: %v\n", err)
		fmt.Printf("\n可能的原因:\n")
		fmt.Printf("1. Redis 服务未启动\n")
		fmt.Printf("2. 防火墙阻止了端口 %s\n", port)
		fmt.Printf("3. Redis 配置需要密码但未设置\n")
		fmt.Printf("4. Redis 运行在不同的主机或端口\n")

		// 尝试使用更简单的配置
		fmt.Printf("\n尝试使用默认配置连接 (localhost:6379, 无密码)...\n")
		defaultClient := redis.NewClient(&redis.Options{
			Addr:        "localhost:6379",
			DB:          0,
			DialTimeout: 5 * time.Second,
		})

		err = defaultClient.Ping(ctx).Err()
		if err != nil {
			fmt.Printf("❌ 默认配置也失败: %v\n", err)
			fmt.Printf("\n建议的操作:\n")
			fmt.Printf("1. 检查 Redis 是否安装并运行: redis-cli ping\n")
			fmt.Printf("2. 如果是 Docker 环境，检查容器是否运行: docker ps | grep redis\n")
			fmt.Printf("3. 检查防火墙设置: sudo ufw status (Linux)\n")
		} else {
			fmt.Printf("✅ 默认配置连接成功！可能是 .env 配置有问题\n")
			defaultClient.Close()
		}
	} else {
		fmt.Printf("✅ Redis 连接成功 (耗时: %v)\n", duration)
		fmt.Printf("✅ Redis 客户端已准备好使用\n")
	}

	client.Close()
}

func getRedisEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
