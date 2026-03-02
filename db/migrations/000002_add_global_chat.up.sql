-- 创建全局聊天室表
CREATE TABLE IF NOT EXISTS global_chat_rooms (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    max_members INTEGER DEFAULT 1000,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- 创建全局消息表
CREATE TABLE IF NOT EXISTS global_messages (
    id VARCHAR(255) PRIMARY KEY,
    room_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    username VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_global_messages_room_id ON global_messages(room_id);
CREATE INDEX IF NOT EXISTS idx_global_messages_created_at ON global_messages(created_at);
CREATE INDEX IF NOT EXISTS idx_global_messages_user_id ON global_messages(user_id);

-- 插入默认全局聊天室
INSERT INTO global_chat_rooms (id, name, description, max_members, created_at, updated_at)
VALUES ('global', 'Alike大家庭', '欢迎来到Alike大家庭！这里是所有用户的聊天空间，认识新朋友，分享生活点滴。', 1000, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
