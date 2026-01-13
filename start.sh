#!/bin/bash

# Elysia Backend 启动脚本
# 在运行 main.go 之前自动启动 MySQL 和 Redis 服务

echo "🚀 启动 Elysia Backend 服务..."

# 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker 未运行，请先启动 Docker"
    exit 1
fi

echo "📦 检查并启动数据库服务..."

# 检查 MySQL 服务状态
MYSQL_RUNNING=$(docker ps --filter "name=elysia-mysql" --format "{{.Names}}")
if [ -z "$MYSQL_RUNNING" ]; then
    echo "🔧 启动 MySQL 服务..."
    docker-compose up -d mysql
    
    # 等待 MySQL 健康检查通过
    echo "⏳ 等待 MySQL 服务就绪..."
    for i in {1..30}; do
        if docker ps --filter "name=elysia-mysql" --filter "health=healthy" --format "{{.Names}}" | grep -q "elysia-mysql"; then
            echo "✅ MySQL 服务已就绪"
            break
        fi
        if [ $i -eq 30 ]; then
            echo "❌ MySQL 服务启动超时"
            exit 1
        fi
        sleep 2
    done
else
    echo "✅ MySQL 服务已在运行"
fi

# 检查 Redis 服务状态
REDIS_RUNNING=$(docker ps --filter "name=elysia-redis" --format "{{.Names}}")
if [ -z "$REDIS_RUNNING" ]; then
    echo "🔧 启动 Redis 服务..."
    docker-compose up -d redis
    
    # 等待 Redis 健康检查通过
    echo "⏳ 等待 Redis 服务就绪..."
    for i in {1..15}; do
        if docker ps --filter "name=elysia-redis" --filter "health=healthy" --format "{{.Names}}" | grep -q "elysia-redis"; then
            echo "✅ Redis 服务已就绪"
            break
        fi
        if [ $i -eq 15 ]; then
            echo "❌ Redis 服务启动超时"
            exit 1
        fi
        sleep 2
    done
else
    echo "✅ Redis 服务已在运行"
fi

echo "🎯 所有依赖服务已就绪，启动 Go 应用..."

# 运行 Go 应用
go run main.go

echo "👋 应用已停止"