# 基于 Redis 的短信验证码系统使用指南

## 📋 目录
- [系统架构](#系统架构)
- [快速开始](#快速开始)
- [API 接口说明](#api-接口说明)
- [Redis Key 设计](#redis-key-设计)
- [配置说明](#配置说明)
- [测试指南](#测试指南)
- [常见问题](#常见问题)

## 🏗️ 系统架构

### 技术栈
- **Redis**: 存储验证码，支持自动过期
- **腾讯云短信**: 发送短信验证码
- **Go**: 后端服务
- **tRPC-Go**: RPC 框架

### 数据流程
```
用户请求 → Router → Service → Redis + 腾讯云短信
                              ↓
                         验证码存储（5分钟）
                              ↓
                         用户验证 → 删除验证码
```

## 🚀 快速开始

### 1. 安装 Redis

#### 使用 Docker（推荐）
```bash
docker run -d --name redis -p 6379:6379 redis:latest
```

#### 使用 Homebrew（macOS）
```bash
brew install redis
brew services start redis
```

#### 验证 Redis 运行
```bash
redis-cli ping
# 应该返回: PONG
```

### 2. 配置腾讯云短信

#### 2.1 登录腾讯云控制台
访问：https://console.cloud.tencent.com/smsv2

#### 2.2 创建短信应用
1. 进入"应用管理" → 点击"创建应用"
2. 填写应用名称，如"Elysia Backend"
3. 记录 **SDK AppID**

#### 2.3 创建短信签名
1. 进入"国内短信" → "签名管理" → "创建签名"
2. 签名类型：选择"网站"
3. 签名用途：选择"自用"
4. 提交审核（通常1-2小时）
5. 记录 **签名名称**

#### 2.4 创建短信模板
1. 进入"国内短信" → "正文模板管理" → "创建正文模板"
2. 模板名称：验证码短信
3. 短信类型：验证码
4. 短信内容示例：
   ```
   您的验证码是{1}，{2}分钟内有效，请勿泄露给他人。
   ```
5. 提交审核（通常1-2小时）
6. 记录 **模板ID**

#### 2.5 获取 API 密钥
1. 访问：https://console.cloud.tencent.com/cam/capi
2. 点击"新建密钥"
3. 记录 **SecretId** 和 **SecretKey**

### 3. 配置环境变量

复制 `.env.example` 为 `.env`：
```bash
cp .env.example .env
```

编辑 `.env` 文件：
```bash
# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# 腾讯云短信配置
TENCENT_SMS_SECRET_ID=你的SecretId
TENCENT_SMS_SECRET_KEY=你的SecretKey
TENCENT_SMS_SDK_APP_ID=你的SDK_AppID
TENCENT_SMS_SIGN_NAME=你的签名名称
TENCENT_SMS_TEMPLATE_ID=你的模板ID
TENCENT_SMS_REGION=ap-guangzhou
```

### 4. 安装依赖

```bash
go get github.com/redis/go-redis/v9
go get github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111
go mod tidy
```

### 5. 启动服务

```bash
go run main.go
```

服务启动后，你应该看到：
```
Redis客户端初始化成功: localhost:6379 (DB: 0)
数据库初始化成功
服务器启动成功，监听端口: 8080
```

## 📡 API 接口说明

### 1. 发送验证码

**接口**: `POST /api/auth/send-code`

**请求参数**:
```json
{
  "phone_number": "13800138000",
  "code_type": "register"
}
```

**响应示例**:
```json
{
  "data": {
    "message": "验证码发送成功"
  },
  "error": null
}
```

**错误响应**:
```json
{
  "data": null,
  "error": {
    "code": 400,
    "message": "发送过于频繁，请59秒后再试"
  }
}
```

### 2. 手机号+验证码注册

**接口**: `POST /api/auth/register-sms`

**请求参数**:
```json
{
  "phone_number": "13800138000",
  "code": "123456"
}
```

**响应示例**:
```json
{
  "data": {
    "user_info": {
      "user_id": "usr_xxx",
      "user_name": "user_13800138000",
      "phone_number": "13800138000",
      "status": 2,
      "create_time": "2025-12-29 20:00:00"
    },
    "token": "abc123..."
  },
  "error": null
}
```

### 3. 手机号+验证码登录

**接口**: `POST /api/auth/login-sms`

**请求参数**:
```json
{
  "phone_number": "13800138000",
  "code": "123456"
}
```

**响应示例**: 同注册接口

### 4. 手机号+密码登录

**接口**: `POST /api/auth/login-password`

**请求参数**:
```json
{
  "phone_number": "13800138000",
  "password": "your_password"
}
```

**响应示例**: 同注册接口

## 🔑 Redis Key 设计

### Key 格式

| 用途 | Key 格式 | Value | TTL |
|------|---------|-------|-----|
| 注册验证码 | `sms:code:register:{phone}` | 6位数字 | 5分钟 |
| 登录验证码 | `sms:code:login:{phone}` | 6位数字 | 5分钟 |
| 发送频率限制 | `sms:frequency:{phone}` | "1" | 60秒 |

### 示例

```
sms:code:register:13800138000 = "123456" (TTL: 300秒)
sms:code:login:13800138000 = "654321" (TTL: 300秒)
sms:frequency:13800138000 = "1" (TTL: 60秒)
```

### 查看 Redis 数据

```bash
# 连接 Redis
redis-cli

# 查看所有验证码相关的 key
KEYS sms:*

# 查看某个 key 的值
GET sms:code:register:13800138000

# 查看某个 key 的剩余过期时间（秒）
TTL sms:code:register:13800138000

# 删除某个 key
DEL sms:code:register:13800138000
```

## ⚙️ 配置说明

### Redis 配置

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| REDIS_HOST | Redis 主机地址 | localhost |
| REDIS_PORT | Redis 端口 | 6379 |
| REDIS_PASSWORD | Redis 密码 | 空 |
| REDIS_DB | Redis 数据库编号 | 0 |

### 腾讯云短信配置

| 配置项 | 说明 | 获取方式 |
|--------|------|----------|
| TENCENT_SMS_SECRET_ID | API 密钥 ID | 控制台 → 访问管理 → API密钥 |
| TENCENT_SMS_SECRET_KEY | API 密钥 Key | 同上 |
| TENCENT_SMS_SDK_APP_ID | 短信应用 ID | 控制台 → 短信 → 应用管理 |
| TENCENT_SMS_SIGN_NAME | 短信签名 | 控制台 → 短信 → 签名管理 |
| TENCENT_SMS_TEMPLATE_ID | 短信模板 ID | 控制台 → 短信 → 正文模板 |
| TENCENT_SMS_REGION | 地域 | ap-guangzhou |

## 🧪 测试指南

### 使用 GoLand 测试

1. 打开 `test/auth_sms.http` 文件
2. 点击请求旁边的绿色运行按钮
3. 查看响应结果

### 使用 curl 测试

```bash
# 1. 发送注册验证码
curl -X POST http://localhost:8080/api/auth/send-code \
  -H "Content-Type: application/json" \
  -d '{"phone_number":"13800138000","code_type":"register"}'

# 2. 使用验证码注册（替换为实际验证码）
curl -X POST http://localhost:8080/api/auth/register-sms \
  -H "Content-Type: application/json" \
  -d '{"phone_number":"13800138000","code":"123456"}'

# 3. 发送登录验证码
curl -X POST http://localhost:8080/api/auth/send-code \
  -H "Content-Type: application/json" \
  -d '{"phone_number":"13800138000","code_type":"login"}'

# 4. 使用验证码登录
curl -X POST http://localhost:8080/api/auth/login-sms \
  -H "Content-Type: application/json" \
  -d '{"phone_number":"13800138000","code":"123456"}'
```

### 测试场景

#### ✅ 正常场景
1. 发送注册验证码 → 收到短信 → 使用验证码注册成功
2. 发送登录验证码 → 收到短信 → 使用验证码登录成功
3. 使用密码登录成功

#### ❌ 异常场景
1. 手机号格式错误 → 返回错误
2. 验证码类型错误 → 返回错误
3. 验证码错误 → 返回错误
4. 验证码过期（5分钟后）→ 返回错误
5. 重复注册 → 返回错误
6. 未注册用户登录 → 返回错误
7. 频繁发送验证码（60秒内）→ 返回错误

## ❓ 常见问题

### 1. Redis 连接失败

**错误**: `Redis连接失败: dial tcp 127.0.0.1:6379: connect: connection refused`

**解决方案**:
```bash
# 检查 Redis 是否运行
redis-cli ping

# 如果没有运行，启动 Redis
# Docker
docker start redis

# Homebrew
brew services start redis
```

### 2. 短信发送失败

**错误**: `发送短信失败: ...`

**可能原因**:
1. 签名或模板未审核通过
2. API 密钥错误
3. 账户余额不足
4. 手机号格式错误

**解决方案**:
1. 检查腾讯云控制台，确认签名和模板已审核通过
2. 检查 `.env` 文件中的配置是否正确
3. 检查腾讯云账户余额
4. 确保手机号是11位数字

### 3. 验证码过期

**问题**: 验证码5分钟后自动过期

**解决方案**: 这是正常行为，用户需要重新发送验证码

### 4. 频繁发送限制

**问题**: 60秒内只能发送一次验证码

**解决方案**: 
- 这是防止短信轰炸的安全机制
- 如需调整，修改 `service/auth.go` 中的 `60*time.Second`

### 5. 验证码一次性使用

**问题**: 验证码使用后立即失效

**解决方案**: 这是安全机制，每个验证码只能使用一次

## 🔒 安全建议

1. **生产环境**:
   - 使用 HTTPS
   - 设置 Redis 密码
   - 限制 API 访问频率
   - 使用 JWT 替代简单 Token

2. **验证码安全**:
   - 验证码长度：6位数字
   - 有效期：5分钟
   - 一次性使用
   - 发送频率限制：60秒

3. **密钥安全**:
   - 不要将 `.env` 文件提交到 Git
   - 定期更换 API 密钥
   - 使用环境变量管理敏感信息

## 📊 监控建议

### Redis 监控

```bash
# 查看 Redis 信息
redis-cli INFO

# 查看内存使用
redis-cli INFO memory

# 查看连接数
redis-cli INFO clients

# 实时监控命令
redis-cli MONITOR
```

### 日志监控

关注以下日志：
- Redis 连接成功/失败
- 短信发送成功/失败
- 验证码验证成功/失败
- 频繁发送告警

## 🎯 下一步优化

1. **JWT Token**: 替代简单的随机 Token
2. **图形验证码**: 防止机器人攻击
3. **IP 限流**: 防止恶意攻击
4. **短信模板优化**: 支持多种场景
5. **监控告警**: 接入监控系统
6. **日志分析**: 统计发送成功率

---

**文档版本**: v1.0  
**最后更新**: 2025-12-29  
**维护者**: Elysia Backend Team
