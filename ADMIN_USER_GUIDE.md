# 管理员用户使用指南

## 管理员用户特性

### 与普通用户的区别
| 特性 | 普通用户 | 管理员用户 |
|------|----------|------------|
| 登录方式 | 手机号/验证码、账号密码 | 仅账号密码 |
| 账号格式 | 手机号为主 | 自定义账号（不一定是手机号） |
| 验证码登录 | 支持 | 不支持 |
| 密码更新 | 支持 | 支持 |
| 账号删除 | 支持 | 不支持（仅状态控制） |
| 注册方式 | 多种渠道 | 手动创建 |

### 管理员用户表字段说明
- `admin_id`: 管理员唯一标识
- `username`: 登录账号（自定义，不一定是手机号）
- `password`: 加密密码
- `real_name`: 真实姓名
- `email`: 邮箱地址
- `role`: 角色权限（admin/super_admin）
- `status`: 账号状态（0-禁用，1-启用）

## 默认管理员账号

系统初始化时会自动创建默认超级管理员：
- **账号**: `admin`
- **密码**: `admin123`（建议首次登录后修改）
- **角色**: `super_admin`
- **邮箱**: `admin@elysia.com`

## 管理员登录流程

### 登录接口要求
1. **仅支持账号密码登录**，不支持验证码登录
2. 账号可以是任意字符串，不限制为手机号格式
3. 密码需要加密传输
4. 登录成功后记录操作日志

### 密码安全策略
1. 密码必须加密存储（使用 bcrypt 等算法）
2. 支持密码更新功能
3. 记录密码最后更新时间
4. 支持登录失败次数限制

## 管理员权限管理

### 角色定义
- **admin**: 普通管理员，具备基本管理权限
- **super_admin**: 超级管理员，具备所有权限，包括用户管理

### 权限控制建议
1. 管理员只能管理普通用户，不能管理其他管理员
2. 超级管理员可以管理所有用户（包括管理员）
3. 管理员操作需要记录操作日志
4. 支持管理员账号的启用/禁用

## 数据库表结构

### 管理员用户表 (admin_user)
```sql
CREATE TABLE `admin_user` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `admin_id` VARCHAR(64) UNIQUE NOT NULL,
  `username` VARCHAR(128) UNIQUE NOT NULL,
  `password` VARCHAR(512) NOT NULL,
  `real_name` VARCHAR(128) NOT NULL,
  `email` VARCHAR(128) UNIQUE NOT NULL,
  `role` VARCHAR(32) DEFAULT 'admin',
  `status` INT DEFAULT 1,
  -- 登录相关字段
  `last_login_time` TIMESTAMP NULL,
  `last_login_ip` VARCHAR(64),
  `login_fail_count` INT DEFAULT 0,
  `password_update_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  -- 时间戳
  `create_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `update_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `remark` VARCHAR(512)
);
```

### 管理员操作日志表 (admin_operation_log)
```sql
CREATE TABLE `admin_operation_log` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `admin_id` VARCHAR(64) NOT NULL,
  `operation_type` VARCHAR(32) NOT NULL,
  `operation_detail` TEXT,
  `ip_address` VARCHAR(64),
  `user_agent` VARCHAR(512),
  `operation_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 开发建议

### 接口设计
1. 管理员登录接口（仅账号密码）
2. 管理员密码更新接口
3. 管理员信息查询接口
4. 管理员操作日志查询接口

### 安全考虑
1. 管理员接口需要严格的权限验证
2. 操作日志需要详细记录
3. 支持登录失败次数限制和账号锁定
4. 密码需要强加密存储

### 业务逻辑
1. 管理员账号由系统管理员手动创建
2. 不支持管理员自助注册
3. 管理员账号状态控制代替删除功能
4. 支持管理员角色的灵活配置

## 部署说明

### 数据库初始化
系统启动时会自动创建管理员相关表结构，并插入默认管理员账号。

### 密码重置
如果忘记管理员密码，可以通过数据库直接重置：
```sql
UPDATE admin_user SET password = '加密后的新密码' WHERE username = 'admin';
```

### 监控建议
建议定期检查管理员操作日志，监控异常登录行为。