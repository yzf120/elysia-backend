package dao

import (
	"fmt"

	"github.com/yzf120/elysia-backend/model/intent"
	"gorm.io/gorm"
)

// IntentDictDAO 意图字典数据访问接口
type IntentDictDAO interface {
	// CreateIntentDict 创建意图字典
	CreateIntentDict(dict *intent.IntentDict) error
	// GetIntentDictById 根据ID查询
	GetIntentDictById(id int) (*intent.IntentDict, error)
	// GetIntentDictByCode 根据意图编码查询
	GetIntentDictByCode(intentCode string) (*intent.IntentDict, error)
	// ListIntentDicts 查询意图字典列表
	ListIntentDicts(page, pageSize int, intentLevel1 string, isValid int) ([]*intent.IntentDict, int64, error)
	// ListValidIntentDicts 查询所有有效的意图字典
	ListValidIntentDicts() ([]*intent.IntentDict, error)
	// UpdateIntentDict 更新意图字典
	UpdateIntentDict(dict *intent.IntentDict) error
	// UpdateIntentDictStatus 更新意图字典状态
	UpdateIntentDictStatus(id int, isValid int) error
	// DeleteIntentDict 删除意图字典（软删除：设为无效）
	DeleteIntentDict(id int) error
}

// intentDictDAOImpl 意图字典数据访问实现
type intentDictDAOImpl struct {
	db *gorm.DB
}

// NewIntentDictDAO 创建意图字典DAO实例
func NewIntentDictDAO() IntentDictDAO {
	return &intentDictDAOImpl{
		db: GetDB(),
	}
}

// CreateIntentDict 创建意图字典
func (d *intentDictDAOImpl) CreateIntentDict(dict *intent.IntentDict) error {
	if dict == nil {
		return fmt.Errorf("intent dict cannot be nil")
	}
	return d.db.Create(dict).Error
}

// GetIntentDictById 根据ID查询
func (d *intentDictDAOImpl) GetIntentDictById(id int) (*intent.IntentDict, error) {
	var dict intent.IntentDict
	err := d.db.Where("id = ?", id).First(&dict).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询意图字典失败: %v", err)
	}
	return &dict, nil
}

// GetIntentDictByCode 根据意图编码查询
func (d *intentDictDAOImpl) GetIntentDictByCode(intentCode string) (*intent.IntentDict, error) {
	var dict intent.IntentDict
	err := d.db.Where("intent_code = ?", intentCode).First(&dict).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询意图字典失败: %v", err)
	}
	return &dict, nil
}

// ListIntentDicts 查询意图字典列表（支持按一级意图和状态过滤）
func (d *intentDictDAOImpl) ListIntentDicts(page, pageSize int, intentLevel1 string, isValid int) ([]*intent.IntentDict, int64, error) {
	var dicts []*intent.IntentDict
	var total int64

	db := d.db.Model(&intent.IntentDict{})

	// 添加过滤条件
	if intentLevel1 != "" {
		db = db.Where("intent_level1 = ?", intentLevel1)
	}
	if isValid >= 0 {
		db = db.Where("is_valid = ?", isValid)
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询意图字典总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := db.Order("priority DESC, create_time DESC").Offset(offset).Limit(pageSize).Find(&dicts).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询意图字典列表失败: %v", err)
	}

	return dicts, total, nil
}

// ListValidIntentDicts 查询所有有效的意图字典
func (d *intentDictDAOImpl) ListValidIntentDicts() ([]*intent.IntentDict, error) {
	var dicts []*intent.IntentDict
	err := d.db.Where("is_valid = ?", 1).Order("priority DESC").Find(&dicts).Error
	if err != nil {
		return nil, fmt.Errorf("查询有效意图字典失败: %v", err)
	}
	return dicts, nil
}

// UpdateIntentDict 更新意图字典
func (d *intentDictDAOImpl) UpdateIntentDict(dict *intent.IntentDict) error {
	if dict == nil {
		return fmt.Errorf("intent dict cannot be nil")
	}
	return d.db.Save(dict).Error
}

// UpdateIntentDictStatus 更新意图字典状态
func (d *intentDictDAOImpl) UpdateIntentDictStatus(id int, isValid int) error {
	return d.db.Model(&intent.IntentDict{}).
		Where("id = ?", id).
		Update("is_valid", isValid).Error
}

// DeleteIntentDict 删除意图字典（软删除：设为无效）
func (d *intentDictDAOImpl) DeleteIntentDict(id int) error {
	return d.UpdateIntentDictStatus(id, 0)
}
