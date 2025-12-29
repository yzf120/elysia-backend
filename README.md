# Elysia Backend - tRPC-Go 微服务

基于 tRPC-Go 框架构建的微服务后端项目。

## 项目结构

```
elysia-backend/
├── main.go                 # 主入口文件
├── trpc_go.yaml           # tRPC 配置文件
├── go.mod                 # Go 模块依赖
├── dao/                   # 数据访问层
│   ├── dao.go
│   └── client.go
├── service/               # 业务逻辑层
│   └── service.go
├── handler/               # tRPC 服务处理器
│   ├── conversation/      # 会话服务处理器
│   │   └── conversation.go
│   └── agent/            # 智能体服务处理器
│       └── agent.go
├── proto/                # Protobuf 定义和生成代码
│   ├── conversation/
│   │   ├── conversation.proto
│   │   ├── conversation.pb.go
│   │   └── conversation.trpc.go
│   ├── agent/
│   │   ├── agent.proto
│   │   ├── agent.pb.go
│   │   └── agent.trpc.go
│   └── helloworld/
│       ├── helloworld.proto
│       ├── helloworld.pb.go
│       └── helloworld.trpc.go
└── config/               # 配置管理
    └── config.go
```

## 快速开始

### 1. 安装依赖

```bash
# 下载项目依赖
go mod tidy

# 安装 protoc 编译器（如果还没安装）
# macOS
brew install protobuf

# 安装 protoc-gen-go 和 trpc 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install trpc.group/trpc-go/trpc-cmdline/protoc-gen-go-trpc@latest
```

### 2. 配置环境变量

创建 `.env` 文件：

```bash
cp .env.example .env
```

编辑 `.env` 文件，配置数据库连接等信息：

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=elysia
```

### 3. 生成 Protobuf 代码

```bash
# 进入 proto 目录
cd proto

# 生成所有 proto 文件
make generate-all

# 或者单独生成某个服务
make generate-conversation
make generate-agent
```

### 4. 运行服务

```bash
# 返回项目根目录
cd ..

# 运行服务
go run main.go

# 或者先编译再运行
go build -o elysia-backend
./elysia-backend
```

服务将在 `127.0.0.1:8000` 启动。

## 服务说明

### Conversation Service (会话服务)

提供会话管理功能：

- `CreateConversation` - 创建新会话
- `GetConversation` - 获取会话详情
- `SendMessage` - 发送消息

### Agent Service (智能体服务)

提供智能体管理功能：

- `CreateAgent` - 创建智能体
- `GetAgent` - 获取智能体详情
- `UpdateAgent` - 更新智能体
- `DeleteAgent` - 删除智能体
- `ExecuteAgent` - 执行智能体

## 开发指南

### 添加新服务

1. 在 `proto/` 目录下创建新的 `.proto` 文件
2. 定义服务接口和消息类型
3. 运行 `make generate-xxx` 生成代码
4. 在 `handler/` 目录下实现服务接口
5. 在 `main.go` 中注册服务

### Proto 文件示例

```protobuf
syntax = "proto3";

package trpc.myservice;
option go_package = "github.com/yzf120/elysia-backend/proto/myservice";

service MyService {
  rpc MyMethod (MyRequest) returns (MyResponse) {}
}

message MyRequest {
  string id = 1;
}

message MyResponse {
  string result = 1;
}
```


## 配置说明

### trpc_go.yaml

tRPC 框架的主配置文件，包含：

- 服务监听地址和端口
- 协议配置
- 超时设置
- 日志配置
- 插件配置

示例：

```yaml
server:
  service:
    - name: trpc.conversation.ConversationService
      ip: 127.0.0.1
      port: 8000
      protocol: trpc
      timeout: 1000

plugins:
  log:
    default:
      - writer: console
        level: debug
        format: json

global:
  namespace: Development
  env_name: dev
```

## 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./handler/conversation

# 运行测试并显示覆盖率
go test -cover ./...
```

## 部署

### Docker 部署

```bash
# 构建镜像
docker build -t elysia-backend .

# 运行容器
docker run -p 8000:8000 --env-file .env elysia-backend
```

### Docker Compose

```bash
docker-compose up -d
```

## 常见问题

### 1. 依赖下载失败

如果遇到 `trpc.group` 相关依赖下载失败，请确保：

- 配置了正确的 Go proxy：`go env -w GOPROXY=https://goproxy.cn,direct`
- 网络连接正常

### 2. Proto 生成失败

确保已安装：
- protoc 编译器
- protoc-gen-go 插件
- protoc-gen-go-trpc 插件

### 3. 服务启动失败

检查：
- 端口 8000 是否被占用
- trpc_go.yaml 配置是否正确
- 数据库连接是否正常

## 技术栈

- **框架**: tRPC-Go v1.0.3
- **协议**: Protocol Buffers
- **数据库**: MySQL
- **日志**: Zap
- **配置**: godotenv

## License

MIT