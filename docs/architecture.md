# Alike Project Structure

## Directory Layout

```
Alike/
├── cmd/                    # 应用入口
│   ├── api/               # API服务器
│   │   └── main.go
│   ├── worker/            # 后台任务处理器
│   │   └── main.go
│   └── migrate/           # 数据库迁移工具
│       └── main.go
│
├── internal/              # 私有应用代码
│   ├── api/              # HTTP handlers
│   │   ├── handler/      # 请求处理器
│   │   │   ├── auth.go
│   │   │   ├── chat.go
│   │   │   ├── global_chat.go
│   │   │   ├── match.go
│   │   │   ├── notification.go
│   │   │   └── user.go
│   │   ├── middleware/   # 中间件
│   │   │   ├── auth.go
│   │   │   └── common.go
│   │   └── router/       # 路由配置
│   │       └── router.go
│   │
│   ├── auth/             # 认证服务
│   │   └── service.go
│   │
│   ├── chat/             # 聊天服务
│   │   └── service.go
│   │
│   ├── match/            # 匹配服务
│   │   └── service.go
│   │
│   ├── user/             # 用户服务
│   │   ├── repository.go
│   │   └── service.go
│   │
│   ├── notification/     # 推送通知
│   │   └── service.go
│   │
│   ├── upload/           # 文件上传
│   │   └── (待实现)
│   │
│   ├── pkg/              # 内部共享包
│   │   ├── errors/       # 错误定义
│   │   ├── logger/       # 日志
│   │   ├── response/     # API响应
│   │   │   └── response.go
│   │   └── validator/    # 验证器
│   │
│   ├── repository/       # 数据访问层
│   │   ├── global_chat.go
│   │   └── user.go
│   │
│   └── domain/           # 领域模型
│       ├── chat.go
│       ├── common.go
│       ├── global_chat.go
│       ├── match.go
│       ├── notification.go
│       └── user.go
│
├── pkg/                  # 公共库（可被外部导入）
│   ├── database/         # 数据库连接
│   │   ├── postgres.go
│   │   ├── redis.go
│   │   └── migrations/
│   │
│   ├── config/           # 配置管理
│   │   ├── config.go
│   │   └── config.yaml
│   │
│   └── utils/            # 工具函数
│       ├── crypto.go
│       ├── geo.go
│       └── time.go
│
├── web/                  # Web前端（可选）
│   ├── admin/            # 管理后台
│   └── landing/          # 落地页
│
├── scripts/              # 脚本
│   ├── deploy.sh
│   ├── migrate.sh
│   └── seed.sh
│
├── deployments/          # 部署配置
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── k8s/              # Kubernetes配置
│
├── docs/                 # 文档
│   ├── api/              # API文档
│   ├── architecture.md   # 架构设计
│   └── database.md       # 数据库设计
│
├── .github/              # GitHub配置
│   └── workflows/        # CI/CD
│
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Tech Stack

### 后端
- **语言**: Go 1.21+
- **框架**: Gin ✅
- **数据库**: PostgreSQL 14+ ✅
- **缓存**: Redis 7+ (计划中)
- **消息队列**: Redis Streams (计划中)
- **WebSocket**: gorilla/websocket (计划中)
- **验证**: go-playground/validator ✅
- **配置**: Viper ✅
- **日志**: Zap (计划中)

### 前端
- **技术栈**: 纯 HTML/CSS/JavaScript ✅
- **UI 设计**: Material Icons + 自定义设计系统 ✅
- **响应式**: 桌面端/移动端适配 ✅
- **PWA 支持**: 计划中

### 基础设施
- **容器**: Docker ✅
- **编排**: Docker Compose ✅
- **反向代理**: Nginx ✅
- **监控**: Prometheus + Grafana (计划中)
- **CI/CD**: GitHub Actions (计划中)

### 存储 (计划中)
- **对象存储**: 阿里云OSS / AWS S3
- **CDN**: 阿里云CDN / Cloudflare

### 推送 (计划中)
- **iOS**: APNs
- **Android**: Firebase Cloud Messaging

## 实现状态

| 模块 | 状态 | 完成度 | 说明 |
|------|------|--------|------|
| 认证系统 | ✅ 已实现 | 80% | 注册、登录、JWT 完成，缺少 SMS 验证 |
| 用户系统 | ✅ 已实现 | 70% | CRUD 完成，缺少头像上传 |
| 匹配系统 | 🚧 开发中 | 40% | 基础匹配完成，缺少推荐算法 |
| 聊天系统 | 🚧 开发中 | 60% | 全局聊天完成，私聊部分完成 |
| 通知系统 | ✅ 已实现 | 50% | 基础通知完成，缺少推送 |
| 数据库 | ✅ 已实现 | 90% | 表结构完成，需要优化索引 |
| API 文档 | 📋 计划中 | 0% | 待生成 Swagger/OpenAPI |
| 前端页面 | 🚧 开发中 | 40% | 基础页面完成，缺少交互 |
| WebSocket | 📋 计划中 | 0% | 待实现 |
| 文件上传 | 📋 计划中 | 0% | 待实现 |

## API Design

### 健康检查
```
GET    /health                     # ✅ 已实现
```

### 认证相关 (✅ 已实现)
```
POST   /api/v1/auth/register       # 用户注册
POST   /api/v1/auth/login          # 用户登录
POST   /api/v1/auth/refresh        # 刷新Token
POST   /api/v1/auth/logout         # 用户登出
GET    /api/v1/auth/me             # 获取当前用户信息（需认证）
```

### 用户相关 (✅ 部分实现)
```
GET    /api/v1/users/me            # 获取当前用户信息 ✅
PUT    /api/v1/users/me            # 更新当前用户信息 ✅
GET    /api/v1/users/nearby        # 获取附近用户 ✅
POST   /api/v1/users/:id/view      # 记录用户访问 📋
POST   /api/v1/users/:id/like      # 点赞用户 📋
DELETE /api/v1/users/:id/like      # 取消点赞 📋
POST   /api/v1/users/:id/block     # 屏蔽用户 📋
DELETE /api/v1/users/:id/block     # 取消屏蔽 📋
GET    /api/v1/users/:id           # 获取用户详情 📋
POST   /api/v1/users/me/avatar     # 上传头像 📋
GET    /api/v1/users/recommend     # 推荐用户 📋
GET    /api/v1/users/blocked       # 获取屏蔽列表 📋
```

### 匹配相关 (✅ 部分实现)
```
GET    /api/v1/matches             # 获取匹配列表 ✅
GET    /api/v1/matches/:id         # 获取匹配详情 ✅
POST   /api/v1/matches/:id/like    # 点赞用户 ✅
POST   /api/v1/matches             # 创建匹配 📋
DELETE /api/v1/matches/:id         # 删除匹配 📋
POST   /api/v1/matches/:id/read    # 标记为已读 📋
```

### 聊天相关 (✅ 部分实现)
```
GET    /api/v1/chats               # 获取聊天列表 ✅
GET    /api/v1/chats/:id           # 获取聊天详情 ✅
GET    /api/v1/chats/:id/messages  # 获取消息列表 ✅
POST   /api/v1/chats/:id/messages  # 发送消息 ✅
POST   /api/v1/chats               # 创建聊天 📋
DELETE /api/v1/chats/:id           # 删除聊天 📋
POST   /api/v1/chats/:id/messages/:id/read  # 标记消息已读 📋
DELETE /api/v1/chats/:id/messages/:id       # 删除消息 📋
WS     /api/v1/chats/:id/ws        # WebSocket 连接 📋
```

### 全局聊天室 (✅ 已实现)
```
GET    /api/v1/global/room         # 获取全局聊天室信息 ✅
GET    /api/v1/global/messages     # 获取全局消息 ✅
POST   /api/v1/global/messages     # 发送全局消息 ✅
POST   /api/v1/global/join         # 加入聊天室 ✅
```

### 通知相关 (✅ 已实现)
```
GET    /api/v1/notifications       # 获取通知列表 ✅
POST   /api/v1/notifications/:id/read    # 标记为已读 ✅
POST   /api/v1/notifications/read-all    # 全部标记为已读 ✅
PUT    /api/v1/notifications/settings    # 更新通知设置 📋
GET    /api/v1/notifications/settings    # 获取通知设置 📋
```

### 上传相关 (📋 计划中)
```
POST   /api/v1/upload/image        # 上传图片
POST   /api/v1/upload/avatar       # 上传头像
DELETE /api/v1/upload/:id          # 删除文件
```

**图例说明**：
- ✅ 已实现：功能已完成并可用
- 📋 计划中：功能已规划，待实现

## Database Schema

### 核心表

#### users (用户表)
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone VARCHAR(20) UNIQUE NOT NULL,
    nickname VARCHAR(50) NOT NULL,
    avatar_url VARCHAR(500),
    birth_date DATE,
    height SMALLINT,
    weight SMALLINT,
    role VARCHAR(10), -- '1', '0', '0.5'
    bio TEXT,
    location_lat DECIMAL(10, 8),
    location_lng DECIMAL(11, 8),
    location_name VARCHAR(200),
    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    last_online_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_location ON users(location_lat, location_lng);
CREATE INDEX idx_users_last_online ON users(last_online_at DESC);
```

#### user_tags (用户标签)
```sql
CREATE TABLE user_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    tag VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_tags_user_id ON user_tags(user_id);
CREATE INDEX idx_user_tags_tag ON user_tags(tag);
```

#### user_images (用户图片)
```sql
CREATE TABLE user_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    image_url VARCHAR(500) NOT NULL,
    order_index SMALLINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_images_user_id ON user_images(user_id);
```

#### likes (喜欢记录)
```sql
CREATE TABLE likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    liker_id UUID REFERENCES users(id) ON DELETE CASCADE,
    liked_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(liker_id, liked_id)
);

CREATE INDEX idx_likes_liker_id ON likes(liker_id);
CREATE INDEX idx_likes_liked_id ON likes(liked_id);
CREATE INDEX idx_likes_created_at ON likes(created_at DESC);
```

#### matches (匹配记录)
```sql
CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user1_id UUID REFERENCES users(id) ON DELETE CASCADE,
    user2_id UUID REFERENCES users(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT TRUE,
    last_message_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user1_id, user2_id)
);

CREATE INDEX idx_matches_user1_id ON matches(user1_id);
CREATE INDEX idx_matches_user2_id ON matches(user2_id);
CREATE INDEX idx_matches_last_message ON matches(last_message_at DESC);
```

#### chats (聊天会话)
```sql
CREATE TABLE chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID REFERENCES matches(id) ON DELETE CASCADE,
    user1_id UUID REFERENCES users(id) ON DELETE CASCADE,
    user2_id UUID REFERENCES users(id) ON DELETE CASCADE,
    last_message_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_chats_match_id ON chats(match_id);
CREATE INDEX idx_chats_user1_id ON chats(user1_id);
CREATE INDEX idx_chats_user2_id ON chats(user2_id);
```

#### messages (消息)
```sql
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id UUID REFERENCES chats(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id) ON DELETE CASCADE,
    content TEXT,
    message_type VARCHAR(20) DEFAULT 'text', -- 'text', 'image', 'location', 'voice'
    metadata JSONB,
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_messages_chat_id ON messages(chat_id, created_at DESC);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);
CREATE INDEX idx_messages_is_read ON messages(is_read);
```

#### blocks (屏蔽记录)
```sql
CREATE TABLE blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blocker_id UUID REFERENCES users(id) ON DELETE CASCADE,
    blocked_id UUID REFERENCES users(id) ON DELETE CASCADE,
    reason VARCHAR(200),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(blocker_id, blocked_id)
);

CREATE INDEX idx_blocks_blocker_id ON blocks(blocker_id);
CREATE INDEX idx_blocks_blocked_id ON blocks(blocked_id);
```

#### views (访问记录)
```sql
CREATE TABLE views (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    viewer_id UUID REFERENCES users(id) ON DELETE CASCADE,
    viewed_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_views_viewer_id ON views(viewer_id);
CREATE INDEX idx_views_viewed_id ON views(viewed_id);
CREATE INDEX idx_views_created_at ON views(created_at DESC);
```

#### notifications (通知)
```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- 'like', 'match', 'message', 'view'
    title VARCHAR(200),
    content TEXT,
    data JSONB,
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id, created_at DESC);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
```

#### notification_settings (通知设置)
```sql
CREATE TABLE notification_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    like_notification BOOLEAN DEFAULT TRUE,
    match_notification BOOLEAN DEFAULT TRUE,
    message_notification BOOLEAN DEFAULT TRUE,
    view_notification BOOLEAN DEFAULT FALSE,
    email_notification BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMP DEFAULT NOW()
);
```

#### refresh_tokens (刷新令牌)
```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(500) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
```

## 环境变量

```env
# 服务器配置
SERVER_PORT=8080
SERVER_MODE=development # development, production

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_NAME=alike_db
DB_USER=alike_user
DB_PASSWORD=your_password
DB_SSL_MODE=disable

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT配置
JWT_SECRET=your_jwt_secret
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=7d

# 文件上传
UPLOAD_MAX_SIZE=10485760 # 10MB
ALLOWED_IMAGE_TYPES=image/jpeg,image/png,image/gif

# OSS配置
OSS_ENDPOINT=
OSS_ACCESS_KEY_ID=
OSS_ACCESS_KEY_SECRET=
OSS_BUCKET_NAME=
OSS_REGION=

# 推送配置
FCM_SERVER_KEY=
APNS_KEY_ID=
APNS_TEAM_ID=
APNS_BUNDLE_ID=

# 短信配置
SMS_PROVIDER=
SMS_ACCESS_KEY=
SMS_SECRET_KEY=
SMS_SIGN_NAME=

# 日志配置
LOG_LEVEL=info
LOG_OUTPUT=stdout

# 地图配置
MAP_API_KEY=
```

## 下一步计划

1. ✅ UI/UX设计完成
2. ✅ 项目结构规划完成
3. ⏳ 创建项目骨架
4. ⏳ 实现数据库模型和迁移
5. ⏳ 实现认证模块
6. ⏳ 实现用户模块
7. ⏳ 实现匹配算法
8. ⏳ 实现聊天功能（WebSocket）
9. ⏳ 实现推送通知
10. ⏳ 前端开发（iOS/Android）

---

*创建时间: 2026-03-02*
*状态: 规划完成，准备开发*
