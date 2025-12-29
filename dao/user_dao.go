package dao

import (
	"fmt"

	"github.com/yzf120/elysia-backend/model"
	"gorm.io/gorm"
)

// UserDAO 用户数据访问接口
type UserDAO interface {
	// CreateUser 创建用户
	CreateUser(user *model.User) error

	// GetUserById 根据用户ID查询用户
	GetUserById(userId string) (*model.User, error)

	// GetUserByPhoneNumber 根据手机号查询用户
	GetUserByPhoneNumber(phoneNumber string) (*model.User, error)

	// GetUserByWxOpenId 根据微信openID查询用户
	GetUserByWxOpenId(wxOpenId string) (*model.User, error)

	// CheckPhoneExists 检查手机号是否存在
	CheckPhoneExists(phoneNumber string) (bool, error)

	// UpdateUser 更新用户信息
	UpdateUser(userID string, updates map[string]interface{}) error

	// DeleteUser 删除用户（软删除）
	DeleteUser(userID string) error

	// ListUsers 查询用户列表
	ListUsers(whereClause string, args []interface{}, limit, offset int32) ([]*model.User, error)

	// CountUsers 统计用户数量
	CountUsers(whereClause string, args []interface{}) (int32, error)
}

// userDAOImpl 用户数据访问实现
type userDAOImpl struct {
	db *gorm.DB
}

// NewUserDAO 创建用户DAO实例
func NewUserDAO() UserDAO {
	return &userDAOImpl{
		db: GetDB(),
	}
}

// CreateUser 创建用户
func (d *userDAOImpl) CreateUser(user *model.User) error {
	if err := d.db.Create(user).Error; err != nil {
		return fmt.Errorf("创建用户失败: %v", err)
	}
	return nil
}

// GetUserById 根据用户ID查询用户
func (d *userDAOImpl) GetUserById(userId string) (*model.User, error) {
	var user model.User
	err := d.db.Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	return &user, nil
}

// GetUserByPhoneNumber 根据手机号查询用户
func (d *userDAOImpl) GetUserByPhoneNumber(phoneNumber string) (*model.User, error) {
	var user model.User
	err := d.db.Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	return &user, nil
}

// GetUserByWxOpenId 根据微信openID查询用户
func (d *userDAOImpl) GetUserByWxOpenId(wxOpenId string) (*model.User, error) {
	var user model.User
	err := d.db.Where("wx_mini_app_open_id = ?", wxOpenId).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	return &user, nil
}

// CheckPhoneExists 检查手机号是否存在
func (d *userDAOImpl) CheckPhoneExists(phoneNumber string) (bool, error) {
	var count int64
	err := d.db.Model(&model.User{}).Where("phone_number = ?", phoneNumber).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查手机号失败: %v", err)
	}
	return count > 0, nil
}

// UpdateUser 更新用户信息
func (d *userDAOImpl) UpdateUser(userID string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("没有需要更新的字段")
	}

	err := d.db.Model(&model.User{}).Where("user_id = ?", userID).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("更新用户失败: %v", err)
	}
	return nil
}

// DeleteUser 删除用户（软删除）
func (d *userDAOImpl) DeleteUser(userID string) error {
	err := d.db.Model(&model.User{}).Where("user_id = ?", userID).Update("status", 4).Error
	if err != nil {
		return fmt.Errorf("删除用户失败: %v", err)
	}
	return nil
}

// ListUsers 查询用户列表
func (d *userDAOImpl) ListUsers(whereClause string, args []interface{}, limit, offset int32) ([]*model.User, error) {
	var users []*model.User
	query := d.db.Model(&model.User{})

	// 如果有 where 条件，添加条件
	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Order("create_time DESC").Limit(int(limit)).Offset(int(offset)).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %v", err)
	}
	return users, nil
}

// CountUsers 统计用户数量
func (d *userDAOImpl) CountUsers(whereClause string, args []interface{}) (int32, error) {
	var count int64
	query := d.db.Model(&model.User{})

	// 如果有 where 条件，添加条件
	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("统计用户数量失败: %v", err)
	}
	return int32(count), nil
}
