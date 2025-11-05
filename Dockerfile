# 构建阶段
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache git gcc musl-dev

# 复制 go.mod 和 go.sum
COPY go.mod go.sum* ./
RUN go mod download || true

# 复制源代码
COPY . .

# 构建应用
ENV GOTOOLCHAIN=auto
RUN CGO_ENABLED=1 GOOS=linux go build -o mygallery .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

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
  CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# 运行应用
CMD ["./mygallery"]

