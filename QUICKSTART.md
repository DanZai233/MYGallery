# 🚀 MYGallery 快速开始指南

本指南将帮助你在 5 分钟内启动 MYGallery。

## 📋 前置要求

- Docker 和 Docker Compose（推荐）
- 或 Go 1.21+（本地开发）

## 🎯 方式一：Docker 快速部署（推荐）

### 1. 克隆项目

```bash
git clone https://github.com/yourusername/mygallery.git
cd mygallery
```

### 2. 一键启动

```bash
# 使用安装脚本
bash scripts/install.sh

# 或手动启动
cp config.example.yaml config.yaml
docker-compose up -d
```

### 3. 访问应用

- 📷 前台展示：http://localhost:8080
- ⚙️ 后台管理：http://localhost:8080/admin
- 👤 默认账号：`admin` / `admin123`

就这么简单！🎉

## 💻 方式二：本地开发

### 1. 安装依赖

```bash
cd mygallery
go mod download
```

### 2. 创建配置

```bash
cp config.example.yaml config.yaml
```

### 3. 运行应用

```bash
go run main.go
```

### 4. 访问应用

- 前台：http://localhost:8080
- 后台：http://localhost:8080/admin

## 🎨 快速测试

### 上传你的第一张照片

1. 访问 http://localhost:8080/admin
2. 使用 `admin` / `admin123` 登录
3. 拖拽照片到上传区域
4. 等待上传完成
5. 访问 http://localhost:8080 查看效果

### 查看 EXIF 信息

- 在前台页面，鼠标悬停在照片上查看基本信息
- 点击照片打开灯箱，查看完整的 EXIF 元数据

### 编辑照片信息

1. 在后台管理页面，点击照片上的"编辑"按钮
2. 填写标题、描述、标签等信息
3. 保存后在前台查看效果

## ⚙️ 常用命令

### 使用 Makefile

```bash
# 查看所有命令
make help

# 初始化项目
make init

# 编译应用
make build

# 运行应用
make run

# Docker 构建
make docker-build

# Docker 启动
make docker-run

# 查看日志
make docker-logs

# 停止服务
make docker-stop
```

### 使用 Docker Compose

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 重启服务
docker-compose restart

# 停止服务
docker-compose down

# 重新构建
docker-compose build
```

## 🔧 基础配置

### 修改管理员密码

编辑 `config.yaml`：

```yaml
admin:
  username: "admin"
  password: "your-new-password"
```

**重要**：首次启动后，密码会被加密存储，之后修改需要在后台管理页面操作。

### 修改端口

编辑 `config.yaml`：

```yaml
server:
  port: 8080  # 改为你想要的端口
```

或修改 `docker-compose.yml`：

```yaml
ports:
  - "3000:8080"  # 外部端口:内部端口
```

### 更换数据库

编辑 `config.yaml`：

```yaml
database:
  type: "mysql"  # sqlite, mysql, postgres
  mysql:
    host: "localhost"
    port: 3306
    username: "root"
    password: "password"
    database: "mygallery"
```

### 使用对象存储

编辑 `config.yaml`：

```yaml
storage:
  type: "s3"  # local, s3, minio, aliyun
  s3:
    region: "us-east-1"
    bucket: "mygallery"
    access_key: "your-key"
    secret_key: "your-secret"
```

## 📝 下一步

- 📖 阅读 [完整文档](README.md)
- 🚀 查看 [部署指南](DEPLOYMENT.md)
- 🏗️ 了解 [项目结构](STRUCTURE.md)
- 📋 查看 [更新日志](CHANGELOG.md)

## 💡 小提示

### 提高上传速度
- 使用有线网络
- 选择离你近的对象存储区域
- 批量上传时分批处理

### 优化存储空间
- 上传前压缩照片
- 定期清理不需要的照片
- 使用对象存储的生命周期策略

### 备份数据
```bash
# 备份 SQLite 数据库
cp data/mygallery.db data/backup.db

# 备份上传文件
tar -czf uploads-backup.tar.gz uploads/
```

## ❓ 遇到问题？

### 端口被占用
```bash
# 修改 docker-compose.yml 中的端口
ports:
  - "8888:8080"  # 使用其他端口
```

### 权限问题
```bash
# 给予目录写入权限
sudo chown -R $USER:$USER data/ uploads/
```

### 看不到照片
```bash
# 检查日志
docker-compose logs mygallery

# 检查配置
cat config.yaml
```

## 📞 获取帮助

- 🐛 [提交 Issue](https://github.com/yourusername/mygallery/issues)
- 💬 [参与讨论](https://github.com/yourusername/mygallery/discussions)
- 📧 发邮件：your-email@example.com

---

开始享受你的照片墙之旅吧！📷✨

