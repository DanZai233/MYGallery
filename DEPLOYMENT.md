# MYGallery 部署指南

本文档提供了 MYGallery 的详细部署步骤和最佳实践。

## 📋 目录

- [快速部署](#快速部署)
- [生产环境部署](#生产环境部署)
- [配置详解](#配置详解)
- [常见问题](#常见问题)

## 🚀 快速部署

### 方式一：一键安装脚本（推荐）

```bash
# 克隆仓库
git clone https://github.com/yourusername/mygallery.git
cd mygallery

# 运行安装脚本
bash scripts/install.sh
```

### 方式二：Docker Compose

```bash
# 1. 克隆仓库
git clone https://github.com/yourusername/mygallery.git
cd mygallery

# 2. 创建配置文件
cp config.example.yaml config.yaml
# 编辑配置文件（可选）
nano config.yaml

# 3. 启动服务
docker-compose up -d

# 4. 查看日志
docker-compose logs -f
```

### 方式三：Makefile

```bash
# 初始化项目
make init

# 启动 Docker 容器
make docker-run

# 查看日志
make docker-logs
```

## 🏭 生产环境部署

### 1. 使用 Nginx 反向代理

创建 Nginx 配置文件 `/etc/nginx/sites-available/mygallery`：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 强制使用 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL 证书配置
    ssl_certificate /path/to/your/certificate.crt;
    ssl_certificate_key /path/to/your/private.key;
    
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # 上传文件大小限制
    client_max_body_size 50M;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

启用配置：

```bash
sudo ln -s /etc/nginx/sites-available/mygallery /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 2. 使用对象存储（S3/MinIO/阿里云OSS）

修改 `config.yaml`：

#### MinIO 配置

```yaml
storage:
  type: "minio"
  minio:
    endpoint: "your-minio-server:9000"
    bucket: "mygallery"
    access_key: "your-access-key"
    secret_key: "your-secret-key"
    use_ssl: true
    url_prefix: "https://your-minio-server:9000/mygallery"
```

#### AWS S3 配置

```yaml
storage:
  type: "s3"
  s3:
    region: "us-east-1"
    bucket: "mygallery"
    access_key: "your-aws-access-key"
    secret_key: "your-aws-secret-key"
    url_prefix: "https://mygallery.s3.amazonaws.com"
```

#### 阿里云 OSS 配置

```yaml
storage:
  type: "aliyun"
  aliyun:
    endpoint: "oss-cn-hangzhou.aliyuncs.com"
    bucket: "mygallery"
    access_key: "your-aliyun-access-key"
    secret_key: "your-aliyun-secret-key"
    url_prefix: "https://mygallery.oss-cn-hangzhou.aliyuncs.com"
```

### 3. 使用 PostgreSQL 数据库

修改 `docker-compose.yml`，取消 PostgreSQL 服务的注释，然后修改 `config.yaml`：

```yaml
database:
  type: "postgres"
  postgres:
    host: "postgres"  # 或外部数据库地址
    port: 5432
    username: "mygallery"
    password: "your-secure-password"
    database: "mygallery"
    sslmode: "disable"
```

### 4. 设置环境变量

创建 `.env` 文件（如果使用敏感配置）：

```bash
# JWT 密钥（强烈建议修改）
JWT_SECRET=your-super-secret-jwt-key-change-this

# 管理员密码（首次启动后会被加密）
ADMIN_PASSWORD=your-secure-admin-password
```

## ⚙️ 配置详解

### 服务器配置

```yaml
server:
  host: "0.0.0.0"      # 监听地址
  port: 8080           # 监听端口
  mode: "release"      # 模式：debug, release, test
```

### 数据库配置

支持三种数据库：
- **SQLite**：适合个人使用，无需额外配置
- **MySQL**：适合中等规模，需要更好的并发性能
- **PostgreSQL**：适合大规模部署，功能最强大

### 存储配置

支持四种存储方式：
- **local**：本地存储，简单直接
- **s3**：AWS S3，全球 CDN
- **minio**：自托管对象存储
- **aliyun**：阿里云 OSS，国内访问速度快

### JWT 配置

```yaml
jwt:
  secret: "your-secret-key"  # ⚠️ 生产环境必须修改
  expire_hours: 168          # Token 有效期（小时）
```

### 图片处理配置

```yaml
image:
  max_upload_size: 52428800  # 最大上传大小（字节）
  allowed_types:
    - "image/jpeg"
    - "image/png"
    - "image/gif"
    - "image/webp"
  thumbnail:
    width: 400    # 缩略图宽度
    height: 400   # 缩略图高度
    quality: 85   # 缩略图质量
```

## 🔒 安全最佳实践

### 1. 修改默认密码

首次登录后立即修改管理员密码：
1. 登录后台管理：http://your-domain/admin
2. 使用默认账号登录（admin/admin123）
3. 在用户设置中修改密码

### 2. 修改 JWT 密钥

编辑 `config.yaml`，修改 JWT 密钥：

```yaml
jwt:
  secret: "use-a-long-random-string-here-at-least-32-characters"
```

生成随机密钥：

```bash
openssl rand -base64 32
```

### 3. 启用 HTTPS

使用 Let's Encrypt 免费证书：

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

### 4. 配置防火墙

```bash
# 允许 HTTP 和 HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# 如果使用 SSH，确保允许
sudo ufw allow 22/tcp

# 启用防火墙
sudo ufw enable
```

## 🔧 维护操作

### 备份数据

#### 备份 SQLite 数据库

```bash
# 备份数据库
cp data/mygallery.db data/mygallery.db.backup

# 或使用 SQLite 备份命令
sqlite3 data/mygallery.db ".backup data/mygallery.db.backup"
```

#### 备份上传文件

```bash
# 打包上传目录
tar -czf uploads-backup-$(date +%Y%m%d).tar.gz uploads/
```

### 更新应用

```bash
# 拉取最新代码
git pull origin main

# 重新构建镜像
docker-compose build

# 重启服务
docker-compose down
docker-compose up -d
```

### 查看日志

```bash
# 查看实时日志
docker-compose logs -f

# 查看最近 100 行日志
docker-compose logs --tail=100

# 查看特定服务日志
docker-compose logs mygallery
```

### 重启服务

```bash
# 重启所有服务
docker-compose restart

# 重启特定服务
docker-compose restart mygallery
```

## 🐛 常见问题

### Q: 上传照片失败

**A:** 检查以下几点：
1. 文件大小是否超过限制（默认 50MB）
2. 存储目录是否有写入权限
3. 磁盘空间是否充足
4. 如果使用对象存储，检查配置是否正确

### Q: 无法登录后台

**A:** 
1. 确认使用默认账号：admin / admin123
2. 检查数据库是否正常初始化
3. 查看日志：`docker-compose logs mygallery`

### Q: 图片无法显示

**A:**
1. 检查存储配置是否正确
2. 如果使用对象存储，确认 bucket 权限为公开读
3. 检查 URL 前缀配置是否正确

### Q: 如何迁移数据

**A:**
1. 备份数据库文件和上传目录
2. 在新服务器上部署 MYGallery
3. 恢复数据库和上传文件
4. 更新配置文件
5. 重启服务

### Q: 如何更改端口

**A:** 修改 `docker-compose.yml`：

```yaml
ports:
  - "your-port:8080"  # 例如："3000:8080"
```

## 📊 性能优化

### 1. 使用 CDN

如果使用对象存储，配置 CDN 加速：

```yaml
storage:
  s3:
    url_prefix: "https://cdn.your-domain.com"
```

### 2. 开启 Gzip 压缩

在 Nginx 配置中添加：

```nginx
gzip on;
gzip_vary on;
gzip_types text/plain text/css application/json application/javascript text/xml application/xml;
```

### 3. 配置缓存

在 Nginx 配置中添加：

```nginx
location /uploads/ {
    expires 30d;
    add_header Cache-Control "public, immutable";
}
```

## 📞 获取帮助

- 📖 文档：查看 [README.md](README.md)
- 🐛 问题：提交 [GitHub Issues](https://github.com/yourusername/mygallery/issues)
- 💬 讨论：加入 [Discussions](https://github.com/yourusername/mygallery/discussions)

---

祝你部署愉快！🎉

