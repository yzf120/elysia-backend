package dao

import (
	"fmt"

	"github.com/yzf120/elysia-backend/model/intent"
	"gorm.io/gorm"
)

// IntentPromptTemplateDAO 意图提示词模板数据访问接口
type IntentPromptTemplateDAO interface {
	// CreateTemplate 创建提示词模板
	CreateTemplate(tpl *intent.IntentPromptTemplate) error
	// GetTemplateById 根据ID查询
	GetTemplateById(id int) (*intent.IntentPromptTemplate, error)
	// ListTemplatesByIntentCode 根据意图编码查询模板列表
	ListTemplatesByIntentCode(intentCode string) ([]*intent.IntentPromptTemplate, error)
	// ListTemplates 查询模板列表（支持分页和过滤）
	ListTemplates(page, pageSize int, intentCode, templateType string) ([]*intent.IntentPromptTemplate, int64, error)
	// GetActiveTemplate 获取某意图某类型的启用模板
	GetActiveTemplate(intentCode, templateType string) (*intent.IntentPromptTemplate, error)
	// UpdateTemplate 更新模板
	UpdateTemplate(tpl *intent.IntentPromptTemplate) error
	// UpdateTemplateStatus 更新模板状态
	UpdateTemplateStatus(id int, isActive int) error
	// DeleteTemplate 删除模板
	DeleteTemplate(id int) error
}

// intentPromptTemplateDAOImpl 意图提示词模板数据访问实现
type intentPromptTemplateDAOImpl struct {
	db *gorm.DB
}

// NewIntentPromptTemplateDAO 创建意图提示词模板DAO实例
func NewIntentPromptTemplateDAO() IntentPromptTemplateDAO {
	return &intentPromptTemplateDAOImpl{
		db: GetDB(),
	}
}

// CreateTemplate 创建提示词模板
func (d *intentPromptTemplateDAOImpl) CreateTemplate(tpl *intent.IntentPromptTemplate) error {
	if tpl == nil {
		return fmt.Errorf("template cannot be nil")
	}
	return d.db.Create(tpl).Error
}

// GetTemplateById 根据ID查询
func (d *intentPromptTemplateDAOImpl) GetTemplateById(id int) (*intent.IntentPromptTemplate, error) {
	var tpl intent.IntentPromptTemplate
	err := d.db.Where("id = ?", id).First(&tpl).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询模板失败: %v", err)
	}
	return &tpl, nil
}

// ListTemplatesByIntentCode 根据意图编码查询模板列表
func (d *intentPromptTemplateDAOImpl) ListTemplatesByIntentCode(intentCode string) ([]*intent.IntentPromptTemplate, error) {
	var tpls []*intent.IntentPromptTemplate
	err := d.db.Where("intent_code = ?", intentCode).Order("version DESC").Find(&tpls).Error
	if err != nil {
		return nil, fmt.Errorf("查询模板列表失败: %v", err)
	}
	return tpls, nil
}

// ListTemplates 查询模板列表（支持分页和过滤）
func (d *intentPromptTemplateDAOImpl) ListTemplates(page, pageSize int, intentCode, templateType string) ([]*intent.IntentPromptTemplate, int64, error) {
	var tpls []*intent.IntentPromptTemplate
	var total int64

	db := d.db.Model(&intent.IntentPromptTemplate{})

	if intentCode != "" {
		db = db.Where("intent_code = ?", intentCode)
	}
	if templateType != "" {
		db = db.Where("template_type = ?", templateType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询模板总数失败: %v", err)
	}

	offset := (page - 1) * pageSize
	err := db.Order("intent_code ASC, template_type ASC, version DESC").Offset(offset).Limit(pageSize).Find(&tpls).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询模板列表失败: %v", err)
	}

	return tpls, total, nil
}

// GetActiveTemplate 获取某意图某类型的启用模板
func (d *intentPromptTemplateDAOImpl) GetActiveTemplate(intentCode, templateType string) (*intent.IntentPromptTemplate, error) {
	var tpl intent.IntentPromptTemplate
	err := d.db.Where("intent_code = ? AND template_type = ? AND is_active = 1", intentCode, templateType).
		Order("version DESC").
		First(&tpl).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询启用模板失败: %v", err)
	}
	return &tpl, nil
}

// UpdateTemplate 更新模板
func (d *intentPromptTemplateDAOImpl) UpdateTemplate(tpl *intent.IntentPromptTemplate) error {
	if tpl == nil {
		return fmt.Errorf("template cannot be nil")
	}
	return d.db.Save(tpl).Error
}

// UpdateTemplateStatus 更新模板状态
func (d *intentPromptTemplateDAOImpl) UpdateTemplateStatus(id int, isActive int) error {
	return d.db.Model(&intent.IntentPromptTemplate{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}

// DeleteTemplate 删除模板（物理删除）
func (d *intentPromptTemplateDAOImpl) DeleteTemplate(id int) error {
	return d.db.Where("id = ?", id).Delete(&intent.IntentPromptTemplate{}).Error
}
