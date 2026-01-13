# Elysia Backend 启动指南

## 启动方式

### 方式一：使用启动脚本（推荐）
```bash
# 使用启动脚本，自动检查并启动数据库服务
make run-with-services
```

### 方式二：直接运行（需要手动启动数据库服务）
```bash
# 直接运行 Go 应用（需要先确保数据库服务已启动）
make run
```

### 方式三：手动使用启动脚本
```bash
# 直接执行启动脚本
chmod +x start.sh
./start.sh
```

## 数据库服务管理

### 启动所有服务
```bash
docker-compose up -d
```

### 停止所有服务
```bash
docker-compose down
```

### 查看服务状态
```bash
docker-compose ps
```

### 查看服务日志
```bash
# 查看 MySQL 日志
docker-compose logs mysql

# 查看 Redis 日志
docker-compose logs redis

# 查看所有服务日志
docker-compose logs -f
```

## 服务配置

### MySQL 服务
- 容器名称: `elysia-mysql`
- 端口: `3306`
- 数据库: `elysia`
- 用户名: `elysia`
- 密码: `elysia123`

### Redis 服务
- 容器名称: `elysia-redis`
- 端口: `6379`
- 无密码
- 数据库编号: `0`

## 环境变量配置

确保 `.env` 文件中的配置与 docker-compose.yml 中的服务配置匹配：

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=elysia
DB_PASSWORD=elysia123
DB_NAME=elysia

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

## 故障排除

### 服务启动失败
如果服务启动失败，可以尝试：
1. 检查 Docker 是否运行：`docker info`
2. 查看服务日志：`docker-compose logs`
3. 重启服务：`docker-compose restart`

### 端口冲突
如果端口 3306 或 6379 已被占用：
1. 停止占用端口的其他服务
2. 或者修改 docker-compose.yml 中的端口映射

### 健康检查失败
如果服务健康检查失败：
1. 等待更长时间让服务完全启动
2. 检查服务配置是否正确
3. 查看服务日志了解具体错误

## 开发建议

1. **推荐使用 `make run-with-services`**：自动处理数据库服务依赖
2. **开发时保持服务运行**：避免频繁启动停止服务
3. **使用 Docker Desktop**：便于管理容器服务
4. **定期备份数据**：重要数据定期备份到本地