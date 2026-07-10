-- 003_perf_indexes.sql — 性能索引（架构审视 #1/#2 P1 性能欠账）
-- 针对三类热点查询补齐索引：全文搜索(ILIKE 前后通配)、排行榜排序、活跃度聚合。
-- 全部 IF NOT EXISTS，重复执行安全。

-- ============ 1. 全文搜索：pg_trgm GIN 索引 ============
-- search 用 ILIKE '%q%'（前导通配符），普通 B-tree 无法命中，需三元组 GIN 索引。
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- 消息内容搜索（仅未软删的可搜，但 GIN 建全量，查询条件里再过滤 deleted_at）
CREATE INDEX IF NOT EXISTS idx_messages_content_trgm
    ON messages USING gin (content gin_trgm_ops);

-- 日记标题/内容搜索（仅公开日记可搜）
CREATE INDEX IF NOT EXISTS idx_diaries_content_trgm
    ON diaries USING gin (content gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_diaries_title_trgm
    ON diaries USING gin (title gin_trgm_ops);

-- 频道名/slug/描述搜索
CREATE INDEX IF NOT EXISTS idx_channels_name_trgm
    ON channels USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_channels_desc_trgm
    ON channels USING gin (description gin_trgm_ops);

-- 用户昵称/简介搜索
CREATE INDEX IF NOT EXISTS idx_users_nickname_trgm
    ON users USING gin (nickname gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_users_bio_trgm
    ON users USING gin (bio gin_trgm_ops);

-- ============ 2. 排行榜：排序覆盖索引 ============
-- 最受共情帖子榜：ORDER BY empathy_count DESC，仅统计未软删且 empathy_count>0
CREATE INDEX IF NOT EXISTS idx_messages_empathy_rank
    ON messages(empathy_count DESC, created_at DESC)
    WHERE deleted_at IS NULL AND empathy_count > 0;

-- 最暖牛马榜：ORDER BY empathy_given DESC，仅统计 empathy_given>0
CREATE INDEX IF NOT EXISTS idx_users_empathy_given
    ON users(empathy_given DESC)
    WHERE empathy_given > 0;

-- 最受共情牛马（备用）：empathy_received DESC
CREATE INDEX IF NOT EXISTS idx_users_empathy_received
    ON users(empathy_received DESC)
    WHERE empathy_received > 0;

-- ============ 3. 活跃度榜：按用户+时间聚合 ============
-- RankingActive 按 messages(user_id, created_at) JOIN + 近7天过滤 + GROUP BY user_id
CREATE INDEX IF NOT EXISTS idx_messages_user_created
    ON messages(user_id, created_at)
    WHERE deleted_at IS NULL;
