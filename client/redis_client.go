package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient Redis客户端
type RedisClient struct {
	Client *redis.Client
	ctx    context.Context
}

var defaultRedisClient *RedisClient

// InitRedisClient 初始化Redis客户端
func InitRedisClient() error {
	host := getRedisEnv("REDIS_HOST", "localhost")
	port := getRedisEnv("REDIS_PORT", "6379")
	password := getRedisEnv("REDIS_PASSWORD", "")
	dbStr := getRedisEnv("REDIS_DB", "0")
	
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0
	}

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

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis连接失败: %v", err)
	}

	defaultRedisClient = &RedisClient{
		Client: client,
		ctx:    ctx,
	}

	log.Printf("Redis客户端初始化成功: %s:%s (DB: %d)", host, port, db)
	return nil
}

// GetRedisClient 获取默认Redis客户端
func GetRedisClient() *RedisClient {
	return defaultRedisClient
}

// Close 关闭Redis连接
func (c *RedisClient) Close() error {
	if c.Client != nil {
		return c.Client.Close()
	}
	return nil
}

// Set 设置键值对
func (c *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(c.ctx, key, value, expiration).Err()
}

// Get 获取值
func (c *RedisClient) Get(key string) (string, error) {
	return c.Client.Get(c.ctx, key).Result()
}

// Del 删除键
func (c *RedisClient) Del(keys ...string) error {
	return c.Client.Del(c.ctx, keys...).Err()
}

// Exists 检查键是否存在
func (c *RedisClient) Exists(keys ...string) (int64, error) {
	return c.Client.Exists(c.ctx, keys...).Result()
}

// Expire 设置过期时间
func (c *RedisClient) Expire(key string, expiration time.Duration) error {
	return c.Client.Expire(c.ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func (c *RedisClient) TTL(key string) (time.Duration, error) {
	return c.Client.TTL(c.ctx, key).Result()
}

// getRedisEnv 获取环境变量，如果不存在则返回默认值
func getRedisEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
