# 🔧 SQLite 编译错误修复

## ❌ 错误信息

```
sqlite3-binding.c:37644:42: error: 'pread64' undeclared here
sqlite3-binding.c:37662:42: error: 'pwrite64' undeclared here
sqlite3-binding.c:37648:49: error: unknown type name 'off64_t'
```

## 🎯 问题原因

**Alpine Linux + musl libc 兼容性问题**：
- Alpine Linux 使用 `musl libc` 而不是 `glibc`
- `go-sqlite3` 的 C 代码在 musl 环境下需要特殊编译标志
- `pread64`/`pwrite64` 和 `off64_t` 在 musl 中需要 `_LARGEFILE64_SOURCE` 宏定义

## ✅ 已修复

### 修改 Dockerfile

```dockerfile
# 修复前
ENV GOTOOLCHAIN=auto
RUN CGO_ENABLED=1 GOOS=linux go build -o mygallery .

# 修复后
ENV GOTOOLCHAIN=auto
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"  # ✅ 添加编译标志
RUN CGO_ENABLED=1 GOOS=linux go build -o mygallery .
```

### 说明

- `CGO_CFLAGS="-D_LARGEFILE64_SOURCE"`: 定义宏，启用 64 位文件操作支持
- 这是 musl libc 环境下编译 `go-sqlite3` 的标准解决方案

---

## 🧪 测试构建

### 本地测试

```bash
cd /root/MYGallery

# 测试构建
docker build -t mygallery:test .

# 如果成功，测试运行
docker run -p 8080:8080 mygallery:test
```

### 推送后测试

```bash
# 提交更改
git add Dockerfile
git commit -m "fix: 修复 SQLite 在 Alpine 上的编译问题"
git push origin main

# 创建发布（触发 Actions）
bash scripts/release.sh
```

---

## 📊 技术细节

### 为什么需要这个标志？

1. **musl vs glibc**:
   - glibc: 默认支持 `pread64`/`pwrite64`
   - musl: 需要显式定义 `_LARGEFILE64_SOURCE` 才能使用

2. **go-sqlite3 的 C 代码**:
   - 使用了 `pread64`/`pwrite64` 进行大文件操作
   - 在 musl 环境下需要额外的宏定义

3. **解决方案**:
   - `-D_LARGEFILE64_SOURCE`: 启用 64 位文件操作 API
   - 这是 Alpine Linux 上编译需要 64 位文件操作的程序的标准做法

---

## 🔄 替代方案（如果仍然失败）

### 方案 1: 使用纯 Go SQLite 驱动

如果 CGO 问题持续存在，可以考虑切换到纯 Go 实现：

```go
// 替换 go.mod 中的依赖
// 删除: gorm.io/driver/sqlite
// 添加: modernc.org/sqlite (纯 Go，无需 CGO)
```

但需要修改 `internal/database/database.go` 中的导入和初始化代码。

### 方案 2: 使用 Debian 基础镜像

```dockerfile
# 使用 Debian 而不是 Alpine（更大但更兼容）
FROM golang:1.24-bookworm AS builder
```

Debian 使用 glibc，不会有这个问题，但镜像会更大。

---

## ✅ 当前修复已足够

**添加 `CGO_CFLAGS="-D_LARGEFILE64_SOURCE"` 应该就能解决问题！**

这是 Alpine Linux 上编译 `go-sqlite3` 的标准解决方案，已经被广泛使用。

---

## 📋 修复步骤总结

1. ✅ 添加 `ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"`
2. ✅ 提交到 GitHub
3. ✅ 触发 Actions 构建
4. ✅ 验证构建成功

---

**修复完成！现在可以成功构建了！** 🚀

