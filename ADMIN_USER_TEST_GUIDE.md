# 管理员用户功能测试指南

## 功能概述

管理员用户功能为系统提供了专门的管理员账号管理能力，与普通用户账号分离，具有以下特点：

- **独立账号体系**：管理员账号与普通用户账号完全分离
- **简化字段**：减少普通用户相关字段，专注于管理员核心功能
- **账号密码登录**：仅支持账号密码登录，不支持验证码登录
- **账号不限制**：管理员账号不一定是手机号，可以是任意字符串
- **权限控制**：支持角色权限管理（admin/super_admin）
- **操作日志**：记录管理员操作历史

## 数据库表结构

### 管理员用户表 (admin_user)
```sql
CREATE TABLE `admin_user` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `admin_id` varchar(64) NOT NULL COMMENT '管理员id',
  `username` varchar(128) NOT NULL COMMENT '管理员账号',
  `password` varchar(512) NOT NULL COMMENT '密码',
  `real_name` varchar(128) NOT NULL COMMENT '真实姓名',
  `email` varchar(128) NOT NULL COMMENT '邮箱',
  `role` varchar(32) NOT NULL DEFAULT 'admin' COMMENT '角色',
  `status` int NOT NULL DEFAULT '1' COMMENT '状态',
  `last_login_time` timestamp NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(64) COMMENT '最后登录IP',
  `login_fail_count` int NOT NULL DEFAULT '0' COMMENT '登录失败次数',
  `password_update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '密码最后更新时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
  `remark` varchar(512) COMMENT '备注'
);
```

### 管理员操作日志表 (admin_operation_log)
```sql
CREATE TABLE `admin_operation_log` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `admin_id` varchar(64) NOT NULL COMMENT '管理员id',
  `operation_type` varchar(32) NOT NULL COMMENT '操作类型',
  `operation_detail` text COMMENT '操作详情',
  `ip_address` varchar(64) COMMENT '操作IP',
  `user_agent` varchar(512) COMMENT '用户代理',
  `operation_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间'
);
```

## 默认管理员账号

系统初始化时会自动创建默认超级管理员账号：

- **账号**: `admin`
- **密码**: `admin123`
- **角色**: `super_admin`
- **邮箱**: `admin@elysia.com`
- **状态**: 启用 (1)

## API 接口测试

### 1. 管理员用户登录

**接口**: `POST /api/admin/login`

**请求示例**:
```json
{
  "username": "admin",
  "password": "admin123",
  "ip_address": "192.168.1.100"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "登录成功",
  "admin_user": {
    "id": 1,
    "admin_id": "admin_001",
    "username": "admin",
    "real_name": "系统管理员",
    "email": "admin@elysia.com",
    "role": "super_admin",
    "status": 1,
    "last_login_time": "2024-01-10 15:34:30",
    "last_login_ip": "192.168.1.100",
    "create_time": "2024-01-10 10:00:00"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 2. 创建管理员用户

**接口**: `POST /api/admin/users`

**请求示例**:
```json
{
  "username": "test_admin",
  "password": "Test@123",
  "real_name": "测试管理员",
  "email": "test@elysia.com",
  "role": "admin",
  "remark": "测试账号"
}
```

### 3. 查询管理员用户列表

**接口**: `GET /api/admin/users?page=1&page_size=10&role=admin&status=1`

### 4. 更新管理员用户密码

**接口**: `PUT /api/admin/users/{admin_id}/password`

**请求示例**:
```json
{
  "old_password": "admin123",
  "new_password": "NewAdmin@123"
}
```

## 数据库测试

### 1. 连接数据库
```bash
mysql -h localhost -u root -p elysia_db
```

### 2. 执行测试脚本
```sql
-- 查看表结构
DESCRIBE admin_user;

-- 查询默认管理员账号
SELECT * FROM admin_user WHERE username = 'admin';

-- 测试管理员用户列表查询
SELECT admin_id, username, real_name, email, role, status, create_time 
FROM admin_user 
ORDER BY create_time DESC;

-- 查看操作日志
SELECT * FROM admin_operation_log ORDER BY operation_time DESC LIMIT 10;
```

### 3. 手动测试数据操作
```sql
-- 创建测试管理员用户
INSERT INTO admin_user (
  admin_id, username, password, real_name, email, role, status, remark
) VALUES (
  'admin_test_001', 
  'test_admin', 
  '$2a$10$r3z7t8v9w0x1y2z3a4b5c6d7e8f9g0h1i2j3k4l5m6n7o8p9q0r1s2t3u4',
  '测试管理员',
  'test@elysia.com',
  'admin',
  1,
  '测试账号'
);

-- 验证数据插入
SELECT * FROM admin_user WHERE username = 'test_admin';

-- 清理测试数据
DELETE FROM admin_user WHERE username = 'test_admin';
```

## 代码层测试

### 1. 服务层测试

// 测试管理员用户登录
//func TestAdminUserLogin(t *testing.T) {
//    service := NewAdminUserService()
//    
//    // 测试成功登录
//    adminUserInfo, token, err := service.LoginAdminUser(
//        context.Background(), 
//        "admin", 
//        "admin123", 
//        "127.0.0.1"
//    )
//    
//    assert.NoError(t, err)
//    assert.NotNil(t, adminUserInfo)
//    assert.NotEmpty(t, token)
//    
//    // 测试密码错误
//    _, _, err = service.LoginAdminUser(
//        context.Background(), 
//        "admin", 
//        "wrongpassword", 
//        "127.0.0.1"
//    )
//    
//    assert.Error(t, err)
//}
```

### 2. DAO层测试

// 测试管理员用户查询
//func TestGetAdminUserByUsername(t *testing.T) {
//    dao := NewAdminUserDAO()
//    
//    adminUser, password, err := dao.GetAdminUserByUsername("admin")
//    
//    assert.NoError(t, err)
//    assert.NotNil(t, adminUser)
//    assert.NotEmpty(t, password)
//    assert.Equal(t, "admin", adminUser.Username)
//}
```

## 安全注意事项

1. **密码安全**：首次登录后应立即修改默认密码
2. **账号权限**：根据实际需求分配适当的角色权限
3. **操作日志**：定期检查操作日志，监控异常行为
4. **账号状态**：及时禁用不再使用的管理员账号
5. **密码策略**：强制使用强密码策略

## 故障排查

### 常见问题

1. **登录失败**：检查账号状态、密码是否正确
2. **账号不存在**：确认用户名拼写正确
3. **权限不足**：检查用户角色权限设置
4. **数据库连接失败**：检查数据库服务是否正常启动

### 日志查看

查看应用日志和数据库日志，定位问题：
```bash
# 查看应用日志
tail -f /var/log/elysia-backend/app.log

# 查看数据库日志
tail -f /var/log/mysql/mysql.log
```

## 性能优化建议

1. **索引优化**：确保关键字段（username, email, admin_id）有索引
2. **查询优化**：避免全表扫描，使用分页查询
3. **缓存策略**：对频繁查询的管理员信息进行缓存
4. **日志清理**：定期清理历史操作日志

## 部署检查清单

- [ ] 数据库表结构已创建
- [ ] 默认管理员账号已初始化
- [ ] API接口已正确配置
- [ ] 权限验证中间件已启用
- [ ] 操作日志记录功能正常
- [ ] 密码加密功能正常
- [ ] 登录失败次数限制生效
- [ ] 账号状态管理功能正常