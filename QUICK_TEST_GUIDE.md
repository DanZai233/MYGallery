# 🧪 快速测试指南

测试新功能的快速指南。

## 🚀 启动应用

```bash
cd /root/MYGallery
go run main.go
```

等待看到：
```
✅ 发布完成！MYGallery v1.0.0 启动成功

📷 前台页面: http://localhost:8080
⚙️ 后台管理: http://localhost:8080/admin
```

## ✅ 测试清单

### 1. 测试黑夜模式 🌙

**步骤**：
1. 访问 http://localhost:8080
2. 点击右上角月亮图标 🌙
3. 查看主题切换效果
4. 点击太阳图标 ☀️ 切回

**预期结果**：
- 背景变暗
- 文字高对比度
- 磨砂玻璃效果调整
- 刷新页面保持设置

### 2. 测试照片 Hover 效果 🖼️

**步骤**：
1. 鼠标悬停在照片上
2. 观察底部滑入的描述

**预期结果**：
- 描述从底部平滑滑入
- 显示标题、描述、EXIF、位置
- 毛玻璃背景效果
- 移开鼠标后滑出

### 3. 测试灯箱效果 🔍

**步骤**：
1. 点击任意照片
2. 观察灯箱打开效果
3. 使用 ← → 键切换
4. 按 ESC 关闭

**预期结果**：
- 背景毛玻璃模糊
- 照片居中显示
- EXIF 信息面板有毛玻璃效果
- 键盘导航流畅

### 4. 测试后台访问 🔒

**步骤**：
1. 在首页查看导航栏
2. 确认没有"后台管理"按钮
3. 直接访问 http://localhost:8080/admin
4. 使用 admin / admin123 登录

**预期结果**：
- 首页没有后台入口
- `/admin` 可以正常访问
- 登录成功进入后台

### 5. 测试分类 API 📁

**创建分类**：
```bash
# 先登录获取 token
TOKEN=$(curl -s http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.token')

# 创建分类
curl -X POST http://localhost:8080/api/categories \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "风景",
    "slug": "landscape",
    "description": "风景照片分类",
    "sort_order": 1
  }'
```

**获取分类列表**：
```bash
curl http://localhost:8080/api/categories
```

**预期结果**：
```json
[
  {
    "id": 1,
    "name": "风景",
    "slug": "landscape",
    "description": "风景照片分类",
    "sort_order": 1,
    "photo_count": 0
  }
]
```

### 6. 测试网站设置 API ⚙️

**获取设置**：
```bash
curl http://localhost:8080/api/settings
```

**更新设置**：
```bash
curl -X PUT http://localhost:8080/api/settings \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "site_title": "我的照片墙",
    "site_description": "记录生活的美好瞬间",
    "icp_record": "京ICP备12345678号",
    "contact_email": "admin@example.com"
  }'
```

**预期结果**：
```json
{
  "success": true,
  "message": "设置更新成功",
  "settings": {
    "site_title": "我的照片墙",
    ...
  }
}
```

### 7. 测试照片分类筛选 🔍

**按分类获取照片**：
```bash
# 获取所有照片
curl http://localhost:8080/api/photos

# 获取风景分类的照片
curl http://localhost:8080/api/photos?category=风景

# 分页获取
curl http://localhost:8080/api/photos?category=风景&page=1&size=10
```

### 8. 测试图片显示 🖼️

**上传一张测试图片**：
1. 登录后台
2. 点击"上传图片"
3. 选择一张照片
4. 填写标题、描述、选择分类
5. 保存

**验证**：
1. 在首页查看
2. 确认缩略图正常显示
3. 点击查看大图
4. Hover 查看描述

## 🐛 常见问题

### Q: 缩略图显示不出来？

**A**: 检查：
1. `uploads/thumbnails/` 目录是否存在
2. 目录权限是否正确
3. 查看浏览器控制台错误

```bash
# 创建目录
mkdir -p uploads/thumbnails

# 设置权限
chmod 755 uploads uploads/thumbnails
```

### Q: 黑夜模式不保存？

**A**: 检查浏览器是否禁用了 localStorage

```javascript
// 在浏览器控制台测试
localStorage.setItem('theme', 'dark');
localStorage.getItem('theme');
```

### Q: API 返回 401 未授权？

**A**: Token 可能过期，重新登录获取

```bash
# 重新登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### Q: 照片 Hover 没有描述？

**A**: 检查照片是否有描述字段

```bash
# 更新照片添加描述
curl -X PUT http://localhost:8080/api/photos/1 \
  -H "Authorization: Bearer $TOKEN" \
  -F "description=这是一张美丽的风景照"
```

## 📊 性能测试

### 加载速度
```bash
# 测试首页加载时间
curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8080/
```

### 并发测试
```bash
# 使用 Apache Bench
ab -n 1000 -c 10 http://localhost:8080/

# 使用 wrk
wrk -t4 -c100 -d30s http://localhost:8080/
```

## ✅ 测试完成清单

- [ ] 黑夜模式切换
- [ ] 照片 Hover 效果
- [ ] 灯箱毛玻璃效果
- [ ] 后台访问限制
- [ ] 创建分类
- [ ] 更新网站设置
- [ ] 按分类筛选照片
- [ ] 上传照片并显示
- [ ] 缩略图正常显示
- [ ] 响应式布局正常

## 🎉 全部通过？

恭喜！所有功能正常工作！

接下来可以：
1. 上传更多照片
2. 创建更多分类
3. 自定义网站设置
4. 部署到生产环境

---

**需要帮助？** 查看 [TROUBLESHOOTING.md](TROUBLESHOOTING.md)

