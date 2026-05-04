package service

import (
	"fmt"
	"log"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/model/ai_model"
)

// AIModelConfigService AI模型配置管理服务
type AIModelConfigService struct {
	modelConfigDAO dao.AIModelConfigDAO
}

// NewAIModelConfigService 创建AI模型配置服务
func NewAIModelConfigService() *AIModelConfigService {
	return &AIModelConfigService{
		modelConfigDAO: dao.NewAIModelConfigDAO(),
	}
}

// InitDefaultModels 初始化默认模型配置（如果表中没有数据则插入默认记录）
func (s *AIModelConfigService) InitDefaultModels() {
	defaultModels := []*ai_model.AIModelConfig{
		{
			ModelId:     "doubao-seed-2-0-lite-260215",
			ModelName:   "Doubao-Seed-2.0-lite",
			Provider:    "doubao",
			Description: "豆包多模态模型，支持深度思考，适合快速响应场景",
			IsEnabled:   1,
		},
		{
			ModelId:     "qwen3-omni-flash",
			ModelName:   "Qwen3-Omni-Flash",
			Provider:    "qwen",
			Description: "通义千问全模态模型，Thinker–Talker 架构，支持深度思考",
			IsEnabled:   1,
		},
	}

	for _, m := range defaultModels {
		if err := s.modelConfigDAO.Upsert(m); err != nil {
			log.Printf("[AIModelConfig] 初始化默认模型失败: model_id=%s, err=%v", m.ModelId, err)
		}
	}
	log.Printf("[AIModelConfig] 默认模型配置初始化完成")
}

// ListAllModels 查询所有模型配置
func (s *AIModelConfigService) ListAllModels() ([]*ai_model.AIModelConfig, error) {
	return s.modelConfigDAO.ListAll()
}

// ListEnabledModels 查询所有启用的模型配置
func (s *AIModelConfigService) ListEnabledModels() ([]*ai_model.AIModelConfig, error) {
	return s.modelConfigDAO.ListEnabled()
}

// ListDisabledModelIds 查询所有禁用的模型ID列表
func (s *AIModelConfigService) ListDisabledModelIds() ([]string, error) {
	return s.modelConfigDAO.ListDisabledModelIds()
}

// ToggleModelStatus 切换模型启用/禁用状态
func (s *AIModelConfigService) ToggleModelStatus(modelId string, enabled bool) error {
	// 检查模型是否存在
	existing, err := s.modelConfigDAO.GetByModelId(modelId)
	if err != nil {
		return fmt.Errorf("查询模型配置失败: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("模型不存在: %s", modelId)
	}

	isEnabled := 0
	if enabled {
		isEnabled = 1
	}

	if err := s.modelConfigDAO.UpdateStatus(modelId, isEnabled); err != nil {
		return fmt.Errorf("更新模型状态失败: %v", err)
	}

	action := "禁用"
	if enabled {
		action = "启用"
	}
	log.Printf("[AIModelConfig] 模型状态已更新: model_id=%s, action=%s", modelId, action)
	return nil
}

// GetModelConfig 获取单个模型配置
func (s *AIModelConfigService) GetModelConfig(modelId string) (*ai_model.AIModelConfig, error) {
	config, err := s.modelConfigDAO.GetByModelId(modelId)
	if err != nil {
		return nil, fmt.Errorf("查询模型配置失败: %v", err)
	}
	if config == nil {
		return nil, fmt.Errorf("模型不存在: %s", modelId)
	}
	return config, nil
}

// IsModelEnabled 检查模型是否启用
func (s *AIModelConfigService) IsModelEnabled(modelId string) (bool, error) {
	config, err := s.modelConfigDAO.GetByModelId(modelId)
	if err != nil {
		return false, err
	}
	// 如果数据库中没有该模型的配置记录，默认启用
	if config == nil {
		return true, nil
	}
	return config.IsEnabled == 1, nil
}
