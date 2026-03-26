.PHONY: help build run dev docker docker-build docker-run clean

help: ## 显示帮助信息
	@echo "MYGallery - 个人照片墙系统"
	@echo ""
	@echo "可用命令:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## 编译应用
	@echo "🔨 编译应用..."
	go build -o bin/mygallery main.go
	@echo "✅ 编译完成: bin/mygallery"

run: ## 运行应用
	@echo "🚀 启动应用..."
	go run main.go

dev: ## 开发模式（自动重载）
	@echo "🔧 开发模式..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "请先安装 air: go install github.com/cosmtrek/air@latest"; \
		go run main.go; \
	fi

test: ## 运行测试
	@echo "🧪 运行测试..."
	go test -v ./...

# Docker 相关命令
docker-build: ## 构建 Docker 镜像
	@echo "🐳 构建 Docker 镜像..."
	docker build -t mygallery:latest .
	@echo "✅ 镜像构建完成"

docker-build-multi: ## 构建多平台 Docker 镜像 (amd64/arm64)
	@echo "🐳 构建多平台 Docker 镜像..."
	docker buildx build --platform linux/amd64,linux/arm64 -t mygallery:latest .

docker-build-multi-push: ## 构建并推送多平台镜像到 Docker Hub
	@echo "🐳 构建并推送多平台镜像..."
	docker buildx build --platform linux/amd64,linux/arm64 -t danzai233/mygallery:$(VERSION) -t danzai233/mygallery:latest --push .

docker-tag: ## 为镜像打标签
	@echo "🏷️  为镜像打标签..."
	docker tag mygallery:latest danzai233/mygallery:$(VERSION)
	docker tag mygallery:latest danzai233/mygallery:latest
	@echo "✅ 标签完成"

docker-push: ## 推送镜像到 Docker Hub
	@echo "📤 推送镜像到 Docker Hub..."
	docker push danzai233/mygallery:$(VERSION)
	docker push danzai233/mygallery:latest
	@echo "✅ 推送完成"

docker-login: ## 登录 Docker Hub
	@echo "🔐 登录 Docker Hub..."
	docker login

docker-up-simple: ## 启动简化版 Docker Compose (SQLite)
	@echo "🚀 启动简化版 Docker Compose..."
	docker compose -f docker-compose.simple.yml up -d
	@echo "✅ 容器已启动"
	@echo "📷 前台: http://localhost:8080"
	@echo "⚙️  后台: http://localhost:8080/admin"

docker-up-full: ## 启动完整版 Docker Compose (PostgreSQL + MinIO)
	@echo "🚀 启动完整版 Docker Compose..."
	docker compose -f docker-compose.full.yml up -d
	@echo "✅ 容器已启动"
	@echo "📷 前台: http://localhost:8080"
	@echo "⚙️  后台: http://localhost:8080/admin"

docker-down-simple: ## 停止简化版 Docker Compose
	@echo "🛑 停止简化版 Docker Compose..."
	docker compose -f docker-compose.simple.yml down

docker-down-full: ## 停止完整版 Docker Compose
	@echo "🛑 停止完整版 Docker Compose..."
	docker compose -f docker-compose.full.yml down

docker-logs-simple: ## 查看简化版 Docker 日志
	docker compose -f docker-compose.simple.yml logs -f

docker-logs-full: ## 查看完整版 Docker 日志
	docker compose -f docker-compose.full.yml logs -f

docker-ps: ## 查看 Docker 容器状态
	docker compose -f docker-compose.simple.yml ps

docker-restart-simple: ## 重启简化版 Docker Compose
	docker compose -f docker-compose.simple.yml restart

docker-restart-full: ## 重启完整版 Docker Compose
	docker compose -f docker-compose.full.yml restart

docker-exec: ## 进入容器调试
	docker exec -it mygallery sh

init: ## 初始化项目
	@echo "🎉 初始化项目..."
	@if [ ! -f config.yaml ]; then \
		cp config.example.yaml config.yaml; \
		echo "✅ 配置文件已创建: config.yaml"; \
	else \
		echo "⚠️  配置文件已存在"; \
	fi
	@mkdir -p data uploads uploads/thumbnails public/assets
	@echo "✅ 目录结构已创建"
	@echo ""
	@echo "下一步："
	@echo "1. 编辑 config.yaml 配置文件"
	@echo "2. 运行 make run 或 make docker-up-simple"

clean: ## 清理构建文件
	@echo "🧹 清理构建文件..."
	rm -rf bin/
	rm -rf uploads/*
	@echo "✅ 清理完成"

deps: ## 安装依赖
	@echo "📦 安装依赖..."
	go mod download
	go mod tidy
	@echo "✅ 依赖安装完成"

