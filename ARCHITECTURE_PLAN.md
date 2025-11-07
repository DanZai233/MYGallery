## ChronoFrame-Style Revamp Plan

### Goals
- Match the ChronoFrame user experience and feature breadth while keeping the Go backend.
- Support both SQLite (default) and MySQL databases with hot-swappable configuration.
- Preserve SSR deployment by wiring the Nuxt frontend to the new Go REST API.
- Deliver first-class local storage: deep directory layouts, unlimited quota, backup workflows.

### Frontend Integration
- Fork or vendor the ChronoFrame Nuxt app and keep its Tailwind-based UI.
- Replace API calls with a typed client talking to the Go backend (REST or tRPC-style JSON over HTTP).
- Maintain SSR mode: Nuxt runs as its own service, but shares auth via HTTP-only JWT cookies.
- Expose a system settings dashboard screen for live configuration edits (DB, storage, site metadata).

### Backend Architecture (Go)
- Entry: `cmd/server/main.go` launching Gin with modular routers.
- Domain slices:
  - `internal/app` – bootstrap, DI wiring, background jobs.
  - `internal/config` – load + persist settings, change notifications.
  - `internal/database` – GORM setup, repository factories for SQLite/MySQL.
  - `internal/media` – upload pipeline, EXIF parsing, thumbnail tasks, dedupe.
  - `internal/storage` – pluggable drivers (`local`, future `s3`).
  - `internal/backup` – scheduled sync jobs and manual triggers.
  - `internal/api` – handlers grouped by resource; use DTOs aligned with ChronoFrame schema.
- Authentication: JWT with refresh tokens, optional OAuth ready hook.
- Background jobs: worker pool powered by Go routines + channel queue; ready to swap for RabbitMQ/Redis later.

### Database Strategy
- ORM: GORM with separate DSN builders per engine.
- Migration layer: goose or atlas; store migration metadata in DB itself.
- Hot swap flow:
  1. Admin submits new DB settings via `/api/admin/settings/database`.
  2. Backend validates connection, runs migrations, updates persistent config.
  3. Connection pool rotated without server restart.
- Default SQLite path: `data/mygallery.db`; MySQL example DSN emitted in UI helper.

### Local Storage Roadmap
- Interface: `storage.Driver` with methods for `Put`, `Get`, `Delete`, `Walk`, `GenerateSignedURL` (no-op for local).
- Local driver features:
  - Hierarchical layout rules (per user/album/date) configurable in admin UI.
  - Thumbnail + derivative directories adjacent to originals.
  - Integrity manifests (`sha256`) stored in DB for backup verification.
- Backup module:
  - Targets: secondary filesystem path (rsync), optional S3-compatible bucket.
  - Scheduling: cron expressions saved in config, executed by background scheduler.
  - Manual trigger endpoint with progress events (SSE/websocket optional later).

### Settings & Config Persistence
- Replace static YAML reliance with dynamic config store table (`system_settings`).
- On startup: load defaults, overlay database values, keep YAML as bootstrap only.
- Provide audit log of config changes with actor, diff, timestamp.
- Settings API secured to admin role; frontend uses form and tests connection before commit.

### Deployment Layout
- Docker Compose services:
  - `api`: Go server, port 8080, mounts uploads + config volume.
  - `frontend`: Nuxt SSR container, depends on `api`.
  - `mysql`: optional; default stack uses SQLite so container disabled by default.
  - `backup-runner`: future cron-style container reusing API binary with `backup run` command.
- Environment variables minimal; runtime config stored via API.
- CI pipeline builds multi-arch Docker images, runs Go tests + lint, Nuxt build tests.

### Immediate Next Actions
1. Scaffold new Go project layout under `cmd/` + `internal/` as outlined.
2. Implement config service with SQLite persistence and HTTP endpoints.
3. Adapt existing handlers to new layering, beginning with auth + settings.
4. Start Nuxt extraction and document API contract gaps versus ChronoFrame.
