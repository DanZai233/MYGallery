# MYGallery - 个人照片墙系统

<div align="center">
  <h1>📷 MYGallery</h1>
  <p>一个简约、美观、功能完整的个人照片墙系统</p>
  <p>人人都可以自部署自己的照片展示空间</p>
  
  <p>
    <img src="https://img.shields.io/badge/Version-1.1.2-blue" alt="Version">
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go" alt="Go Version">
    <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
  </p>
</div>

---

## ✨ 特性

### 🎨 精美的前端展示
- **瀑布流布局**：响应式设计，自适应各种屏幕尺寸
- **灯箱效果**：点击查看大图，支持键盘导航和手势操作
- **磨砂玻璃质感**：现代化的UI设计，圆润的界面语言
- **EXIF 元数据展示**：自动显示相机型号、拍摄参数、位置等信息
- **分类筛选**：按分类快速筛选照片，毛玻璃按钮效果
- **搜索功能**：支持搜索标题、描述、标签、位置等
- **黑夜模式**：一键切换白天/黑夜主题

### ⚙️ 强大的后台管理
- **用户登录验证**：JWT token 认证，安全可靠
- **图片上传**：支持拖拽上传，批量上传，实时进度显示
- **元数据编辑**：编辑照片标题、描述、标签、位置、版权等信息
- **EXIF 自动提取**：自动读取照片的相机参数、GPS 位置等元数据
- **EXIF 手动编辑**：所有相机参数都可以手动修改
- **分类管理**：创建、编辑、删除照片分类
- **网站设置**：自定义网站信息、备案信息、Header/Footer等

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

### 本地运行（最简单）

```bash
# 1. 克隆仓库
git clone https://github.com/yourusername/mygallery.git
cd mygallery

# 2. 安装依赖
go mod download

# 3. 创建配置
cp config.example.yaml config.yaml

# 4. 运行应用
go run main.go
```

### Docker 部署

```bash
# 1. 克隆仓库
git clone https://github.com/yourusername/mygallery.git
cd mygallery

# 2. 创建配置
cp config.example.yaml config.yaml

# 3. 启动服务（需要更新 Docker 和 Buildx）
docker compose build
docker compose up -d

# 或者直接本地运行
go run main.go
```

### 访问应用

- 📷 前台展示：http://localhost:8080
- ⚙️ 后台管理：http://localhost:8080/admin
- 👤 默认账号：admin / admin123

⚠️ **首次登录后请立即修改默认密码！**

## 📖 文档

- 📚 [快速开始指南](QUICKSTART.md) - 5 分钟快速入门

## 🎯 核心功能

### 照片展示
- ✅ 瀑布流布局（响应式4/3/2/1列）
- ✅ 灯箱大图预览（键盘导航）
- ✅ 元数据完整展示（相机/参数/位置）
- ✅ Hover 简介效果
- ✅ 黑夜模式切换

### 照片管理
- ✅ 拖拽上传
- ✅ EXIF 自动提取
- ✅ 元数据编辑（所有字段）
- ✅ 分类设置
- ✅ 批量操作

### 分类系统
- ✅ 分类管理（创建/编辑/删除）
- ✅ 首页分类筛选
- ✅ 自动统计照片数量

### 搜索功能
- ✅ 多字段搜索（标题/描述/标签/位置）
- ✅ 实时搜索结果
- ✅ 毛玻璃UI

### 网站设置
- ✅ 基本信息（标题/描述/Logo）
- ✅ 备案信息（ICP/公安备案）
- ✅ 联系方式（邮箱/电话/微信）
- ✅ 社交媒体（GitHub/Twitter/微博）
- ✅ 自定义代码（HTML/CSS/JS/统计）

## 🛠️ 技术栈

**后端**
- Gin Web Framework
- GORM ORM
- JWT 认证
- goexif EXIF 解析
- imaging 图片处理

**前端**
- 原生 JavaScript
- Tailwind CSS
- lightGallery.js

**部署**
- Docker & Docker Compose
- Nginx 反向代理支持

## 📝 配置

编辑 `config.yaml` 文件：

```yaml
# 服务器配置
server:
  port: 8080

# 数据库配置（默认 SQLite）
database:
  type: "sqlite"
  sqlite:
    path: "./data/mygallery.db"

# 存储配置（默认本地）
storage:
  type: "local"
  local:
    upload_dir: "./uploads"
```

更多配置选项请参考 `config.example.yaml`

## 🔒 安全建议

1. **修改默认密码**
2. **修改 JWT 密钥**
3. **启用 HTTPS**（生产环境）
4. **配置防火墙**
5. **定期备份数据**

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

---

<div align="center">
  <p>Made with ❤️ by MYGallery</p>
  <p>如果这个项目对你有帮助，请给个 ⭐️ Star 吧！</p>
</div>
