package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/model/content_security"
	"github.com/yzf120/elysia-backend/rpc"
)

// 违规次数阈值，超过此值禁止使用AI对话
const ViolationBanThreshold = 5

// ModerationResult 审核结果
type ModerationResult struct {
	Passed     bool   // 是否通过审核
	Suggestion string // Block / Review / Pass
	Label      string // 违规标签
	Score      int32  // 违规置信度
	Message    string // 给用户的提示信息
}

// ContentModerationService 内容审核业务服务
type ContentModerationService struct {
	violationDAO dao.ViolationRecordDAO
}

// NewContentModerationService 创建内容审核服务
func NewContentModerationService() *ContentModerationService {
	return &ContentModerationService{
		violationDAO: dao.NewViolationRecordDAO(),
	}
}

// IsUserBanned 检查用户是否因违规被禁止使用AI对话
func (s *ContentModerationService) IsUserBanned(userId string) (bool, int64, error) {
	count, err := s.violationDAO.CountByUserId(userId)
	if err != nil {
		return false, 0, fmt.Errorf("查询违规记录失败: %w", err)
	}
	return count >= ViolationBanThreshold, count, nil
}

// ModerateText 同步文本审核
// 返回审核结果，如果违规同时写入违规记录表
func (s *ContentModerationService) ModerateText(ctx context.Context, userId, userRole, sessionId, senderType, content string) (*ModerationResult, error) {
	if strings.TrimSpace(content) == "" {
		return &ModerationResult{Passed: true, Suggestion: "Pass"}, nil
	}

	accessClient := rpc.GetAccessClient()

	textReq := &rpc.TextModerationRequest{
		Content: content,
		BizType: "default",
		UserId:  userId,
		DataId:  sessionId,
	}

	textResp, err := accessClient.TextModeration(ctx, textReq)
	if err != nil {
		// 审核服务调用失败时（含超时），默认拦截，因为敏感内容往往导致审核耗时更长
		log.Printf("[ContentModeration] 文本审核调用失败，默认拦截，userId: %s, sender: %s, err: %v", userId, senderType, err)
		return &ModerationResult{
			Passed:     false,
			Suggestion: "Block",
			Label:      "ServiceError",
			Message:    "内容审核服务暂时不可用，请稍后再试。",
		}, nil
	}

	result := &ModerationResult{
		Suggestion: textResp.Suggestion,
		Label:      textResp.Label,
		Score:      textResp.Score,
	}

	if textResp.Suggestion == "Pass" {
		result.Passed = true
		log.Printf("[ContentModeration] 文本审核通过，userId: %s, sender: %s, label: %s", userId, senderType, textResp.Label)
		return result, nil
	}

	// 违规：写入违规记录
	result.Passed = false
	result.Message = buildViolationMessage(textResp.Label, senderType)

	log.Printf("[ContentModeration] ⚠️ 文本审核未通过！userId: %s, sender: %s, sessionId: %s, suggestion: %s, label: %s, score: %d",
		userId, senderType, sessionId, textResp.Suggestion, textResp.Label, textResp.Score)

	// 截取违规内容（最多500字符）
	truncatedContent := content
	if len([]rune(truncatedContent)) > 500 {
		truncatedContent = string([]rune(truncatedContent)[:500]) + "..."
	}

	record := &content_security.ViolationRecord{
		UserId:      userId,
		UserRole:    userRole,
		SessionId:   sessionId,
		SenderType:  senderType,
		Content:     truncatedContent,
		Suggestion:  textResp.Suggestion,
		Label:       textResp.Label,
		Score:       textResp.Score,
		RequestId:   textResp.RequestId,
		ContentType: "text",
		CreateTime:  time.Now(),
	}

	if err := s.violationDAO.CreateRecord(record); err != nil {
		log.Printf("[ContentModeration] 写入违规记录失败，userId: %s, err: %v", userId, err)
	}

	return result, nil
}

// ModerateImage 同步图片审核
func (s *ContentModerationService) ModerateImage(ctx context.Context, userId, userRole, sessionId, senderType string, imgReq *rpc.ImageModerationRequest) (*ModerationResult, error) {
	accessClient := rpc.GetAccessClient()

	imgResp, err := accessClient.ImageModeration(ctx, imgReq)
	if err != nil {
		// 审核服务调用失败时（含超时），默认拦截
		log.Printf("[ContentModeration] 图片审核调用失败，默认拦截，userId: %s, sender: %s, err: %v", userId, senderType, err)
		return &ModerationResult{
			Passed:     false,
			Suggestion: "Block",
			Label:      "ServiceError",
			Message:    "内容审核服务暂时不可用，请稍后再试。",
		}, nil
	}

	result := &ModerationResult{
		Suggestion: imgResp.Suggestion,
		Label:      imgResp.Label,
		Score:      imgResp.Score,
	}

	if imgResp.Suggestion == "Pass" {
		result.Passed = true
		log.Printf("[ContentModeration] 图片审核通过，userId: %s, sender: %s", userId, senderType)
		return result, nil
	}

	result.Passed = false
	result.Message = buildViolationMessage(imgResp.Label, senderType)

	log.Printf("[ContentModeration] ⚠️ 图片审核未通过！userId: %s, sender: %s, sessionId: %s, suggestion: %s, label: %s, score: %d",
		userId, senderType, sessionId, imgResp.Suggestion, imgResp.Label, imgResp.Score)

	record := &content_security.ViolationRecord{
		UserId:      userId,
		UserRole:    userRole,
		SessionId:   sessionId,
		SenderType:  senderType,
		Content:     "[图片内容]",
		Suggestion:  imgResp.Suggestion,
		Label:       imgResp.Label,
		Score:       imgResp.Score,
		RequestId:   imgResp.RequestId,
		ContentType: "image",
		CreateTime:  time.Now(),
	}

	if err := s.violationDAO.CreateRecord(record); err != nil {
		log.Printf("[ContentModeration] 写入图片违规记录失败，userId: %s, err: %v", userId, err)
	}

	return result, nil
}

// buildViolationMessage 根据违规标签构建给用户的提示信息
func buildViolationMessage(label, senderType string) string {
	labelMap := map[string]string{
		"Porn":          "色情",
		"Abuse":         "谩骂",
		"Ad":            "广告",
		"Custom":        "自定义违规",
		"Contraband":    "违禁品",
		"Flood":         "灌水",
		"Meaningless":   "无意义",
		"UnsafeContent": "不安全内容",
	}

	labelCN, ok := labelMap[label]
	if !ok {
		labelCN = "违规"
	}

	if senderType == "user" {
		return fmt.Sprintf("您发送的消息包含%s内容，已被系统拦截。请规范用语，多次违规将被禁止使用AI对话功能。", labelCN)
	}
	return ""
}

// GetUserViolationCount 获取用户违规次数
func (s *ContentModerationService) GetUserViolationCount(userId string) (int64, error) {
	return s.violationDAO.CountByUserId(userId)
}
