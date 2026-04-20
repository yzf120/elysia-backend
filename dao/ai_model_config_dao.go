package dao

import (
	"fmt"

	"github.com/yzf120/elysia-backend/model/ai_model"
	"gorm.io/gorm"
)

// AIModelConfigDAO AI模型配置数据访问接口
type AIModelConfigDAO interface {
	// GetByModelId 根据模型ID查询配置
	GetByModelId(modelId string) (*ai_model.AIModelConfig, error)
	// ListAll 查询所有模型配置
	ListAll() ([]*ai_model.AIModelConfig, error)
	// ListEnabled 查询所有启用的模型配置
	ListEnabled() ([]*ai_model.AIModelConfig, error)
	// ListDisabledModelIds 查询所有禁用的模型ID列表
	ListDisabledModelIds() ([]string, error)
	// UpdateStatus 更新模型启用状态
	UpdateStatus(modelId string, isEnabled int) error
	// Upsert 插入或更新模型配置
	Upsert(config *ai_model.AIModelConfig) error
}

// aiModelConfigDAOImpl AI模型配置数据访问实现
type aiModelConfigDAOImpl struct {
	db *gorm.DB
}

// NewAIModelConfigDAO 创建AI模型配置DAO实例
func NewAIModelConfigDAO() AIModelConfigDAO {
	return &aiModelConfigDAOImpl{
		db: GetDB(),
	}
}

// GetByModelId 根据模型ID查询配置
func (d *aiModelConfigDAOImpl) GetByModelId(modelId string) (*ai_model.AIModelConfig, error) {
	var config ai_model.AIModelConfig
	err := d.db.Where("model_id = ?", modelId).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询模型配置失败: %v", err)
	}
	return &config, nil
}

// ListAll 查询所有模型配置
func (d *aiModelConfigDAOImpl) ListAll() ([]*ai_model.AIModelConfig, error) {
	var configs []*ai_model.AIModelConfig
	err := d.db.Order("id ASC").Find(&configs).Error
	if err != nil {
		return nil, fmt.Errorf("查询模型配置列表失败: %v", err)
	}
	return configs, nil
}

// ListEnabled 查询所有启用的模型配置
func (d *aiModelConfigDAOImpl) ListEnabled() ([]*ai_model.AIModelConfig, error) {
	var configs []*ai_model.AIModelConfig
	err := d.db.Where("is_enabled = ?", 1).Order("id ASC").Find(&configs).Error
	if err != nil {
		return nil, fmt.Errorf("查询启用模型列表失败: %v", err)
	}
	return configs, nil
}

// ListDisabledModelIds 查询所有禁用的模型ID列表
func (d *aiModelConfigDAOImpl) ListDisabledModelIds() ([]string, error) {
	var modelIds []string
	err := d.db.Model(&ai_model.AIModelConfig{}).
		Where("is_enabled = ?", 0).
		Pluck("model_id", &modelIds).Error
	if err != nil {
		return nil, fmt.Errorf("查询禁用模型ID列表失败: %v", err)
	}
	return modelIds, nil
}

// UpdateStatus 更新模型启用状态
func (d *aiModelConfigDAOImpl) UpdateStatus(modelId string, isEnabled int) error {
	result := d.db.Model(&ai_model.AIModelConfig{}).
		Where("model_id = ?", modelId).
		Update("is_enabled", isEnabled)
	if result.Error != nil {
		return fmt.Errorf("更新模型状态失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("模型不存在: %s", modelId)
	}
	return nil
}

// Upsert 插入或更新模型配置（如果model_id已存在则更新，否则插入）
func (d *aiModelConfigDAOImpl) Upsert(config *ai_model.AIModelConfig) error {
	existing, err := d.GetByModelId(config.ModelId)
	if err != nil {
		return err
	}
	if existing != nil {
		// 更新已有记录
		return d.db.Model(&ai_model.AIModelConfig{}).
			Where("model_id = ?", config.ModelId).
			Updates(map[string]interface{}{
				"model_name":  config.ModelName,
				"provider":    config.Provider,
				"description": config.Description,
				"is_enabled":  config.IsEnabled,
			}).Error
	}
	// 插入新记录
	return d.db.Create(config).Error
}
