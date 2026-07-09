# ============================================================
# Alike Makefile — 统一开发命令
# 使用：make <target>，无参数或 make help 查看全部命令
# ============================================================

# --- 配置 ---
COMPOSE        := docker compose
ENV_FILE       := .env
BACKEND_DIR    := backend
FRONTEND_DIR   := frontend
MIGRATIONS_DIR := backend/migrations
PG_SERVICE     := postgres

# 从 .env 读取 DB 参数（存在则加载），提供默认值
ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export
endif
POSTGRES_USER ?= alike
POSTGRES_DB   ?= alike

.DEFAULT_GOAL := help
.PHONY: help check-docker check-env dev up down build test test-backend \
        test-frontend migrate lint lint-backend lint-frontend seed clean logs ps

# ---------- 辅助检查 ----------
check-docker:
	@command -v docker >/dev/null 2>&1 || { echo "✗ 未找到 docker，请先安装 Docker Desktop"; exit 1; }
	@docker compose version >/dev/null 2>&1 || { echo "✗ 未找到 docker compose 插件"; exit 1; }

check-env:
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "⚠ 未找到 $(ENV_FILE)，正在从 .env.example 复制..."; \
		cp .env.example $(ENV_FILE); \
		echo "→ 已生成 $(ENV_FILE)，请检查其中的占位密码后重试"; \
	fi

# ---------- 开发环境 ----------
dev: check-docker check-env ## 启动完整开发环境（构建并后台运行）
	@echo "→ 启动 Alike 开发环境..."
	$(COMPOSE) up -d --build
	@echo "✓ 已启动。入口: http://localhost:$${NGINX_HTTP_PORT:-80}  MinIO 控制台: http://localhost:$${MINIO_CONSOLE_PORT:-9001}"

up: check-docker check-env ## 启动服务（不强制重建）
	$(COMPOSE) up -d

down: check-docker ## 停止并移除容器（保留数据卷）
	$(COMPOSE) down

logs: check-docker ## 跟踪所有服务日志
	$(COMPOSE) logs -f --tail=100

ps: check-docker ## 查看服务状态
	$(COMPOSE) ps

# ---------- 构建 ----------
build: check-docker ## 构建前后端镜像
	$(COMPOSE) build

# ---------- 测试 ----------
test: test-backend test-frontend ## 运行前后端全部测试

test-backend: ## 运行 Go 后端测试
	@echo "→ 运行后端测试..."
	@cd $(BACKEND_DIR) && go test ./... || { echo "✗ 后端测试失败"; exit 1; }

test-frontend: ## 运行前端测试
	@echo "→ 运行前端测试..."
	@cd $(FRONTEND_DIR) && npm test

# ---------- 数据库迁移 ----------
migrate: check-docker ## 执行 backend/migrations 下的 SQL 迁移（按文件名排序）
	@echo "→ 执行数据库迁移..."
	@if [ -z "$$(ls -A $(MIGRATIONS_DIR)/*.sql 2>/dev/null)" ]; then \
		echo "⚠ $(MIGRATIONS_DIR) 下暂无 .sql 迁移文件，跳过"; \
	else \
		for f in $$(ls $(MIGRATIONS_DIR)/*.sql | sort); do \
			echo "  → 应用 $$f"; \
			$(COMPOSE) exec -T $(PG_SERVICE) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -v ON_ERROR_STOP=1 < $$f || exit 1; \
		done; \
		echo "✓ 迁移完成"; \
	fi

# ---------- 代码检查 ----------
lint: lint-backend lint-frontend ## 前后端代码检查

lint-backend: ## Go 代码检查（优先 golangci-lint，回退 go vet）
	@echo "→ 后端 lint..."
	@cd $(BACKEND_DIR) && \
		if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run ./...; \
		else echo "  (未装 golangci-lint，回退 go vet)"; go vet ./...; fi

lint-frontend: ## 前端代码检查
	@echo "→ 前端 lint..."
	@cd $(FRONTEND_DIR) && npm run lint

# ---------- 种子数据 ----------
seed: check-docker ## 灌入种子数据（预设频道等）
	@echo "→ 灌入种子数据..."
	@if [ -f $(MIGRATIONS_DIR)/seed.sql ]; then \
		$(COMPOSE) exec -T $(PG_SERVICE) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -v ON_ERROR_STOP=1 < $(MIGRATIONS_DIR)/seed.sql && echo "✓ 种子数据完成"; \
	else \
		echo "⚠ 未找到 $(MIGRATIONS_DIR)/seed.sql，跳过（待后端角色提供）"; \
	fi

# ---------- 清理 ----------
clean: check-docker ## 停止并删除容器、数据卷、构建产物
	@echo "→ 清理容器与数据卷..."
	$(COMPOSE) down -v --remove-orphans
	@rm -rf $(FRONTEND_DIR)/.output $(FRONTEND_DIR)/.nuxt $(BACKEND_DIR)/server 2>/dev/null || true
	@echo "✓ 清理完成"

# ---------- 帮助 ----------
help: ## 显示本帮助
	@echo "Alike 开发命令："
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'
