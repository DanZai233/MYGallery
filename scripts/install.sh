#!/bin/bash

# MYGallery 安装脚本

set -e

echo "╔══════════════════════════════════════════╗"
echo "║   📷 MYGallery 安装向导                   ║"
echo "╚══════════════════════════════════════════╝"
echo ""

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    echo "   访问 https://docs.docker.com/get-docker/"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

echo "✅ Docker 环境检查通过"
echo ""

# 创建配置文件
if [ ! -f "config.yaml" ]; then
    echo "📝 创建配置文件..."
    cp config.example.yaml config.yaml
    echo "✅ 配置文件已创建: config.yaml"
    echo ""
    echo "⚠️  请编辑 config.yaml 修改默认配置（特别是管理员密码）"
    echo ""
    read -p "是否现在编辑配置文件？(y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        ${EDITOR:-nano} config.yaml
    fi
fi

# 创建必要目录
echo "📁 创建必要目录..."
mkdir -p data uploads uploads/thumbnails public/assets
echo "✅ 目录创建完成"
echo ""

# 构建镜像
echo "🐳 构建 Docker 镜像..."
if docker compose build 2>/dev/null || docker-compose build 2>/dev/null; then
    echo "✅ 镜像构建完成"
else
    echo "❌ 镜像构建失败，请检查错误信息"
    exit 1
fi
echo ""

# 启动服务
echo "🚀 启动服务..."
if docker compose up -d 2>/dev/null || docker-compose up -d 2>/dev/null; then
    echo "✅ 服务已启动"
else
    echo "❌ 服务启动失败，请检查错误信息"
    exit 1
fi
echo ""

# 等待服务就绪
echo "⏳ 等待服务就绪..."
sleep 3

# 显示访问信息
echo "╔══════════════════════════════════════════╗"
echo "║   ✅ MYGallery 安装成功！                 ║"
echo "╚══════════════════════════════════════════╝"
echo ""
echo "📷 前台页面: http://localhost:8080"
echo "⚙️  后台管理: http://localhost:8080/admin"
echo ""
echo "默认管理员账号:"
echo "  用户名: admin"
echo "  密码: admin123"
echo ""
echo "⚠️  首次登录后请立即修改密码！"
echo ""
echo "查看日志: docker-compose logs -f"
echo "停止服务: docker-compose down"
echo ""

