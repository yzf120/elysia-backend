#!/bin/bash

# 阿里云短信服务依赖安装脚本
# 用途：自动安装阿里云短信服务所需的Go依赖包

set -e

echo "🚀 开始安装阿里云短信服务依赖..."
echo ""

# 检查是否在正确的目录
if [ ! -f "go.mod" ]; then
    echo "❌ 错误：未找到 go.mod 文件"
    echo "请在项目根目录下运行此脚本"
    exit 1
fi

echo "📦 清理旧依赖..."
go clean -modcache 2>/dev/null || true

echo ""
echo "📥 下载阿里云SDK依赖..."

# 安装阿里云号码认证服务SDK
echo "  - 安装 dypnsapi-20170525..."
go get github.com/alibabacloud-go/dypnsapi-20170525/v3@latest

# 安装阿里云OpenAPI SDK
echo "  - 安装 darabonba-openapi..."
go get github.com/alibabacloud-go/darabonba-openapi/v2@latest

# 安装Tea SDK
echo "  - 安装 tea..."
go get github.com/alibabacloud-go/tea@latest

# 安装Tea Utils
echo "  - 安装 tea-utils..."
go get github.com/alibabacloud-go/tea-utils/v2@latest

# 安装阿里云凭证管理
echo "  - 安装 credentials-go..."
go get github.com/aliyun/credentials-go@latest

echo ""
echo "🔧 整理依赖..."
go mod tidy

echo ""
echo "✅ 依赖安装完成！"
echo ""
echo "📋 已安装的阿里云SDK："
echo "  ✓ github.com/alibabacloud-go/dypnsapi-20170525/v3"
echo "  ✓ github.com/alibabacloud-go/darabonba-openapi/v2"
echo "  ✓ github.com/alibabacloud-go/tea"
echo "  ✓ github.com/alibabacloud-go/tea-utils/v2"
echo "  ✓ github.com/aliyun/credentials-go"
echo ""
echo "📝 下一步："
echo "  1. 配置 .env 文件中的阿里云凭证"
echo "  2. 运行 go build 编译项目"
echo "  3. 启动服务并测试短信发送功能"
echo ""
echo "📖 详细说明请查看：ALIYUN_SMS_MIGRATION.md"
echo ""
