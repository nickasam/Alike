<p align="center">
  <img src="./docs/brand/banner.svg" alt="Alike — 汇聚天下牛马" height="128"/>
</p>

<p align="center">
  面向打工人的情感共鸣型 Web 聊天社区。<br/>
  <b>共鸣 · 归属 · 宣泄 · 互助</b>
</p>

---

## 技术栈

- **前端** — Vue 3 / Nuxt 3 / Tailwind CSS
- **后端** — Go / Gin
- **数据** — PostgreSQL · Redis · MinIO
- **实时** — WebSocket
- **部署** — Docker Compose + Nginx

---

## 快速开始

```bash
git clone git@github.com:nickasam/Alike.git
cd Alike

cp .env.example .env    # 填入密码 / 密钥
make dev                # 拉起完整环境
make migrate            # 应用数据库迁移
```

浏览器访问 http://localhost —— 服务由 Nginx 反代统一暴露在 80 端口。

其它常用命令：`make down` · `make test` · `make lint` · `make seed` · `make logs` · `make ps` · `make clean`。完整清单：`make help`。

---

## 项目结构

```
backend/       Go 后端（Gin） · 业务模块在 internal/  · SQL 迁移在 migrations/
frontend/      Nuxt 3 前端    · 页面在 pages/  · 组件在 components/
nginx/         Nginx 反代配置
docs/          设计与规划文档（架构 · PRD · UI 稿 · 测试计划）
docker-compose.yml
Makefile
.env.example
CLAUDE.md      项目上下文入口
```

---

## 编码约定

- **Go 后端** — 每个业务模块四段式：`model.go` / `repository.go` / `handler.go` / `handler_test.go`；`database/sql` + `pgx`；API 统一响应 `{code, message, data}`；JWT 走 `Authorization: Bearer <token>`。
- **前端** — `<script setup lang="ts">` + Pinia + Tailwind；API/WS 封装在 `composables/`；暗色主题用 `.light` class 变量覆盖。
- **Git** — 分支 `feat/<阶段>-<模块>-<描述>`；commit `feat(模块): 描述` / `fix(模块): 描述` / `docs(模块): 描述`；提交前 `make test` 必须通过；不直接推 `main`，走 PR review。

---

## 更多

- 架构、模块、API、实施进度 → `CLAUDE.md`
- 设计与规划文档 → `docs/`
- 品牌资源（logo · banner · 图标）→ `docs/brand/`

---

## License

私有项目，未开源。© Alike Team.
