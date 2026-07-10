# Alike 聊天软件 — 架构设计实施计划

> **产品定位：** 汇聚天下牛马，总有人懂你的辛苦
>
> **一句话描述：** 面向打工人的情感共鸣聊天社区，用频道+线程让牛马们吐槽、互助、抱团取暖。

部署的主机信息， ssh root@39.107.58.169 密码 zheng199512!


---

## 一、产品概述

### 1.1 定位

**Alike** 是一个面向广大"牛马"（打工人）的情感共鸣型 Web 聊天社区。不同行业、不同岗位的打工人可以在这里找到同路人，吐槽日常、分享经验、互相鼓励。

### 1.2 核心价值

- **共鸣** — 总有人懂你的辛苦，你不是一个人在扛
- **归属** — 按行业/岗位分频道，快速找到自己人
- **宣泄** — 匿名吐槽、情绪标签，健康释放压力
- **互助** — 线程式讨论，经验分享，问答互助

### 1.3 技术栈

| 层 | 技术 | 说明 |
|---|------|------|
| 前端 | Vue 3 + Nuxt 3 | SSR + SPA 混合，SEO 友好 |
| 后端 | Go (Gin) | 高性能、并发友好 |
| 数据库 | PostgreSQL | 源数据存储 |
| 缓存 | Redis | 会话管理 + WebSocket Pub/Sub 广播 |
| 实时通信 | WebSocket | 消息实时推送 |
| 对象存储 | MinIO / S3 | 文件/图片存储 |
| 部署 | Docker + Docker Compose | 容器化部署 |

---

## 二、核心功能

### 2.1 用户系统

| 功能 | 描述 |
|------|------|
| 注册/登录 | 邮箱注册，JWT 认证 |
| 个人主页 | 昵称、头像、行业、岗位、工龄、牛马等级 |
| 匿名模式 | 发帖/回复可选择匿名，隐藏身份 |
| 牛马等级 | 根据活跃度、被共情次数累积，自动升级（小牛马 → 老牛马 → 牛马之王） |

### 2.2 频道系统

| 功能 | 描述 |
|------|------|
| 行业频道 | #互联网、#工厂、#外卖骑手、#医护、#教培、#服务业… |
| 岗位频道 | #程序员、#产品经理、#设计师、#流水线、#骑手… |
| 主题频道 | #吐槽大会、#打工日记、#薪资揭秘、#维权互助、#摸鱼技巧… |
| 自建频道 | 用户可申请创建频道（需审核） |
| 频道管理 | 频道描述、规则、管理员、成员数 |

### 2.3 消息与线程

| 功能 | 描述 |
|------|------|
| 发布消息 | 文本、图片、表情、情绪标签 |
| 线程回复 | 在任意消息下开启线程讨论，类似 Slack |
| 实时推送 | WebSocket 实时推送新消息 |
| 消息历史 | 分页加载历史消息 |
| @提及 | @用户通知 |

### 2.4 情绪系统（特色）

| 功能 | 描述 |
|------|------|
| 情绪标签 | 😮‍💨疲惫、😡愤怒、😢委屈、🤯崩溃、😴麻木、🔥想润、😰焦虑、💪加油（共 8 种，与 PRD 一致） |
| 情绪看板 | 频道内实时情绪热力图，看看大家今天都什么状态 |
| 情绪趋势 | 个人/频道情绪趋势统计 |

### 2.5 共情机制（特色）

| 功能 | 描述 |
|------|------|
| 抱团取暖 | 不只是"赞"，是"我懂你" — 共情按钮 |
| 共情统计 | 帖子被共情次数、用户被共情总数、用户给出共情总数 |
| 牛马排行榜 | 最受共情帖子 / 最暖牛马（按给出共情数） / 连续打卡牛马 / 本周最活跃牛马 |

### 2.6 打工日记（v1.5）

> **版本边界：** 打工日记为 v1.5 功能，不在 v1.0 MVP 范围内。以下表结构与 API 预留设计，实施排期见第七章。

| 功能 | 描述 |
|------|------|
| 每日日记 | 记录今天打工感受，可设为私密/公开 |
| 连续打卡 | 打卡日历，连续记录可获得徽章（连续 + 累计打卡天数） |
| 日记广场 | 公开日记流，可共情可评论 |

### 2.7 其他

| 功能 | 描述 |
|------|------|
| 文件上传 | 图片、文档分享 |
| 通知系统 | @提及、共情、回复通知 |
| 搜索 | 频道内/全局搜索消息 |
| 管理 | 用户举报、内容审核、管理员后台 |

---

## 三、系统架构设计

### 3.1 架构总览

```
┌─────────────────────────────────────────────────┐
│                   用户浏览器                      │
│              Vue 3 + Nuxt 3 (SSR)                │
├────────────┬────────────────────────────────────┤
│  HTTP/REST │         WebSocket                  │
├────────────┴────────────────────────────────────┤
│                  Nginx (反代)                     │
├─────────────────────────────────────────────────┤
│               Go Backend (Gin)                   │
│  ┌──────────┬──────────┬──────────┬──────────┐  │
│  │ Auth模块  │ Channel  │ Message  │ Emotion  │  │
│  │          │  模块     │  模块     │  模块    │  │
│  ├──────────┼──────────┼──────────┼──────────┤  │
│  │ User模块  │ Diary    │ Empathy  │  Search  │  │
│  │          │  模块     │  模块    │  模块    │  │
│  └──────────┴──────────┴──────────┴──────────┘  │
│         WebSocket Hub (实时连接管理)               │
├─────────┬──────────────────────┬────────────────┤
│PostgreSQL│       Redis          │  MinIO/S3      │
│ (源数据)  │ (会话+Pub/Sub广播)    │  (文件存储)     │
└─────────┴──────────────────────┴────────────────┘
```

### 3.2 模块划分

#### 后端模块（Go）

| 模块 | 职责 | 关键包 |
|------|------|--------|
| `auth` | 注册、登录、JWT 签发与验证、匿名令牌 | `/internal/auth` |
| `user` | 用户信息、牛马等级、个人主页 | `/internal/user` |
| `channel` | 频道 CRUD、成员管理、频道列表 | `/internal/channel` |
| `message` | 消息发布、线程回复、历史查询 | `/internal/message` |
| `emotion` | 情绪标签、情绪看板、趋势统计 | `/internal/emotion` |
| `empathy` | 共情（抱团取暖）、排行榜 | `/internal/empathy` |
| `diary` | 打工日记 CRUD、打卡日历 | `/internal/diary` |
| `websocket` | WebSocket Hub、连接管理、事件广播 | `/internal/ws` |
| `notification` | @提及、共情、回复通知 | `/internal/notification` |
| `search` | 消息/日记/频道/用户全文搜索 | `/internal/search` |
| `storage` | 文件上传到 MinIO/S3 | `/internal/storage` |
| `middleware` | JWT 鉴权、限流、CORS、请求日志 | `/internal/middleware` |

#### 前端模块（Vue 3 / Nuxt 3）

| 模块 | 职责 |
|------|------|
| `layouts/default` | 顶部导航 + 左侧频道列表 + 右侧内容区 |
| `pages/index` | 首页：热门频道 + 今日牛马榜 |
| `pages/channel/[id]` | 频道页：消息流 + 情绪看板 |
| `pages/diary` | 打工日记广场 |
| `pages/diary/[id]` | 日记详情 |
| `pages/profile/[id]` | 个人主页 |
| `pages/ranking` | 牛马排行榜 |
| `components/chat` | 消息列表、输入框、线程面板 |
| `components/emotion` | 情绪标签选择器、情绪看板 |
| `components/empathy` | 抱团取暖按钮、共情动画 |
| `composables/useWebSocket` | WebSocket 连接管理 |
| `composables/useAuth` | 认证状态管理 |
| `stores/` | Pinia 状态管理（channel、message、user） |

### 3.3 数据流

#### 消息发布流程

```
用户输入消息 → 前端 → REST API (/api/messages)
  → Go 后端验证 + 写入 PostgreSQL
  → 发布事件到 Redis Pub/Sub
  → WebSocket Hub 订阅 → 推送到所有在线频道成员
  → 前端实时显示新消息
```

#### 共情（抱团取暖）流程

```
用户点击"我懂你" → REST API (/api/empathy)
  → 写入 empathy 表 + 更新用户共情统计
  → Redis Pub/Sub 通知
  → WebSocket 推送给被共情者（实时通知）
  → 前端动画 + 数字更新
```

### 3.4 WebSocket 事件协议

> **连接鉴权：** JWT 不放在 URL query（会写入 Nginx access log 泄露），改用首帧鉴权 —
> 连接建立后客户端首帧发送 `{ "type": "auth", "data": { "token": "<JWT>" } }`，服务端校验通过前不接受其他消息，
> 超时（5s）未鉴权则断开。鉴权成功后服务端回 `auth_ok`。
>
> **心跳与重连：** 服务端每 30s 发 `ping`，客户端回 `pong`；客户端断线后指数退避重连（1s→2s→4s…上限 30s），重连后重新 `auth` 并 `join_channel`。
>
> **发送鉴权：** `send_message` / `join_channel` 服务端必须校验该用户是否为目标频道成员，非成员拒绝，防止越权。
>
> **幂等与去重：** 客户端为 `send_message` 附带 `client_msg_id`（UUID）用于去重；服务端多实例经 Redis Pub/Sub 广播时，发送方实例本地已推送的连接跳过再投，避免同实例双投。

```json
// 客户端 → 服务端（业务字段统一放在 data 内）
{ "type": "auth", "data": { "token": "<JWT>" } }
{ "type": "join_channel", "data": { "channel_id": 123 } }
{ "type": "leave_channel", "data": { "channel_id": 123 } }
{ "type": "typing", "data": { "channel_id": 123 } }
{ "type": "send_message", "data": { "channel_id": 123, "content": "...", "emotion": "tired", "is_anonymous": false, "client_msg_id": "uuid" } }
{ "type": "pong" }

// 服务端 → 客户端
{ "type": "ping" }
{ "type": "auth_ok", "data": { "user_id": 1 } }
{ "type": "new_message", "data": { "id": 1, "channel_id": 123, "content": "...", "emotion": "tired", "author": {...}, "client_msg_id": "uuid" } }
{ "type": "thread_reply", "data": { "parent_id": 1, "reply": {...} } }
{ "type": "message_deleted", "data": { "message_id": 1 } }
{ "type": "empathy", "data": { "message_id": 1, "empathy_count": 42 } }
{ "type": "user_joined", "data": { "user_id": 1, "channel_id": 123 } }
{ "type": "emotion_update", "data": { "channel_id": 123, "board": {...} } }
{ "type": "notification", "data": { "type": "mention", "content": "..." } }
{ "type": "error", "data": { "message": "..." } }
```

---

## 四、数据库模型设计

### 4.1 核心表结构

```sql
-- 扩展：email 大小写不敏感唯一，避免 Foo@x.com 与 foo@x.com 重复注册
CREATE EXTENSION IF NOT EXISTS citext;

-- 用户表
CREATE TABLE users (
    id          BIGSERIAL PRIMARY KEY,
    email       CITEXT UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,  -- bcrypt hash
    nickname    VARCHAR(100) NOT NULL,
    avatar_url  VARCHAR(500),
    bio         VARCHAR(200),           -- 个人简介（≤200 字）
    industry    VARCHAR(100),           -- 行业
    job_title   VARCHAR(100),           -- 岗位
    work_years  INT DEFAULT 0 CHECK (work_years >= 0),  -- 工龄
    level       INT DEFAULT 1,          -- 牛马等级
    empathy_received INT DEFAULT 0,     -- 被共情总数
    empathy_given    INT DEFAULT 0,     -- 给出共情总数（支撑"最暖牛马榜"）
    total_check_in_days INT DEFAULT 0,  -- 累计打卡天数（连续天数由 diaries 计算）
    is_anonymous BOOLEAN DEFAULT FALSE,  -- 全局匿名开关
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

-- 频道表
CREATE TABLE channels (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,   -- 频道名（如 #互联网）
    slug        VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    category    VARCHAR(50) NOT NULL,    -- industry/job/topic/custom
    icon        VARCHAR(50),             -- 图标标识
    status      VARCHAR(20) DEFAULT 'active', -- pending(待审核)/active/archived
    member_count INT DEFAULT 0,
    created_by  BIGINT REFERENCES users(id),
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- 频道成员表
CREATE TABLE channel_members (
    id          BIGSERIAL PRIMARY KEY,
    channel_id  BIGINT REFERENCES channels(id) ON DELETE CASCADE,
    user_id     BIGINT REFERENCES users(id) ON DELETE CASCADE,
    role        VARCHAR(20) DEFAULT 'member', -- member/admin
    joined_at   TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(channel_id, user_id)
);

-- 消息表
CREATE TABLE messages (
    id          BIGSERIAL PRIMARY KEY,
    channel_id  BIGINT REFERENCES channels(id) ON DELETE CASCADE,
    user_id     BIGINT REFERENCES users(id),
    parent_id   BIGINT REFERENCES messages(id),  -- NULL=主消息, 非NULL=线程回复
    content     TEXT NOT NULL,
    emotion     VARCHAR(50),             -- 情绪标签
    is_anonymous BOOLEAN DEFAULT FALSE,
    empathy_count INT DEFAULT 0,         -- 共情次数（冗余，加速查询）
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ             -- 软删除：保留线程外键完整性，删除消息时置时间戳
);

-- 共情（抱团取暖）表
CREATE TABLE empathies (
    id          BIGSERIAL PRIMARY KEY,
    message_id  BIGINT REFERENCES messages(id) ON DELETE CASCADE,
    user_id     BIGINT REFERENCES users(id),
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(message_id, user_id)  -- 每人每条消息只能共情一次
);

-- 打工日记表（v1.5）
CREATE TABLE diaries (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT REFERENCES users(id),
    title       VARCHAR(200),
    content     TEXT NOT NULL,
    mood        VARCHAR(50),             -- 今日心情
    is_public   BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMPTZ DEFAULT NOW()
    -- 连续/累计打卡天数在 users(total_check_in_days) 及应用层按 created_at 计算，
    -- 不放在每条日记行上，避免断签重算逻辑无处落地
);

-- 日记评论表（v1.5，支撑 2.6"日记广场可评论"）
CREATE TABLE diary_comments (
    id          BIGSERIAL PRIMARY KEY,
    diary_id    BIGINT REFERENCES diaries(id) ON DELETE CASCADE,
    user_id     BIGINT REFERENCES users(id),
    content     TEXT NOT NULL,
    is_anonymous BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

-- 通知表
CREATE TABLE notifications (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT REFERENCES users(id),  -- 接收者
    type        VARCHAR(50) NOT NULL,    -- mention/empathy/reply/system
    content     TEXT,
    ref_id      BIGINT,                  -- 关联的消息/日记 ID
    is_read     BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- 文件附件表
CREATE TABLE attachments (
    id          BIGSERIAL PRIMARY KEY,
    message_id  BIGINT REFERENCES messages(id) ON DELETE CASCADE,
    file_url    VARCHAR(500) NOT NULL,
    file_type   VARCHAR(50),             -- image/document/video
    file_size   BIGINT,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);
```

### 4.2 索引

```sql
-- 频道主消息流：仅索引主消息（parent_id IS NULL）且未软删，避免线程回复污染列表分页
CREATE INDEX idx_messages_channel_created ON messages(channel_id, created_at DESC)
    WHERE parent_id IS NULL AND deleted_at IS NULL;
CREATE INDEX idx_messages_parent ON messages(parent_id) WHERE parent_id IS NOT NULL;
CREATE INDEX idx_empathies_message ON empathies(message_id);
CREATE INDEX idx_notifications_user_unread ON notifications(user_id) WHERE is_read = FALSE;
CREATE INDEX idx_channel_members_user ON channel_members(user_id);
CREATE INDEX idx_diaries_user_created ON diaries(user_id, created_at DESC);
CREATE INDEX idx_diaries_public_created ON diaries(created_at DESC) WHERE is_public = TRUE;
CREATE INDEX idx_diary_comments_diary ON diary_comments(diary_id, created_at) WHERE deleted_at IS NULL;
```

---

## 五、API 接口设计

### 5.1 RESTful API

#### 认证

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/auth/register` | 注册 |
| POST | `/api/auth/login` | 登录 |
| POST | `/api/auth/refresh` | 刷新 JWT（refresh token 轮换） |
| POST | `/api/auth/logout` | 登出（撤销 refresh token / 加入黑名单） |
| GET | `/api/auth/me` | 获取当前用户信息 |

#### 用户

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/users/:id` | 用户主页 |
| PUT | `/api/users/:id` | 更新资料 |
| GET | `/api/users/:id/diaries` | 用户打工日记列表 |
| GET | `/api/users/:id/stats` | 用户统计（共情数、打卡天数等） |

#### 频道

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/channels` | 频道列表（按分类） |
| POST | `/api/channels` | 申请创建频道（status=pending 待审核） |
| GET | `/api/channels/:id` | 频道详情 |
| POST | `/api/channels/:id/join` | 加入频道 |
| POST | `/api/channels/:id/leave` | 离开频道 |
| GET | `/api/channels/:id/members` | 频道成员列表 |
| GET | `/api/channels/:id/emotion-board` | 频道情绪看板 |

#### 消息

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/channels/:id/messages` | 频道消息列表（分页） |
| GET | `/api/messages/:id/threads` | 线程回复列表 |
| POST | `/api/channels/:id/messages` | 发布消息 |
| POST | `/api/messages/:id/replies` | 线程回复 |
| DELETE | `/api/messages/:id` | 删除消息 |

#### 共情

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/messages/:id/empathy` | 抱团取暖（共情） |
| DELETE | `/api/messages/:id/empathy` | 取消共情 |
| GET | `/api/messages/:id/empathy-users` | 共情用户列表 |

#### 打工日记（v1.5）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/diaries` | 日记广场（公开流） |
| POST | `/api/diaries` | 写日记/打卡 |
| GET | `/api/diaries/:id` | 日记详情 |
| GET | `/api/diaries/:id/comments` | 日记评论列表 |
| POST | `/api/diaries/:id/comments` | 发表日记评论 |
| GET | `/api/diaries/streak/:user_id` | 打卡日历（连续 + 累计天数） |

#### 排行榜

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/ranking/empathy` | 最受共情帖子榜 |
| GET | `/api/ranking/warmest` | 最暖牛马榜（按给出共情数排序） |
| GET | `/api/ranking/streak` | 连续打卡牛马榜（v1.5） |
| GET | `/api/ranking/active` | 本周最活跃牛马 |

#### 通知

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/notifications` | 通知列表 |
| PUT | `/api/notifications/:id/read` | 标记已读 |
| PUT | `/api/notifications/read-all` | 全部已读 |

#### 搜索

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/search?q=xxx&type=message` | 搜索消息 |
| GET | `/api/search?q=xxx&type=diary` | 搜索日记（v1.5） |
| GET | `/api/search?q=xxx&type=channel` | 搜索频道 |
| GET | `/api/search?q=xxx&type=user` | 搜索用户 |

### 5.3 统一响应格式与错误码

**成功响应：**
```json
{ "code": 0, "message": "success", "data": { ... } }
```

**分页响应：**
```json
{ "code": 0, "message": "success", "data": { "list": [], "total": 100, "page": 1, "page_size": 20 } }
```

**失败响应：** `code` 为业务错误码（非 0），HTTP 状态码与之对应。

| HTTP | code | 含义 |
|------|------|------|
| 400 | 40000 | 请求参数错误 |
| 401 | 40100 | 未认证 / token 失效 |
| 403 | 40300 | 无权限（如非频道成员发消息） |
| 404 | 40400 | 资源不存在 |
| 409 | 40900 | 冲突（如邮箱已注册、重复共情） |
| 422 | 42200 | 校验失败（如密码强度、字段超长） |
| 429 | 42900 | 触发限流 |
| 500 | 50000 | 服务端内部错误 |

> WebSocket 错误码使用 4xxx 段（见 3.4 节 `error` 事件），如 4001 未鉴权、4003 越权、4029 限流。

### 5.2 WebSocket 端点

```
WS /api/ws
```

> 鉴权走首帧 `auth` 消息（见 3.4 节），不放 URL query。

事件类型见上方 3.4 节。

---

## 六、项目目录结构

```
Alike/
├── CLAUDE.md
├── README.md
├── docs/
│   └── plans/
│       └── architecture-design.md      ← 本文档
│   └── architecture-diagram.html       ← 架构图（待生成）
├── docker-compose.yml
├── Makefile
│
├── backend/                             ← Go 后端
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── auth/
│   │   ├── user/
│   │   ├── channel/
│   │   ├── message/
│   │   ├── emotion/
│   │   ├── empathy/
│   │   ├── diary/
│   │   ├── ws/                          ← WebSocket Hub
│   │   ├── notification/
│   │   ├── search/
│   │   ├── storage/
│   │   └── middleware/
│   ├── pkg/
│   │   ├── database/                    ← PostgreSQL 连接
│   │   ├── redis/                       ← Redis 连接
│   │   ├── jwt/                         ← JWT 工具
│   │   └── config/                      ← 配置加载
│   ├── migrations/                      ← 数据库迁移
│   │   └── 001_init.sql
│   ├── go.mod
│   └── go.sum
│
├── frontend/                            ← Vue 3 / Nuxt 3
│   ├── nuxt.config.ts
│   ├── package.json
│   ├── pages/
│   │   ├── index.vue
│   │   ├── channel/[id].vue
│   │   ├── diary/index.vue
│   │   ├── diary/[id].vue
│   │   ├── profile/[id].vue
│   │   └── ranking.vue
│   ├── components/
│   │   ├── chat/
│   │   │   ├── MessageList.vue
│   │   │   ├── MessageInput.vue
│   │   │   └── ThreadPanel.vue
│   │   ├── emotion/
│   │   │   ├── EmotionPicker.vue
│   │   │   └── EmotionBoard.vue
│   │   ├── empathy/
│   │   │   └── EmpathyButton.vue
│   │   └── layout/
│   │       ├── ChannelSidebar.vue
│   │       └── TopNav.vue
│   ├── composables/
│   │   ├── useWebSocket.ts
│   │   ├── useAuth.ts
│   │   └── useApi.ts
│   ├── stores/
│   │   ├── auth.ts
│   │   ├── channel.ts
│   │   └── message.ts
│   └── assets/
│       └── css/
│
└── nginx/
    └── nginx.conf
```

---

## 七、分阶段实施计划

### 阶段一：项目脚手架（预计 2-3 小时）

| 任务 | 说明 |
|------|------|
| 1.1 初始化 Go 后端 | `go mod init`, Gin 框架, 目录结构, 配置加载 |
| 1.2 初始化 Nuxt 前端 | `npx nuxi init`, 目录结构, Tailwind CSS |
| 1.3 Docker Compose | PostgreSQL + Redis + MinIO + Nginx 开发环境 |
| 1.4 数据库迁移 | 创建所有表 + 索引 |
| 1.5 Makefile | 统一的开发命令（dev/build/test/migrate） |

### 阶段二：用户认证系统（预计 2-3 小时）

| 任务 | 说明 |
|------|------|
| 2.1 注册接口 | 邮箱注册 + bcrypt 密码 |
| 2.2 登录接口 | JWT 签发 + Redis 会话 |
| 2.3 JWT 中间件 | 请求鉴权 |
| 2.4 前端登录/注册页 | 表单 + Token 存储 |
| 2.5 前端认证守卫 | 路由守卫 + 自动刷新 Token |

### 阶段三：频道系统（预计 2-3 小时）

| 任务 | 说明 |
|------|------|
| 3.1 频道 CRUD | 创建/列表/详情 |
| 3.2 频道成员 | 加入/离开/成员列表 |
| 3.3 预设频道种子数据 | 行业/岗位/主题频道初始化 |
| 3.4 前端频道侧边栏 | 频道列表 + 分类筛选 |
| 3.5 前端频道页 | 频道信息 + 成员数 |

### 阶段四：消息与线程（预计 3-4 小时）

| 任务 | 说明 |
|------|------|
| 4.1 消息发布 API | REST 发布 + 写入 DB |
| 4.2 消息列表 API | 分页查询 + 线程回复数 |
| 4.3 线程回复 API | 线程列表 + 回复 |
| 4.4 WebSocket Hub | 连接管理 + 频道广播 |
| 4.5 实时消息推送 | Redis Pub/Sub → WebSocket |
| 4.6 前端消息列表 | 消息流 + 分页加载 |
| 4.7 前端消息输入 | 输入框 + 表情 + 情绪标签 |
| 4.8 前端线程面板 | 线程回复 UI |

### 阶段五：情绪系统（预计 2 小时）

| 任务 | 说明 |
|------|------|
| 5.1 情绪标签 | 消息附带情绪 + 统计 |
| 5.2 情绪看板 API | 频道内情绪热力图数据 |
| 5.3 前端情绪选择器 | 发消息时选择情绪 |
| 5.4 前端情绪看板 | 实时情绪可视化 |

### 阶段六：共情机制（预计 2 小时）

| 任务 | 说明 |
|------|------|
| 6.1 共情 API | 抱团取暖 + 取消 + 统计 |
| 6.2 共情实时推送 | WebSocket 通知被共情者 |
| 6.3 前端共情按钮 | 动画效果 + 数字更新 |

### 阶段七：打工日记（预计 2 小时）

| 任务 | 说明 |
|------|------|
| 7.1 日记 CRUD | 发布/列表/详情 |
| 7.2 打卡日历 | 连续打卡 + 徽章 |
| 7.3 前端日记页 | 日记广场 + 写日记 + 打卡日历 |

### 阶段八：排行榜与通知（预计 2 小时）

| 任务 | 说明 |
|------|------|
| 8.1 排行榜 API | 最受共情/连续打卡/最活跃 |
| 8.2 通知系统 | @提及/共情/回复通知 |
| 8.3 前端排行榜页 | 排行榜 UI |
| 8.4 前端通知面板 | 通知列表 + 已读 |

### 阶段九：打磨与部署（预计 2-3 小时）

| 任务 | 说明 |
|------|------|
| 9.1 全文搜索 | PostgreSQL 全文索引 + 搜索 API |
| 9.2 文件上传 | MinIO 集成 + 图片上传 |
| 9.3 匿名模式 | 匿名发帖/回复 |
| 9.4 前端 UI 打磨 | 响应式 + 暗色主题 + 动画 |
| 9.5 Docker 部署 | 生产 Dockerfile + Compose |
| 9.6 Nginx 配置 | 反向代理 + WebSocket 代理 |

---

## 八、风险与注意事项

| 风险 | 应对 |
|------|------|
| WebSocket 多实例广播 | 必须用 Redis Pub/Sub，不能仅靠内存；发送方实例本地已推送连接跳过再投，防同实例双投 |
| WebSocket 鉴权泄露 | JWT 不放 URL query（进 Nginx 日志），改首帧 `auth` 或子协议头；发消息校验频道成员身份防越权 |
| 消息量大时查询慢 | 主消息流用 `(channel_id, created_at DESC) WHERE parent_id IS NULL AND deleted_at IS NULL` 部分索引，游标分页不用深 OFFSET |
| 匿名模式隐私泄露 | 匿名消息接口不返回 user_id；`messages.user_id` 仍存真实 ID 仅用于举报审计，需数据库层访问控制 + 审计日志，避免直接查询穿帮 |
| 冗余计数漂移 | `empathy_received`/`empathy_given`/`member_count`/`empathy_count` 等冗余字段用原子 `UPDATE ... SET x = x ± 1` 或触发器维护，并发下保持一致 |
| 文件上传安全 | 限制类型（白名单）和大小（图片 ≤5MB、文档 ≤10MB）；服务端 MIME 嗅探校验；SVG 走下载而非内联渲染，防 XSS；图片压缩 |
| 情绪看板实时性 | WebSocket 推送情绪更新，Redis 缓存聚合数据 |
| 限流 | 按 IP + 用户 + 接口三维度限流；WebSocket 发消息独立限流，防刷屏 |

---

## 九、今晚执行清单

| # | 任务 | 产出 |
|---|------|------|
| 1 | Claude Code `/init` | `CLAUDE.md` |
| 2 | 撰写本文档 | `docs/plans/architecture-design.md` |
| 3 | 生成架构设计图 | `docs/architecture-diagram.html` |

> 后续阶段实施可在 Claude Code 中按本文档逐步执行。
