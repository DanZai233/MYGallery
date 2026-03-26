# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache git gcc musl-dev make

# 复制依赖文件
COPY go.mod go.sum* ./

# 下载依赖（利用 Docker 缓存）
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
# 修复 go-sqlite3 在 Alpine/musl 上的编译问题
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
ENV GOOS=linux
RUN go build -ldflags="-s -w" -o mygallery .

# 运行阶段
FROM alpine:latest

LABEL maintainer="danzai233 <danzai233@gmail.com>"
LABEL description="MYGallery - 个人照片墙系统"
LABEL version="1.1.4"

WORKDIR /app

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata curl

# 从构建阶段复制二进制文件
COPY --from=builder /app/mygallery .
COPY --from=builder /app/config.example.yaml ./config.example.yaml
COPY --from=builder /app/public ./public

# 创建必要的目录
RUN mkdir -p /app/data /app/uploads /app/uploads/thumbnails

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# 运行应用
CMD ["./mygallery"]

