# MYGallery - 个人照片墙系统

<div align="center">
  <h1>📷 MYGallery</h1>
  <p>一个简约、美观、功能完整的个人照片墙系统</p>
  <p>人人都可以自部署自己的照片展示空间</p>
  
  <p>
    <img src="https://img.shields.io/badge/Version-1.1.2-blue" alt="Version">
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go" alt="Go Version">
    <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
    <img src="https://img.shields.io/github/workflow/status/yourusername/mygallery/Docker%20Build%20and%20Push" alt="Build Status">
    <img src="https://img.shields.io/docker/image-size/yourusername/mygallery" alt="Docker Image Size">
    <img src="https://img.shields.io/docker/pulls/yourusername/mygallery" alt="Docker Pulls">
    <img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg" alt="PRs Welcome">
  </p>
</div>

---

## 🎯 项目特点

### 为什么选择 MYGallery？

- 🚀 **开箱即用**：一键部署，5分钟启动你的照片墙
- 🎨 **精美设计**：瀑布流布局 + 磨砂玻璃质感 + 灯箱效果
- 📸 **智能识别**：自动提取照片 EXIF 信息（相机、参数、GPS）
- 🔧 **灵活配置**：支持多种数据库和对象存储
- 🐳 **容器化部署**：Docker 一键部署，省心省力
- 🌍 **人人可用**：详细文档，无需专业知识

## ✨ 核心功能

### 📷 前台展示
- ✅ 瀑布流布局，自适应响应式设计
- ✅ 灯箱效果，支持键盘和手势导航
- ✅ EXIF 元数据展示（相机型号、光圈、快门、ISO、GPS）
- ✅ 磨砂玻璃质感的现代 UI 设计
- ✅ 图片懒加载，优化加载速度

### ⚙️ 后台管理
- ✅ JWT 安全认证
- ✅ 拖拽上传，支持批量处理
- ✅ 元数据编辑（标题、描述、标签、位置、版权）
- ✅ 照片管理（查看、编辑、删除）
- ✅ 统计面板（照片数、浏览量、存储空间）

### 🔧 技术特性
- ✅ **多数据库支持**：SQLite、MySQL、PostgreSQL
- ✅ **多存储支持**：本地、AWS S3、MinIO、阿里云 OSS
- ✅ **EXIF 自动提取**：相机信息、拍摄参数、GPS 位置
- ✅ **自动缩略图**：加快页面加载速度
- ✅ **Docker 部署**：容器化部署，环境一致

## 🚀 快速开始

### 方式一：一键安装（推荐）

```bash
# 1. 克隆仓库
git clone https://github.com/yourusername/mygallery.git
cd mygallery

# 2. 运行安装脚本
bash scripts/install.sh
```

### 方式二：Docker Compose

```bash
# 1. 克隆项目
git clone https://github.com/yourusername/mygallery.git
cd mygallery

# 2. 创建配置文件
cp config.example.yaml config.yaml

# 3. 启动服务
docker compose up -d
# 或使用旧版命令
docker-compose up -d
```

### 方式三：本地开发

```bash
# 1. 确保安装了 Go 1.21+
go version

# 2. 克隆项目
git clone https://github.com/yourusername/mygallery.git
cd mygallery

# 3. 安装依赖
go mod download

# 4. 创建配置
cp config.example.yaml config.yaml

# 5. 运行应用
go run main.go
```

### 访问应用

- 📷 **前台展示**：http://localhost:8080
- ⚙️ **后台管理**：http://localhost:8080/admin
- 👤 **默认账号**：admin / admin123

⚠️ **重要**：首次登录后请立即修改默认密码！

## 📖 文档

- 📚 [快速开始指南](QUICKSTART.md) - 5 分钟快速入门

## ⚙️ 配置说明

### 基础配置

编辑 `config.yaml` 文件：

```yaml
# 服务器配置
server:
  host: "0.0.0.0"
  port: 8080
  mode: "release"

# 管理员配置
admin:
  username: "admin"
  password: "admin123"  # ⚠️ 请修改

# JWT 配置
jwt:
  secret: "change-this-secret"  # ⚠️ 请修改
  expire_hours: 168
```

### 数据库配置

**SQLite（默认，适合个人使用）**
```yaml
database:
  type: "sqlite"
  sqlite:
    path: "./data/mygallery.db"
```

**MySQL（适合中等规模）**
```yaml
database:
  type: "mysql"
  mysql:
    host: "localhost"
    port: 3306
    username: "mygallery"
    password: "your_password"
    database: "mygallery"
```

**PostgreSQL（适合大规模）**
```yaml
database:
  type: "postgres"
  postgres:
    host: "localhost"
    port: 5432
    username: "mygallery"
    password: "your_password"
    database: "mygallery"
```

### 存储配置

**本地存储（默认）**
```yaml
storage:
  type: "local"
  local:
    upload_dir: "./uploads"
    thumbnail_dir: "./uploads/thumbnails"
```

**AWS S3**
```yaml
storage:
  type: "s3"
  s3:
    region: "us-east-1"
    bucket: "mygallery"
    access_key: "your_access_key"
    secret_key: "your_secret_key"
```

**MinIO（自托管对象存储）**
```yaml
storage:
  type: "minio"
  minio:
    endpoint: "localhost:9000"
    bucket: "mygallery"
    access_key: "minioadmin"
    secret_key: "minioadmin"
```

**阿里云 OSS**
```yaml
storage:
  type: "aliyun"
  aliyun:
    endpoint: "oss-cn-hangzhou.aliyuncs.com"
    bucket: "mygallery"
    access_key: "your_access_key"
    secret_key: "your_secret_key"
```

## 🛠️ 常用命令

### 使用 Makefile

```bash
make help          # 显示帮助
make init          # 初始化项目
make build         # 编译应用
make run           # 运行应用
make docker-build  # 构建 Docker 镜像
make docker-run    # 启动 Docker 容器
make docker-logs   # 查看日志
make docker-stop   # 停止容器
```

### 使用 Docker Compose

```bash
docker compose up -d          # 启动服务
docker compose down           # 停止服务
docker compose logs -f        # 查看日志
docker compose restart        # 重启服务
docker compose build          # 重新构建
docker compose ps             # 查看状态
```

## 📊 系统架构

```
┌─────────────┐
│   用户端    │
│  (浏览器)   │
└──────┬──────┘
       │
       ↓
┌─────────────────────────────┐
│      Gin Web Server         │
│   ┌──────────────────────┐  │
│   │  JWT 认证中间件       │  │
│   │  CORS 中间件          │  │
│   └──────────────────────┘  │
│   ┌──────────────────────┐  │
│   │  路由层 (Router)      │  │
│   └──────────────────────┘  │
│   ┌──────────────────────┐  │
│   │  处理器 (Handlers)    │  │
│   │  - AuthHandler        │  │
│   │  - PhotoHandler       │  │
│   └──────────────────────┘  │
└─────────┬───────────┬───────┘
          │           │
          ↓           ↓
    ┌─────────┐  ┌─────────┐
    │ 数据库  │  │ 存储层  │
    │ SQLite  │  │ Local   │
    │ MySQL   │  │ S3      │
    │Postgres │  │ MinIO   │
    └─────────┘  │ Aliyun  │
                 └─────────┘
```

## 🎨 界面预览

*（待添加截图）*

### 前台展示
- 瀑布流照片墙
- 灯箱大图预览
- EXIF 信息展示

### 后台管理
- 登录界面
- 照片上传
- 元数据编辑

## 🔐 安全建议

### 生产环境部署必做

1. **修改默认密码**
   ```yaml
   admin:
     password: "strong-password-here"
   ```

2. **修改 JWT 密钥**
   ```yaml
   jwt:
     secret: "long-random-secret-at-least-32-chars"
   ```

3. **启用 HTTPS**
   ```bash
   # 使用 Let's Encrypt
   sudo certbot --nginx -d your-domain.com
   ```

4. **配置防火墙**
   ```bash
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw enable
   ```

5. **定期备份数据**
   ```bash
   # 备份数据库
   cp data/mygallery.db data/backup-$(date +%Y%m%d).db
   
   # 备份照片
   tar -czf uploads-backup-$(date +%Y%m%d).tar.gz uploads/
   ```

## 📈 性能优化

### 1. 使用对象存储 + CDN
- 减轻服务器负载
- 加快图片加载速度
- 支持全球 CDN 加速

### 2. 启用缩略图
```yaml
image:
  thumbnail:
    width: 400
    height: 400
    quality: 85
```

### 3. 配置 Nginx 缓存
```nginx
location /uploads/ {
    expires 30d;
    add_header Cache-Control "public, immutable";
}
```

## 🤝 贡献

欢迎贡献代码！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 开源协议

本项目采用 [MIT](LICENSE) 协议。

## 🙏 致谢

感谢以下开源项目：

- [Gin](https://github.com/gin-gonic/gin) - Go Web 框架
- [GORM](https://gorm.io/) - Go ORM 库
- [lightGallery](https://www.lightgalleryjs.com/) - 灯箱插件
- [Tailwind CSS](https://tailwindcss.com/) - CSS 框架
- [goexif](https://github.com/rwcarlsen/goexif) - EXIF 解析库

## 📞 联系方式

- 📧 邮箱：your-email@example.com
- 🐛 问题反馈：[GitHub Issues](https://github.com/yourusername/mygallery/issues)
- 💬 讨论交流：[GitHub Discussions](https://github.com/yourusername/mygallery/discussions)

## ⭐ Star History

如果这个项目对你有帮助，请给个 ⭐️ Star 支持一下！

---

<div align="center">
  <p>Made with ❤️ by MYGallery Contributors</p>
  <p>© 2025 MYGallery. All rights reserved.</p>
</div>

