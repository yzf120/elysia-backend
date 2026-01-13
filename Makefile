# Elysia Backend Makefile

.PHONY: help build run test clean deps setup proto install-tools

# 默认目标
help:
	@echo "Elysia Backend 构建工具"
	@echo ""
	@echo "可用命令:"
	@echo "  make deps          - 下载依赖"
	@echo "  make install-tools - 安装trpc工具"
	@echo "  make proto         - 生成protobuf代码"
	@echo "  make build         - 构建项目"
	@echo "  make run           - 运行项目"
	@echo "  make test          - 运行测试"
	@echo "  make clean         - 清理构建文件"
	@echo "  make setup         - 初始化项目环境"
	@echo ""

# 下载依赖
deps:
	@echo "下载依赖..."
	go mod tidy
	go mod download

# 安装trpc工具
install-tools:
	@echo "安装trpc工具..."
	go install trpc.group/trpc-go/trpc-cmdline/trpc@latest
	@echo "trpc工具安装完成"

# 生成protobuf代码
proto: install-tools
	@echo "生成protobuf代码..."
	cd proto && make all
	@echo "protobuf代码生成完成"

# 构建项目
build: deps proto
	@echo "构建项目..."
	go build -o bin/elysia-backend main.go
	@echo "构建完成: bin/elysia-backend"

# 运行项目
run: deps
	@echo "启动开发服务器..."
	go run main.go

# 使用启动脚本运行项目（推荐）
run-with-services: deps
	@echo "使用启动脚本运行项目（自动启动数据库服务）..."
	chmod +x start.sh
	./start.sh

# 运行测试
test: deps
	@echo "运行测试..."
	go test ./... -v

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	rm -f elysia-backend
	rm -rf proto/*/*.pb.go
	rm -rf proto/*/*.trpc.go

# 初始化项目环境
setup:
	@echo "初始化项目环境..."
	@if [ ! -f .env ]; then \
		echo "创建.env文件..."; \
		cp .env.example .env 2>/dev/null || echo "请手动创建.env文件"; \
	fi
	@echo "请编辑.env文件配置数据库连接信息"
	@echo "初始化完成"

# 生产环境构建
build-prod: deps proto
	@echo "构建生产版本..."
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/elysia-backend-prod main.go
	@echo "生产版本构建完成: bin/elysia-backend-prod"

# 代码检查
lint: deps
	@echo "运行代码检查..."
	go vet ./...
	golangci-lint run ./... || echo "请安装golangci-lint: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"