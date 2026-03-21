package dao

import (
	"fmt"

	"github.com/yzf120/elysia-backend/model/intent"
	"gorm.io/gorm"
)

// UserIntentRecordDAO 用户意图记录数据访问接口
type UserIntentRecordDAO interface {
	// CreateRecord 创建用户意图记录
	CreateRecord(record *intent.UserIntentRecord) error
	// ListRecordsByUserId 根据用户ID查询意图记录列表
	ListRecordsByUserId(userId string, page, pageSize int) ([]*intent.UserIntentRecord, int64, error)
	// ListRecords 查询意图记录列表（支持多条件过滤）
	ListRecords(page, pageSize int, userId, intentCode, intentLevel1 string) ([]*intent.UserIntentRecord, int64, error)
	// CountByIntentCode 按意图编码统计次数
	CountByIntentCode() (map[string]int64, error)
}

// userIntentRecordDAOImpl 用户意图记录数据访问实现
type userIntentRecordDAOImpl struct {
	db *gorm.DB
}

// NewUserIntentRecordDAO 创建用户意图记录DAO实例
func NewUserIntentRecordDAO() UserIntentRecordDAO {
	return &userIntentRecordDAOImpl{
		db: GetDB(),
	}
}

// CreateRecord 创建用户意图记录
func (d *userIntentRecordDAOImpl) CreateRecord(record *intent.UserIntentRecord) error {
	if record == nil {
		return fmt.Errorf("record cannot be nil")
	}
	return d.db.Create(record).Error
}

// ListRecordsByUserId 根据用户ID查询意图记录列表
func (d *userIntentRecordDAOImpl) ListRecordsByUserId(userId string, page, pageSize int) ([]*intent.UserIntentRecord, int64, error) {
	var records []*intent.UserIntentRecord
	var total int64

	db := d.db.Model(&intent.UserIntentRecord{}).Where("user_id = ?", userId)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询用户意图记录总数失败: %v", err)
	}

	offset := (page - 1) * pageSize
	err := db.Order("create_time DESC").Offset(offset).Limit(pageSize).Find(&records).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户意图记录列表失败: %v", err)
	}

	return records, total, nil
}

// ListRecords 查询意图记录列表（支持多条件过滤）
func (d *userIntentRecordDAOImpl) ListRecords(page, pageSize int, userId, intentCode, intentLevel1 string) ([]*intent.UserIntentRecord, int64, error) {
	var records []*intent.UserIntentRecord
	var total int64

	db := d.db.Model(&intent.UserIntentRecord{})

	if userId != "" {
		db = db.Where("user_id = ?", userId)
	}
	if intentCode != "" {
		db = db.Where("intent_code = ?", intentCode)
	}
	if intentLevel1 != "" {
		db = db.Where("intent_level1 = ?", intentLevel1)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询意图记录总数失败: %v", err)
	}

	offset := (page - 1) * pageSize
	err := db.Order("create_time DESC").Offset(offset).Limit(pageSize).Find(&records).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询意图记录列表失败: %v", err)
	}

	return records, total, nil
}

// CountByIntentCode 按意图编码统计次数
func (d *userIntentRecordDAOImpl) CountByIntentCode() (map[string]int64, error) {
	type Result struct {
		IntentCode string `gorm:"column:intent_code"`
		Count      int64  `gorm:"column:count"`
	}

	var results []Result
	err := d.db.Model(&intent.UserIntentRecord{}).
		Select("intent_code, COUNT(*) as count").
		Group("intent_code").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("统计意图记录失败: %v", err)
	}

	countMap := make(map[string]int64)
	for _, r := range results {
		countMap[r.IntentCode] = r.Count
	}
	return countMap, nil
}
