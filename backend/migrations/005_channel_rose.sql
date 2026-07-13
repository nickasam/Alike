-- 005_channel_rose.sql — 新增「玫瑰少年」频道（性少数打工人的职场同温层）
-- 零改表，纯 seed 数据。ON CONFLICT 保证幂等，重复执行安全。

INSERT INTO channels (name, slug, description, category, icon, status) VALUES
    ('#玫瑰少年', 'rose', '每个不被理解的少年都值得被温柔以待 —— 性少数打工人的树洞：找到懂你的人', 'topic', 'rose', 'active')
ON CONFLICT (slug) DO NOTHING;
