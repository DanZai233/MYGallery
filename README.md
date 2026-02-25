# MYGallery - 个人照片墙系统

<div align="center">
  <h1>📷 MYGallery</h1>
  <p>一个简约、美观、功能完整的个人照片墙系统</p>
  <p>人人都可以自部署自己的照片展示空间</p>
  
  <p>
    <img src="https://img.shields.io/badge/Version-1.1.4-blue" alt="Version">
    <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go" alt="Go Version">
    <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
  </p>
</div>

---

## ✨ 特性

### 🎨 精美的前端展示
- **瀑布流布局**：响应式设计，自适应各种屏幕尺寸（4/3/2 列自动切换）
- **灯箱效果**：点击查看大图，支持键盘导航和缩放
- **毛玻璃质感**：现代化 Glassmorphism UI 设计
- **EXIF 元数据展示**：自动显示相机/手机型号、拍摄参数、位置等信息
- **分类筛选**：按分类快速筛选照片，移动端支持横向滑动
- **搜索功能**：支持搜索标题、描述、标签、位置等
- **黑夜模式**：一键切换白天/黑夜主题，记忆用户偏好
- **Live Photo 支持**：支持 Apple Live Photo 标识展示

### 📱 移动端友好
- **响应式导航**：移动端自动切换为汉堡菜单
- **触摸友好**：照片卡片支持触摸交互，长按显示详情
- **横向滑动分类**：分类过多时自动支持横向滚动
- **滚动回顶**：一键回到顶部按钮
- **PWA 优化**：支持添加到主屏幕

### ⚙️ 强大的后台管理
- **用户登录验证**：JWT token 认证，安全可靠
- **图片上传**：支持拖拽上传，批量上传，实时进度显示
- **EXIF 自动提取**：自动读取照片的相机参数、GPS 位置等元数据
- **增强 EXIF 解析**：支持相机和手机拍摄的照片，提取更多元数据字段
- **Live Photo 上传**：支持同时上传照片和配套视频
- **分类管理**：创建、编辑、删除照片分类
- **网站设置**：自定义网站信息、备案信息、Header/Footer 等

### 🔧 灵活的配置系统
- **多数据库支持**：SQLite（默认）、MySQL、PostgreSQL
- **多存储支持**：本地存储（默认）、AWS S3、MinIO、阿里云 OSS
- **云存储缩略图**：所有存储后端均支持自动生成和上传缩略图
- **YAML 配置**：简单直观的配置文件
- **Docker 部署**：一键部署，开箱即用

---

## 📦 快速开始

### 环境要求

- **Go 1.24+**（需要 CGO 支持，用于 SQLite）
- **GCC**（C 编译器，SQLite 驱动需要）
- **Git**

### 本地运行（推荐开发使用）

```bash
# 1. 克隆仓库
git clone https://github.com/danzai233/mygallery.git
cd mygallery

# 2. 安装 Go 依赖
go mod download

# 3. 创建配置文件
cp config.example.yaml config.yaml

# 4. 初始化项目目录
make init

# 5. 运行应用
go run main.go
```

访问应用：
- 📷 **前台展示**：http://localhost:8080
- ⚙️ **后台管理**：http://localhost:8080/admin
- 👤 **默认账号**：`admin` / `admin123`

> ⚠️ **首次登录后请立即修改默认密码！**

### Docker 部署（推荐生产使用）

```bash
# 1. 克隆仓库
git clone https://github.com/danzai233/mygallery.git
cd mygallery

# 2. 创建配置
cp config.example.yaml config.yaml
# 编辑 config.yaml，修改管理员密码和 JWT 密钥

# 3. 构建并启动
docker compose build
docker compose up -d

# 查看日志
docker compose logs -f

# 停止服务
docker compose down
```

### 使用 Makefile

```bash
make help       # 查看所有可用命令
make init       # 初始化项目（创建配置文件和目录）
make deps       # 安装依赖
make build      # 编译二进制
make run        # 运行应用
make dev        # 开发模式（需安装 air 热重载）
make test       # 运行测试
make clean      # 清理构建文件
```

---

## 🔧 配置详解

配置文件为项目根目录下的 `config.yaml`，从 `config.example.yaml` 复制修改。

### 服务器配置

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "release"  # debug(开发), release(生产), test(测试)
```

### 数据库配置

支持三种数据库，通过 `database.type` 切换：

#### SQLite（默认，零配置）

```yaml
database:
  type: "sqlite"
  sqlite:
    path: "./data/mygallery.db"
```

#### MySQL

```yaml
database:
  type: "mysql"
  mysql:
    host: "localhost"
    port: 3306
    username: "root"
    password: "your-password"
    database: "mygallery"
    charset: "utf8mb4"
```

#### PostgreSQL

```yaml
database:
  type: "postgres"
  postgres:
    host: "localhost"
    port: 5432
    username: "postgres"
    password: "your-password"
    database: "mygallery"
    sslmode: "disable"
```

### 存储配置

支持四种存储后端，通过 `storage.type` 切换。所有后端均支持自动缩略图生成。

#### 本地存储（默认）

```yaml
storage:
  type: "local"
  local:
    upload_dir: "./uploads"
    thumbnail_dir: "./uploads/thumbnails"
    url_prefix: "/uploads"
```

#### AWS S3

```yaml
storage:
  type: "s3"
  s3:
    region: "us-east-1"
    bucket: "mygallery"
    access_key: "YOUR_ACCESS_KEY"
    secret_key: "YOUR_SECRET_KEY"
    endpoint: ""              # 留空使用 AWS 默认，或自定义
    url_prefix: "https://your-bucket.s3.amazonaws.com"
```

#### MinIO（S3 兼容的自建对象存储）

```yaml
storage:
  type: "minio"
  minio:
    endpoint: "localhost:9000"
    bucket: "mygallery"
    access_key: "minioadmin"
    secret_key: "minioadmin"
    use_ssl: false
    url_prefix: "http://localhost:9000/mygallery"
```

#### 阿里云 OSS

```yaml
storage:
  type: "aliyun"
  aliyun:
    endpoint: "oss-cn-hangzhou.aliyuncs.com"
    bucket: "mygallery"
    access_key: "YOUR_ACCESS_KEY"
    secret_key: "YOUR_SECRET_KEY"
    url_prefix: "https://mygallery.oss-cn-hangzhou.aliyuncs.com"
```

### JWT 与管理员配置

```yaml
jwt:
  secret: "change-this-to-a-random-string"  # 生产环境务必修改！
  expire_hours: 168  # Token 有效期（7天）

admin:
  username: "admin"
  password: "admin123"  # 首次启动后自动加密
  email: "admin@example.com"
```

### 图片处理配置

```yaml
image:
  max_upload_size: 52428800  # 50MB
  allowed_types:
    - "image/jpeg"
    - "image/png"
    - "image/gif"
    - "image/webp"
  thumbnail:
    width: 400
    height: 400
    quality: 85
```

---

## 📸 EXIF 元数据支持

MYGallery 自动从上传的照片中提取 EXIF 元数据，兼容相机和手机拍摄的照片：

| 字段 | 说明 | 相机 | iPhone | Android |
|------|------|:----:|:------:|:-------:|
| 相机品牌 | Camera Make | ✅ | ✅ | ✅ |
| 相机型号 | Camera Model | ✅ | ✅ | ✅ |
| 镜头型号 | Lens Model | ✅ | ✅ | 部分 |
| 焦距 | Focal Length | ✅ | ✅ | ✅ |
| 光圈 | F-Number | ✅ | ✅ | ✅ |
| 快门速度 | Exposure Time | ✅ | ✅ | ✅ |
| ISO | ISO Speed | ✅ | ✅ | ✅ |
| 拍摄时间 | Date Taken | ✅ | ✅ | ✅ |
| GPS 位置 | GPS Coordinates | ✅ | ✅ | ✅ |
| 拍摄软件 | Software | ✅ | ✅ | ✅ |
| 图片方向 | Orientation | ✅ | ✅ | ✅ |
| 白平衡 | White Balance | ✅ | ✅ | ✅ |
| 闪光灯 | Flash | ✅ | ✅ | ✅ |
| 曝光模式 | Exposure Program | ✅ | ✅ | 部分 |
| 测光模式 | Metering Mode | ✅ | ✅ | ✅ |
| 曝光补偿 | Exposure Bias | ✅ | ✅ | ✅ |
| 色彩空间 | Color Space | ✅ | ✅ | ✅ |

> 缩略图生成时会自动根据 EXIF 方向标记旋转图片，确保显示方向正确。

---

## 🎬 Live Photo 支持

MYGallery 支持 Apple Live Photo 的上传和展示：

- **上传**：在后台上传照片时，可同时上传配套的 `.mov` 视频文件
- **展示**：前端照片卡片左上角显示 `LIVE` 徽章
- **API**：通过 `POST /api/photos` 的 `live_photo` 字段上传视频

```bash
# 上传 Live Photo 示例
curl -X POST http://localhost:8080/api/photos \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "photo=@IMG_1234.jpg" \
  -F "live_photo=@IMG_1234.mov" \
  -F "title=我的 Live Photo"
```

---

## 🌐 API 接口

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/photos` | 获取照片列表（支持 `page`、`size`、`category`、`search` 参数） |
| GET | `/api/photos/:id` | 获取单张照片详情 |
| GET | `/api/categories` | 获取分类列表 |
| GET | `/api/settings` | 获取网站设置 |
| GET | `/health` | 健康检查 |

### 需认证接口（需要 `Authorization: Bearer <token>` 头）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/auth/login` | 管理员登录 |
| POST | `/api/auth/change-password` | 修改密码 |
| POST | `/api/photos` | 上传照片（支持 Live Photo） |
| PUT | `/api/photos/:id` | 更新照片信息 |
| DELETE | `/api/photos/:id` | 删除照片 |
| POST | `/api/categories` | 创建分类 |
| PUT | `/api/categories/:id` | 更新分类 |
| DELETE | `/api/categories/:id` | 删除分类 |
| PUT | `/api/settings` | 更新网站设置 |

---

## 🚀 部署指南

### 方式一：直接运行二进制

```bash
# 编译
go build -o mygallery main.go

# 创建配置
cp config.example.yaml config.yaml
# 编辑 config.yaml

# 创建目录
mkdir -p data uploads uploads/thumbnails

# 运行
./mygallery
```

### 方式二：Docker Compose（推荐）

```bash
cp config.example.yaml config.yaml
docker compose build
docker compose up -d
```

### 方式三：Nginx 反向代理

```nginx
server {
    listen 80;
    server_name gallery.example.com;
    client_max_body_size 50M;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 方式四：Systemd 服务

```ini
# /etc/systemd/system/mygallery.service
[Unit]
Description=MYGallery Photo Gallery
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/mygallery
ExecStart=/opt/mygallery/mygallery
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable mygallery
sudo systemctl start mygallery
```

---

## 🔒 安全建议

1. **修改默认密码** — 首次登录后立即修改
2. **修改 JWT 密钥** — `config.yaml` 中的 `jwt.secret`
3. **启用 HTTPS** — 生产环境使用 Nginx + Let's Encrypt
4. **配置防火墙** — 只开放必要端口
5. **定期备份数据** — 备份 `data/` 和 `uploads/` 目录

---

## 🛠️ 技术栈

| 组件 | 技术 |
|------|------|
| 后端框架 | Go + Gin |
| ORM | GORM |
| 认证 | JWT (golang-jwt) |
| EXIF 解析 | rwcarlsen/goexif |
| 图片处理 | disintegration/imaging |
| 数据库 | SQLite / MySQL / PostgreSQL |
| 对象存储 | Local / AWS S3 / MinIO / 阿里云 OSS |
| 前端 | 原生 JavaScript + Tailwind CSS + lightGallery |

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

---

<div align="center">
  <p>Made with ❤️ by MYGallery</p>
  <p>如果这个项目对你有帮助，请给个 ⭐ Star 吧！</p>
</div>
