# Proto 文件管理说明

## 概述

本目录包含项目的 Protocol Buffers 定义文件，用于定义 API 接口和消息结构。

## 文件结构

```
proto/
├── Makefile              # 代码生成脚本
├── README.md             # 说明文档
├── admin/                # 管理员用户相关
│   ├── admin.proto       # 管理员用户服务定义
│   ├── admin.pb.go       # 生成的 protobuf 代码
│   └── admin.trpc.go     # 生成的 tRPC 代码
├── auth/                 # 认证相关
│   ├── auth.proto
│   ├── auth.pb.go
│   └── auth.trpc.go
├── user/                 # 用户相关
│   ├── user.proto
│   ├── user.pb.go
│   └── user.trpc.go
└── ...                  # 其他服务
```

## 新添加的 Admin 服务

### 功能描述
- **管理员用户管理**：独立的管理员账号体系，与普通用户分离
- **账号密码登录**：仅支持账号密码登录，不支持验证码登录
- **权限管理**：支持角色权限控制（admin/super_admin）
- **操作日志**：记录管理员操作历史

### 主要接口

1. **管理员用户登录**

   rpc LoginAdminUser(LoginAdminUserRequest) returns (LoginAdminUserResponse);
   ```

2. **创建管理员用户**

   rpc CreateAdminUser(CreateAdminUserRequest) returns (CreateAdminUserResponse);
   ```

3. **查询管理员用户列表**
   rpc ListAdminUsers(ListAdminUsersRequest) returns (ListAdminUsersResponse);
   ```

4. **更新管理员用户密码**
   rpc UpdateAdminUserPassword(UpdateAdminUserPasswordRequest) returns (UpdateAdminUserPasswordResponse);
   ```

## 代码生成

### 生成所有 proto 代码
```bash
cd proto
make generate-all
```

### 清理生成的代码
```bash
cd proto
make clean
```

### 手动生成单个 proto 文件
```bash
# 生成 admin.proto
cd proto
trpc create \
    -p admin/admin.proto \
    -o admin \
    --rpconly \
    --nogomod \
    --mock=false
```

## 依赖要求

- **trpc-cmdline**: tRPC 命令行工具
- **protoc**: Protocol Buffers 编译器
- **protoc-gen-go**: Go 语言插件

## 开发流程

### 1. 修改 proto 文件
编辑相应的 `.proto` 文件，添加新的接口或修改现有定义。

### 2. 生成代码
运行 `make generate-all` 重新生成所有 Go 代码。

### 3. 实现服务
在 `service_impl/` 目录下实现对应的服务接口。

### 4. 注册服务
在服务注册文件中添加新的服务实例。

## 注意事项

1. **保持一致性**：修改 proto 文件后必须重新生成代码
2. **版本控制**：生成的 `.pb.go` 和 `.trpc.go` 文件不应手动修改
3. **接口兼容性**：修改现有接口时要考虑向后兼容性
4. **命名规范**：遵循 protobuf 的命名约定

## 常见问题

### Q: 为什么生成的代码没有更新？
A: 确保运行了 `make generate-all`，并且没有语法错误。

### Q: 如何添加新的 proto 文件？
A: 1. 创建新的 `.proto` 文件
   2. 在 `Makefile` 的 `generate-all` 目标中添加生成规则
   3. 运行 `make generate-all`

### Q: 生成的代码有编译错误？
A: 检查 proto 文件语法，确保所有消息和服务定义正确。

## 文件说明

### admin.proto
管理员用户服务定义，包含：
- 管理员用户登录接口
- 管理员用户管理接口
- 管理员操作日志记录

### auth.proto
认证服务定义，包含：
- 普通用户注册/登录接口
- 管理员用户登录接口（通过 auth 服务统一入口）
- 认证令牌管理

## 版本历史

- **v1.0.0**: 初始版本，包含基础服务定义
- **v1.1.0**: 添加管理员用户服务（admin.proto）

## 相关文档

- [Protocol Buffers 官方文档](https://developers.google.com/protocol-buffers)
- [tRPC 框架文档](https://github.com/trpc-group/trpc-go)
- [项目 API 文档](../docs/api.md)