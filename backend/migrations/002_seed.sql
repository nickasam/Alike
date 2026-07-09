-- 002_seed.sql — 预置频道种子数据
-- 覆盖行业/岗位/主题三类初始频道（含 #维权互助）。
-- 重复执行安全：slug 唯一，冲突则跳过。

INSERT INTO channels (name, slug, description, category, icon, status) VALUES
    -- 行业频道 industry
    ('#互联网',   'internet',    '互联网牛马集中营，996 的日常', 'industry', 'code',      'active'),
    ('#工厂',     'factory',     '流水线上的兄弟姐妹们',         'industry', 'factory',   'active'),
    ('#外卖骑手', 'delivery',    '风里雨里，路上见',             'industry', 'bike',      'active'),
    ('#医护',     'medical',     '白衣天使也是牛马',             'industry', 'hospital',  'active'),
    ('#教培',     'education',   '教书育人，也想下班',           'industry', 'book',      'active'),
    ('#服务业',   'service',     '微笑服务背后的疲惫',           'industry', 'store',     'active'),

    -- 岗位频道 job
    ('#程序员',   'developer',   'CRUD 到秃头',                 'job',      'terminal',  'active'),
    ('#产品经理', 'product',     '需求改到天荒地老',             'job',      'clipboard', 'active'),
    ('#设计师',   'designer',    '甲方说再大一点',               'job',      'palette',   'active'),
    ('#流水线',   'assembly',    '重复的动作，重复的一天',       'job',      'gear',      'active'),

    -- 主题频道 topic
    ('#吐槽大会', 'complain',    '有苦水尽管倒',                 'topic',    'megaphone', 'active'),
    ('#打工日记', 'diary',       '记录每一个搬砖的日子',         'topic',    'notebook',  'active'),
    ('#薪资揭秘', 'salary',      '同行不同命，聊聊工资',         'topic',    'coins',     'active'),
    ('#维权互助', 'rights',      '劳动法在手，维权不孤单',       'topic',    'shield',    'active'),
    ('#摸鱼技巧', 'slacking',    '带薪如厕的艺术',               'topic',    'fish',      'active')
ON CONFLICT (slug) DO NOTHING;
