package dao

import (
	"fmt"

	"github.com/yzf120/elysia-backend/model"
	"gorm.io/gorm"
)

// AuthDAO 认证数据访问接口
type AuthDAO interface {
	// GetUserByPhoneNumberWithPassword 根据手机号查询用户（包含密码）
	GetUserByPhoneNumberWithPassword(phoneNumber string) (*model.User, string, error)
}

// authDAOImpl 认证数据访问实现
type authDAOImpl struct {
	db *gorm.DB
}

// NewAuthDAO 创建认证DAO实例
func NewAuthDAO() AuthDAO {
	return &authDAOImpl{
		db: GetDB(),
	}
}

// GetUserByPhoneNumberWithPassword 根据手机号查询用户（包含密码）
func (d *authDAOImpl) GetUserByPhoneNumberWithPassword(phoneNumber string) (*model.User, string, error) {
	var user model.User
	err := d.db.Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", nil
		}
		return nil, "", fmt.Errorf("查询用户失败: %v", err)
	}

	// 返回用户信息和密码
	password := user.Password
	return &user, password, nil
}
