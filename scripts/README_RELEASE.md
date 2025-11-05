# 🚀 自动发布系统使用指南

## 📦 已创建的文件

### 1. 发布脚本
- `scripts/release.sh` - 自动发布脚本

### 2. GitHub Actions 工作流
- `.github/workflows/docker-publish.yml` - Docker 镜像构建和推送
- `.github/workflows/update-badge.yml` - 自动更新 README 徽章
- `.github/workflows/test.yml` - 代码测试工作流

### 3. 版本文件
- `VERSION` - 当前版本号
- `RELEASE_GUIDE.md` - 详细发布指南

## 🎯 快速开始

### 方式一：使用自动脚本（推荐）⭐

```bash
cd /root/MYGallery
bash scripts/release.sh
```

脚本会交互式引导你完成整个发布流程：

```
╔══════════════════════════════════════════╗
║   📦 MYGallery 自动发布工具              ║
╚══════════════════════════════════════════╝

📌 当前版本: 2.0.0

请选择版本类型:
  1) major (大版本更新，如 1.0.0 -> 2.0.0)
  2) minor (功能更新，如 1.0.0 -> 1.1.0)
  3) patch (修复更新，如 1.0.0 -> 1.0.1)
  4) 手动输入版本号

请选择 [1-4]:
```

选择版本类型后，脚本会自动：
- ✅ 更新 VERSION 文件
- ✅ 更新 config.example.yaml
- ✅ 更新 README.md 和 README_CN.md
- ✅ 更新 CHANGELOG.md
- ✅ 提交更改到 Git
- ✅ 创建 Git tag
- ✅ 推送到 GitHub
- ✅ 触发 GitHub Actions 自动构建

### 方式二：手动发布

```bash
# 1. 更新版本号
echo "2.1.0" > VERSION

# 2. 创建 tag
git add .
git commit -m "chore: bump version to 2.1.0"
git tag -a v2.1.0 -m "Release 2.1.0"

# 3. 推送
git push origin main
git push origin v2.1.0
```

## ⚙️ GitHub Secrets 配置

在使用自动发布功能前，需要在 GitHub 仓库设置中配置以下 Secrets：

### 1. Docker Hub 凭据

```
Settings -> Secrets and variables -> Actions -> New repository secret
```

添加两个 Secret：
- **Name**: `DOCKER_USERNAME`
  **Value**: 你的 Docker Hub 用户名

- **Name**: `DOCKER_PASSWORD`
  **Value**: 你的 Docker Hub 访问令牌（不是密码）

### 2. 获取 Docker Hub 访问令牌

1. 登录 [Docker Hub](https://hub.docker.com/)
2. 点击右上角头像 -> Account Settings
3. 进入 Security 标签
4. 点击 "New Access Token"
5. 输入描述（如 "GitHub Actions"）
6. 选择权限（Read, Write, Delete）
7. 生成并复制 Token
8. ⚠️ **重要**: 立即保存到 GitHub Secrets（关闭后无法再查看）

## 🔄 自动化工作流程

### 流程图

```
开发者运行 release.sh
         ↓
脚本自动更新版本号和文档
         ↓
创建 Git tag 并推送
         ↓
触发 GitHub Actions (docker-publish.yml)
         ↓
构建多平台 Docker 镜像
         ↓
推送到 Docker Hub
         ↓
触发徽章更新 (update-badge.yml)
         ↓
自动更新 README 中的徽章
         ↓
✅ 发布完成
```

### GitHub Actions 构建的镜像标签

每次发布会自动创建以下 Docker 镜像标签：

```bash
yourusername/mygallery:2.1.0    # 完整版本号
yourusername/mygallery:2.1      # 主版本 + 次版本
yourusername/mygallery:2        # 主版本
yourusername/mygallery:latest   # 最新版本
```

## 📋 发布检查清单

在运行发布脚本前：

- [ ] 所有代码已提交
- [ ] 测试已通过（`go test ./...`）
- [ ] 编译成功（`go build`）
- [ ] 本地测试运行正常
- [ ] CHANGELOG 准备好更新内容
- [ ] GitHub Secrets 已配置
- [ ] Docker Hub 账号正常

## 📊 查看发布状态

### 1. GitHub Actions

访问你的仓库：
```
https://github.com/yourusername/mygallery/actions
```

查看工作流运行状态：
- ✅ 绿色勾：构建成功
- ❌ 红色叉：构建失败
- 🟡 黄色圈：正在构建

### 2. Docker Hub

访问：
```
https://hub.docker.com/r/yourusername/mygallery
```

查看：
- 镜像标签列表
- 镜像大小
- 拉取次数
- 最后更新时间

### 3. 本地测试新版本

```bash
# 拉取最新镜像
docker pull yourusername/mygallery:latest

# 运行容器
docker run -d -p 8080:8080 yourusername/mygallery:latest

# 测试健康检查
curl http://localhost:8080/health
```

## 🔧 常见问题

### Q: 推送 tag 后没有触发 Actions？

**A**: 检查：
1. `.github/workflows/docker-publish.yml` 文件是否存在
2. 工作流文件是否在 main 分支
3. 仓库是否启用了 Actions

### Q: Docker 镜像构建失败？

**A**: 检查：
1. GitHub Secrets 是否正确配置
2. Docker Hub 用户名和仓库名是否匹配
3. Dockerfile 是否有语法错误
4. 查看 Actions 日志获取详细错误

### Q: 徽章没有更新？

**A**: 
1. 等待 1-2 分钟（缓存刷新）
2. 检查 update-badge.yml 工作流是否运行
3. 手动触发工作流：
   ```
   Actions -> Update README Badges -> Run workflow
   ```

### Q: 如何回滚版本？

**A**:
```bash
# 删除本地 tag
git tag -d v2.1.0

# 删除远程 tag
git push origin :refs/tags/v2.1.0

# 恢复 VERSION 文件
echo "2.0.0" > VERSION
git add VERSION
git commit -m "chore: rollback to 2.0.0"
git push origin main
```

## 📚 相关文档

- [RELEASE_GUIDE.md](../RELEASE_GUIDE.md) - 详细发布指南
- [CHANGELOG.md](../CHANGELOG.md) - 版本更新记录
- [GitHub Actions 文档](https://docs.github.com/actions)
- [Docker Hub 文档](https://docs.docker.com/docker-hub/)

## 🎉 示例：发布 v2.1.0

```bash
# 1. 运行发布脚本
bash scripts/release.sh

# 2. 选择版本类型（选择 2 - minor）
请选择 [1-4]: 2

# 3. 确认发布
🎯 新版本: 2.1.0
确认发布版本 2.1.0？(y/n) y

# 4. 等待自动完成
✅ 发布完成！
📦 版本: 2.1.0
🏷️  标签: v2.1.0
🌿 分支: main

# 5. 访问 GitHub Actions 查看构建
https://github.com/yourusername/mygallery/actions

# 6. 等待几分钟后，新版本镜像即可使用
docker pull yourusername/mygallery:2.1.0
```

## 💡 最佳实践

1. **小步快跑**: 频繁发布小版本而不是积累大量更改
2. **语义化版本**: 严格遵循版本规范
3. **详细日志**: 在 CHANGELOG 中记录所有更改
4. **测试充分**: 发布前确保所有测试通过
5. **备份重要**: 发布前备份重要数据

---

**祝你发布顺利！** 🚀✨

