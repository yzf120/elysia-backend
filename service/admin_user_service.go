package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	"github.com/yzf120/elysia-backend/model/admin"
	adminPb "github.com/yzf120/elysia-backend/proto/admin"
	"github.com/yzf120/elysia-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

// AdminUserService 管理员用户服务
type AdminUserService struct {
	adminUserDAO dao.AdminUserDAO
	jwtService   *utils.JWTService
}

// NewAdminUserService 创建管理员用户服务
func NewAdminUserService() *AdminUserService {
	return &AdminUserService{
		adminUserDAO: dao.NewAdminUserDAO(),
		jwtService:   utils.NewJWTService(),
	}
}

// CreateAdminUser 创建管理员用户
func (s *AdminUserService) CreateAdminUser(req *adminPb.CreateAdminUserRequest) (*adminPb.AdminUserInfo, error) {
	// 参数校验
	if err := s.validateCreateAdminUserRequest(req); err != nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, err.Error())
	}

	// 检查用户名是否已存在
	existingUser, _, err := s.adminUserDAO.GetAdminUserByUsername(req.Username)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "检查用户名失败: "+err.Error())
	}
	if existingUser != nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "用户名已存在")
	}

	// 检查邮箱是否已存在
	existingEmail, err := s.adminUserDAO.GetAdminUserByEmail(req.Email)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "检查邮箱失败: "+err.Error())
	}
	if existingEmail != nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "密码加密失败: "+err.Error())
	}

	// 生成管理员ID
	adminId := utils.GenerateAdminId()

	// 创建管理员用户
	adminUser := &admin.AdminUser{
		AdminId:    adminId,
		Username:   req.Username,
		Password:   string(hashedPassword),
		RealName:   req.RealName,
		Email:      req.Email,
		Role:       req.Role,
		Status:     1, // 默认启用
		Remark:     req.Remark,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := s.adminUserDAO.CreateAdminUser(adminUser); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "创建管理员用户失败: "+err.Error())
	}

	// 转换为响应格式
	adminUserInfo := s.convertModelToAdminUserInfo(adminUser)
	return adminUserInfo, nil
}

// LoginAdminUser 管理员用户登录
func (s *AdminUserService) LoginAdminUser(ctx context.Context, username, password, ipAddress string) (*adminPb.AdminUserInfo, string, error) {
	// 参数校验
	if username == "" {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "用户名不能为空")
	}
	if password == "" {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "密码不能为空")
	}

	// 查询管理员用户
	adminUser, hashedPassword, err := s.adminUserDAO.GetAdminUserByUsername(username)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "查询管理员用户失败: "+err.Error())
	}
	if adminUser == nil {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "管理员用户不存在")
	}

	// 检查用户状态
	if adminUser.Status != 1 {
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "管理员账号已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		// 密码错误，增加登录失败次数
		s.handleLoginFailure(adminUser.AdminId)
		return nil, "", errs.NewCommonError(errs.ErrBadRequest, "密码错误")
	}

	// 登录成功，更新登录信息
	if err := s.adminUserDAO.UpdateAdminUserLoginInfo(adminUser.AdminId, ipAddress, time.Now()); err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "更新登录信息失败: "+err.Error())
	}

	// 生成登录令牌
	token, err := s.jwtService.GenerateToken(adminUser.AdminId)
	if err != nil {
		return nil, "", errs.NewCommonError(errs.ErrInternal, "生成令牌失败: "+err.Error())
	}

	// 转换为响应格式
	adminUserInfo := s.convertModelToAdminUserInfo(adminUser)
	return adminUserInfo, token, nil
}

// UpdateAdminUserPassword 更新管理员用户密码
func (s *AdminUserService) UpdateAdminUserPassword(adminId, oldPassword, newPassword string) error {
	// 参数校验
	if adminId == "" {
		return errs.NewCommonError(errs.ErrBadRequest, "管理员ID不能为空")
	}
	if oldPassword == "" || newPassword == "" {
		return errs.NewCommonError(errs.ErrBadRequest, "密码不能为空")
	}
	if len(newPassword) < 6 {
		return errs.NewCommonError(errs.ErrBadRequest, "新密码长度不能少于6位")
	}

	// 查询管理员用户
	adminUser, hashedPassword, err := s.adminUserDAO.GetAdminUserByUsername("")
	if err != nil {
		return errs.NewCommonError(errs.ErrInternal, "查询管理员用户失败: "+err.Error())
	}
	if adminUser == nil {
		return errs.NewCommonError(errs.ErrBadRequest, "管理员用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(oldPassword)); err != nil {
		return errs.NewCommonError(errs.ErrBadRequest, "旧密码错误")
	}

	// 加密新密码
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errs.NewCommonError(errs.ErrInternal, "密码加密失败: "+err.Error())
	}

	// 更新密码
	if err := s.adminUserDAO.UpdateAdminUserPassword(adminUser.AdminId, string(newHashedPassword)); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "更新密码失败: "+err.Error())
	}

	return nil
}

// GetAdminUserById 根据ID查询管理员用户
func (s *AdminUserService) GetAdminUserById(id int64) (*adminPb.AdminUserInfo, error) {
	adminUser, err := s.adminUserDAO.GetAdminUserById(id)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "查询管理员用户失败: "+err.Error())
	}
	if adminUser == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "管理员用户不存在")
	}

	return s.convertModelToAdminUserInfo(adminUser), nil
}

// GetAdminUserByAdminId 根据管理员ID查询
func (s *AdminUserService) GetAdminUserByAdminId(adminId string) (*adminPb.AdminUserInfo, error) {
	adminUser, err := s.adminUserDAO.GetAdminUserByAdminId(adminId)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "查询管理员用户失败: "+err.Error())
	}
	if adminUser == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "管理员用户不存在")
	}

	return s.convertModelToAdminUserInfo(adminUser), nil
}

// ListAdminUsers 查询管理员用户列表
func (s *AdminUserService) ListAdminUsers(page, pageSize int, role, status string) ([]*adminPb.AdminUserInfo, int64, error) {
	adminUsers, total, err := s.adminUserDAO.ListAdminUsers(page, pageSize, role, status)
	if err != nil {
		return nil, 0, errs.NewCommonError(errs.ErrInternal, "查询管理员用户列表失败: "+err.Error())
	}

	// 转换为响应格式
	var adminUserInfos []*adminPb.AdminUserInfo
	for _, adminUser := range adminUsers {
		adminUserInfos = append(adminUserInfos, s.convertModelToAdminUserInfo(adminUser))
	}

	return adminUserInfos, total, nil
}

// UpdateAdminUserStatus 更新管理员用户状态
func (s *AdminUserService) UpdateAdminUserStatus(adminId string, status int32) error {
	if adminId == "" {
		return errs.NewCommonError(errs.ErrBadRequest, "管理员ID不能为空")
	}

	if err := s.adminUserDAO.UpdateAdminUserStatus(adminId, status); err != nil {
		return errs.NewCommonError(errs.ErrInternal, "更新管理员用户状态失败: "+err.Error())
	}

	return nil
}

// handleLoginFailure 处理登录失败
func (s *AdminUserService) handleLoginFailure(adminId string) {
	// 这里可以添加登录失败次数限制逻辑
	// 例如：连续失败5次锁定账号
	// 当前先简单记录失败次数
	s.adminUserDAO.UpdateAdminUserLoginInfo(adminId, "", time.Now())
}

// validateCreateAdminUserRequest 校验创建管理员用户请求
func (s *AdminUserService) validateCreateAdminUserRequest(req *adminPb.CreateAdminUserRequest) error {
	if req.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if req.Password == "" {
		return fmt.Errorf("密码不能为空")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("密码长度不能少于6位")
	}
	if req.Email == "" {
		return fmt.Errorf("邮箱不能为空")
	}
	if req.RealName == "" {
		return fmt.Errorf("真实姓名不能为空")
	}
	if req.Role == "" {
		return fmt.Errorf("角色不能为空")
	}
	return nil
}

// convertModelToAdminUserInfo 将数据模型转换为管理员用户信息
func (s *AdminUserService) convertModelToAdminUserInfo(model *admin.AdminUser) *adminPb.AdminUserInfo {
	if model == nil {
		return nil
	}
	return &adminPb.AdminUserInfo{
		Id:                 model.Id,
		AdminId:            model.AdminId,
		Username:           model.Username,
		RealName:           model.RealName,
		Email:              model.Email,
		Role:               model.Role,
		Status:             model.Status,
		LastLoginTime:      model.LastLoginTime.Format("2006-01-02 15:04:05"),
		LastLoginIp:        model.LastLoginIp,
		LoginFailCount:     model.LoginFailCount,
		PasswordUpdateTime: model.PasswordUpdateTime.Format("2006-01-02 15:04:05"),
		CreateTime:         model.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:         model.UpdateTime.Format("2006-01-02 15:04:05"),
		Remark:             model.Remark,
	}
}
