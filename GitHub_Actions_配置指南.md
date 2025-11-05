# 🤖 GitHub Actions 配置指南

## ❌ 错误原因

```
Error: Username and password required
```

**原因**：GitHub Actions 需要 Docker Hub 的登录凭据才能推送镜像。

---

## ✅ 解决方案：配置 GitHub Secrets

### 步骤 1: 获取 Docker Hub 访问令牌

1. **登录 Docker Hub**
   - 访问：https://hub.docker.com/
   - 使用你的账号登录（DanZai233）

2. **进入安全设置**
   - 点击右上角头像
   - Account Settings → Security

3. **创建访问令牌**
   - 点击 **"New Access Token"**
   - 描述：`GitHub Actions`
   - 权限：选择 **Read, Write, Delete**
   - 点击 **"Generate"**

4. **复制令牌**
   - ⚠️ **重要**：立即复制并保存令牌
   - 关闭后无法再查看！

### 步骤 2: 在 GitHub 添加 Secrets

1. **打开你的 GitHub 仓库**
   ```
   https://github.com/DanZai233/mygallery
   ```

2. **进入设置**
   - 点击仓库顶部的 **Settings** 标签

3. **打开 Secrets 设置**
   - 左侧菜单找到 **Secrets and variables**
   - 点击 **Actions**

4. **添加第一个 Secret**
   - 点击 **"New repository secret"**
   - Name: `DOCKER_USERNAME`
   - Secret: `DanZai233`（你的 Docker Hub 用户名）
   - 点击 **"Add secret"**

5. **添加第二个 Secret**
   - 再次点击 **"New repository secret"**
   - Name: `DOCKER_PASSWORD`
   - Secret: 粘贴刚才复制的 Docker Hub 访问令牌
   - 点击 **"Add secret"**

### 步骤 3: 验证配置

配置完成后，你应该看到：

```
Secrets / Actions secrets / Repository secrets

DOCKER_USERNAME    Updated now
DOCKER_PASSWORD    Updated now
```

---

## 🔄 重新运行 Workflow

### 方法 1: 推送新标签

```bash
cd /root/MYGallery

# 创建新标签
git tag -a v1.0.1 -m "Test release"
git push origin v1.0.1

# 或使用发布脚本
bash scripts/release.sh
```

### 方法 2: 手动触发

1. 打开仓库的 **Actions** 标签
2. 选择 **Docker Build and Push** workflow
3. 点击 **"Run workflow"**
4. 选择分支（main）
5. 点击 **"Run workflow"**

---

## 📋 配置检查清单

在推送标签前，确认：

- [ ] Docker Hub 账号已登录
- [ ] 访问令牌已创建
- [ ] 令牌权限包含 Read, Write, Delete
- [ ] GitHub Secrets 已添加（2个）
  - [ ] DOCKER_USERNAME
  - [ ] DOCKER_PASSWORD
- [ ] Secrets 名称拼写正确（区分大小写）
- [ ] 令牌已正确复制（没有多余空格）

---

## 🐛 常见问题

### Q1: 仍然提示 "Username and password required"

**检查**：
1. Secret 名称是否正确（DOCKER_USERNAME, DOCKER_PASSWORD）
2. 是否有多余的空格
3. 访问令牌是否有效

**重新添加**：
1. 删除旧的 Secrets
2. 重新创建 Docker Hub 访问令牌
3. 重新添加到 GitHub

### Q2: 推送失败 "denied: requested access to the resource is denied"

**原因**：仓库名称不匹配

**检查 workflow 文件**：
```yaml
env:
  IMAGE_NAME: DanZai233/mygallery  # 确保与 Docker Hub 仓库名一致
```

**Docker Hub 仓库**：
- 确保已创建仓库：`DanZai233/mygallery`
- 或在 Docker Hub 创建新仓库

### Q3: workflow 触发失败

**检查**：
1. workflow 文件是否在 main 分支
2. 标签格式是否正确（v1.0.0）
3. 是否推送了标签：`git push origin v1.0.0`

---

## 📝 完整的发布流程

### 使用自动发布脚本（推荐）

```bash
cd /root/MYGallery

# 1. 运行发布脚本
bash scripts/release.sh

# 2. 选择版本类型
#    选择 3 (patch) 创建 v1.0.1

# 3. 确认发布
#    输入 y

# 4. 等待推送完成
#    脚本会自动：
#    - 更新版本号
#    - 创建 Git tag
#    - 推送到 GitHub
#    - 触发 Actions

# 5. 查看构建进度
#    访问: https://github.com/DanZai233/mygallery/actions
```

### 手动发布

```bash
# 1. 创建标签
git tag -a v1.0.1 -m "Release v1.0.1"

# 2. 推送标签
git push origin v1.0.1

# 3. 查看 Actions
#    访问 GitHub Actions 页面
```

---

## 🎯 预期的 Actions 流程

配置正确后，推送标签会触发：

```
1. Checkout 代码
2. 设置 QEMU（多平台支持）
3. 设置 Docker Buildx
4. 登录 Docker Hub ← 使用你配置的 Secrets
5. 提取元数据（版本标签）
6. 构建 Docker 镜像（amd64, arm64）
7. 推送到 Docker Hub
8. 更新仓库描述
```

**完成后，可以使用**：
```bash
docker pull DanZai233/mygallery:1.0.1
docker pull DanZai233/mygallery:latest
```

---

## 📖 参考资料

### Docker Hub 访问令牌
- 文档：https://docs.docker.com/docker-hub/access-tokens/
- 管理：https://hub.docker.com/settings/security

### GitHub Secrets
- 文档：https://docs.github.com/en/actions/security-guides/encrypted-secrets
- 管理：`仓库设置 → Secrets and variables → Actions`

---

## 🎉 配置完成后

1. **推送一个测试标签**
   ```bash
   bash scripts/release.sh
   ```

2. **查看构建状态**
   ```
   https://github.com/DanZai233/mygallery/actions
   ```

3. **等待构建完成**（约 5-10 分钟）

4. **使用新镜像**
   ```bash
   docker pull DanZai233/mygallery:latest
   docker run -p 8080:8080 DanZai233/mygallery:latest
   ```

---

## 💡 小提示

### 保护你的令牌
- ❌ 不要提交到代码
- ❌ 不要分享给他人
- ✅ 只存储在 GitHub Secrets
- ✅ 定期轮换令牌

### 令牌权限
- **Read, Write, Delete** - 完整权限（推荐）
- **Read, Write** - 基本权限（够用）
- **Read** - 只能拉取（不够）

---

**按照以上步骤配置后，GitHub Actions 就能正常工作了！** 🚀

