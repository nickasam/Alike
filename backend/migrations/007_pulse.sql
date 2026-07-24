-- 007_pulse.sql — 最近发生（Pulse / 世界脉搏）
-- 首发聚焦：只搬运、只读展示，纯外部数据镜像。
-- 见 docs/plans/pulse-module-design.md 第 3.2 章。

-- 专题：一个专题 = 一个外部数据源的封装（collector 注册键 + 可调 JSON 配置）。
CREATE TABLE IF NOT EXISTS pulse_topics (
    id                    BIGSERIAL PRIMARY KEY,
    slug                  VARCHAR(32) UNIQUE NOT NULL,   -- 'github-trending' / 'ai-news'
    name                  VARCHAR(64) NOT NULL,          -- 'GitHub 每日 Star 榜'
    emoji                 VARCHAR(8),                    -- '🌟' / '🤖'
    description           TEXT,                          -- 卡片下方灰字副标题
    collector_kind        VARCHAR(32) NOT NULL,          -- 'github_trending' / 'hackernews_ai'
    collector_config      JSONB,                         -- collector 可调参数
    sort_order            INT DEFAULT 0,                 -- 前端 Tab 顺序
    is_active             BOOLEAN DEFAULT TRUE,          -- 关闭抓取但保留数据
    refresh_interval_min  INT DEFAULT 60,                -- 抓取间隔（分钟）
    last_fetched_at       TIMESTAMPTZ,                   -- 最近一次成功抓取时间
    last_error            TEXT,                          -- 最近一次抓取错误（成功则清空）
    created_at            TIMESTAMPTZ DEFAULT NOW(),
    updated_at            TIMESTAMPTZ DEFAULT NOW()
);

-- 热点条目：一条 = 一个外部事件（GitHub 仓库当日快照 / HN 帖子）。
-- 不做软删除：外部数据镜像不需要"删除保护"，7 天窗口后 scheduler 直接 DELETE。
CREATE TABLE IF NOT EXISTS pulse_items (
    id             BIGSERIAL PRIMARY KEY,
    topic_id       BIGINT NOT NULL REFERENCES pulse_topics(id) ON DELETE CASCADE,
    source         VARCHAR(32) NOT NULL,               -- 'github' / 'hackernews'
    source_id      VARCHAR(255) NOT NULL,              -- 去重键（owner/repo / HN story_id）
    title          TEXT NOT NULL,
    summary        TEXT,
    url            TEXT NOT NULL,
    author         VARCHAR(128),
    score          INT DEFAULT 0,                      -- 主分：today_stars / hn_points
    score_delta    INT DEFAULT 0,                      -- 相对上次抓取的增量
    extra          JSONB,                              -- 源特有字段：{"lang":"Go","total_stars":12345,"forks":234}
    published_at   TIMESTAMPTZ,                        -- 原始发布/首次出现时间
    captured_at    TIMESTAMPTZ DEFAULT NOW(),          -- 本次抓取时间
    UNIQUE(topic_id, source, source_id)
);

-- 索引
-- 主查询：某专题 top-N（按 score 降序）
CREATE INDEX IF NOT EXISTS idx_pulse_items_topic_score
    ON pulse_items(topic_id, score DESC, captured_at DESC);
-- 清理任务：按抓取时间清 7 天前的老数据
CREATE INDEX IF NOT EXISTS idx_pulse_items_topic_captured
    ON pulse_items(topic_id, captured_at DESC);

-- 种子专题（M0）：GitHub 每日 Star 榜 + AI 圈大事
INSERT INTO pulse_topics (slug, name, emoji, description, collector_kind, collector_config, sort_order, refresh_interval_min) VALUES
  ('github-trending', 'GitHub 每日 Star 榜', '🌟', '今天全世界程序员都在给谁星星', 'github_trending',
   '{"since":"daily","limit":25}'::jsonb, 1, 60),
  ('ai-news',         'AI 圈大事',         '🤖', 'HN 上关于 AI 的最新高赞讨论',    'hackernews_ai',
   '{"keywords":["openai","anthropic","claude","gpt","llm","gemini","mistral","deepseek"],"min_points":50,"window_hours":72,"limit":20}'::jsonb,
   2, 30)
ON CONFLICT (slug) DO NOTHING;
