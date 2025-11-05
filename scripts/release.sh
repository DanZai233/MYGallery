#!/bin/bash

# MYGallery 自动发布脚本
# 功能：自动打 tag、推送到 GitHub、触发 Actions 构建、更新文档

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔══════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   📦 MYGallery 自动发布工具              ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════╝${NC}"
echo ""

# 检查是否在 git 仓库中
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}❌ 错误: 当前目录不是 git 仓库${NC}"
    exit 1
fi

# 检查是否有未提交的更改
if [[ -n $(git status -s) ]]; then
    echo -e "${YELLOW}⚠️  警告: 有未提交的更改${NC}"
    git status -s
    echo ""
    read -p "是否继续？(y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${RED}❌ 已取消${NC}"
        exit 1
    fi
fi

# 获取当前版本
if [ -f "VERSION" ]; then
    CURRENT_VERSION=$(cat VERSION)
else
    CURRENT_VERSION="0.0.0"
fi

echo -e "${GREEN}📌 当前版本: ${CURRENT_VERSION}${NC}"
echo ""

# 选择版本类型
echo "请选择版本类型:"
echo "  1) major (大版本更新，如 1.0.0 -> 2.0.0)"
echo "  2) minor (功能更新，如 1.0.0 -> 1.1.0)"
echo "  3) patch (修复更新，如 1.0.0 -> 1.0.1)"
echo "  4) 手动输入版本号"
echo ""
read -p "请选择 [1-4]: " choice

case $choice in
    1)
        VERSION_TYPE="major"
        ;;
    2)
        VERSION_TYPE="minor"
        ;;
    3)
        VERSION_TYPE="patch"
        ;;
    4)
        read -p "请输入新版本号 (如 2.1.0): " NEW_VERSION
        if [[ ! $NEW_VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            echo -e "${RED}❌ 错误: 版本号格式不正确 (应为 x.y.z)${NC}"
            exit 1
        fi
        VERSION_TYPE="manual"
        ;;
    *)
        echo -e "${RED}❌ 错误: 无效的选择${NC}"
        exit 1
        ;;
esac

# 自动计算新版本号
if [ "$VERSION_TYPE" != "manual" ]; then
    IFS='.' read -r -a VERSION_PARTS <<< "$CURRENT_VERSION"
    MAJOR="${VERSION_PARTS[0]}"
    MINOR="${VERSION_PARTS[1]}"
    PATCH="${VERSION_PARTS[2]}"

    case $VERSION_TYPE in
        major)
            MAJOR=$((MAJOR + 1))
            MINOR=0
            PATCH=0
            ;;
        minor)
            MINOR=$((MINOR + 1))
            PATCH=0
            ;;
        patch)
            PATCH=$((PATCH + 1))
            ;;
    esac

    NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"
fi

echo ""
echo -e "${GREEN}🎯 新版本: ${NEW_VERSION}${NC}"
echo ""

# 确认发布
read -p "确认发布版本 ${NEW_VERSION}？(y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${RED}❌ 已取消${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}📝 开始发布流程...${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# 1. 更新 VERSION 文件
echo -e "${YELLOW}1/7${NC} 更新 VERSION 文件..."
echo "$NEW_VERSION" > VERSION
echo -e "${GREEN}✓ 已更新 VERSION${NC}"

# 2. 更新 config.example.yaml 中的版本
echo -e "${YELLOW}2/7${NC} 更新配置文件版本..."
if [ -f "config.example.yaml" ]; then
    sed -i "s/version: \".*\"/version: \"$NEW_VERSION\"/" config.example.yaml
    echo -e "${GREEN}✓ 已更新 config.example.yaml${NC}"
fi

# 3. 更新 README.md 中的版本
echo -e "${YELLOW}3/7${NC} 更新 README 文档..."
for readme in README.md README_CN.md; do
    if [ -f "$readme" ]; then
        # 更新版本号
        sed -i "s/v[0-9]\+\.[0-9]\+\.[0-9]\+/v$NEW_VERSION/g" "$readme"
        # 更新 Docker 镜像标签
        sed -i "s/mygallery:[0-9]\+\.[0-9]\+\.[0-9]\+/mygallery:$NEW_VERSION/g" "$readme"
        sed -i "s/mygallery:v[0-9]\+\.[0-9]\+\.[0-9]\+/mygallery:v$NEW_VERSION/g" "$readme"
        echo -e "${GREEN}✓ 已更新 $readme${NC}"
    fi
done

# 4. 更新 CHANGELOG.md
echo -e "${YELLOW}4/7${NC} 更新 CHANGELOG..."
if [ -f "CHANGELOG.md" ]; then
    RELEASE_DATE=$(date +%Y-%m-%d)
    # 在文件开头插入新版本信息
    sed -i "/^## \[.*\]/i\\
\\
## [${NEW_VERSION}] - ${RELEASE_DATE}\\
\\
### 更新内容\\
\\
- 发布版本 ${NEW_VERSION}\\
- 请在此添加具体的更新内容\\
" CHANGELOG.md
    echo -e "${GREEN}✓ 已更新 CHANGELOG.md${NC}"
fi

# 5. 提交更改
echo -e "${YELLOW}5/7${NC} 提交更改到 Git..."
git add VERSION config.example.yaml README*.md CHANGELOG.md 2>/dev/null || true
git commit -m "chore: bump version to ${NEW_VERSION}" || true
echo -e "${GREEN}✓ 已提交更改${NC}"

# 6. 创建 Git tag
echo -e "${YELLOW}6/7${NC} 创建 Git tag..."
TAG_NAME="v${NEW_VERSION}"
if git rev-parse "$TAG_NAME" >/dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Tag ${TAG_NAME} 已存在，将删除并重新创建${NC}"
    git tag -d "$TAG_NAME"
fi

git tag -a "$TAG_NAME" -m "Release ${NEW_VERSION}"
echo -e "${GREEN}✓ 已创建 tag: ${TAG_NAME}${NC}"

# 7. 推送到 GitHub
echo -e "${YELLOW}7/7${NC} 推送到 GitHub..."
CURRENT_BRANCH=$(git branch --show-current)
git push origin "$CURRENT_BRANCH"
git push origin "$TAG_NAME"
echo -e "${GREEN}✓ 已推送到 GitHub${NC}"

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ 发布完成！${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${GREEN}📦 版本: ${NEW_VERSION}${NC}"
echo -e "${GREEN}🏷️  标签: ${TAG_NAME}${NC}"
echo -e "${GREEN}🌿 分支: ${CURRENT_BRANCH}${NC}"
echo ""
echo -e "${YELLOW}📋 接下来的步骤:${NC}"
echo -e "  1. GitHub Actions 将自动构建 Docker 镜像"
echo -e "  2. 访问 GitHub Actions 查看构建进度:"
echo -e "     ${BLUE}https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]\(.*\)\.git/\1/')/actions${NC}"
echo -e "  3. 构建完成后，镜像将推送到 Docker Hub"
echo -e "  4. 更新 CHANGELOG.md 中的具体更新内容"
echo ""
echo -e "${GREEN}🎉 感谢使用 MYGallery！${NC}"

