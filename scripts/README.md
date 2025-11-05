# 脚本说明

本目录包含 MYGallery 的辅助脚本。

## install.sh

一键安装脚本，自动完成以下操作：

1. 检查 Docker 环境
2. 创建配置文件
3. 创建必要目录
4. 构建 Docker 镜像
5. 启动服务

### 使用方法

```bash
bash scripts/install.sh
```

### 要求

- Docker 20.10+
- Docker Compose 1.29+

## 其他脚本（待添加）

- `backup.sh` - 数据备份脚本
- `restore.sh` - 数据恢复脚本
- `update.sh` - 应用更新脚本
- `migrate.sh` - 数据迁移脚本
