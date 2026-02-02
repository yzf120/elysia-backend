# User相关代码删除总结

## 删除原因

项目已经使用 `student`、`teacher`、`admin` 三个独立的表和认证系统，不再需要统一的 `user` 表和相关逻辑。

## 删除的文件清单

### 1. Router层
- ✅ `router/user.go` - 用户路由处理器
- ✅ `router/login.go` - 旧的登录路由
- ✅ `router/router.go` 中的 `registerUserRoutes` 函数

### 2. Service层
- ✅ `service/user_service.go` - 用户服务
- ✅ `service/auth_service.go` - 旧的认证服务（已被student_auth_service、teacher_auth_service、admin_auth_service替代）

### 3. Service实现层
- ✅ `service_impl/user_service_impl.go` - 用户服务实现

### 4. DAO层
- ✅ `dao/user_dao.go` - 用户数据访问对象

### 5. Model层
- ✅ `model/user.go` - 用户模型

### 6. Proto层
- ✅ `proto/user/` - 整个user proto目录
  - `user.proto`
  - `user.pb.go`
  - `user.trpc.go`

### 7. 测试文件
- ✅ `test/user_test.go` - 用户测试文件

### 8. SQL文件
- ✅ `sql/user.sql` - 用户表SQL
- ✅ `sql/init.sql` 中的user表定义

### 9. 其他修改
- ✅ `main.go` - 删除user proto导入和RPC服务注册

## 保留的文件

以下文件保留，因为它们是新的认证系统的一部分：

### Student相关
- `service/student_service.go`
- `service/student_auth_service.go`
- `dao/student_dao.go`
- `model/student/`
- `router/student.go`

### Teacher相关
- `service/teacher_service.go`
- `service/teacher_auth_service.go`
- `dao/teacher_dao.go`
- `model/teacher/`
- `router/teacher.go`

### Admin相关
- `service/admin_user_service.go`
- `service/admin_auth_service.go`
- `dao/admin_user_dao.go`
- `model/admin/`
- `router/admin.go`

## 新的认证架构

### 路由结构
```
/api/student/auth/
  - POST /send-code (发送验证码)
  - POST /register-sms (短信注册)
  - POST /login-sms (短信登录)
  - POST /login-password (密码登录)
  - POST /logout (登出)

/api/teacher/auth/
  - POST /send-code (发送验证码)
  - POST /register-sms (短信注册)
  - POST /login-sms (短信登录)
  - POST /login-password (密码登录)
  - POST /logout (登出)

/api/admin/auth/
  - POST /send-code (发送验证码)
  - POST /register-sms (短信注册)
  - POST /login-sms (短信登录)
  - POST /login-password (密码登录)
  - POST /login-username (用户名密码登录)
  - POST /logout (登出)
```

### 数据库表结构
```
student (学生表)
  - student_id (主键)
  - phone_number (手机号，唯一)
  - password (密码)
  - chinese_name (姓名)
  - ...

teacher (教师表)
  - teacher_id (主键)
  - phone_number (手机号，唯一)
  - password (密码)
  - chinese_name (姓名)
  - employee_number (工号)
  - ...

admin_user (管理员表)
  - admin_id (主键)
  - username (用户名，唯一)
  - password (密码)
  - real_name (真实姓名)
  - ...
```

## 优势

1. **职责分离**：每个角色有独立的表和服务，职责更清晰
2. **灵活扩展**：可以为不同角色添加特定字段和功能
3. **安全性提升**：不同角色的认证逻辑独立，更容易管理权限
4. **代码可维护性**：避免了大量的if-else判断用户类型
5. **性能优化**：查询时不需要关联多个表

## 迁移注意事项

如果数据库中已经有user表的数据，需要：

1. **备份数据**：
```sql
CREATE TABLE `user_backup` AS SELECT * FROM `user`;
```

2. **迁移数据到对应表**：
```sql
-- 迁移学生数据
INSERT INTO student (student_id, phone_number, password, chinese_name, ...)
SELECT user_id, phone_number, password, chinese_name, ...
FROM user WHERE user_type = 'student';

-- 迁移教师数据
INSERT INTO teacher (teacher_id, phone_number, password, chinese_name, ...)
SELECT user_id, phone_number, password, chinese_name, ...
FROM user WHERE user_type = 'teacher';

-- 迁移管理员数据
INSERT INTO admin_user (admin_id, username, password, real_name, ...)
SELECT user_id, phone_number, password, chinese_name, ...
FROM user WHERE user_type = 'admin';
```

3. **删除user表**（确认数据迁移成功后）：
```sql
DROP TABLE IF EXISTS `user`;
```

## 验证

可以通过以下命令验证是否还有user相关的引用：

```bash
# 查找user相关的导入
grep -r "proto/user" --include="*.go" .

# 查找UserService的引用
grep -r "UserService\|user_service\|userService" --include="*.go" .

# 查找user表的SQL
grep -r "CREATE TABLE.*user\|INSERT INTO.*user" --include="*.sql" .
```

## 完成时间

2026-02-01

## 相关文档

- [常量重构文档](./CONSTANTS_REFACTOR.md)
- [管理员用户指南](./ADMIN_USER_GUIDE.md)
