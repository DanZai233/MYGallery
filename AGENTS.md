# AGENTS.md

## Cursor Cloud specific instructions

### Project Overview

MYGallery is a self-hosted personal photo gallery system built with Go (Gin framework) + SQLite + static HTML/JS frontend. Single binary serves both the REST API and the frontend on port 8080.

### Running the Application

Standard commands are in `Makefile` and `README.md`. Key commands:

| Task | Command |
|------|---------|
| Run dev server | `go run main.go` (serves on `:8080`) |
| Build binary | `make build` (output: `bin/mygallery`) |
| Run tests | `make test` or `go test -v ./...` |
| Lint | `go vet ./...` |
| Dev mode (hot reload) | `make dev` (requires `air`) |
| Init project dirs | `make init` |

### Non-obvious Notes

- **CGO is required**: SQLite driver (`mattn/go-sqlite3`) requires CGO. Ensure `gcc` is available. The system Go 1.24 and GCC are pre-installed.
- **Config file**: The app requires `config.yaml` at the project root. Copy from `config.example.yaml` if missing. For development, set `server.mode` to `"debug"`.
- **Auto-created directories**: The app auto-creates `data/`, `uploads/`, and `uploads/thumbnails/` on startup, but `make init` can pre-create them.
- **Default credentials**: Admin login is `admin` / `admin123` (configured in `config.yaml`).
- **No external services needed**: Default config uses embedded SQLite and local filesystem storage. No Docker, database server, or object storage required for development.
- **Frontend is static HTML**: The `public/` directory contains static `.html` files served by the Go server. There is no Node.js/npm build step for the frontend.
- **Planned SSR frontend**: `docker-compose.yml` references a `frontend/` Nuxt.js app that does not exist yet. Ignore it for now.
- **Test data script**: `scripts/create-test-data.sh` creates sample categories via the API (requires `jq` and a running server).
