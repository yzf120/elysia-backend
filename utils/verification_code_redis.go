package utils

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yzf120/elysia-backend/client"
)

// VerificationCodeService 验证码服务（基于Redis）
type VerificationCodeService struct {
	redisClient *client.RedisClient
}

// NewVerificationCodeService 创建验证码服务
func NewVerificationCodeService() *VerificationCodeService {
	return &VerificationCodeService{
		redisClient: client.GetRedisClient(),
	}
}

// SaveVerificationCode 保存验证码到Redis
// key格式: sms:code:{codeType}:{phoneNumber}
// value: 验证码
// expiration: 过期时间（默认5分钟）
func (s *VerificationCodeService) SaveVerificationCode(phoneNumber, code, codeType string, expiration time.Duration) error {
	if expiration == 0 {
		expiration = 5 * time.Minute // 默认5分钟
	}

	key := fmt.Sprintf("sms:code:%s:%s", codeType, phoneNumber)
	return s.redisClient.Set(key, code, expiration)
}

// GetVerificationCode 获取验证码
func (s *VerificationCodeService) GetVerificationCode(phoneNumber, codeType string) (string, error) {
	key := fmt.Sprintf("sms:code:%s:%s", codeType, phoneNumber)
	code, err := s.redisClient.Get(key)
	if err == redis.Nil {
		return "", fmt.Errorf("验证码不存在或已过期")
	}
	return code, err
}

// VerifyCode 验证验证码
func (s *VerificationCodeService) VerifyCode(phoneNumber, code, codeType string) error {
	// 获取存储的验证码
	storedCode, err := s.GetVerificationCode(phoneNumber, codeType)
	if err != nil {
		return err
	}

	// 验证验证码
	if storedCode != code {
		return fmt.Errorf("验证码错误")
	}

	// 验证成功后删除验证码（一次性使用）
	key := fmt.Sprintf("sms:code:%s:%s", codeType, phoneNumber)
	if err := s.redisClient.Del(key); err != nil {
		// 删除失败不影响验证结果，只记录日志
		fmt.Printf("删除验证码失败: %v\n", err)
	}
	return nil
}

// DeleteVerificationCode 删除验证码
func (s *VerificationCodeService) DeleteVerificationCode(phoneNumber, codeType string) error {
	key := fmt.Sprintf("sms:code:%s:%s", codeType, phoneNumber)
	return s.redisClient.Del(key)
}

// GetTTL 获取验证码剩余有效时间
func (s *VerificationCodeService) GetTTL(phoneNumber, codeType string) (time.Duration, error) {
	key := fmt.Sprintf("sms:code:%s:%s", codeType, phoneNumber)
	return s.redisClient.TTL(key)
}

// CheckSendFrequency 检查发送频率（防止频繁发送）
// 返回true表示可以发送，false表示需要等待
func (s *VerificationCodeService) CheckSendFrequency(phoneNumber string, interval time.Duration) (bool, time.Duration, error) {
	key := fmt.Sprintf("sms:frequency:%s", phoneNumber)

	// 检查是否存在
	exists, err := s.redisClient.Exists(key)
	if err != nil {
		return false, 0, err
	}

	if exists > 0 {
		// 获取剩余时间
		ttl, err := s.redisClient.TTL(key)
		if err != nil {
			return false, 0, err
		}
		return false, ttl, nil
	}

	// 设置频率限制
	if err := s.redisClient.Set(key, "1", interval); err != nil {
		return false, 0, err
	}

	return true, 0, nil
}
