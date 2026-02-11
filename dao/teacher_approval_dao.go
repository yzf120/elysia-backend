package dao

import (
	"github.com/yzf120/elysia-backend/model/teacher"
)

// TeacherApprovalDAO 教师审批单数据访问接口
type TeacherApprovalDAO interface {
	// CreateApproval 创建审批单
	CreateApproval(approval *teacher.TeacherApproval) error

	// GetApprovalById 根据审批单ID查询
	GetApprovalById(approvalId string) (*teacher.TeacherApproval, error)

	// GetApprovalByTeacherId 根据教师ID查询审批单
	GetApprovalByTeacherId(teacherId string) (*teacher.TeacherApproval, error)

	// UpdateApproval 更新审批单
	UpdateApproval(approvalId string, updates map[string]interface{}) error

	// DeleteApproval 删除审批单
	DeleteApproval(approvalId string) error

	// ListApprovals 查询审批单列表
	ListApprovals(whereClause string, args []interface{}, limit, offset int32) ([]*teacher.TeacherApproval, error)

	// CountApprovals 统计审批单数量
	CountApprovals(whereClause string, args []interface{}) (int32, error)
}

// teacherApprovalDAOImpl 教师审批单DAO实现
type teacherApprovalDAOImpl struct{}

// NewTeacherApprovalDAO 创建教师审批单DAO实例
func NewTeacherApprovalDAO() TeacherApprovalDAO {
	return &teacherApprovalDAOImpl{}
}

// CreateApproval 创建审批单
func (d *teacherApprovalDAOImpl) CreateApproval(approval *teacher.TeacherApproval) error {
	db := DB
	return db.Create(approval).Error
}

// GetApprovalById 根据审批单ID查询
func (d *teacherApprovalDAOImpl) GetApprovalById(approvalId string) (*teacher.TeacherApproval, error) {
	db := DB
	var approval teacher.TeacherApproval
	err := db.Where("approval_id = ?", approvalId).First(&approval).Error
	if err != nil {
		return nil, err
	}
	return &approval, nil
}

// GetApprovalByTeacherId 根据教师ID查询审批单
func (d *teacherApprovalDAOImpl) GetApprovalByTeacherId(teacherId string) (*teacher.TeacherApproval, error) {
	db := DB
	var approval teacher.TeacherApproval
	err := db.Where("teacher_id = ?", teacherId).First(&approval).Error
	if err != nil {
		return nil, err
	}
	return &approval, nil
}

// UpdateApproval 更新审批单
func (d *teacherApprovalDAOImpl) UpdateApproval(approvalId string, updates map[string]interface{}) error {
	db := DB
	return db.Model(&teacher.TeacherApproval{}).Where("approval_id = ?", approvalId).Updates(updates).Error
}

// DeleteApproval 删除审批单
func (d *teacherApprovalDAOImpl) DeleteApproval(approvalId string) error {
	db := DB
	return db.Where("approval_id = ?", approvalId).Delete(&teacher.TeacherApproval{}).Error
}

// ListApprovals 查询审批单列表
func (d *teacherApprovalDAOImpl) ListApprovals(whereClause string, args []interface{}, limit, offset int32) ([]*teacher.TeacherApproval, error) {
	db := DB
	var approvals []*teacher.TeacherApproval
	query := db.Model(&teacher.TeacherApproval{})

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Order("create_time DESC").Limit(int(limit)).Offset(int(offset)).Find(&approvals).Error
	if err != nil {
		return nil, err
	}

	return approvals, nil
}

// CountApprovals 统计审批单数量
func (d *teacherApprovalDAOImpl) CountApprovals(whereClause string, args []interface{}) (int32, error) {
	db := DB
	var count int64
	query := db.Model(&teacher.TeacherApproval{})

	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int32(count), nil
}
