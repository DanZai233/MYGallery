# MYGallery - 个人照片墙系统

<div align="center">
  <h1>📷 MYGallery</h1>
  <p>一个简约、美观、功能完整的个人照片墙系统</p>
  <p>
    <img src="https://img.shields.io/badge/Version-1.1.4-blue" alt="Version">
    <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go" alt="Go Version">
    <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
  </p>
</div>

---

## 📸 界面预览

### 暗色模式 · 瀑布流画廊

默认暗色主题，5 列紧凑瀑布流布局，侧栏 SVG 简笔画图标导航。

![暗色模式画廊](docs/screenshots/gallery-dark.webp)

### 亮色模式

一键切换亮色/暗色主题，偏好自动记忆。

![亮色模式画廊](docs/screenshots/gallery-light.webp)

### 灯箱查看器 · EXIF 信息面板

ChronoFrame 风格全屏查看器，右侧信息面板分区展示 EXIF 元数据，底部缩略图条可快速跳转。

![灯箱查看器](docs/screenshots/lightbox-info.webp)

### 灯箱 · GPS 迷你地图

当照片包含 GPS 坐标时，信息面板顶部嵌入 Leaflet 迷你地图，实时显示拍摄位置。

![GPS 迷你地图](docs/screenshots/lightbox-minimap.webp)

### 灯箱 · 全屏模式

点击 ℹ 按钮关闭信息面板后，照片自动铺满全屏居中展示。

![全屏模式](docs/screenshots/lightbox-fullwidth.webp)

### 表态系统

访客可对照片进行 emoji 表态（👍❤️😍😂😮😢🔥✨），无需登录，基于指纹识别。

| 表态选择器 | 表态已添加 |
|:---:|:---:|
| ![表态选择器](docs/screenshots/reactions.webp) | ![表态已添加](docs/screenshots/reaction-added.webp) |

### 照片地图

基于 Leaflet + OpenStreetMap 的全屏地图，展示带 GPS 坐标的照片位置，支持聚合标记。

| 地图标记 | 标记弹窗 |
|:---:|:---:|
| ![地图标记](docs/screenshots/map-markers.webp) | ![标记弹窗](docs/screenshots/map-popup.webp) |

### 后台管理

ChronoFrame 风格侧栏导航，仪表盘统计，拖拽上传，照片/分类管理。

![后台管理](docs/screenshots/admin.webp)

---

## ✨ 特性

### 🎨 前端展示
- **瀑布流布局** — 5/4/3/2 列自适应，4px 紧凑间距
- **暗色模式默认** — ChronoFrame 风格中性配色
- **侧栏图标导航** — SVG 简笔画图标，极简设计
- **灯箱查看器** — 全屏查看 + 右侧信息面板 + 底部缩略图导航
- **GPS 迷你地图** — 信息面板嵌入 Leaflet 地图显示拍摄位置
- **表态系统** — 8 种 emoji 表态，指纹识别匿名用户
- **分享链接** — 一键复制照片直链（`?photo=ID`），支持 Deep Link
- **分类筛选** — 按分类快速筛选，横向滑动支持
- **搜索功能** — 搜索标题、描述、标签、位置
- **亮暗切换** — 一键切换，记忆用户偏好

### 📱 移动端友好
- 移动端底部栏导航
- 卡片触摸交互，灯箱触摸滑动
- 移动端信息面板自适应为底部弹出
- 分类横向滚动

### ⚙️ 后台管理
- 侧栏导航：仪表盘、分类管理、网站设置、地图、前台
- 拖拽上传，批量上传，实时进度
- EXIF 自动提取 + 手动编辑
- Live Photo 上传（照片 + 配套视频）

### 📍 照片地图
- Leaflet + OpenStreetMap，无需 API Key
- 照片缩略图标记 + 弹出详情卡
- 密集区域聚合标记
- 缩放控制 + 重置

### 📸 EXIF 元数据

兼容相机（Nikon / Canon / Sony / Fujifilm）和手机（iPhone / Android），提取字段：

| 字段 | 说明 |
|------|------|
| 品牌 / 型号 | Camera Make / Model |
| 镜头型号 | Lens Model |
| 焦距 / 光圈 / 快门 / ISO | Focal / Aperture / Shutter / ISO |
| GPS 坐标 | 纬度 / 经度（自动过滤 0,0） |
| 拍摄时间 | DateTimeOriginal |
| 软件 / 固件 | Software |
| 白平衡 / 闪光灯 | White Balance / Flash |
| 曝光模式 / 测光模式 | Exposure / Metering |
| 曝光补偿 / 色彩空间 | Exposure Bias / Color Space |
| 场景类型 / 方向 | Scene Type / Orientation |

缩略图生成时自动根据 EXIF 方向旋转。

### 🎬 Live Photo
- 上传时附加配套 `.mov` 视频
- 前端 LIVE 徽章标识
- API 支持 `live_photo` 字段

### 🔧 灵活配置
- **多数据库** — SQLite（默认）、MySQL、PostgreSQL
- **多存储** — 本地（默认）、AWS S3、MinIO、阿里云 OSS
- **云端缩略图** — 所有存储后端支持缩略图自动上传

---

## 📦 快速开始

### 环境要求

- **Go 1.24+**（需 CGO，SQLite 驱动依赖）
- **GCC**（C 编译器）

### 本地运行

```bash
git clone https://github.com/danzai233/mygallery.git
cd mygallery
go mod download
cp config.example.yaml config.yaml
make init
go run main.go
```

访问：
| 页面 | 地址 |
|------|------|
| 📷 前台画廊 | http://localhost:8080 |
| ⚙️ 后台管理 | http://localhost:8080/admin |
| 🗺️ 照片地图 | http://localhost:8080/map.html |
| 👤 默认账号 | `admin` / `admin123` |

### Docker 部署

#### 方式一：一键部署（推荐，使用预构建镜像）

使用官方 Docker Hub 镜像，本地 SQLite 数据库，最简单快速。

```bash
# 1. 克隆项目
git clone https://github.com/danzai233/mygallery.git
cd mygallery

# 2. 复制配置文件
cp config.example.yaml config.yaml

# 3. 启动服务
docker compose -f docker-compose.simple.yml up -d

# 4. 查看日志
docker compose -f docker-compose.simple.yml logs -f
```

访问 http://localhost:8080，默认账号：`admin` / `admin123`。

#### 方式二：Docker Compose 完整部署（可选 PostgreSQL + MinIO）

使用 PostgreSQL 数据库和 MinIO 对象存储，适合生产环境。

```bash
# 1. 复制配置文件
cp config.example.yaml config.yaml

# 2. 编辑 docker-compose.full.yml，取消注释需要的服务（postgres/minio）

# 3. 创建 .env 文件设置密码（可选）
cat > .env <<EOF
POSTGRES_PASSWORD=your_strong_password
MINIO_ROOT_USER=your_minio_user
MINIO_ROOT_PASSWORD=your_minio_password
TZ=Asia/Shanghai
EOF

# 4. 启动服务
docker compose -f docker-compose.full.yml up -d
```

如果启用了 PostgreSQL，需要修改 `config.yaml`：

```yaml
database:
  type: "postgres"
  postgres:
    host: "postgres"          # Docker Compose 网络中的服务名
    port: 5432
    username: "mygallery"
    password: "your_strong_password"
    database: "mygallery"
    sslmode: "disable"
```

如果启用了 MinIO，需要修改 `config.yaml`：

```yaml
storage:
  type: "minio"
  minio:
    endpoint: "minio:9000"   # Docker Compose 网络中的服务名
    bucket: "mygallery"
    access_key: "your_minio_user"
    secret_key: "your_minio_password"
    use_ssl: false
    url_prefix: "http://localhost:9000/mygallery"
```

首次启动 MinIO 后，需要通过 http://localhost:9001 进入控制台手动创建 bucket。

#### 方式三：自定义镜像构建

如果需要自定义构建（例如修改源码）：

```bash
# 构建镜像
docker build -t mygallery:custom .

# 运行容器
docker run -d \
  --name mygallery \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/uploads:/app/uploads \
  -e TZ=Asia/Shanghai \
  mygallery:custom
```

##### 自动化构建和发布

项目配置了 GitHub Actions，支持自动构建和发布 Docker 镜像：

- **触发条件**：推送带有 `v*.*.*` 格式的 tag 时自动构建（如 `v1.1.5`）
- **支持平台**：`linux/amd64` 和 `linux/arm64`
- **镜像仓库**：[Docker Hub](https://hub.docker.com/r/danzai233/mygallery)

手动发布新版本：

```bash
# 使用 Makefile
make docker-login  # 登录 Docker Hub
make docker-build-multi-push VERSION=v1.1.5

# 或使用发布脚本
./scripts/publish-docker.sh v1.1.5
```

#### Docker 部署常用命令

```bash
# 查看容器状态
docker compose -f docker-compose.simple.yml ps

# 查看日志
docker compose -f docker-compose.simple.yml logs -f

# 停止服务
docker compose -f docker-compose.simple.yml stop

# 启动服务
docker compose -f docker-compose.simple.yml start

# 重启服务
docker compose -f docker-compose.simple.yml restart

# 停止并删除容器
docker compose -f docker-compose.simple.yml down

# 停止并删除容器及数据卷（慎用，会删除所有数据）
docker compose -f docker-compose.simple.yml down -v

# 进入容器调试
docker exec -it mygallery sh

# 更新镜像到最新版本
docker pull danzai233/mygallery:latest
docker compose -f docker-compose.simple.yml up -d
```

#### Docker 数据持久化

以下目录通过 volume 挂载持久化数据：

- `./data` - SQLite 数据库文件（如果使用 SQLite）
- `./uploads` - 上传的图片和缩略图
- `./config.yaml` - 配置文件（建议使用 `:ro` 只读挂载）

**备份重要数据**：

```bash
# 备份数据
tar -czf mygallery-backup-$(date +%Y%m%d).tar.gz data/ uploads/ config.yaml

# 恢复数据
tar -xzf mygallery-backup-20240326.tar.gz
```

### Makefile

```bash
make help                      # 显示所有命令
make run                       # 运行应用
make build                     # 编译应用
make test                      # 运行测试
make dev                       # 热重载开发模式

# Docker 相关命令
make docker-build              # 构建 Docker 镜像
make docker-build-multi        # 构建多平台镜像
make docker-build-multi-push   # 构建并推送多平台镜像
make docker-login              # 登录 Docker Hub
make docker-up-simple          # 启动简化版部署
make docker-up-full            # 启动完整版部署
make docker-down-simple        # 停止简化版部署
make docker-logs-simple        # 查看简化版日志
```

---

## 🔧 配置详解

配置文件为项目根目录下的 `config.yaml`，首次使用请从示例文件复制：

```bash
cp config.example.yaml config.yaml
```

以下逐项说明每个配置块。**只需修改你用到的部分**，未用到的保持默认即可。

### 服务器配置

```yaml
server:
  host: "0.0.0.0"       # 监听地址，0.0.0.0 表示所有网卡
  port: 8080             # 监听端口
  mode: "release"        # 运行模式：debug（开发，输出详细日志）/ release（生产）/ test
```

> 开发时建议设为 `debug`，可以看到每个请求的详细日志。

### 数据库配置

通过 `database.type` 切换数据库类型，**只需配置你选用的那一种**。

#### SQLite（默认，零依赖）

```yaml
database:
  type: "sqlite"
  sqlite:
    path: "./data/mygallery.db"   # 数据库文件路径，会自动创建
```

> SQLite 不需要额外安装数据库服务，适合个人部署。数据库文件保存在 `data/` 目录下。

#### MySQL

```yaml
database:
  type: "mysql"
  mysql:
    host: "localhost"        # MySQL 主机地址
    port: 3306               # 端口
    username: "root"         # 用户名
    password: "your-pass"    # 密码
    database: "mygallery"    # 数据库名（需提前创建）
    charset: "utf8mb4"       # 字符集，建议 utf8mb4 以支持 emoji
```

> 需要提前创建数据库：`CREATE DATABASE mygallery CHARACTER SET utf8mb4;`

#### PostgreSQL

```yaml
database:
  type: "postgres"
  postgres:
    host: "localhost"        # PostgreSQL 主机地址
    port: 5432               # 端口
    username: "postgres"     # 用户名
    password: "your-pass"    # 密码
    database: "mygallery"    # 数据库名（需提前创建）
    sslmode: "disable"       # SSL 模式：disable / require / verify-full
```

> 需要提前创建数据库：`CREATE DATABASE mygallery;`

### 存储配置

通过 `storage.type` 切换存储后端。所有后端均支持自动生成和上传缩略图。

#### 本地存储（默认）

```yaml
storage:
  type: "local"
  local:
    upload_dir: "./uploads"              # 原图存放目录
    thumbnail_dir: "./uploads/thumbnails" # 缩略图目录
    url_prefix: "/uploads"               # 前端访问 URL 前缀
```

#### AWS S3

```yaml
storage:
  type: "s3"
  s3:
    region: "us-east-1"                              # AWS 区域
    bucket: "mygallery"                              # S3 Bucket 名称
    access_key: "AKIAIOSFODNN7EXAMPLE"               # Access Key ID
    secret_key: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLE"  # Secret Access Key
    endpoint: ""                                     # 留空使用 AWS 默认端点，或填自定义端点
    url_prefix: "https://mygallery.s3.amazonaws.com" # 图片访问 URL 前缀
```

> `endpoint` 留空会使用 AWS 默认端点。如果使用兼容 S3 协议的服务（如 Cloudflare R2），需要填写对应的 endpoint。

#### MinIO（自建 S3 兼容存储）

```yaml
storage:
  type: "minio"
  minio:
    endpoint: "localhost:9000"                      # MinIO 服务地址（不含 http://）
    bucket: "mygallery"                             # Bucket 名称，不存在会自动创建
    access_key: "minioadmin"                        # Access Key
    secret_key: "minioadmin"                        # Secret Key
    use_ssl: false                                  # 是否使用 HTTPS
    url_prefix: "http://localhost:9000/mygallery"   # 图片访问 URL 前缀
```

#### 阿里云 OSS

```yaml
storage:
  type: "aliyun"
  aliyun:
    endpoint: "oss-cn-hangzhou.aliyuncs.com"              # OSS Endpoint（不含 Bucket 名）
    bucket: "mygallery"                                   # Bucket 名称
    access_key: "LTAI5tExample"                           # AccessKey ID
    secret_key: "HpMGhExample"                            # AccessKey Secret
    url_prefix: "https://mygallery.oss-cn-hangzhou.aliyuncs.com"  # 访问 URL 前缀
```

> 阿里云 Endpoint 列表参考：[OSS 各地域 Endpoint](https://help.aliyun.com/document_detail/31837.html)

### JWT 认证配置

```yaml
jwt:
  secret: "your-secret-key-change-this"   # JWT 签名密钥，⚠️ 生产环境务必修改！
  expire_hours: 168                       # Token 有效期（小时），默认 7 天
```

> **安全提示**：`secret` 建议使用随机字符串，至少 32 位。可以用 `openssl rand -hex 32` 生成。

### 管理员配置

```yaml
admin:
  username: "admin"          # 管理员用户名
  password: "admin123"       # 初始密码，首次启动后自动加密存储
  email: "admin@example.com" # 管理员邮箱
```

> 首次启动时会用此配置创建管理员账号。密码会自动 bcrypt 加密，修改配置文件中的密码不会影响已创建的账号。如需重置密码，删除 `data/mygallery.db` 重新初始化。

### 图片处理配置

```yaml
image:
  max_upload_size: 52428800   # 单张图片最大上传大小（字节），默认 50MB
  allowed_types:              # 允许上传的图片 MIME 类型
    - "image/jpeg"
    - "image/png"
    - "image/gif"
    - "image/webp"
  thumbnail:
    width: 400                # 缩略图最大宽度（像素）
    height: 400               # 缩略图最大高度（像素）
    quality: 85               # 缩略图 JPEG 压缩质量（1-100）
```

### CORS 跨域配置

```yaml
cors:
  enabled: true               # 是否启用 CORS
  allow_origins:              # 允许的来源域名
    - "*"                     # * 表示允许所有域名，生产环境建议指定具体域名
  allow_methods:              # 允许的 HTTP 方法
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allow_headers:              # 允许的请求头
    - "Origin"
    - "Content-Type"
    - "Authorization"
```

### 应用配置

```yaml
app:
  name: "MYGallery"                                     # 应用名称
  version: "1.1.4"                                      # 版本号
  site_title: "MYGallery"                               # 网站标题
  site_description: "A simple and beautiful photo gallery system"  # 网站描述
  pagination:
    page_size: 20                                       # 默认每页照片数
```

### 完整配置示例

<details>
<summary>点击查看完整 config.yaml 示例（MySQL + S3）</summary>

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "release"

database:
  type: "mysql"
  mysql:
    host: "db.example.com"
    port: 3306
    username: "mygallery"
    password: "StrongP@ssw0rd"
    database: "mygallery"
    charset: "utf8mb4"

storage:
  type: "s3"
  s3:
    region: "ap-northeast-1"
    bucket: "my-gallery-photos"
    access_key: "AKIAIOSFODNN7EXAMPLE"
    secret_key: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLE"
    endpoint: ""
    url_prefix: "https://my-gallery-photos.s3.ap-northeast-1.amazonaws.com"

jwt:
  secret: "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
  expire_hours: 72

admin:
  username: "admin"
  password: "MySecureP@ss123"
  email: "admin@example.com"

image:
  max_upload_size: 104857600  # 100MB
  allowed_types:
    - "image/jpeg"
    - "image/png"
    - "image/gif"
    - "image/webp"
  thumbnail:
    width: 600
    height: 600
    quality: 80

cors:
  enabled: true
  allow_origins:
    - "https://gallery.example.com"
  allow_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  allow_headers: ["Origin", "Content-Type", "Authorization"]

app:
  name: "My Photo Gallery"
  version: "1.1.4"
  site_title: "My Photo Gallery"
  site_description: "Personal photography showcase"
  pagination:
    page_size: 30
```

</details>

---

## 🌐 API

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/photos` | 照片列表（`page`, `size`, `category`, `search`） |
| GET | `/api/photos/:id` | 照片详情 |
| GET | `/api/categories` | 分类列表 |
| GET | `/api/settings` | 网站设置 |
| GET | `/api/photos/:id/reactions` | 表态统计 |
| POST | `/api/photos/:id/reactions` | 添加/更新表态 |
| DELETE | `/api/photos/:id/reactions` | 删除表态 |
| GET | `/health` | 健康检查 |

### 需认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/auth/login` | 登录 |
| POST | `/api/photos` | 上传照片（支持 `live_photo` 字段） |
| PUT | `/api/photos/:id` | 更新照片 |
| DELETE | `/api/photos/:id` | 删除照片 |
| POST | `/api/categories` | 创建分类 |
| PUT / DELETE | `/api/categories/:id` | 更新 / 删除分类 |

---

## 🐳 Docker 部署

### 快速开始

#### 最简配置（SQLite + 本地存储）

```bash
git clone https://github.com/danzai233/mygallery.git
cd mygallery
cp config.example.yaml config.yaml
docker compose -f docker-compose.simple.yml up -d
```

访问 http://localhost:8080，默认账号：`admin` / `admin123`。

### 多种部署模式

| 配置文件 | 数据库 | 存储 | 适用场景 |
|---------|-------|------|---------|
| `docker-compose.simple.yml` | SQLite | 本地 | 个人部署、开发环境 |
| `docker-compose.full.yml` | PostgreSQL | MinIO | 生产环境、高并发 |

### 环境变量

在 `.env` 文件中配置（可选）：

```bash
# 时区
TZ=Asia/Shanghai

# PostgreSQL 配置（仅完整模式）
POSTGRES_PASSWORD=your_strong_password

# MinIO 配置（仅完整模式）
MINIO_ROOT_USER=your_minio_user
MINIO_ROOT_PASSWORD=your_minio_password
```

### 多平台支持

官方镜像支持以下架构：

- `linux/amd64`
- `linux/arm64` (Apple Silicon / ARM 服务器)
- `linux/arm/v7` (树莓派等 ARM 设备)

无需额外配置，Docker 会自动选择合适架构的镜像。

### 健康检查

容器内置健康检查，可通过以下命令查看：

```bash
docker inspect mygallery --format='{{.State.Health.Status}}'
```

---

## 🚀 部署

<details>
<summary>二进制运行</summary>

```bash
go build -o mygallery main.go
cp config.example.yaml config.yaml
mkdir -p data uploads uploads/thumbnails
./mygallery
```
</details>

<details>
<summary>Nginx 反向代理</summary>

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
    }
}
```
</details>

<details>
<summary>Systemd 服务</summary>

```ini
[Unit]
Description=MYGallery
After=network.target
[Service]
Type=simple
WorkingDirectory=/opt/mygallery
ExecStart=/opt/mygallery/mygallery
Restart=on-failure
[Install]
WantedBy=multi-user.target
```
</details>

---

## 🔒 安全建议

1. 修改默认密码和 JWT 密钥
2. 生产环境启用 HTTPS
3. 配置防火墙
4. 定期备份 `data/` 和 `uploads/`

## 🛠️ 技术栈

| 组件 | 技术 |
|------|------|
| 后端 | Go + Gin + GORM |
| 认证 | JWT |
| EXIF | rwcarlsen/goexif |
| 图片 | disintegration/imaging |
| 数据库 | SQLite / MySQL / PostgreSQL |
| 存储 | Local / S3 / MinIO / Aliyun OSS |
| 前端 | Vanilla JS + Tailwind CSS |
| 地图 | Leaflet + OpenStreetMap |

## 📄 License

MIT

---

<div align="center">
  <p>Made with ❤️ by MYGallery</p>
</div>
