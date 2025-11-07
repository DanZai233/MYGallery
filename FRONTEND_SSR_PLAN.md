## ChronoFrame-Style Frontend & SSR Integration

参考项目：[ChronoFrame](https://github.com/HoshinoSuzumi/chronoframe)

### Nuxt 应用策略
- 直接 fork ChronoFrame 的 Nuxt 代码，保留 TailwindCSS + TypeScript 结构。
- 替换数据层：集中在 `app/plugins/api`、`app/composables` 等位置重写 HTTP 客户端，指向 Go 后端 REST 接口。
- 维护 SSR：Nuxt 以 `npm run dev` / `pnpm dev` 方式支持本地调试；生产环境运行 `node .output/server/index.mjs`，并通过环境变量指向 API。
- 认证：采用 HttpOnly Cookie 携带 Access/Refresh Token；Nuxt `server/plugins/auth.ts` 捕获请求预处理并在 SSR 场景注入用户态。

### API 协议映射
- 根据 ChronoFrame 前端的 `server/api` 路由签名，梳理所需接口（照片、相册、地图、设置、用户）。
- 产出 TypeScript DTO，对应 Go API 的 JSON Schema；两侧共享字段命名与分页结构。
- 为上传和批处理提供 `multipart/form-data` 与事件轮询（SSE）端点，以兼容 ChronoFrame 的上传界面和进度反馈。

### Docker Compose 布局
- `frontend-ssr` 服务：
  - 基于 `node:20-alpine` 构建镜像，安装 pnpm，复制 Nuxt 源码。
  - 构建阶段运行 `pnpm install && pnpm build`，运行阶段执行 `node .output/server/index.mjs`。
  - 暴露 3000 端口，代理到 API（通过 `NUXT_API_BASE_URL=http://api:8080`）。
- `api` 服务：
  - 复用新的 Go 程序，监听 8080，提供 REST/上传接口。
- `mysql`（可选）：采用 `depends_on`，默认注释，可通过配置 UI 切换。
- 共享卷：
  - `uploads`、`thumbnails` 通过 Bind 挂载在 API；前端仅通过 URL 访问。
  - `config-data` 保存持久化设置信息。

### CI/CD 流程
- 前端流水线：`pnpm lint && pnpm test`，然后 `pnpm build` 生成 `.output`。
- 后端流水线：`go test ./...`、`go vet`、`staticcheck`（后续添加）。
- Docker 多阶段构建生成 `frontend`、`api` 两个镜像，由 `docker-compose` 或 K8s 组合部署。
- 支持多架构（linux/amd64, linux/arm64），使用 BuildKit v0.17+ 或分层构建脚本。

### 配置中心与 UI 集成
- 在 Nuxt 后台新增 “系统设置” 页面：
  - 读取 `/api/admin/settings`，展示数据库、本地存储、站点信息。
  - 修改后调用 `PUT` 接口，实时刷新前端状态（使用 Pinia store）。
- SSR 期间（`server/plugins/settings.server.ts`）读取缓存的设置信息，减少重复请求。
- 若配置变化涉及重建连接（如数据库切换），Go 端返回任务 ID；前端轮询等待完成并提示。

### 路线图
1. 提取 ChronoFrame Nuxt 代码并建立独立包管理（pnpm workspace）。
2. 编写 API 客户端层，与 Go 后端约定接口。
3. 搭建 docker-compose 新版（api + frontend + mysql 可选）。
4. 构建 SSR 运行镜像，验证登录、上传、地图等核心流程。
5. 在管理后台实现在线配置修改与热重载联动。
