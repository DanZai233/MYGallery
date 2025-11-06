#!/bin/bash
# 使用传统构建方式启动服务（兼容旧版 buildx）
# 禁用 BuildKit 和 Bake 以避免兼容性问题
export DOCKER_BUILDKIT=0
export COMPOSE_DOCKER_CLI_BUILD=0
export COMPOSE_EXPERIMENTAL_DOCKER_BUILDKIT=0
docker compose up -d "$@"

