#!/bin/bash

# 阿里云短信服务测试脚本
# 用途：测试所有短信发送接口

set -e

# 配置
BASE_URL="http://localhost:8001"
PHONE_NUMBER="18873197041"  # 修改为你的测试手机号

echo "📱 阿里云短信服务测试"
echo "===================="
echo ""
echo "📞 测试手机号: $PHONE_NUMBER"
echo "🌐 服务地址: $BASE_URL"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_sms() {
    local name=$1
    local endpoint=$2
    
    echo -e "${YELLOW}测试: $name${NC}"
    echo "接口: POST $endpoint"
    
    response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL$endpoint" \
        -H "Content-Type: application/json" \
        -d "{\"phone_number\": \"$PHONE_NUMBER\"}")
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)
    
    if [ "$http_code" = "200" ]; then
        echo -e "${GREEN}✅ 成功${NC}"
        echo "响应: $body"
    else
        echo -e "${RED}❌ 失败 (HTTP $http_code)${NC}"
        echo "响应: $body"
    fi
    
    echo ""
    echo "⏳ 等待60秒（避免频率限制）..."
    sleep 60
    echo ""
}

# 检查服务是否运行
echo "🔍 检查后端服务..."
if ! curl -s "$BASE_URL/health" > /dev/null 2>&1; then
    echo -e "${RED}❌ 后端服务未运行${NC}"
    echo "请先启动后端服务: ./elysia-backend"
    exit 1
fi
echo -e "${GREEN}✅ 后端服务正常运行${NC}"
echo ""

# 开始测试
echo "🧪 开始测试短信发送功能..."
echo ""

# 1. 学生注册验证码
test_sms "学生注册验证码" "/api/student/auth/send-register-code"

# 2. 学生登录验证码
test_sms "学生登录验证码" "/api/student/auth/send-login-code"

# 3. 教师注册验证码
test_sms "教师注册验证码" "/api/teacher/auth/send-register-code"

# 4. 教师登录验证码
test_sms "教师登录验证码" "/api/teacher/auth/send-login-code"

# 5. 管理员登录验证码
test_sms "管理员登录验证码" "/api/admin/auth/send-login-code"

echo "===================="
echo -e "${GREEN}🎉 测试完成！${NC}"
echo ""
echo "📝 注意事项："
echo "  1. 请检查手机是否收到验证码短信"
echo "  2. 验证码有效期为5分钟"
echo "  3. 发送频率限制为60秒/次"
echo ""
