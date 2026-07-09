-- Postgres 首次初始化脚本（挂载至 /docker-entrypoint-initdb.d）
-- 仅在数据卷为空时执行一次。启用 citext 扩展供 users.email 大小写不敏感唯一。
-- 建表 / 索引由 `make migrate` 执行 backend/migrations 完成，此处只装扩展。
CREATE EXTENSION IF NOT EXISTS citext;
