.PHONY: help build run dev docker docker-build docker-run clean

help: ## 显示帮助信息
	@echo "MYGallery - 个人照片墙系统"
	@echo ""
	@echo "可用命令:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

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

docker-build: ## 构建 Docker 镜像
	@echo "🐳 构建 Docker 镜像..."
	docker build -t mygallery:latest .
	@echo "✅ 镜像构建完成"

docker-run: ## 运行 Docker 容器
	@echo "🚀 启动 Docker 容器..."
	docker-compose up -d
	@echo "✅ 容器已启动"
	@echo "📷 前台: http://localhost:8080"
	@echo "⚙️  后台: http://localhost:8080/admin"

docker-stop: ## 停止 Docker 容器
	@echo "🛑 停止 Docker 容器..."
	docker-compose down

docker-logs: ## 查看 Docker 日志
	docker-compose logs -f

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
	@echo "2. 运行 make run 或 make docker-run"

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

