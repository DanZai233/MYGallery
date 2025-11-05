#!/bin/bash

# 创建测试数据脚本

echo "📦 创建测试数据..."
echo ""

# 检查服务是否运行
if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "❌ 错误: MYGallery 服务未运行"
    echo "请先运行: go run main.go"
    exit 1
fi

echo "✅ 服务正在运行"
echo ""

# 登录获取 token
echo "🔐 登录获取 token..."
TOKEN=$(curl -s http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
    echo "❌ 登录失败"
    exit 1
fi

echo "✅ 登录成功"
echo ""

# 创建分类
echo "📁 创建测试分类..."

categories=(
  '{"name":"风景","slug":"landscape","description":"自然风光照片","sort_order":1}'
  '{"name":"人像","slug":"portrait","description":"人物摄影","sort_order":2}'
  '{"name":"城市","slug":"urban","description":"城市街拍","sort_order":3}'
  '{"name":"美食","slug":"food","description":"美食摄影","sort_order":4}'
  '{"name":"旅行","slug":"travel","description":"旅行记录","sort_order":5}'
)

for category in "${categories[@]}"; do
    RESULT=$(curl -s -X POST http://localhost:8080/api/categories \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "$category")
    
    NAME=$(echo "$RESULT" | jq -r '.category.name // empty')
    if [ -n "$NAME" ]; then
        echo "  ✓ 创建分类: $NAME"
    fi
done

echo ""
echo "✅ 测试数据创建完成"
echo ""
echo "📋 已创建的分类:"
curl -s http://localhost:8080/api/categories | jq -r '.[] | "  - \(.name) (\(.slug))"'
echo ""
echo "🎉 现在可以上传照片并设置分类了！"
echo ""
echo "访问:"
echo "  📷 前台: http://localhost:8080"
echo "  ⚙️  后台: http://localhost:8080/admin"

