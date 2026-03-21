package dao

import (
	"github.com/yzf120/elysia-backend/model/content_security"
)

// ViolationRecordDAO 违规记录数据访问对象
type ViolationRecordDAO interface {
	// CreateRecord 创建违规记录
	CreateRecord(record *content_security.ViolationRecord) error
	// CountByUserId 统计用户的违规次数
	CountByUserId(userId string) (int64, error)
	// ListByUserId 查询用户的违规记录列表
	ListByUserId(userId string, limit, offset int) ([]*content_security.ViolationRecord, error)
}

type violationRecordDAOImpl struct{}

// NewViolationRecordDAO 创建违规记录DAO
func NewViolationRecordDAO() ViolationRecordDAO {
	return &violationRecordDAOImpl{}
}

// CreateRecord 创建违规记录
func (d *violationRecordDAOImpl) CreateRecord(record *content_security.ViolationRecord) error {
	return DB.Create(record).Error
}

// CountByUserId 统计用户的违规次数
func (d *violationRecordDAOImpl) CountByUserId(userId string) (int64, error) {
	var count int64
	err := DB.Model(&content_security.ViolationRecord{}).Where("user_id = ?", userId).Count(&count).Error
	return count, err
}

// ListByUserId 查询用户的违规记录列表
func (d *violationRecordDAOImpl) ListByUserId(userId string, limit, offset int) ([]*content_security.ViolationRecord, error) {
	var records []*content_security.ViolationRecord
	err := DB.Model(&content_security.ViolationRecord{}).
		Where("user_id = ?", userId).
		Order("create_time DESC").
		Limit(limit).Offset(offset).
		Find(&records).Error
	return records, err
}
