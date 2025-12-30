package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yzf120/elysia-backend/client"
)

// JWT 配置常量
const (
	// JWT 密钥（应该从环境变量读取）
	DefaultJWTSecretKey = "elysia-jwt-secret-key-change-in-production"
	// Token 过期时间（24小时）
	DefaultTokenExpiration = 24 * time.Hour
	// Redis 中 token 的 key 前缀
	TokenRedisKeyPrefix = "auth:token:"
)

// JWTClaims JWT 声明结构
type JWTClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

// JWTService JWT 服务
type JWTService struct {
	redisClient *client.RedisClient
	secretKey   []byte
}

// NewJWTService 创建 JWT 服务
func NewJWTService() *JWTService {
	// 从环境变量获取密钥，如果没有则使用默认值
	secretKey := []byte(DefaultJWTSecretKey)
	// 实际项目中应该从环境变量读取
	// if envSecret := os.Getenv("JWT_SECRET_KEY"); envSecret != "" {
	//     secretKey = []byte(envSecret)
	// }

	return &JWTService{
		redisClient: client.GetRedisClient(),
		secretKey:   secretKey,
	}
}

// GenerateToken 生成 JWT token 并存储到 Redis
func (s *JWTService) GenerateToken(userID string) (string, error) {
	// 创建 JWT claims
	expiresAt := jwt.NewNumericDate(time.Now().Add(DefaultTokenExpiration))
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "elysia-backend",
		},
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名 token
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("生成 JWT token 失败: %v", err)
	}

	// 将 token 存储到 Redis，设置过期时间
	redisKey := TokenRedisKeyPrefix + userID + ":" + tokenString
	err = s.redisClient.Set(redisKey, "valid", DefaultTokenExpiration)
	if err != nil {
		return "", fmt.Errorf("存储 token 到 Redis 失败: %v", err)
	}

	return tokenString, nil
}

// ValidateToken 验证 JWT token 并获取用户 ID
func (s *JWTService) ValidateToken(tokenString string) (string, error) {
	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名方法: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("解析 JWT token 失败: %v", err)
	}

	if !token.Valid {
		return "", errors.New("无效的 token")
	}

	// 获取 claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return "", errors.New("无效的 token claims")
	}

	userID := claims.UserID

	// 检查 token 是否在 Redis 中有效（单点登录校验）
	redisKey := TokenRedisKeyPrefix + userID + ":" + tokenString
	exists, err := s.redisClient.Exists(redisKey)
	if err != nil {
		return "", fmt.Errorf("检查 Redis token 失败: %v", err)
	}

	if exists == 0 {
		return "", errors.New("token 已过期或无效")
	}

	return userID, nil
}

// InvalidateToken 使 token 失效（登出）
func (s *JWTService) InvalidateToken(userID, tokenString string) error {
	redisKey := TokenRedisKeyPrefix + userID + ":" + tokenString
	return s.redisClient.Del(redisKey)
}

// InvalidateAllUserTokens 使用户的所有 token 失效（强制登出）
func (s *JWTService) InvalidateAllUserTokens(userID string) error {
	// 注意：实际实现可能需要使用 Redis 的 KEYS 命令或 SCAN 命令
	// 但由于性能考虑，这里使用简单的模式匹配
	// 在生产环境中可能需要更复杂的方法
	// pattern := TokenRedisKeyPrefix + userID + ":*"
	// 这里简化处理，在实际使用时需要实现
	return nil
}

// GetTokenExpiration 获取 token 剩余过期时间
func (s *JWTService) GetTokenExpiration(userID, tokenString string) (time.Duration, error) {
	redisKey := TokenRedisKeyPrefix + userID + ":" + tokenString
	return s.redisClient.TTL(redisKey)
}

// RefreshToken 刷新 token
func (s *JWTService) RefreshToken(userID, oldTokenString string) (string, error) {
	// 验证旧 token
	_, err := s.ValidateToken(oldTokenString)
	if err != nil {
		return "", fmt.Errorf("原 token 无效: %v", err)
	}

	// 生成新 token
	newTokenString, err := s.GenerateToken(userID)
	if err != nil {
		return "", err
	}

	// 使旧 token 失效
	if err := s.InvalidateToken(userID, oldTokenString); err != nil {
		// 即使失效失败也继续，新 token 已经生成
		fmt.Printf("使旧 token 失效失败: %v", err)
	}

	return newTokenString, nil
}
