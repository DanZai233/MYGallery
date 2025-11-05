# 🚀 发布指南

本文档说明如何发布 MYGallery 的新版本。

## 📋 准备工作

### 1. 配置 GitHub Secrets

在 GitHub 仓库设置中添加以下 Secrets：

- `DOCKER_USERNAME`: Docker Hub 用户名
- `DOCKER_PASSWORD`: Docker Hub 访问令牌

获取 Docker Hub 访问令牌：
1. 登录 Docker Hub
2. 访问 Account Settings -> Security
3. 创建新的 Access Token
4. 复制 Token 到 GitHub Secrets

### 2. 确保代码已提交

```bash
git status
git add .
git commit -m "feat: your feature description"
```

## 🎯 发布新版本

### 自动发布（推荐）

运行发布脚本：

```bash
bash scripts/release.sh
```

脚本会引导你：
1. 选择版本类型（major/minor/patch）
2. 自动计算新版本号
3. 更新所有相关文件
4. 创建 Git tag
5. 推送到 GitHub
6. 触发自动构建

### 手动发布

#### 1. 更新版本号

```bash
# 修改 VERSION 文件
echo "2.1.0" > VERSION

# 更新 config.example.yaml
sed -i 's/version: ".*"/version: "2.1.0"/' config.example.yaml

# 更新 README
sed -i 's/v[0-9]\+\.[0-9]\+\.[0-9]\+/v2.1.0/g' README*.md
```

#### 2. 更新 CHANGELOG

编辑 `CHANGELOG.md`，添加新版本的更新内容：

```markdown
## [2.1.0] - 2025-11-05

### 新增
- 添加新功能 A
- 添加新功能 B

### 改进
- 优化性能
- 改进 UI

### 修复
- 修复 Bug X
- 修复 Bug Y
```

#### 3. 提交更改

```bash
git add VERSION config.example.yaml README*.md CHANGELOG.md
git commit -m "chore: bump version to 2.1.0"
```

#### 4. 创建 Tag

```bash
git tag -a v2.1.0 -m "Release 2.1.0"
```

#### 5. 推送到 GitHub

```bash
git push origin main
git push origin v2.1.0
```

## 🔄 自动化流程

### 1. Docker 镜像构建

推送 tag 后，GitHub Actions 会自动：
- 构建多平台 Docker 镜像（amd64, arm64）
- 推送到 Docker Hub
- 更新镜像描述
- 添加版本标签

查看构建进度：
```
https://github.com/yourusername/mygallery/actions
```

### 2. 徽章更新

构建成功后，会自动更新 README 中的徽章：
- 版本号徽章
- 镜像大小徽章
- 构建状态徽章

## 📦 Docker 镜像标签

每次发布会创建以下标签：

```bash
# 完整版本号
yourusername/mygallery:2.1.0

# 主版本号 + 次版本号
yourusername/mygallery:2.1

# 主版本号
yourusername/mygallery:2

# 最新版本
yourusername/mygallery:latest
```

## 🧪 发布前测试

### 1. 本地测试

```bash
# 运行测试
go test -v ./...

# 编译检查
go build -v ./...

# 运行应用
go run main.go
```

### 2. Docker 测试

```bash
# 构建镜像
docker build -t mygallery:test .

# 运行容器
docker run -p 8080:8080 mygallery:test

# 测试访问
curl http://localhost:8080/health
```

## 📝 版本规范

遵循 [语义化版本](https://semver.org/lang/zh-CN/)：

- **MAJOR (主版本号)**: 不兼容的 API 修改
  - 示例: 1.0.0 -> 2.0.0
  - 场景: 重大架构变更、破坏性更新

- **MINOR (次版本号)**: 向后兼容的功能新增
  - 示例: 1.0.0 -> 1.1.0
  - 场景: 新功能、新特性

- **PATCH (修订号)**: 向后兼容的问题修正
  - 示例: 1.0.0 -> 1.0.1
  - 场景: Bug 修复、小改进

## 🔍 检查发布状态

### 查看 GitHub Actions

```bash
# 打开 Actions 页面
https://github.com/yourusername/mygallery/actions
```

### 查看 Docker Hub

```bash
# 检查镜像是否推送成功
docker pull yourusername/mygallery:latest

# 查看镜像信息
docker images | grep mygallery
```

### 测试新版本

```bash
# 使用新版本运行
docker run -d -p 8080:8080 yourusername/mygallery:2.1.0

# 检查健康状态
curl http://localhost:8080/health
```

## 🛠️ 故障排除

### Actions 构建失败

1. 检查 GitHub Secrets 是否配置正确
2. 查看 Actions 日志找出错误原因
3. 修复问题后重新推送 tag：
   ```bash
   git tag -d v2.1.0
   git push origin :refs/tags/v2.1.0
   git tag -a v2.1.0 -m "Release 2.1.0"
   git push origin v2.1.0
   ```

### Docker Hub 推送失败

1. 验证 Docker Hub 凭据
2. 检查网络连接
3. 确认仓库权限

### 徽章未更新

1. 等待几分钟（缓存刷新）
2. 手动触发工作流：
   ```
   GitHub -> Actions -> Update README Badges -> Run workflow
   ```

## 📚 相关文档

- [CHANGELOG.md](CHANGELOG.md) - 版本更新记录
- [CONTRIBUTING.md](CONTRIBUTING.md) - 贡献指南
- [GitHub Actions 文档](https://docs.github.com/actions)
- [Docker Hub 文档](https://docs.docker.com/docker-hub/)

## 🎉 发布清单

发布新版本前的检查清单：

- [ ] 所有测试通过
- [ ] 代码已审查
- [ ] CHANGELOG 已更新
- [ ] 版本号已更新
- [ ] 文档已更新
- [ ] GitHub Secrets 已配置
- [ ] 本地测试通过
- [ ] Docker 构建成功

---

**需要帮助？** 查看 [TROUBLESHOOTING.md](TROUBLESHOOTING.md) 或提交 [Issue](https://github.com/yourusername/mygallery/issues)

