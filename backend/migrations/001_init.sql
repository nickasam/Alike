-- 001_init.sql — Alike 初始建表与索引
-- 严格对齐 docs/plans/architecture-design.md 第四章。

-- 扩展：email 大小写不敏感唯一，避免 Foo@x.com 与 foo@x.com 重复注册
CREATE EXTENSION IF NOT EXISTS citext;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
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
    is_anonymous BOOLEAN DEFAULT FALSE, -- 全局匿名开关
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

-- 频道表
CREATE TABLE IF NOT EXISTS channels (
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
CREATE TABLE IF NOT EXISTS channel_members (
    id          BIGSERIAL PRIMARY KEY,
    channel_id  BIGINT REFERENCES channels(id) ON DELETE CASCADE,
    user_id     BIGINT REFERENCES users(id) ON DELETE CASCADE,
    role        VARCHAR(20) DEFAULT 'member', -- member/admin
    joined_at   TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(channel_id, user_id)
);

-- 消息表
CREATE TABLE IF NOT EXISTS messages (
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
    deleted_at  TIMESTAMPTZ              -- 软删除：保留线程外键完整性，删除消息时置时间戳
);

-- 共情（抱团取暖）表
CREATE TABLE IF NOT EXISTS empathies (
    id          BIGSERIAL PRIMARY KEY,
    message_id  BIGINT REFERENCES messages(id) ON DELETE CASCADE,
    user_id     BIGINT REFERENCES users(id),
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(message_id, user_id)  -- 每人每条消息只能共情一次
);

-- 打工日记表（v1.5）
CREATE TABLE IF NOT EXISTS diaries (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT REFERENCES users(id),
    title       VARCHAR(200),
    content     TEXT NOT NULL,
    mood        VARCHAR(50),             -- 今日心情
    is_public   BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMPTZ DEFAULT NOW()
    -- 连续/累计打卡天数在 users(total_check_in_days) 及应用层按 created_at 计算
);

-- 日记评论表（v1.5，支撑 2.6"日记广场可评论"）
CREATE TABLE IF NOT EXISTS diary_comments (
    id          BIGSERIAL PRIMARY KEY,
    diary_id    BIGINT REFERENCES diaries(id) ON DELETE CASCADE,
    user_id     BIGINT REFERENCES users(id),
    content     TEXT NOT NULL,
    is_anonymous BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

-- 通知表
CREATE TABLE IF NOT EXISTS notifications (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT REFERENCES users(id),  -- 接收者
    type        VARCHAR(50) NOT NULL,    -- mention/empathy/reply/system
    content     TEXT,
    ref_id      BIGINT,                  -- 关联的消息/日记 ID
    is_read     BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- 文件附件表
CREATE TABLE IF NOT EXISTS attachments (
    id          BIGSERIAL PRIMARY KEY,
    message_id  BIGINT REFERENCES messages(id) ON DELETE CASCADE,
    file_url    VARCHAR(500) NOT NULL,
    file_type   VARCHAR(50),             -- image/document/video
    file_size   BIGINT,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- ============ 索引（对齐 4.2 节）============

-- 频道主消息流：仅索引主消息（parent_id IS NULL）且未软删，避免线程回复污染列表分页
CREATE INDEX IF NOT EXISTS idx_messages_channel_created ON messages(channel_id, created_at DESC)
    WHERE parent_id IS NULL AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_messages_parent ON messages(parent_id) WHERE parent_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_empathies_message ON empathies(message_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_unread ON notifications(user_id) WHERE is_read = FALSE;
CREATE INDEX IF NOT EXISTS idx_channel_members_user ON channel_members(user_id);
CREATE INDEX IF NOT EXISTS idx_diaries_user_created ON diaries(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_diaries_public_created ON diaries(created_at DESC) WHERE is_public = TRUE;
CREATE INDEX IF NOT EXISTS idx_diary_comments_diary ON diary_comments(diary_id, created_at) WHERE deleted_at IS NULL;
