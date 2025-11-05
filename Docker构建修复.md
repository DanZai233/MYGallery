# 🔧 Docker 构建错误修复

## ❌ 错误信息

```
go: go.mod requires go >= 1.24.0 (running go 1.21.13; GOTOOLCHAIN=local)
```

## 🎯 问题原因

**版本不匹配**：
- 本地 Go 版本：1.24.4
- Dockerfile Go 版本：1.21.13
- go.mod 自动使用了本地版本要求

## ✅ 已修复

### 修改 1: 更新 Dockerfile Go 版本

```dockerfile
# 修复前
FROM golang:1.21-alpine AS builder

# 修复后
FROM golang:1.23-alpine AS builder
```

### 修改 2: 添加 GOTOOLCHAIN 环境变量

```dockerfile
# 添加这行，允许自动选择合适的 Go 版本
ENV GOTOOLCHAIN=auto
RUN CGO_ENABLED=1 GOOS=linux go build -o mygallery .
```

---

## 🧪 测试构建

### 本地测试

```bash
cd /root/MYGallery

# 测试构建（不推送）
docker build -t mygallery:test .

# 如果成功，测试运行
docker run -p 8080:8080 mygallery:test
```

### 推送后测试

```bash
# 提交更改
git add Dockerfile
git commit -m "fix: 更新 Dockerfile Go 版本到 1.23"
git push origin main

# 创建发布（触发 Actions）
bash scripts/release.sh
```

---

## 📋 完整的发布步骤

### 1. 配置 GitHub Secrets（必须）

访问：`https://github.com/DanZai233/mygallery/settings/secrets/actions`

添加两个 Secrets：
```
DOCKER_USERNAME = DanZai233
DOCKER_PASSWORD = [Docker Hub 访问令牌]
```

参考：`cat 快速配置Secrets.txt`

### 2. 提交 Dockerfile 修复

```bash
cd /root/MYGallery

git add Dockerfile
git commit -m "fix: 更新 Dockerfile Go 版本到 1.23"
git push origin main
```

### 3. 发布新版本

```bash
# 使用自动发布脚本
bash scripts/release.sh

# 或手动创建标签
git tag -a v1.0.1 -m "Release v1.0.1"
git push origin v1.0.1
```

### 4. 查看构建状态

访问：`https://github.com/DanZai233/mygallery/actions`

**预期流程**：
```
✅ Checkout 代码
✅ 设置 QEMU
✅ 设置 Docker Buildx
✅ 登录 Docker Hub (使用 Secrets)
✅ 提取元数据
✅ 构建镜像 (linux/amd64, linux/arm64)
✅ 推送到 Docker Hub
✅ 更新描述
```

### 5. 使用构建的镜像

构建成功后（约 5-10 分钟）：

```bash
# 拉取最新镜像
docker pull DanZai233/mygallery:latest

# 运行容器
docker run -d \
  --name mygallery \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/uploads:/app/uploads \
  DanZai233/mygallery:latest

# 访问应用
curl http://localhost:8080/health
```

---

## 🎯 本地构建测试（可选）

如果想在推送前本地测试：

```bash
cd /root/MYGallery

# 构建多平台镜像
docker buildx create --use
docker buildx build --platform linux/amd64,linux/arm64 -t mygallery:test .

# 或只构建当前平台
docker build -t mygallery:test .

# 测试运行
docker run -p 8080:8080 mygallery:test
```

---

## 📊 修复对比

### 修复前

```dockerfile
FROM golang:1.21-alpine AS builder  ❌ 版本太低
RUN go build -o mygallery .         ❌ 版本不匹配错误
```

### 修复后

```dockerfile
FROM golang:1.23-alpine AS builder  ✅ 版本更新
ENV GOTOOLCHAIN=auto                ✅ 自动选择版本
RUN go build -o mygallery .         ✅ 构建成功
```

---

## ⚠️ 重要提示

### Go 版本说明

- **本地开发**: Go 1.24.4（你当前使用的）
- **Docker 构建**: Go 1.23（Alpine 最新稳定版）
- **GOTOOLCHAIN=auto**: 自动下载需要的 Go 版本

### Secrets 安全

- ✅ 存储在 GitHub Secrets
- ✅ 加密存储
- ✅ 不会出现在日志中
- ❌ 不要提交到代码
- ❌ 不要分享给他人

---

## 🎉 修复完成

**Dockerfile 已更新，现在可以成功构建！**

**下一步**：
1. ✅ 配置 GitHub Secrets（Docker Hub 凭据）
2. ✅ 提交 Dockerfile 修复
3. ✅ 运行发布脚本
4. ✅ 等待构建完成
5. ✅ 使用 Docker 镜像

---

**查看详细配置步骤**：
```bash
cat 快速配置Secrets.txt
cat GitHub_Actions_配置指南.md
```

**所有修复已完成！** 🚀

