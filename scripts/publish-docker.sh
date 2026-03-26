#!/bin/bash
# MYGallery Docker 镜像发布脚本
# 使用方法：./scripts/publish-docker.sh [version]

set -e

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 配置
IMAGE_NAME="danzai233/mygallery"
VERSION=${1:-latest}

echo -e "${GREEN}🐳 MYGallery Docker 镜像发布脚本${NC}"
echo -e "${YELLOW}================================${NC}"
echo ""

# 检查是否已登录 Docker Hub
echo -e "${YELLOW}检查 Docker Hub 登录状态...${NC}"
if ! docker info | grep -q "Username"; then
    echo -e "${RED}未登录 Docker Hub，请先运行：docker login${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Docker Hub 登录状态正常${NC}"
echo ""

# 构建 Dockerfile
echo -e "${YELLOW}构建 Docker 镜像...${NC}"
docker build -t ${IMAGE_NAME}:${VERSION} -t ${IMAGE_NAME}:latest .
echo -e "${GREEN}✅ Docker 镜像构建完成${NC}"
echo ""

# 推送镜像
echo -e "${YELLOW}推送镜像到 Docker Hub...${NC}"
docker push ${IMAGE_NAME}:${VERSION}
docker push ${IMAGE_NAME}:latest
echo -e "${GREEN}✅ 镜像推送完成${NC}"
echo ""

# 输出镜像信息
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}发布完成！${NC}"
echo -e "${YELLOW}镜像信息：${NC}"
echo -e "  - 镜像名称: ${IMAGE_NAME}"
echo -e "  - 版本标签: ${VERSION}"
echo -e "  - 最新标签: latest"
echo ""
echo -e "${YELLOW}使用方法：${NC}"
echo -e "  docker pull ${IMAGE_NAME}:${VERSION}"
echo -e "  docker pull ${IMAGE_NAME}:latest"
echo ""
echo -e "${GREEN}================================${NC}"
