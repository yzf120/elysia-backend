# 基于 Redis 的短信验证码系统 - 快速开始

## ✅ 已完成的功能

### 1. **注册方式**
- ✅ 手机号 + 验证码注册（无需密码）

### 2. **登录方式**
- ✅ 手机号 + 验证码登录
- ✅ 手机号 + 密码登录

### 3. **核心特性**
- ✅ 使用 Redis 存储验证码（key: 手机号, value: 验证码）
- ✅ 验证码自动过期（5分钟）
- ✅ 验证码一次性使用
- ✅ 发送频率限制（60秒）
- ✅ 腾讯云短信集成

## 📦 新增文件列表

### 核心功能文件
```
client/redis_client.go                    # Redis 客户端
utils/sms_util.go                         # 腾讯云短信工具
utils/verification_code_redis.go         # 验证码服务（基于Redis）
router/auth.go                            # 认证路由处理器
```

### 测试和文档
```
test/auth_sms.http                        # HTTP 测试文件
docs/SMS_REDIS_GUIDE.md                   # 完整使用指南
```

### 修改的文件
```
service/auth.go                           # 添加验证码相关方法
router/router.go                          # 注册 auth 路由
main.go                                   # 初始化 Redis
.env.example                              # 添加 Redis 和短信配置
```

## 🚀 快速开始（3步）

### 1. 启动 Redis
```bash
docker run -d --name redis -p 6379:6379 redis:latest
```

### 2. 配置环境变量
复制 `.env.example` 为 `.env` 并填写配置：
```bash
# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# 腾讯云短信配置（从腾讯云控制台获取）
TENCENT_SMS_SECRET_ID=你的SecretId
TENCENT_SMS_SECRET_KEY=你的SecretKey
TENCENT_SMS_SDK_APP_ID=你的SDK_AppID
TENCENT_SMS_SIGN_NAME=你的签名名称
TENCENT_SMS_TEMPLATE_ID=你的模板ID
```

### 3. 安装依赖并启动
```bash
go get github.com/redis/go-redis/v9
go get github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111
go mod tidy
go run main.go
```

## 📡 API 接口

| 接口 | 方法 | 路径 | 功能 |
|------|------|------|------|
| 发送验证码 | POST | `/api/auth/send-code` | 发送注册/登录验证码 |
| 验证码注册 | POST | `/api/auth/register-sms` | 手机号+验证码注册 |
| 验证码登录 | POST | `/api/auth/login-sms` | 手机号+验证码登录 |
| 密码登录 | POST | `/api/auth/login-password` | 手机号+密码登录 |

## 🧪 测试

### 方式1: 使用 GoLand
打开 `test/auth_sms.http` 文件，点击运行按钮

### 方式2: 使用 curl
```bash
# 发送验证码
curl -X POST http://localhost:8080/api/auth/send-code \
  -H "Content-Type: application/json" \
  -d '{"phone_number":"13800138000","code_type":"register"}'

# 注册
curl -X POST http://localhost:8080/api/auth/register-sms \
  -H "Content-Type: application/json" \
  -d '{"phone_number":"13800138000","code":"123456"}'
```

## 🔑 Redis Key 设计

| Key 格式 | Value | TTL | 说明 |
|---------|-------|-----|------|
| `sms:code:register:{phone}` | 6位数字 | 5分钟 | 注册验证码 |
| `sms:code:login:{phone}` | 6位数字 | 5分钟 | 登录验证码 |
| `sms:frequency:{phone}` | "1" | 60秒 | 发送频率限制 |

### 查看 Redis 数据
```bash
redis-cli
KEYS sms:*
GET sms:code:register:13800138000
TTL sms:code:register:13800138000
```

## 📚 详细文档

完整的配置和使用指南请查看：[docs/SMS_REDIS_GUIDE.md](docs/SMS_REDIS_GUIDE.md)

包含内容：
- 腾讯云短信配置详细步骤
- Redis 安装和配置
- API 接口详细说明
- 测试场景和用例
- 常见问题解答
- 安全建议

## 🔒 安全特性

- ✅ 验证码6位随机数字
- ✅ 5分钟自动过期
- ✅ 一次性使用（使用后自动删除）
- ✅ 60秒发送频率限制
- ✅ 密码 bcrypt 加密
- ✅ 注册时检查手机号是否已存在
- ✅ 登录时检查用户状态

## ⚠️ 注意事项

1. **腾讯云短信**：签名和模板必须审核通过才能发送
2. **Redis**：必须启动 Redis 服务
3. **环境变量**：必须配置 `.env` 文件
4. **测试手机号**：测试阶段建议使用测试手机号

## 🎯 技术架构

```
HTTP 请求
   ↓
Router (router/auth.go)
   ↓
Service (service/auth.go)
   ↓
├─→ Redis (验证码存储)
│   └─→ Key: sms:code:{type}:{phone}
│       Value: 验证码
│       TTL: 5分钟
│
└─→ 腾讯云短信 (utils/sms_util.go)
    └─→ 发送验证码短信
```

## 📊 优势

使用 Redis 替代数据库存储验证码的优势：

| 特性 | Redis | 数据库 |
|------|-------|--------|
| 性能 | ⚡ 极快（内存） | 🐌 较慢（磁盘） |
| 过期机制 | ✅ 自动过期 | ❌ 需要定时清理 |
| 代码复杂度 | ✅ 简单 | ❌ 需要 DAO 层 |
| 并发性能 | ✅ 高 | ⚠️ 中等 |
| 数据持久化 | ⚠️ 可选 | ✅ 持久化 |

## 🚀 下一步优化建议

1. **JWT Token**: 替代简单的随机 Token
2. **图形验证码**: 防止机器人攻击
3. **IP 限流**: 防止恶意攻击
4. **Redis 集群**: 提高可用性
5. **监控告警**: 接入监控系统

---

**版本**: v1.0  
**日期**: 2025-12-29  
**状态**: ✅ 已完成并测试通过
