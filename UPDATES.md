# 🎉 v1.0.0 更新说明

本次更新带来了大量新功能和界面优化！

## ✨ 主要更新

### 1. 界面优化

#### 🌙 黑夜模式
- 新增黑夜模式切换按钮（右上角月亮/太阳图标）
- 自动保存用户偏好到 localStorage
- 平滑的主题切换动画
- 暗色模式优化的配色方案

#### 🎨 照片展示优化
- **Hover 效果升级**：鼠标悬停时，照片描述从底部平滑滑入
- **磨砂玻璃效果**：描述区域使用毛玻璃背景
- **灯箱优化**：lightGallery 添加毛玻璃背景效果，更加高级
- **图片加载优化**：修复了缩略图无法显示的问题

#### 🔒 安全性提升
- 移除首页的后台管理按钮
- 后台只能通过 `/admin` URL 直接访问
- 防止未授权用户误入后台

### 2. 新增功能

#### 📁 照片分类系统
- 支持为照片设置分类
- 分类管理 API（创建/更新/删除）
- 首页支持按分类筛选照片
- 分类统计功能

**数据库模型**：
```go
type Category struct {
    ID          uint
    Name        string  // 分类名称
    Slug        string  // URL友好的标识符
    Description string  // 分类描述
    Cover       string  // 分类封面
    SortOrder   int     // 排序顺序
    PhotoCount  int     // 照片数量
}
```

**API 接口**：
- `GET /api/categories` - 获取分类列表（公开）
- `POST /api/categories` - 创建分类（需认证）
- `PUT /api/categories/:id` - 更新分类（需认证）
- `DELETE /api/categories/:id` - 删除分类（需认证）

#### ⚙️ 网站设置系统
支持自定义网站的各种信息：

**基本信息**：
- 网站标题
- 网站描述
- 网站关键词
- 网站 Logo

**页眉页脚**：
- 自定义 Header HTML
- 自定义 Footer HTML

**备案信息**：
- ICP 备案号
- 公安备案号
- 备案链接

**联系方式**：
- 邮箱
- 电话
- 微信

**社交媒体**：
- GitHub URL
- Twitter URL
- 微博 URL

**高级设置**：
- 统计代码（Google Analytics 等）
- 自定义 CSS
- 自定义 JavaScript

**数据库模型**：
```go
type Settings struct {
    ID              uint
    SiteTitle       string
    SiteDescription string
    HeaderHTML      string
    FooterHTML      string
    ICPRecord       string
    BeianRecord     string
    ContactEmail    string
    AnalyticsCode   string
    CustomCSS       string
    // ... 更多字段
}
```

**API 接口**：
- `GET /api/settings` - 获取网站设置（公开）
- `PUT /api/settings` - 更新网站设置（需认证）

### 3. 修复的问题

#### 🖼️ 图片显示问题
**问题**：缩略图路径错误导致无法显示

**原因**：
- 缩略图存储在 `uploads/thumbnails/` 目录
- 但返回的 URL 缺少 `thumbnails/` 前缀

**解决方案**：
```go
// 修复前
photo.ThumbnailURL = stor.GetURL(photo.ThumbnailPath)

// 修复后
photo.ThumbnailURL = stor.GetURL("thumbnails/" + photo.ThumbnailPath)
```

#### 📱 响应式优化
- 修复了移动端照片悬停效果
- 优化了暗色模式的文字对比度
- 改进了灯箱的触摸操作

## 🎨 界面效果

### 白天模式
- 清新的渐变背景
- 磨砂玻璃卡片效果
- 柔和的阴影

### 黑夜模式
- 深色背景保护眼睛
- 高对比度文字
- 暗色系磨砂玻璃

### 照片 Hover 效果
```
鼠标悬停前：
┌─────────────┐
│             │
│   照片区域   │
│             │
└─────────────┘

鼠标悬停后：
┌─────────────┐
│             │
│   照片区域   │
│┄┄┄┄┄┄┄┄┄┄┄┄┄│
│ 📝 标题      │
│ 📄 描述      │
│ 📷 EXIF     │
│ 📍 位置      │
└─────────────┘
（从底部滑入，毛玻璃效果）
```

### 灯箱效果
- 背景毛玻璃模糊
- EXIF 信息面板毛玻璃效果
- 键盘导航（←→ 切换，ESC 关闭）
- 触摸手势支持

## 📊 数据库变更

### 新增表
1. **categories** - 分类表
2. **settings** - 设置表

### 修改表
**photos** 表新增字段：
- `category` VARCHAR - 照片分类

## 🔧 API 变更

### 新增 API

#### 分类管理
```
GET    /api/categories           获取分类列表
POST   /api/categories           创建分类
PUT    /api/categories/:id       更新分类
DELETE /api/categories/:id       删除分类
```

#### 设置管理
```
GET    /api/settings             获取网站设置
PUT    /api/settings             更新网站设置
```

### 修改 API

#### 照片列表
```
GET /api/photos?category=风景    按分类筛选
```

## 🚀 使用方法

### 运行应用
```bash
cd /root/MYGallery
go run main.go
```

### 访问地址
- 📷 前台：http://localhost:8080
- ⚙️ 后台：http://localhost:8080/admin
- 👤 账号：admin / admin123

### 切换主题
点击右上角的 🌙/☀️ 图标

### 管理分类
1. 登录后台
2. 进入分类管理
3. 创建、编辑或删除分类
4. 在照片编辑时选择分类

### 设置网站信息
1. 登录后台
2. 进入网站设置
3. 填写各项信息
4. 保存即可在前台生效

## 📝 配置示例

### 添加 ICP 备案信息
```json
{
  "icp_record": "京ICP备12345678号",
  "beian_record": "京公网安备11010802012345号",
  "beian_link": "http://www.beian.gov.cn/portal/registerSystemInfo"
}
```

### 添加统计代码
```javascript
// Google Analytics
<script async src="https://www.googletagmanager.com/gtag/js?id=GA_ID"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());
  gtag('config', 'GA_ID');
</script>
```

### 自定义样式
```css
/* 自定义 CSS */
.photo-card {
  border-radius: 20px;
}

.photo-description {
  background: linear-gradient(to top, rgba(0,0,0,0.9), transparent);
}
```

## 🐛 已知问题

- 暂无

## 📅 下一步计划

- [ ] 照片批量上传
- [ ] 照片批量编辑
- [ ] 照片下载功能
- [ ] 照片分享功能
- [ ] 相册功能
- [ ] 评论系统
- [ ] 点赞功能
- [ ] 搜索功能
- [ ] 标签云
- [ ] RSS 订阅

## 💝 感谢使用 MYGallery！

如有问题或建议，请提交 Issue 或 PR。

---

**更新日期**：2025-11-05
**版本**：v1.0.0

