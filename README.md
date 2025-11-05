# MYGallery - 个人照片墙系统

<div align="center">
  <h1>📷 MYGallery</h1>
  <p>一个简约、美观、功能完整的个人照片墙系统</p>
  <p>人人都可以自部署自己的照片展示空间</p>
</div>

## ✨ 特性

### 🎨 精美的前端展示
- **瀑布流布局**：响应式设计，自适应各种屏幕尺寸
- **灯箱效果**：点击图片查看大图，支持键盘导航和手势操作
- **磨砂玻璃质感**：现代化的UI设计，圆润的界面语言
- **EXIF 元数据展示**：自动显示相机型号、拍摄参数、位置等信息

### ⚙️ 强大的后台管理
- **用户登录验证**：JWT token 认证，安全可靠
- **图片上传**：支持拖拽上传，批量上传，实时进度显示
- **元数据编辑**：编辑照片标题、描述、标签、位置、版权等信息
- **EXIF 自动提取**：自动读取照片的相机参数、GPS 位置等元数据

### 🔧 灵活的配置系统
- **多数据库支持**：SQLite、MySQL、PostgreSQL
- **多存储支持**：本地存储、AWS S3、MinIO、阿里云 OSS
- **YAML 配置**：简单直观的配置文件
- **Docker 部署**：一键部署，开箱即用

### 🚀 性能优化
- **缩略图生成**：自动生成缩略图，加快加载速度
- **懒加载**：图片按需加载，节省带宽
- **CDN 支持**：支持对象存储 CDN 加速

## 📦 快速开始

### Docker 部署（推荐）

1. **克隆仓库**
```bash
git clone https://github.com/yourusername/mygallery.git
cd mygallery
```

2. **创建配置文件**
```bash
cp config.example.yaml config.yaml
# 根据需要修改配置文件
```

3. **启动服务**
```bash
docker-compose up -d
```

4. **访问应用**
- 前台展示：http://localhost:8080
- 后台管理：http://localhost:8080/admin
- 默认账号：admin / admin123

### 本地部署

#### 环境要求
- Go 1.21+
- Git

#### 安装步骤

1. **克隆仓库**
```bash
git clone https://github.com/yourusername/mygallery.git
cd mygallery
```

2. **安装依赖**
```bash
go mod download
```

3. **创建配置文件**
```bash
cp config.example.yaml config.yaml
# 根据需要修改配置文件
```

4. **运行应用**
```bash
go run main.go
```

5. **访问应用**
- 前台展示：http://localhost:8080
- 后台管理：http://localhost:8080/admin
- 默认账号：admin / admin123

## 📝 配置说明

### 数据库配置

#### SQLite（默认）
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
    password: "password"
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
    password: "password"
    database: "mygallery"
    sslmode: "disable"
```

### 存储配置

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
    access_key: "your-access-key"
    secret_key: "your-secret-key"
    url_prefix: "https://your-bucket.s3.amazonaws.com"
```

#### MinIO
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
    access_key: "your-access-key"
    secret_key: "your-secret-key"
    url_prefix: "https://mygallery.oss-cn-hangzhou.aliyuncs.com"
```

## 🎯 功能详解

### EXIF 元数据自动提取
系统会自动从上传的照片中提取以下信息：
- 📷 相机品牌和型号
- 🔭 镜头型号
- ⚙️ 拍摄参数（光圈、快门、ISO、焦距）
- 🕐 拍摄时间
- 📍 GPS 位置信息（如果有）
- 📐 图片尺寸

### 照片管理功能
- ✏️ 编辑照片标题和描述
- 🏷️ 添加标签分类
- 📍 标注拍摄位置
- ©️ 设置版权信息
- 🗑️ 删除照片

## 🔐 安全性

- JWT Token 认证
- 密码 bcrypt 加密
- CORS 跨域配置
- 文件类型验证
- 文件大小限制

## 📊 API 接口

### 公开接口
- `GET /api/photos` - 获取照片列表
- `GET /api/photos/:id` - 获取单张照片详情

### 认证接口
- `POST /api/auth/login` - 用户登录

### 需要认证的接口
- `POST /api/photos` - 上传照片
- `PUT /api/photos/:id` - 更新照片信息
- `DELETE /api/photos/:id` - 删除照片
- `POST /api/auth/change-password` - 修改密码

## 🛠️ 技术栈

### 后端
- **框架**：Gin (Go Web Framework)
- **数据库**：GORM (支持 SQLite/MySQL/PostgreSQL)
- **认证**：JWT
- **图片处理**：imaging (缩略图生成)
- **EXIF 解析**：goexif

### 前端
- **框架**：原生 JavaScript
- **样式**：Tailwind CSS
- **灯箱**：lightGallery.js
- **布局**：CSS 瀑布流

### 部署
- **容器化**：Docker & Docker Compose
- **反向代理**：支持 Nginx

## 📸 截图

*待添加*

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 🙏 致谢

感谢所有开源项目的贡献者！

---

<div align="center">
  <p>Made with ❤️ by MYGallery</p>
  <p>如果这个项目对你有帮助，请给个 ⭐️ Star 吧！</p>
</div>

