# Elysia 后端系统 API 测试指南

## 文件概览

本目录下为您生成了以下测试文件：

1. **`elysia-api-postman-collection.json`** - Postman 集合文件
   - 包含所有 33 个 HTTP 接口的完整定义
   - 支持环境变量和动态参数
   - 可直接导入 Postman 使用

2. **`elysia-api-curl-commands.txt`** - cURL 命令文件
   - 包含所有接口的 cURL 命令
   - 每个命令都有详细注释
   - 可直接在终端中复制使用

3. **`test_http_apis.py`** - Python 自动化测试脚本
   - 自动测试所有接口
   - 提供测试结果统计
   - 支持顺序测试和错误处理

## 快速开始

### 方法一：使用 Postman（推荐）

1. 打开 Postman
2. 点击 "Import" 按钮
3. 选择 `elysia-api-postman-collection.json` 文件
4. 导入后，您将看到一个名为 "Elysia 后端系统 API 集合" 的集合
5. 按以下步骤测试：
   - 先运行管理员/学生/教师的登录接口获取 token
   - 在环境变量中设置获取到的 token 和 ID
   - 然后测试受保护接口
   - 最后测试登出接口

### 方法二：使用 cURL 命令

1. 打开终端
2. 进入项目目录：`cd /Users/sylvainyang/project/elysia/elysia-backend`
3. 查看 cURL 命令：`cat elysia-api-curl-commands.txt`
4. 按顺序执行命令，注意先获取 token 再测试受保护接口

### 方法三：使用 Python 测试脚本

1. 确保已安装 Python 3 和 requests 库：`pip install requests`
2. 启动您的服务器（默认在 http://localhost:8001）
3. 运行测试脚本：`python3 test_http_apis.py`
4. 查看详细的测试结果

## 接口分类

### 管理员接口 (10个接口)
**公共接口（无需认证）：**
- 发送注册验证码
- 发送登录验证码  
- 短信注册
- 短信登录
- 密码登录

**受保护接口（需要认证）：**
- 创建管理员
- 获取管理员信息
- 查询管理员列表
- 更新密码
- 更新状态
- 管理员登出

### 学生接口 (10个接口)
**公共接口（无需认证）：**
- 发送注册验证码
- 发送登录验证码
- 短信注册
- 短信登录
- 密码登录

**受保护接口（需要认证）：**
- 创建学生
- 获取学生信息
- 更新学生信息
- 查询学生列表
- 更新学习进度
- 学生登出

### 教师接口 (9个接口)
**公共接口（无需认证）：**
- 教师注册（注意：在 `/register` 路径）
- 发送登录验证码
- 短信登录
- 密码登录

**受保护接口（需要认证）：**
- 获取教师信息
- 更新教师信息
- 查询教师列表
- 审核教师（需要管理员权限）
- 教师登出

### 教师审批接口 (4个接口)
**教师端（需要教师认证）：**
- 获取教师的审批单
- 获取审批单详情

**管理员端（需要管理员认证）：**
- 查询审批单列表
- 审批教师
- 删除审批单

## 环境变量

在 Postman 中，您可以使用以下环境变量：

| 变量名 | 说明 | 如何获取 |
|--------|------|----------|
| `base_url` | 服务器地址 | 默认为 `http://localhost:8001` |
| `admin_token` | 管理员认证 token | 从管理员登录接口获取 |
| `student_token` | 学生认证 token | 从学生登录接口获取 |
| `teacher_token` | 教师认证 token | 从教师登录接口获取 |
| `admin_id` | 管理员 ID | 从登录响应中获取 |
| `student_id` | 学生 ID | 从登录响应中获取 |
| `teacher_id` | 教师 ID | 从登录响应中获取 |
| `approval_id` | 审批单 ID | 从教师审批接口获取 |

## 测试数据

默认使用以下测试数据，您可以根据实际情况修改：

- **手机号**：`13800138000`
- **验证码**：`123456`（开发环境可能不需要真实验证）
- **密码**：`Test@123456`
- **学号**：`S20240001`
- **工号**：`T20240001`
- **默认管理员**：手机号 `admin`，密码 `Test@123456`

## 测试顺序建议

1. **先测试公共接口**：
   - 管理员、学生、教师的登录接口
   - 获取 token 和用户 ID

2. **再测试受保护接口**：
   - 使用获取到的 token 测试管理接口
   - 注意角色权限：某些接口需要管理员权限

3. **最后测试登出接口**：
   - 使 token 失效
   - 验证登出后接口无法访问

4. **特殊接口测试**：
   - 教师注册接口是独立的（`/register`）
   - 教师审批接口有教师端和管理员端之分

## 常见问题

### 1. 服务器连接失败
- 确认服务器已启动：`http://localhost:8001`
- 检查 `trpc_go.yaml` 中的端口配置
- 如果使用其他地址，请修改所有命令中的 `localhost:8001`

### 2. 认证失败
- 确认 token 已正确设置
- 检查 token 是否过期（登出后 token 会失效）
- 确认请求头中包含了正确的 `Authorization: Bearer <token>`

### 3. 参数错误
- 检查请求体 JSON 格式是否正确
- 确认必填参数都已提供
- 对于 GET 请求，参数需要通过查询字符串传递

### 4. Python 测试脚本问题
- 确保已安装 requests 库：`pip install requests`
- 如果服务器地址不同，修改脚本中的 `BASE_URL`
- 检查 Python 版本：需要 Python 3.6 或更高版本

## 开发调试

### 使用 cURL 调试
```bash
# 显示详细请求信息
curl -v -X POST http://localhost:8001/api/admin/auth/login-password ...

# 只显示响应头
curl -i -X POST http://localhost:8001/api/admin/auth/login-password ...

# 将响应保存到文件
curl -o response.json -X POST http://localhost:8001/api/admin/auth/login-password ...
```

### 使用 Postman 测试环境
- 可以创建不同的环境（开发、测试、生产）
- 使用预请求脚本来自动获取 token
- 使用测试脚本来验证响应

## 接口更新维护

当项目接口有变更时：

1. 更新 `elysia-api-postman-collection.json` 中的接口定义
2. 更新 `elysia-api-curl-commands.txt` 中的 cURL 命令
3. 更新 `test_http_apis.py` 中的测试函数
4. 更新本 README 文档

## 许可证

这些测试资源仅供 Elysia 项目开发使用。

## 支持

如有问题，请参考项目代码或联系开发团队。

---
*最后更新：2026-02-02*
*共包含 33 个 HTTP 接口的完整测试资源*