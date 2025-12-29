package service

import (
	"fmt"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/model"
	pb "github.com/yzf120/elysia-backend/proto/user"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct {
	userDAO dao.UserDAO
}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{
		userDAO: dao.NewUserDAO(),
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(req *pb.CreateUserRequest) (*pb.User, error) {
	// 参数校验
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	// 检查手机号是否已存在
	exists, err := s.userDAO.CheckPhoneExists(req.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("检查手机号失败: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("手机号已被注册")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 生成用户ID
	userID := fmt.Sprintf("user_%d", time.Now().UnixNano())

	// 构建用户模型
	userModel := &model.User{
		UserID:          userID,
		UserName:        req.UserName,
		Password:        string(hashedPassword),
		Email:           req.Email,
		Gender:          req.Gender,
		PhoneNumber:     req.PhoneNumber,
		WxMiniAppOpenID: req.WxMiniAppOpenId,
		ChineseName:     req.ChineseName,
		ImageURL:        req.ImageUrl,
		RegisterSource:  req.RegisterSource,
		UserType:        req.UserType,
		Status:          2, // 默认状态：可用
	}

	// 创建用户
	if err := s.userDAO.CreateUser(userModel); err != nil {
		return nil, err
	}

	// 查询创建的用户
	createdUser, err := s.userDAO.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	return s.convertModelToProto(createdUser), nil
}

// GetUser 获取用户信息
func (s *UserService) GetUser(req *pb.GetUserRequest) (*pb.User, error) {
	// 参数校验
	if err := s.validateGetUserRequest(req); err != nil {
		return nil, err
	}

	var userModel *model.User
	var err error

	// 根据不同条件查询
	if req.UserId != "" {
		userModel, err = s.userDAO.GetUserByID(req.UserId)
	} else if req.PhoneNumber != "" {
		userModel, err = s.userDAO.GetUserByPhoneNumber(req.PhoneNumber)
	} else if req.WxMiniAppOpenId != "" {
		userModel, err = s.userDAO.GetUserByWxOpenID(req.WxMiniAppOpenId)
	}

	if err != nil {
		return nil, err
	}

	if userModel == nil {
		return nil, fmt.Errorf("用户不存在")
	}

	return s.convertModelToProto(userModel), nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(req *pb.UpdateUserRequest) (*pb.User, error) {
	// 参数校验
	if err := s.validateUpdateUserRequest(req); err != nil {
		return nil, err
	}

	// 检查用户是否存在
	existingUser, err := s.userDAO.GetUserByID(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	if existingUser == nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if req.UserName != "" {
		updates["user_name"] = req.UserName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Gender != 0 {
		updates["gender"] = req.Gender
	}
	if req.PhoneNumber != "" {
		updates["phone_number"] = req.PhoneNumber
	}
	if req.ChineseName != "" {
		updates["chinese_name"] = req.ChineseName
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}
	if req.ImageUrl != "" {
		updates["image_url"] = req.ImageUrl
	}
	if req.UserType != "" {
		updates["user_type"] = req.UserType
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("没有需要更新的字段")
	}

	// 执行更新
	if err := s.userDAO.UpdateUser(req.UserId, updates); err != nil {
		return nil, err
	}

	// 查询更新后的用户
	updatedUser, err := s.userDAO.GetUserByID(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	return s.convertModelToProto(updatedUser), nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(req *pb.DeleteUserRequest) error {
	// 参数校验
	if req.UserId == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	// 检查用户是否存在
	existingUser, err := s.userDAO.GetUserByID(req.UserId)
	if err != nil {
		return fmt.Errorf("查询用户失败: %v", err)
	}
	if existingUser == nil {
		return fmt.Errorf("用户不存在")
	}

	// 执行删除
	return s.userDAO.DeleteUser(req.UserId)
}

// ListUsers 查询用户列表
func (s *UserService) ListUsers(req *pb.ListUsersRequest) ([]*pb.User, int32, error) {
	// 参数校验和默认值设置
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 构建查询条件
	whereClause := "WHERE 1=1"
	var args []interface{}

	if req.Status != 0 {
		whereClause += " AND status = ?"
		args = append(args, req.Status)
	}
	if req.RegisterSource != "" {
		whereClause += " AND register_source = ?"
		args = append(args, req.RegisterSource)
	}
	if req.UserType != "" {
		whereClause += " AND user_type = ?"
		args = append(args, req.UserType)
	}

	// 查询总数
	total, err := s.userDAO.CountUsers(whereClause, args)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (req.Page - 1) * req.PageSize
	userModels, err := s.userDAO.ListUsers(whereClause, args, req.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// 转换为proto
	users := make([]*pb.User, 0, len(userModels))
	for _, model := range userModels {
		users = append(users, s.convertModelToProto(model))
	}

	return users, total, nil
}

// validateCreateUserRequest 校验创建用户请求
func (s *UserService) validateCreateUserRequest(req *pb.CreateUserRequest) error {
	if req.UserName == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if req.Password == "" {
		return fmt.Errorf("密码不能为空")
	}
	if req.PhoneNumber == "" {
		return fmt.Errorf("手机号不能为空")
	}
	if req.RegisterSource == "" {
		return fmt.Errorf("注册来源不能为空")
	}
	return nil
}

// validateGetUserRequest 校验获取用户请求
func (s *UserService) validateGetUserRequest(req *pb.GetUserRequest) error {
	if req.UserId == "" && req.PhoneNumber == "" && req.WxMiniAppOpenId == "" {
		return fmt.Errorf("用户ID、手机号、微信openID至少提供一个")
	}
	return nil
}

// validateUpdateUserRequest 校验更新用户请求
func (s *UserService) validateUpdateUserRequest(req *pb.UpdateUserRequest) error {
	if req.UserId == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	return nil
}

// convertModelToProto 将数据模型转换为proto
func (s *UserService) convertModelToProto(model *model.User) *pb.User {
	if model == nil {
		return nil
	}
	return &pb.User{
		UserId:          model.UserID,
		UserName:        model.UserName,
		Email:           model.Email,
		Gender:          model.Gender,
		PhoneNumber:     model.PhoneNumber,
		WxMiniAppOpenId: model.WxMiniAppOpenID,
		ChineseName:     model.ChineseName,
		RiskLevel:       model.RiskLevel,
		Status:          model.Status,
		CreateTime:      model.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:      model.UpdateTime.Format("2006-01-02 15:04:05"),
		ImageUrl:        model.ImageURL,
		RegisterSource:  model.RegisterSource,
		UserType:        model.UserType,
	}
}
