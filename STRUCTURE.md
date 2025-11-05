# MYGallery 项目结构

```
MYGallery/
├── main.go                 # 应用入口
├── go.mod                  # Go 模块定义
├── config.yaml             # 配置文件（从 config.example.yaml 复制）
├── config.example.yaml     # 配置文件示例
├── Dockerfile              # Docker 镜像构建文件
├── docker-compose.yml      # Docker Compose 配置
├── Makefile                # Make 构建命令
├── .dockerignore           # Docker 忽略文件
├── .gitignore              # Git 忽略文件
│
├── README.md               # 项目说明文档
├── DEPLOYMENT.md           # 部署指南
├── CHANGELOG.md            # 更新日志
├── STRUCTURE.md            # 项目结构说明（本文件）
├── LICENSE                 # 开源许可证
│
├── internal/               # 内部包（不对外暴露）
│   ├── config/            # 配置管理
│   │   └── config.go      # 配置加载和解析
│   │
│   ├── models/            # 数据模型
│   │   └── models.go      # User、Photo 等模型定义
│   │
│   ├── database/          # 数据库层
│   │   └── database.go    # 数据库初始化和操作
│   │
│   ├── storage/           # 存储层
│   │   ├── storage.go     # 存储接口和本地存储
│   │   ├── s3.go          # AWS S3 存储实现
│   │   ├── minio.go       # MinIO 存储实现
│   │   └── aliyun.go      # 阿里云 OSS 存储实现
│   │
│   ├── middleware/        # 中间件
│   │   └── auth.go        # JWT 认证中间件、CORS 中间件
│   │
│   ├── handlers/          # HTTP 处理器
│   │   ├── auth.go        # 认证相关处理（登录、修改密码）
│   │   └── photo.go       # 照片相关处理（上传、查询、更新、删除）
│   │
│   ├── router/            # 路由配置
│   │   └── router.go      # 路由设置和中间件注册
│   │
│   └── utils/             # 工具函数
│       ├── image.go       # 图片处理（EXIF 提取、缩略图生成）
│       └── file.go        # 文件处理（文件名生成、类型检查）
│
├── public/                # 前端静态文件
│   ├── index.html         # 前台照片展示页面
│   ├── admin.html         # 后台管理页面
│   └── assets/            # 静态资源
│       ├── css/           # 样式文件
│       └── js/            # JavaScript 文件
│
├── data/                  # 数据目录（SQLite 数据库）
│   ├── .gitkeep
│   └── mygallery.db       # SQLite 数据库文件（运行时生成）
│
├── uploads/               # 上传文件目录
│   ├── .gitkeep
│   ├── thumbnails/        # 缩略图目录
│   └── *.jpg/png/...      # 上传的照片文件
│
└── scripts/               # 脚本文件
    └── install.sh         # 一键安装脚本
```

## 核心组件说明

### 1. 配置系统 (`internal/config/`)
- 支持 YAML 格式配置文件
- 支持多种数据库配置
- 支持多种存储配置
- 提供默认配置

### 2. 数据层 (`internal/database/`, `internal/models/`)
- 使用 GORM ORM
- 支持 SQLite、MySQL、PostgreSQL
- 自动迁移表结构
- 提供数据库操作封装

### 3. 存储层 (`internal/storage/`)
- 统一的存储接口
- 支持本地存储
- 支持 AWS S3
- 支持 MinIO
- 支持阿里云 OSS

### 4. 认证系统 (`internal/middleware/auth.go`)
- JWT Token 认证
- 密码 bcrypt 加密
- CORS 跨域支持

### 5. 业务逻辑 (`internal/handlers/`)
- **认证处理器**：用户登录、密码修改
- **照片处理器**：上传、查询、更新、删除

### 6. 图片处理 (`internal/utils/`)
- EXIF 元数据提取
- 自动缩略图生成
- 图片尺寸获取

### 7. 前端页面 (`public/`)
- **index.html**：瀑布流照片展示，灯箱预览
- **admin.html**：后台管理，上传、编辑、删除

## 数据库表结构

### users 表
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email TEXT,
    role TEXT DEFAULT 'admin'
);
```

### photos 表
```sql
CREATE TABLE photos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    
    -- 基本信息
    filename TEXT NOT NULL,
    original_name TEXT NOT NULL,
    title TEXT,
    description TEXT,
    tags TEXT,
    location TEXT,
    
    -- EXIF 元数据
    camera_make TEXT,
    camera_model TEXT,
    lens_model TEXT,
    focal_length TEXT,
    aperture TEXT,
    shutter_speed TEXT,
    iso TEXT,
    date_taken DATETIME,
    
    -- GPS 信息
    gps_latitude REAL,
    gps_longitude REAL,
    
    -- 文件信息
    width INTEGER,
    height INTEGER,
    file_size INTEGER,
    mime_type TEXT,
    
    -- 存储信息
    storage_type TEXT,
    storage_path TEXT,
    thumbnail_path TEXT,
    
    -- 其他
    copyright TEXT,
    user_id INTEGER,
    views INTEGER DEFAULT 0,
    
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

## API 接口

### 公开接口
- `GET /api/photos` - 获取照片列表
- `GET /api/photos/:id` - 获取单张照片详情
- `GET /health` - 健康检查

### 认证接口
- `POST /api/auth/login` - 用户登录

### 需要认证的接口（需要 Bearer Token）
- `POST /api/photos` - 上传照片
- `PUT /api/photos/:id` - 更新照片信息
- `DELETE /api/photos/:id` - 删除照片
- `POST /api/auth/change-password` - 修改密码

## 文件命名规范

### Go 文件
- 小写字母，下划线分隔：`config.go`, `auth.go`
- 测试文件：`*_test.go`

### 目录命名
- 小写字母，单数形式：`config/`, `model/`, `handler/`

### 变量命名
- 驼峰命名：`userID`, `photoList`
- 常量大写：`API_BASE`, `MAX_SIZE`

## 开发流程

### 添加新功能
1. 在 `internal/models/` 定义模型
2. 在 `internal/handlers/` 实现业务逻辑
3. 在 `internal/router/` 注册路由
4. 更新前端页面（如需要）

### 添加新的存储类型
1. 在 `internal/storage/` 创建新文件
2. 实现 `Storage` 接口
3. 在 `storage.go` 中注册新类型
4. 更新 `config.example.yaml`

### 部署新版本
1. 更新代码
2. 构建镜像：`docker-compose build`
3. 重启服务：`docker-compose down && docker-compose up -d`
4. 检查日志：`docker-compose logs -f`

## 性能考虑

### 数据库
- 照片表按创建时间索引
- 用户表用户名唯一索引

### 存储
- 缩略图自动生成
- 支持 CDN 加速
- 对象存储支持

### 前端
- 图片懒加载
- 瀑布流布局优化
- 静态资源缓存

## 安全性

### 认证
- JWT Token 有效期 7 天
- 密码 bcrypt 加密（cost=10）

### 文件上传
- 文件类型白名单
- 文件大小限制（默认 50MB）
- MIME 类型检查

### API
- CORS 配置
- 认证中间件保护

## 扩展性

### 数据库扩展
- 使用 GORM，容易切换数据库
- 表结构自动迁移

### 存储扩展
- 统一的存储接口
- 易于添加新的存储类型

### 功能扩展
- 模块化设计
- 清晰的代码结构

---

如有疑问，请参考其他文档或提交 Issue。

