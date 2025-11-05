# 更新日志

本文档记录 MYGallery 的所有重要更改。

## [2.0.0] - 2025-11-05

### 🎉 重大更新

#### 后端架构
- ✨ 使用 Golang 重写后端，提升性能和稳定性
- 🗄️ 支持多数据库：SQLite、MySQL、PostgreSQL
- ☁️ 支持多对象存储：本地存储、AWS S3、MinIO、阿里云 OSS
- 🔐 JWT 认证系统
- 📸 自动 EXIF 元数据提取
- 🖼️ 自动缩略图生成

#### 前端功能
- 🎨 全新的瀑布流布局设计
- 💫 磨砂玻璃质感 UI
- 🔍 lightGallery 灯箱效果
- 📱 完全响应式设计
- ⚡ 懒加载优化

#### 管理后台
- 👤 用户登录认证
- 📤 拖拽上传支持
- 📝 元数据编辑
- 🗑️ 照片管理（增删改查）
- 📊 统计面板

#### 部署支持
- 🐳 Docker 容器化
- 📦 Docker Compose 一键部署
- 🔧 Makefile 便捷命令
- 📜 一键安装脚本
- 📖 完整部署文档

### 🛠️ 技术栈

**后端**
- Gin Web Framework
- GORM ORM
- JWT 认证
- goexif EXIF 解析
- imaging 图片处理
- AWS SDK / MinIO SDK / 阿里云 SDK

**前端**
- 原生 JavaScript
- Tailwind CSS
- lightGallery.js

**部署**
- Docker & Docker Compose
- Nginx 反向代理支持

### 📝 配置系统

- YAML 配置文件
- 环境变量支持
- 灵活的配置选项
- 安全的默认配置

### 🔒 安全特性

- JWT Token 认证
- bcrypt 密码加密
- CORS 跨域配置
- 文件类型验证
- 文件大小限制

### 📚 文档

- README.md - 项目介绍和快速开始
- DEPLOYMENT.md - 详细部署指南
- config.example.yaml - 配置文件示例
- API 文档（待完善）

### 🎯 核心功能

- ✅ 照片上传
- ✅ EXIF 元数据自动提取
- ✅ 缩略图自动生成
- ✅ 照片列表展示
- ✅ 照片详情查看
- ✅ 照片信息编辑
- ✅ 照片删除
- ✅ 用户登录
- ✅ 密码修改
- ✅ 瀑布流布局
- ✅ 灯箱预览
- ✅ 响应式设计

### 🚀 性能优化

- 缩略图生成
- 图片懒加载
- 数据库索引优化
- 静态资源缓存
- CDN 支持

### 🐛 已知问题

- 暂无

### 📅 计划功能

- [ ] 图片标签系统
- [ ] 照片分类功能
- [ ] 照片搜索功能
- [ ] 照片排序选项
- [ ] 批量操作
- [ ] 照片下载
- [ ] 评论功能
- [ ] 点赞功能
- [ ] RSS 订阅
- [ ] API 文档
- [ ] 多用户支持
- [ ] 主题切换
- [ ] 国际化支持

---

## 版本说明

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)

版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)

