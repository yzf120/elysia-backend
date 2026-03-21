package service

import (
	"fmt"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/model/intent"
)

// IntentService 意图管理服务
type IntentService struct {
	intentDictDAO     dao.IntentDictDAO
	promptTemplateDAO dao.IntentPromptTemplateDAO
	intentRecordDAO   dao.UserIntentRecordDAO
}

// NewIntentService 创建意图管理服务
func NewIntentService() *IntentService {
	return &IntentService{
		intentDictDAO:     dao.NewIntentDictDAO(),
		promptTemplateDAO: dao.NewIntentPromptTemplateDAO(),
		intentRecordDAO:   dao.NewUserIntentRecordDAO(),
	}
}

// ==================== 意图字典管理 ====================

// CreateIntentDict 创建意图字典
func (s *IntentService) CreateIntentDict(dict *intent.IntentDict) error {
	// 检查编码是否已存在
	existing, err := s.intentDictDAO.GetIntentDictByCode(dict.IntentCode)
	if err != nil {
		return fmt.Errorf("检查意图编码失败: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("意图编码 %s 已存在", dict.IntentCode)
	}
	return s.intentDictDAO.CreateIntentDict(dict)
}

// GetIntentDictById 根据ID查询意图字典
func (s *IntentService) GetIntentDictById(id int) (*intent.IntentDict, error) {
	dict, err := s.intentDictDAO.GetIntentDictById(id)
	if err != nil {
		return nil, fmt.Errorf("查询意图字典失败: %v", err)
	}
	if dict == nil {
		return nil, fmt.Errorf("意图字典不存在 (id=%d)", id)
	}
	return dict, nil
}

// GetIntentDictByCode 根据编码查询意图字典
func (s *IntentService) GetIntentDictByCode(intentCode string) (*intent.IntentDict, error) {
	dict, err := s.intentDictDAO.GetIntentDictByCode(intentCode)
	if err != nil {
		return nil, fmt.Errorf("查询意图字典失败: %v", err)
	}
	if dict == nil {
		return nil, fmt.Errorf("意图字典不存在 (code=%s)", intentCode)
	}
	return dict, nil
}

// ListIntentDicts 查询意图字典列表
func (s *IntentService) ListIntentDicts(page, pageSize int, intentLevel1 string, isValid int) ([]*intent.IntentDict, int64, error) {
	return s.intentDictDAO.ListIntentDicts(page, pageSize, intentLevel1, isValid)
}

// ListValidIntentDicts 查询所有有效的意图字典（供意图Agent使用）
func (s *IntentService) ListValidIntentDicts() ([]*intent.IntentDict, error) {
	return s.intentDictDAO.ListValidIntentDicts()
}

// UpdateIntentDict 更新意图字典
func (s *IntentService) UpdateIntentDict(dict *intent.IntentDict) error {
	// 检查记录是否存在
	existing, err := s.intentDictDAO.GetIntentDictById(dict.Id)
	if err != nil {
		return fmt.Errorf("查询意图字典失败: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("意图字典不存在 (id=%d)", dict.Id)
	}

	// 如果修改了编码，检查新编码是否冲突
	if dict.IntentCode != existing.IntentCode {
		conflict, err := s.intentDictDAO.GetIntentDictByCode(dict.IntentCode)
		if err != nil {
			return fmt.Errorf("检查意图编码失败: %v", err)
		}
		if conflict != nil {
			return fmt.Errorf("意图编码 %s 已存在", dict.IntentCode)
		}
	}

	return s.intentDictDAO.UpdateIntentDict(dict)
}

// UpdateIntentDictStatus 更新意图字典状态
func (s *IntentService) UpdateIntentDictStatus(id int, isValid int) error {
	existing, err := s.intentDictDAO.GetIntentDictById(id)
	if err != nil {
		return fmt.Errorf("查询意图字典失败: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("意图字典不存在 (id=%d)", id)
	}
	return s.intentDictDAO.UpdateIntentDictStatus(id, isValid)
}

// DeleteIntentDict 删除意图字典（软删除）
func (s *IntentService) DeleteIntentDict(id int) error {
	existing, err := s.intentDictDAO.GetIntentDictById(id)
	if err != nil {
		return fmt.Errorf("查询意图字典失败: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("意图字典不存在 (id=%d)", id)
	}
	return s.intentDictDAO.DeleteIntentDict(id)
}

// ==================== 意图提示词模板管理 ====================

// CreatePromptTemplate 创建提示词模板
func (s *IntentService) CreatePromptTemplate(tpl *intent.IntentPromptTemplate) error {
	// 验证关联的意图编码是否存在
	dict, err := s.intentDictDAO.GetIntentDictByCode(tpl.IntentCode)
	if err != nil {
		return fmt.Errorf("查询意图字典失败: %v", err)
	}
	if dict == nil {
		return fmt.Errorf("关联的意图编码 %s 不存在", tpl.IntentCode)
	}
	return s.promptTemplateDAO.CreateTemplate(tpl)
}

// GetPromptTemplateById 根据ID查询模板
func (s *IntentService) GetPromptTemplateById(id int) (*intent.IntentPromptTemplate, error) {
	tpl, err := s.promptTemplateDAO.GetTemplateById(id)
	if err != nil {
		return nil, fmt.Errorf("查询模板失败: %v", err)
	}
	if tpl == nil {
		return nil, fmt.Errorf("模板不存在 (id=%d)", id)
	}
	return tpl, nil
}

// ListPromptTemplates 查询模板列表
func (s *IntentService) ListPromptTemplates(page, pageSize int, intentCode, templateType string) ([]*intent.IntentPromptTemplate, int64, error) {
	return s.promptTemplateDAO.ListTemplates(page, pageSize, intentCode, templateType)
}

// GetActivePromptTemplate 获取某意图某类型的启用模板
func (s *IntentService) GetActivePromptTemplate(intentCode, templateType string) (*intent.IntentPromptTemplate, error) {
	return s.promptTemplateDAO.GetActiveTemplate(intentCode, templateType)
}

// UpdatePromptTemplate 更新模板
func (s *IntentService) UpdatePromptTemplate(tpl *intent.IntentPromptTemplate) error {
	existing, err := s.promptTemplateDAO.GetTemplateById(tpl.Id)
	if err != nil {
		return fmt.Errorf("查询模板失败: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("模板不存在 (id=%d)", tpl.Id)
	}
	return s.promptTemplateDAO.UpdateTemplate(tpl)
}

// UpdatePromptTemplateStatus 更新模板状态
func (s *IntentService) UpdatePromptTemplateStatus(id int, isActive int) error {
	existing, err := s.promptTemplateDAO.GetTemplateById(id)
	if err != nil {
		return fmt.Errorf("查询模板失败: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("模板不存在 (id=%d)", id)
	}
	return s.promptTemplateDAO.UpdateTemplateStatus(id, isActive)
}

// DeletePromptTemplate 删除模板
func (s *IntentService) DeletePromptTemplate(id int) error {
	existing, err := s.promptTemplateDAO.GetTemplateById(id)
	if err != nil {
		return fmt.Errorf("查询模板失败: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("模板不存在 (id=%d)", id)
	}
	return s.promptTemplateDAO.DeleteTemplate(id)
}

// ==================== 用户意图记录管理 ====================

// CreateIntentRecord 创建用户意图记录
func (s *IntentService) CreateIntentRecord(record *intent.UserIntentRecord) error {
	return s.intentRecordDAO.CreateRecord(record)
}

// ListIntentRecords 查询意图记录列表
func (s *IntentService) ListIntentRecords(page, pageSize int, userId, intentCode, intentLevel1 string) ([]*intent.UserIntentRecord, int64, error) {
	return s.intentRecordDAO.ListRecords(page, pageSize, userId, intentCode, intentLevel1)
}

// GetIntentStats 获取意图统计数据
func (s *IntentService) GetIntentStats() (map[string]int64, error) {
	return s.intentRecordDAO.CountByIntentCode()
}
