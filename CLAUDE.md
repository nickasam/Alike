# CLAUDE.md — Alike 项目指南

> **本文件是所有 Claude Code 代理的项目上下文入口。每次执行任务前必读。**

---

## 项目简介

**Alike** — 汇聚天下牛马，总有人懂你的辛苦。

面向打工人的情感共鸣型 Web 聊天社区，用频道+线程让牛马们吐槽、互助、抱团取暖。

### 核心价值
- **共鸣** — 总有人懂你的辛苦，你不是一个人在扛
- **归属** — 按行业/岗位分频道，快速找到自己人
- **宣泄** — 匿名吐槽、情绪标签，健康释放压力
- **互助** — 线程式讨论，经验分享，问答互助

---

## 技术栈

| 层 | 技术 | 说明 |
|---|------|------|
| 前端 | Vue 3 + Nuxt 3 | SSR + SPA 混合，SEO 友好 |
| 后端 | Go (Gin) | 高性能、并发友好 |
| 数据库 | PostgreSQL | 源数据存储 |
| 缓存 | Redis | 会话管理 + WebSocket Pub/Sub 广播 |
| 实时通信 | WebSocket | 消息实时推送 |
| 对象存储 | MinIO / S3 | 文件/图片存储 |
| 部署 | Docker + Docker Compose | 容器化部署 |
| CSS | Tailwind CSS | 原子化 CSS，响应式 + 暗色主题 |

---

## 项目结构

```
Alike/
├── CLAUDE.md                            ← 本文件（项目上下文入口）
├── docs/
│   ├── plans/
│   │   └── architecture-design.md       ← 架构设计实施计划（权威文档）
│   ├── prd/
│   │   └── PRD.md                        ← 产品需求文档
│   ├── stories/
│   │   └── user-stories.md              ← 用户故事
│   ├── qa/
│   │   └── test-plan.md                 ← 测试计划
│   ├── design/
│   │   ├── design-system.md             ← 设计规范（token/主题/组件）
│   │   ├── component-spec.md            ← 组件规格（props/事件）
│   │   ├── interaction-spec.md          ← 交互规范
│   │   ├── 01-home.html                 ← 首页高保真稿
│   │   ├── 02-channel.html              ← 频道页高保真稿
│   │   ├── 03-diary.html                ← 日记页高保真稿
│   │   ├── 04-ranking.html              ← 排行榜高保真稿
│   │   └── 05-profile.html              ← 个人主页高保真稿
│   ├── architecture-diagram.html        ← 架构图
│   └── team-roles.md                    ← 团队角色与分工
│
├── backend/                             ← Go 后端（待创建）
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── auth/          user/         channel/     message/
│   │   ├── emotion/       empathy/      diary/       ws/
│   │   ├── notification/  search/       storage/     middleware/
│   ├── pkg/
│   │   ├── database/      redis/        jwt/         config/
│   ├── migrations/
│   ├── go.mod
│   └── go.sum
│
├── frontend/                            ← Vue3 / Nuxt3 前端（待创建）
│   ├── nuxt.config.ts
│   ├── package.json
│   ├── pages/             components/   composables/  stores/
│   └── assets/css/
│
├── nginx/nginx.conf                     ← Nginx 反代配置（待创建）
├── docker-compose.yml                   ← 开发环境（待创建）
├── Makefile                             ← 统一命令（待创建）
└── .env.example                         ← 环境变量模板（待创建）
```

---

## 团队角色

本项目由 **Hermes（架构师）** 统筹，分配任务给 7 个 Claude Code 代理：

| 角色 | 代号 | 职责 | 输出目录 |
|------|------|------|---------|
| 产品经理 | `claude-pm` | PRD、用户故事、优先级 | `docs/prd/`、`docs/stories/` |
| 设计师 | `claude-designer` | 交互设计、信息架构、流程图 | `docs/design/wireframes/` |
| 美工 | `claude-visual` | 高保真 UI、视觉规范 | `docs/design/` |
| 后端开发 | `claude-backend` | Go API、WebSocket、中间件 | `backend/` |
| 前端开发 | `claude-frontend` | Vue3/Nuxt3 页面、组件 | `frontend/` |
| DevOps | `claude-devops` | Docker、Nginx、CI/CD | 根目录 + `nginx/` |
| QA | `claude-qa` | 测试用例、自动化测试 | `backend/**/*_test.go`、`frontend/test/` |

> **你被分配任务时会被告知你的角色。请只做你角色范围内的工作。**

---

## 架构概览

```
用户浏览器 (Vue3 + Nuxt3)
    │
    ├── HTTP/REST ──────┐
    └── WebSocket ──────┤
                        ▼
                    Nginx (反代)
                        │
                        ▼
                Go Backend (Gin)
    ┌──────────┬──────────┬──────────┐
    │  Auth    │ Channel  │ Message  │  ...其他模块
    └──────────┴──────────┴──────────┘
              WebSocket Hub
    ┌─────────┬──────────────┬────────────┐
    │PostgreSQL│    Redis     │  MinIO/S3  │
    │ (源数据) │(会话+Pub/Sub)│  (文件存储) │
    └─────────┴──────────────┴────────────┘
```

### 后端模块（12个）

| 模块 | 路径 | 职责 |
|------|------|------|
| auth | `/internal/auth` | 注册、登录、JWT 签发与验证 |
| user | `/internal/user` | 用户信息、牛马等级、个人主页 |
| channel | `/internal/channel` | 频道 CRUD、成员管理 |
| message | `/internal/message` | 消息发布、线程回复、历史查询 |
| emotion | `/internal/emotion` | 情绪标签、情绪看板、趋势统计 |
| empathy | `/internal/empathy` | 抱团取暖、共情统计、排行榜 |
| diary | `/internal/diary` | 打工日记 CRUD、打卡日历 |
| ws | `/internal/ws` | WebSocket Hub、连接管理、事件广播 |
| notification | `/internal/notification` | @提及、共情、回复通知 |
| search | `/internal/search` | 全文搜索 |
| storage | `/internal/storage` | MinIO/S3 文件上传 |
| middleware | `/internal/middleware` | JWT 鉴权、限流、CORS、日志 |

### 前端页面（6个）

| 页面 | 路径 | 功能 |
|------|------|------|
| 首页 | `pages/index.vue` | 热门频道 + 今日牛马榜 |
| 频道页 | `pages/channel/[id].vue` | 消息流 + 情绪看板 |
| 日记广场 | `pages/diary/index.vue` | 公开日记流 |
| 日记详情 | `pages/diary/[id].vue` | 日记内容 + 评论 |
| 排行榜 | `pages/ranking.vue` | 牛马排行榜 |
| 个人主页 | `pages/profile/[id].vue` | 用户信息 + 统计 |

### 前端组件

| 组件 | 路径 | 功能 |
|------|------|------|
| 消息列表 | `components/chat/MessageList.vue` | 消息流 + 分页加载 |
| 消息输入 | `components/chat/MessageInput.vue` | 输入框 + 表情 + 情绪标签 |
| 线程面板 | `components/chat/ThreadPanel.vue` | 线程回复 UI |
| 情绪选择器 | `components/emotion/EmotionPicker.vue` | 发消息时选情绪 |
| 情绪看板 | `components/emotion/EmotionBoard.vue` | 实时情绪热力图 |
| 共情按钮 | `components/empathy/EmpathyButton.vue` | 抱团取暖 + 动画 |
| 频道侧边栏 | `components/layout/ChannelSidebar.vue` | 频道列表 + 分类 |
| 顶部导航 | `components/layout/TopNav.vue` | 导航 + 通知 |

---

## 数据库核心表

| 表 | 说明 |
|----|------|
| `users` | 用户：邮箱(CITEXT)、密码、昵称、头像、简介、行业、岗位、工龄、牛马等级、被共情数、给出共情数、累计打卡天数、匿名开关 |
| `channels` | 频道：名称、slug、描述、分类(industry/job/topic/custom)、图标、状态(pending/active/archived)、成员数 |
| `channel_members` | 频道成员：频道ID、用户ID、角色(member/admin) |
| `messages` | 消息：频道ID、用户ID、父消息ID(线程)、内容、情绪标签、匿名、共情数、软删除(deleted_at) |
| `empathies` | 共情：消息ID、用户ID (唯一约束: 每人每条消息一次) |
| `diaries` | 日记(v1.5)：用户ID、标题、内容、心情、公开/私密 (打卡天数在 users 与应用层计算) |
| `diary_comments` | 日记评论(v1.5)：日记ID、用户ID、内容、匿名、软删除 |
| `notifications` | 通知：用户ID、类型(mention/empathy/reply/system)、内容、已读 |
| `attachments` | 附件：消息ID、文件URL、类型、大小 |

> 完整 SQL 建表语句见 `docs/plans/architecture-design.md` 第四章。

---

## API 契约

### REST API 概要

| 模块 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 认证 | POST | `/api/auth/register` | 注册 |
| 认证 | POST | `/api/auth/login` | 登录 |
| 认证 | POST | `/api/auth/refresh` | 刷新 JWT |
| 认证 | GET | `/api/auth/me` | 当前用户信息 |
| 用户 | GET | `/api/users/:id` | 用户主页 |
| 用户 | PUT | `/api/users/:id` | 更新资料 |
| 频道 | GET | `/api/channels` | 频道列表 |
| 频道 | GET | `/api/channels/:id` | 频道详情 |
| 频道 | POST | `/api/channels/:id/join` | 加入频道 |
| 频道 | POST | `/api/channels/:id/leave` | 离开频道 |
| 消息 | GET | `/api/channels/:id/messages` | 消息列表(分页) |
| 消息 | POST | `/api/channels/:id/messages` | 发布消息 |
| 消息 | GET | `/api/messages/:id/threads` | 线程回复 |
| 消息 | POST | `/api/messages/:id/replies` | 线程回复 |
| 消息 | DELETE | `/api/messages/:id` | 删除消息 |
| 共情 | POST | `/api/messages/:id/empathy` | 抱团取暖 |
| 共情 | DELETE | `/api/messages/:id/empathy` | 取消共情 |
| 日记 | GET | `/api/diaries` | 日记广场 |
| 日记 | POST | `/api/diaries` | 写日记 |
| 排行榜 | GET | `/api/ranking/empathy` | 最受共情榜 |
| 通知 | GET | `/api/notifications` | 通知列表 |
| 搜索 | GET | `/api/search?q=xxx&type=message` | 搜索消息 |

### WebSocket 端点

```
WS /ws?token=<JWT>
```

**客户端 → 服务端事件：**
- `join_channel` / `leave_channel` / `typing` / `send_message`

**服务端 → 客户端事件：**
- `new_message` / `thread_reply` / `empathy` / `user_joined` / `emotion_update` / `notification`

> 完整 API 列表见 `docs/plans/architecture-design.md` 第五章。

---

## 开发命令（Makefile）

```bash
make dev        # 启动开发环境 (docker-compose up + 热重载)
make build      # 构建前后端
make test       # 运行所有测试
make migrate    # 执行数据库迁移
make lint       # 代码检查
make seed       # 灌入种子数据
make clean      # 清理构建产物
```

---

## 编码规范

### Go 后端
- 遵循 [Effective Go](https://go.dev/doc/effective_go) 和标准项目布局
- 包名小写单数：`auth`、`channel`、`message`
- 错误处理：返回 `error`，不 panic
- 配置：从环境变量加载，支持 `.env` 文件
- API 响应统一格式：
  ```json
  { "code": 0, "message": "success", "data": { ... } }
  ```
- 分页响应格式：
  ```json
  { "code": 0, "message": "success", "data": { "list": [...], "total": 100, "page": 1, "page_size": 20 } }
  ```
- JWT 认证：Header `Authorization: Bearer <token>`
- 数据库：使用 `database/sql` + `pgx` 驱动
- Redis：使用 `go-redis/redis/v9`
- WebSocket：使用 `gorilla/websocket`
- 密码：bcrypt hash
- 文件上传：限制类型(image/document)和大小(图片 ≤5MB、文档 ≤10MB)，服务端 MIME 校验

### Vue3 / Nuxt3 前端
- `<script setup lang="ts">` 组合式 API
- 组件名 PascalCase：`MessageList.vue`
- Props 用 `defineProps<T>()` 类型标注
- 状态管理用 Pinia
- API 调用封装在 `composables/useApi.ts`
- WebSocket 封装在 `composables/useWebSocket.ts`
- 样式：Tailwind CSS，暗色主题用 `dark:` 前缀
- 路由：Nuxt 文件路由，`pages/` 目录自动生成

---

## 分阶段实施计划

| 阶段 | 内容 | 状态 |
|------|------|------|
| 阶段零 | 架构设计文档 + UI 设计稿 | ✅ 已完成 |
| 阶段一 | 项目脚手架 (Go+Nuxt+Docker+DB迁移) | ✅ 已完成 |
| 阶段二 | 用户认证系统 | ⏳ 待开始 |
| 阶段三 | 频道系统 | ⏳ 待开始 |
| 阶段四 | 消息与线程 + WebSocket | ⏳ 待开始 |
| 阶段五 | 情绪系统 | ⏳ 待开始 |
| 阶段六 | 共情机制 | ⏳ 待开始 |
| 阶段七 | 打工日记 | ⏳ 待开始 |
| 阶段八 | 排行榜与通知 | ⏳ 待开始 |
| 阶段九 | 搜索 + 文件上传 + 打磨部署 | ⏳ 待开始 |

> 完整计划见 `docs/plans/architecture-design.md` 第七章。

---

## 任务执行规则

1. **收到任务卡后** — 仔细阅读目标、输入、输出、验收标准
2. **只做角色范围内的工作** — 不要越界改其他角色的文件
3. **遵循架构设计** — 模块划分、API 路径、表结构以 `architecture-design.md` 为准
4. **增量提交** — 完成一个子任务即可提交，不要攒一大坨
5. **自测通过再交付** — 至少保证编译通过、基本逻辑正确
6. **遇到模糊处** — 参考 `docs/design/` 下的高保真稿和架构文档，不要自行发挥
7. **提交信息格式** — `feat(模块): 描述` / `fix(模块): 描述` / `docs(模块): 描述`
8. **不要修改 `docs/plans/architecture-design.md`** — 这是架构师的权威文档，如需变更请反馈

---

## 关键文档索引

| 文档 | 路径 | 用途 |
|------|------|------|
| 架构设计实施计划 | `docs/plans/architecture-design.md` | 技术架构、数据库、API、实施计划 |
| 产品需求文档 | `docs/prd/PRD.md` | 功能需求、优先级、验收标准（PM 权威输入） |
| 用户故事 | `docs/stories/user-stories.md` | INVEST 用户故事与验收条件 |
| 测试计划 | `docs/qa/test-plan.md` | 测试策略、用例、覆盖率目标（QA 权威输入） |
| 设计规范 | `docs/design/design-system.md` | 设计 token、主题、间距、字体（前端权威输入） |
| 组件规格 | `docs/design/component-spec.md` | 组件 props/事件/结构定义 |
| 交互规范 | `docs/design/interaction-spec.md` | 交互流程与状态定义 |
| 团队角色分工 | `docs/team-roles.md` | 角色定义、协作流程、任务卡格式 |
| 架构图 | `docs/architecture-diagram.html` | 系统架构可视化 |
| 首页设计稿 | `docs/design/01-home.html` | 首页高保真 |
| 频道页设计稿 | `docs/design/02-channel.html` | 频道页高保真 |
| 日记页设计稿 | `docs/design/03-diary.html` | 日记页高保真 |
| 排行榜设计稿 | `docs/design/04-ranking.html` | 排行榜高保真 |
| 个人主页设计稿 | `docs/design/05-profile.html` | 个人主页高保真 |

---

## 注意事项

- **WebSocket 多实例** — 必须用 Redis Pub/Sub 广播，不能仅靠内存
- **匿名隐私** — 匿名消息不返回 `user_id`，JWT 中不暴露匿名消息关联
- **情绪看板实时性** — WebSocket 推送情绪更新，Redis 缓存聚合数据
- **文件上传安全** — 限制类型和大小，图片压缩
- **消息分页** — 按 `channel_id + created_at DESC` 索引分页，不用 OFFSET 翻深页
