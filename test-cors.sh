#!/bin/bash

echo "🧪 测试 CORS 配置"
echo "=================="
echo ""

# 测试后端服务是否运行
echo "1️⃣ 检查后端服务..."
if curl -s http://localhost:8001 > /dev/null 2>&1; then
    echo "   ✅ 后端服务运行正常 (http://localhost:8001)"
else
    echo "   ❌ 后端服务未运行，请先启动后端服务"
    exit 1
fi
echo ""

# 测试 OPTIONS 请求
echo "2️⃣ 测试 OPTIONS 请求（Preflight）..."
response=$(curl -s -X OPTIONS http://localhost:8001/api/admin/auth/login-password \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type, Authorization" \
  -i)

if echo "$response" | grep -q "Access-Control-Allow-Origin"; then
    echo "   ✅ OPTIONS 请求成功"
    echo "   📋 CORS 响应头："
    echo "$response" | grep "Access-Control" | sed 's/^/      /'
else
    echo "   ❌ OPTIONS 请求失败"
    echo "$response"
    exit 1
fi
echo ""

# 测试 POST 请求
echo "3️⃣ 测试 POST 请求..."
response=$(curl -s -X POST http://localhost:8001/api/admin/auth/login-password \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3000" \
  -d '{"phone_number":"13800138000","password":"123456Admin"}' \
  -i)

if echo "$response" | grep -q "Access-Control-Allow-Origin"; then
    echo "   ✅ POST 请求成功"
    echo "   📋 CORS 响应头："
    echo "$response" | grep "Access-Control" | sed 's/^/      /'
else
    echo "   ❌ POST 请求失败"
    echo "$response"
    exit 1
fi
echo ""

# 测试前端服务
echo "4️⃣ 检查前端服务..."
if curl -s http://localhost:3000 > /dev/null 2>&1; then
    echo "   ✅ 前端服务运行正常 (http://localhost:3000)"
else
    echo "   ⚠️  前端服务未运行"
    echo "   💡 启动前端：cd /Users/sylvainyang/project/elysia/elysia-frontend/web-frontend && npm run dev"
fi
echo ""

echo "🎉 CORS 配置测试完成！"
echo ""
echo "📝 下一步："
echo "   1. 在浏览器中访问：http://localhost:3000"
echo "   2. 打开开发者工具（F12）→ Network 标签"
echo "   3. 尝试登录功能"
echo "   4. 检查网络请求是否正常"
echo ""
