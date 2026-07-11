-- 004_diary_empathy.sql — 日记共情（我懂你）
-- 为 diaries 增加冗余共情计数列，并新增日记共情关联表（每人每篇一次）。

ALTER TABLE diaries ADD COLUMN IF NOT EXISTS empathy_count INT DEFAULT 0;

CREATE TABLE IF NOT EXISTS diary_empathies (
    id         BIGSERIAL PRIMARY KEY,
    diary_id   BIGINT REFERENCES diaries(id) ON DELETE CASCADE,
    user_id    BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(diary_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_diary_empathies_diary ON diary_empathies(diary_id);
