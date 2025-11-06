#!/bin/bash
# 使用传统构建方式启动服务（兼容旧版 buildx）
DOCKER_BUILDKIT=0 COMPOSE_DOCKER_CLI_BUILD=0 docker compose up -d "$@"

