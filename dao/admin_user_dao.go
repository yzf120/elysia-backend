package dao

import (
	"fmt"
	"time"

	"github.com/yzf120/elysia-backend/model"
	"gorm.io/gorm"
)

// AdminUserDAO 管理员用户数据访问接口
type AdminUserDAO interface {
	// CreateAdminUser 创建管理员用户
	CreateAdminUser(adminUser *model.AdminUser) error

	// GetAdminUserById 根据ID查询管理员用户
	GetAdminUserById(id int64) (*model.AdminUser, error)

	// GetAdminUserByAdminId 根据管理员ID查询
	GetAdminUserByAdminId(adminId string) (*model.AdminUser, error)

	// GetAdminUserByUsername 根据用户名查询管理员用户（包含密码）
	GetAdminUserByUsername(username string) (*model.AdminUser, string, error)

	// GetAdminUserByEmail 根据邮箱查询管理员用户
	GetAdminUserByEmail(email string) (*model.AdminUser, error)

	// UpdateAdminUser 更新管理员用户信息
	UpdateAdminUser(adminUser *model.AdminUser) error

	// UpdateAdminUserStatus 更新管理员用户状态
	UpdateAdminUserStatus(adminId string, status int32) error

	// UpdateAdminUserPassword 更新管理员用户密码
	UpdateAdminUserPassword(adminId string, password string) error

	// UpdateAdminUserLoginInfo 更新管理员用户登录信息
	UpdateAdminUserLoginInfo(adminId, ipAddress string, loginTime time.Time) error

	// ListAdminUsers 查询管理员用户列表
	ListAdminUsers(page, pageSize int, role, status string) ([]*model.AdminUser, int64, error)

	// DeleteAdminUser 删除管理员用户（软删除）
	DeleteAdminUser(adminId string) error
}

// adminUserDAOImpl 管理员用户数据访问实现
type adminUserDAOImpl struct {
	db *gorm.DB
}

// NewAdminUserDAO 创建管理员用户DAO实例
func NewAdminUserDAO() AdminUserDAO {
	return &adminUserDAOImpl{
		db: GetDB(),
	}
}

// CreateAdminUser 创建管理员用户
func (d *adminUserDAOImpl) CreateAdminUser(adminUser *model.AdminUser) error {
	if adminUser == nil {
		return fmt.Errorf("admin user cannot be nil")
	}

	return d.db.Create(adminUser).Error
}

// GetAdminUserById 根据ID查询管理员用户
func (d *adminUserDAOImpl) GetAdminUserById(id int64) (*model.AdminUser, error) {
	var adminUser model.AdminUser
	err := d.db.Where("id = ?", id).First(&adminUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询管理员用户失败: %v", err)
	}
	return &adminUser, nil
}

// GetAdminUserByAdminId 根据管理员ID查询
func (d *adminUserDAOImpl) GetAdminUserByAdminId(adminId string) (*model.AdminUser, error) {
	var adminUser model.AdminUser
	err := d.db.Where("admin_id = ?", adminId).First(&adminUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询管理员用户失败: %v", err)
	}
	return &adminUser, nil
}

// GetAdminUserByUsername 根据用户名查询管理员用户（包含密码）
func (d *adminUserDAOImpl) GetAdminUserByUsername(username string) (*model.AdminUser, string, error) {
	var adminUser model.AdminUser
	err := d.db.Where("username = ?", username).First(&adminUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", nil
		}
		return nil, "", fmt.Errorf("查询管理员用户失败: %v", err)
	}

	// 返回用户信息和密码
	password := adminUser.Password
	return &adminUser, password, nil
}

// GetAdminUserByEmail 根据邮箱查询管理员用户
func (d *adminUserDAOImpl) GetAdminUserByEmail(email string) (*model.AdminUser, error) {
	var adminUser model.AdminUser
	err := d.db.Where("email = ?", email).First(&adminUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询管理员用户失败: %v", err)
	}
	return &adminUser, nil
}

// UpdateAdminUser 更新管理员用户信息
func (d *adminUserDAOImpl) UpdateAdminUser(adminUser *model.AdminUser) error {
	if adminUser == nil {
		return fmt.Errorf("admin user cannot be nil")
	}

	return d.db.Save(adminUser).Error
}

// UpdateAdminUserStatus 更新管理员用户状态
func (d *adminUserDAOImpl) UpdateAdminUserStatus(adminId string, status int32) error {
	return d.db.Model(&model.AdminUser{}).
		Where("admin_id = ?", adminId).
		Update("status", status).Error
}

// UpdateAdminUserPassword 更新管理员用户密码
func (d *adminUserDAOImpl) UpdateAdminUserPassword(adminId string, password string) error {
	return d.db.Model(&model.AdminUser{}).
		Where("admin_id = ?", adminId).
		Updates(map[string]interface{}{
			"password":             password,
			"password_update_time": time.Now(),
		}).Error
}

// UpdateAdminUserLoginInfo 更新管理员用户登录信息
func (d *adminUserDAOImpl) UpdateAdminUserLoginInfo(adminId, ipAddress string, loginTime time.Time) error {
	return d.db.Model(&model.AdminUser{}).
		Where("admin_id = ?", adminId).
		Updates(map[string]interface{}{
			"last_login_time":  loginTime,
			"last_login_ip":    ipAddress,
			"login_fail_count": 0, // 登录成功重置失败次数
		}).Error
}

// ListAdminUsers 查询管理员用户列表
func (d *adminUserDAOImpl) ListAdminUsers(page, pageSize int, role, status string) ([]*model.AdminUser, int64, error) {
	var adminUsers []*model.AdminUser
	var total int64

	db := d.db.Model(&model.AdminUser{})

	// 添加过滤条件
	if role != "" {
		db = db.Where("role = ?", role)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询管理员用户总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := db.Order("create_time DESC").Offset(offset).Limit(pageSize).Find(&adminUsers).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询管理员用户列表失败: %v", err)
	}

	return adminUsers, total, nil
}

// DeleteAdminUser 删除管理员用户（软删除）
func (d *adminUserDAOImpl) DeleteAdminUser(adminId string) error {
	// 管理员用户不支持删除，只支持状态控制
	return d.UpdateAdminUserStatus(adminId, 0) // 设置为禁用状态
}
