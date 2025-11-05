# 🔧 MYGallery 故障排除指南

本文档列出了常见问题及其解决方案。

## 📋 目录

- [Docker 相关问题](#docker-相关问题)
- [构建问题](#构建问题)
- [运行时问题](#运行时问题)
- [上传问题](#上传问题)
- [数据库问题](#数据库问题)
- [存储问题](#存储问题)

## Docker 相关问题

### ❌ docker-compose.yml: the attribute `version` is obsolete

**问题**：Docker Compose 新版本不再需要 version 字段

**解决方案**：已在最新版本中移除 `version` 字段

```bash
# 如果还有此警告，确保使用最新代码
git pull origin main
```

### ❌ compose build requires buildx 0.17 or later

**问题**：Docker Buildx 版本过低

**解决方案**：

```bash
# 方案 1: 更新 Docker
# Ubuntu/Debian
sudo apt update
sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# macOS
brew upgrade docker

# 方案 2: 使用传统构建方式
# 修改 Dockerfile，移除多阶段构建的某些特性

# 方案 3: 使用旧版 docker-compose 命令
docker-compose build --no-cache
docker-compose up -d
```

### ❌ Cannot connect to the Docker daemon

**问题**：Docker 服务未运行

**解决方案**：

```bash
# Linux
sudo systemctl start docker
sudo systemctl enable docker

# macOS
# 打开 Docker Desktop 应用

# 检查 Docker 状态
docker ps
```

### ❌ Permission denied while trying to connect to Docker daemon

**问题**：当前用户没有 Docker 权限

**解决方案**：

```bash
# 将用户添加到 docker 组
sudo usermod -aG docker $USER

# 重新登录或运行
newgrp docker

# 临时使用 sudo
sudo docker-compose up -d
```

## 构建问题

### ❌ go.mod file not found

**问题**：Go 模块文件缺失

**解决方案**：

```bash
# 初始化 Go 模块
cd /root/MYGallery
go mod init github.com/mygallery/mygallery
go mod tidy

# 或重新克隆仓库
git clone <your-repo> MYGallery-new
cd MYGallery-new
```

### ❌ package not found

**问题**：Go 依赖包缺失

**解决方案**：

```bash
# 下载所有依赖
go mod download

# 清理并重新下载
go clean -modcache
go mod download

# 更新依赖
go mod tidy
```

### ❌ CGO_ENABLED error

**问题**：缺少 C 编译器（SQLite 需要）

**解决方案**：

```bash
# Alpine Linux
apk add gcc musl-dev

# Ubuntu/Debian
apt install build-essential

# macOS
xcode-select --install
```

## 运行时问题

### ❌ 端口已被占用

**问题**：8080 端口被其他程序使用

**解决方案**：

```bash
# 方案 1: 修改端口
# 编辑 docker-compose.yml
ports:
  - "3000:8080"  # 改为其他端口

# 方案 2: 查找并关闭占用端口的进程
# Linux/macOS
sudo lsof -i :8080
sudo kill -9 <PID>

# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

### ❌ 无法访问应用

**问题**：服务启动但无法访问

**解决方案**：

```bash
# 1. 检查容器状态
docker ps
docker logs mygallery

# 2. 检查防火墙
sudo ufw status
sudo ufw allow 8080

# 3. 检查配置文件
cat config.yaml

# 4. 使用 localhost 或 127.0.0.1 访问
curl http://localhost:8080/health

# 5. 如果是远程服务器，检查云服务商安全组规则
```

### ❌ 登录失败

**问题**：无法登录后台管理

**解决方案**：

```bash
# 1. 确认使用默认账号
# 用户名: admin
# 密码: admin123

# 2. 检查数据库
docker exec -it mygallery sh
ls -la /app/data/

# 3. 重置数据库
docker-compose down
rm -rf data/*.db
docker-compose up -d

# 4. 查看日志
docker-compose logs mygallery | grep -i error
```

## 上传问题

### ❌ 照片上传失败

**问题**：无法上传照片

**解决方案**：

```bash
# 1. 检查文件大小（默认限制 50MB）
# 修改 config.yaml
image:
  max_upload_size: 104857600  # 100MB

# 2. 检查文件类型
# 确保是支持的格式: JPG, PNG, GIF, WebP

# 3. 检查磁盘空间
df -h

# 4. 检查目录权限
ls -la uploads/
sudo chown -R 1000:1000 uploads/

# 5. 如果使用 Nginx，检查上传限制
# /etc/nginx/nginx.conf
client_max_body_size 50M;

# 6. 查看详细错误
# 打开浏览器开发者工具 -> Network 标签
```

### ❌ 上传后看不到图片

**问题**：上传成功但图片不显示

**解决方案**：

```bash
# 1. 检查存储配置
cat config.yaml | grep storage

# 2. 检查文件是否存在
ls -la uploads/

# 3. 检查 URL 配置
# config.yaml 中的 url_prefix 要正确

# 4. 检查浏览器控制台错误
# F12 -> Console 标签

# 5. 如果使用对象存储，检查权限
# S3/MinIO/OSS bucket 需要设置公开读权限
```

## 数据库问题

### ❌ SQLite database is locked

**问题**：数据库被锁定

**解决方案**：

```bash
# 1. 重启应用
docker-compose restart

# 2. 检查是否有多个实例运行
docker ps | grep mygallery

# 3. 删除锁文件
rm data/mygallery.db-journal

# 4. 如果频繁出现，考虑使用 MySQL/PostgreSQL
```

### ❌ MySQL connection refused

**问题**：无法连接 MySQL

**解决方案**：

```bash
# 1. 确保 MySQL 服务运行
docker-compose ps

# 2. 检查连接配置
# config.yaml
database:
  mysql:
    host: "mysql"  # 使用服务名，不是 localhost
    port: 3306

# 3. 等待 MySQL 启动完成
# 在 docker-compose.yml 中添加健康检查

# 4. 测试连接
docker exec -it mygallery sh
ping mysql
```

### ❌ PostgreSQL authentication failed

**问题**：PostgreSQL 认证失败

**解决方案**：

```bash
# 1. 检查密码
cat config.yaml | grep postgres

# 2. 确保数据库已创建
docker exec -it postgres psql -U postgres
CREATE DATABASE mygallery;

# 3. 检查环境变量
docker-compose config
```

## 存储问题

### ❌ S3 access denied

**问题**：无法访问 S3

**解决方案**：

```bash
# 1. 检查 Access Key 和 Secret Key
# 2. 检查 bucket 权限
# 3. 检查 IAM 策略
# 4. 确保 region 正确
# 5. 测试连接
aws s3 ls s3://your-bucket --region us-east-1
```

### ❌ MinIO connection error

**问题**：无法连接 MinIO

**解决方案**：

```bash
# 1. 启动 MinIO 服务
docker-compose up -d minio

# 2. 访问 MinIO 控制台
# http://localhost:9001

# 3. 创建 bucket

# 4. 设置公开读权限
mc policy set public myminio/mygallery

# 5. 检查配置
# endpoint 应该是 "minio:9000"（容器内）
# 或 "localhost:9000"（本地）
```

### ❌ 阿里云 OSS 403 Forbidden

**问题**：OSS 访问被拒绝

**解决方案**：

```bash
# 1. 检查 AccessKey 权限
# 2. 检查 bucket 读写权限
# 3. 检查跨域配置（CORS）
# 4. 检查 Referer 白名单
# 5. 使用 ossutil 测试
ossutil ls oss://your-bucket
```

## 性能问题

### ⚠️ 加载速度慢

**解决方案**：

1. **启用缩略图**（已默认启用）
2. **使用 CDN**
   ```yaml
   storage:
     s3:
       url_prefix: "https://cdn.your-domain.com"
   ```
3. **启用 Nginx 缓存**
   ```nginx
   location /uploads/ {
       expires 30d;
       add_header Cache-Control "public, immutable";
   }
   ```
4. **压缩图片**
   - 上传前使用工具压缩
   - 调整缩略图质量

### ⚠️ 内存使用过高

**解决方案**：

```bash
# 1. 限制内存使用
# docker-compose.yml
services:
  mygallery:
    deploy:
      resources:
        limits:
          memory: 512M

# 2. 使用对象存储而不是本地存储
# 3. 定期清理日志
docker-compose logs --tail=100 > /dev/null
```

## 日志和调试

### 查看日志

```bash
# 查看实时日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs mygallery

# 查看最近 100 行
docker-compose logs --tail=100

# 保存日志到文件
docker-compose logs > logs.txt
```

### 进入容器调试

```bash
# 进入容器
docker exec -it mygallery sh

# 查看文件
ls -la
cat config.yaml

# 测试网络
ping google.com
curl http://localhost:8080/health

# 查看进程
ps aux

# 退出
exit
```

### 启用调试模式

```yaml
# config.yaml
server:
  mode: "debug"  # 显示详细日志
```

## 完全重置

如果一切都不工作，尝试完全重置：

```bash
# 1. 停止所有容器
docker-compose down

# 2. 删除数据（⚠️ 谨慎操作）
rm -rf data/*.db
rm -rf uploads/*

# 3. 删除镜像
docker rmi mygallery:latest

# 4. 清理 Docker 缓存
docker system prune -a

# 5. 重新构建
docker-compose build --no-cache

# 6. 重新启动
docker-compose up -d
```

## 获取帮助

如果以上方法都无法解决问题：

1. **检查日志**：`docker-compose logs`
2. **提交 Issue**：https://github.com/yourusername/mygallery/issues
3. **查看文档**：README.md, DEPLOYMENT.md
4. **社区讨论**：GitHub Discussions

提交 Issue 时请包含：
- 操作系统和版本
- Docker 版本
- 错误信息和日志
- 复现步骤

---

祝你早日解决问题！🎉

