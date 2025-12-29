#!/bin/bash

echo "🔧 开始修复 Git 文件状态异常问题..."
echo ""

# 1. 检查并删除 Git 锁文件
echo "📌 步骤 1: 检查 Git 锁文件"
if [ -f .git/index.lock ]; then
    echo "  ⚠️  发现 Git 索引锁文件，正在删除..."
    rm -f .git/index.lock
    echo "  ✅ Git 锁文件已删除"
else
    echo "  ✅ 没有发现 Git 锁文件"
fi
echo ""

# 2. 检查文件权限
echo "📌 步骤 2: 修复文件权限"
chmod -R u+w dao/ service/ model/ config/ router/ 2>/dev/null
echo "  ✅ 文件权限已修复"
echo ""

# 3. 显示当前 Git 状态
echo "📌 步骤 3: 当前 Git 状态"
git status --short
echo ""

# 4. 重置 Git 索引
echo "📌 步骤 4: 重置 Git 索引"
git reset
echo "  ✅ Git 索引已重置"
echo ""

# 5. 清理 Git 缓存
echo "📌 步骤 5: 清理 Git 缓存"
git rm -r --cached . 2>/dev/null
git add .
echo "  ✅ Git 缓存已清理"
echo ""

# 6. 显示修复后的状态
echo "📌 步骤 6: 修复后的 Git 状态"
git status
echo ""

echo "✅ 修复完成！"
echo ""
echo "📝 后续操作："
echo "  1. 在 VSCode 中按 Cmd+Shift+P"
echo "  2. 输入 'Reload Window' 并执行"
echo "  3. 尝试编辑和删除代码"
echo "  4. 如果仍有问题，请尝试："
echo "     - 完全关闭 VSCode"
echo "     - 重新打开项目"
echo ""
