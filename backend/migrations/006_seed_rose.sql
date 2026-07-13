-- 006_seed_rose.sql — 「玫瑰少年」频道种子数据（用户 + 成员 + 消息 + 线程 + 共情）
-- 目的：让 #玫瑰少年 频道有真实用户与氛围内容，不再是空频道。
-- 全程幂等：ON CONFLICT DO NOTHING + 冗余计数用子查询重算，可安全重复执行。
-- 所有种子用户密码统一为明文 alike1234（bcrypt hash 见下）。

-- 1) 种子用户（email 唯一，幂等）
INSERT INTO users (email, password, nickname, bio, industry, job_title, work_years, level) VALUES
    ('rose_lin@seed.alike',  '$2a$10$ezO.DovGm0zl2zHF7YhAzecKp292CnIL.IYV4Mx8WxP2bxGJH6nEy', '阿林',   '刚入职半年，还在学怎么在工位上隐身', '互联网', '前端', 1, 2),
    ('rose_yu@seed.alike',   '$2a$10$ezO.DovGm0zl2zHF7YhAzecKp292CnIL.IYV4Mx8WxP2bxGJH6nEy', '小雨',   '出柜三年，家里还没通气', '教培', '老师', 4, 3),
    ('rose_qi@seed.alike',   '$2a$10$ezO.DovGm0zl2zHF7YhAzecKp292CnIL.IYV4Mx8WxP2bxGJH6nEy', '柒',     '设计狗一枚，甲方虐我千百遍', '互联网', '设计师', 3, 2),
    ('rose_an@seed.alike',   '$2a$10$ezO.DovGm0zl2zHF7YhAzecKp292CnIL.IYV4Mx8WxP2bxGJH6nEy', '阿安',   '医院轮班到怀疑人生', '医护', '护士', 5, 4),
    ('rose_mo@seed.alike',   '$2a$10$ezO.DovGm0zl2zHF7YhAzecKp292CnIL.IYV4Mx8WxP2bxGJH6nEy', '默默',   '习惯了做透明人', '服务业', '店员', 2, 1),
    ('rose_he@seed.alike',   '$2a$10$ezO.DovGm0zl2zHF7YhAzecKp292CnIL.IYV4Mx8WxP2bxGJH6nEy', '河河',   '想找到公司里的同温层', '互联网', '产品', 6, 5)
ON CONFLICT (email) DO NOTHING;

-- 2) 加入 rose 频道（channel_members 唯一约束幂等）
INSERT INTO channel_members (channel_id, user_id, role)
SELECT c.id, u.id, CASE WHEN u.email = 'rose_he@seed.alike' THEN 'admin' ELSE 'member' END
FROM channels c, users u
WHERE c.slug = 'rose'
  AND u.email IN ('rose_lin@seed.alike','rose_yu@seed.alike','rose_qi@seed.alike',
                  'rose_an@seed.alike','rose_mo@seed.alike','rose_he@seed.alike')
ON CONFLICT (channel_id, user_id) DO NOTHING;

-- 3) 主消息（用 content 幂等：同频道同内容不重复插）
--    情绪标签取自 8 种合法值：tired/angry/wronged/break/numb/quit/anxious/cheer
INSERT INTO messages (channel_id, user_id, content, emotion, is_anonymous)
SELECT c.id, u.id, v.content, v.emotion, v.anon
FROM channels c
JOIN (VALUES
    ('rose_he@seed.alike', '开个头吧——这里是给公司里那些不太好说出口的心事留的角落，匿名随意，来聊。', 'cheer',   FALSE),
    ('rose_lin@seed.alike','团建又要玩「介绍对象」的梗，我全程尬笑，散场后一个人走回地铁站。',        'wronged', FALSE),
    ('rose_yu@seed.alike', '同事随口一句「你怎么还不结婚」，我笑着岔开，回工位缓了十分钟。',            'tired',   TRUE),
    ('rose_qi@seed.alike', '发现工位对面的同事也在关注同一个话题，那一刻突然觉得没那么孤单了。',        'cheer',   FALSE),
    ('rose_an@seed.alike', '轮到夜班，更衣室里没人，反而是我一天里最放松的时候。',                    'numb',    FALSE),
    ('rose_mo@seed.alike', '领导知道后态度就变了，绩效沟通时那种客气的疏远，比骂我还难受。',            'break',   TRUE),
    ('rose_lin@seed.alike','其实不奢求被理解，只希望别把我当异类看待。',                              'anxious', FALSE)
) AS v(email, content, emotion, anon) ON TRUE
JOIN users u ON u.email = v.email
WHERE c.slug = 'rose'
  AND NOT EXISTS (
    SELECT 1 FROM messages m WHERE m.channel_id = c.id AND m.content = v.content
  );

-- 4) 线程回复（parent 用「频道+父消息内容」定位，回复用自身 content 幂等）
INSERT INTO messages (channel_id, user_id, parent_id, content, emotion, is_anonymous)
SELECT c.id, u.id, p.id, v.content, v.emotion, v.anon
FROM channels c
JOIN (VALUES
    ('rose_qi@seed.alike', '团建又要玩「介绍对象」的梗，我全程尬笑，散场后一个人走回地铁站。', '懂，那种氛围里最累的就是要不停演。抱抱。',   'wronged', FALSE),
    ('rose_he@seed.alike', '团建又要玩「介绍对象」的梗，我全程尬笑，散场后一个人走回地铁站。', '下次那种局，能溜就溜，别为难自己。',         'cheer',   FALSE),
    ('rose_an@seed.alike', '领导知道后态度就变了，绩效沟通时那种客气的疏远，比骂我还难受。',     '这种冷暴力最伤人。你的能力不该被这个定义。', 'angry',   TRUE)
) AS v(replier_email, parent_content, content, emotion, anon) ON TRUE
JOIN users u ON u.email = v.replier_email
JOIN messages p ON p.channel_id = c.id AND p.content = v.parent_content AND p.parent_id IS NULL
WHERE c.slug = 'rose'
  AND NOT EXISTS (
    SELECT 1 FROM messages m WHERE m.channel_id = c.id AND m.content = v.content
  );

-- 5) 共情（每人每条消息一次，唯一约束幂等）
--    让几条走心的主消息拿到共情，营造「被看见」的氛围。
INSERT INTO empathies (message_id, user_id)
SELECT m.id, u.id
FROM messages m
JOIN channels c ON c.id = m.channel_id AND c.slug = 'rose'
JOIN (VALUES
    ('团建又要玩「介绍对象」的梗，我全程尬笑，散场后一个人走回地铁站。', 'rose_yu@seed.alike'),
    ('团建又要玩「介绍对象」的梗，我全程尬笑，散场后一个人走回地铁站。', 'rose_qi@seed.alike'),
    ('团建又要玩「介绍对象」的梗，我全程尬笑，散场后一个人走回地铁站。', 'rose_he@seed.alike'),
    ('领导知道后态度就变了，绩效沟通时那种客气的疏远，比骂我还难受。',     'rose_lin@seed.alike'),
    ('领导知道后态度就变了，绩效沟通时那种客气的疏远，比骂我还难受。',     'rose_qi@seed.alike'),
    ('领导知道后态度就变了，绩效沟通时那种客气的疏远，比骂我还难受。',     'rose_an@seed.alike'),
    ('其实不奢求被理解，只希望别把我当异类看待。',                       'rose_he@seed.alike'),
    ('其实不奢求被理解，只希望别把我当异类看待。',                       'rose_an@seed.alike')
) AS v(msg_content, user_email) ON m.content = v.msg_content
JOIN users u ON u.email = v.user_email
ON CONFLICT (message_id, user_id) DO NOTHING;

-- 6) 重算冗余计数（幂等：每次按真实关联重新汇总，不做增量累加）
-- 6a) 消息 empathy_count
UPDATE messages m
SET empathy_count = sub.cnt
FROM (
    SELECT message_id, COUNT(*) AS cnt FROM empathies GROUP BY message_id
) sub
WHERE m.id = sub.message_id;

-- 6b) 频道 member_count
UPDATE channels c
SET member_count = (SELECT COUNT(*) FROM channel_members cm WHERE cm.channel_id = c.id)
WHERE c.slug = 'rose';

-- 6c) 用户 empathy_received（被共情总数，跨全站按真实关联重算）
UPDATE users u
SET empathy_received = COALESCE(sub.cnt, 0)
FROM (
    SELECT m.user_id, COUNT(*) AS cnt
    FROM empathies e JOIN messages m ON m.id = e.message_id
    GROUP BY m.user_id
) sub
WHERE u.id = sub.user_id
  AND u.email LIKE '%@seed.alike';
