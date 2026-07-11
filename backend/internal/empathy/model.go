package empathy

import "time"

// User 是共情用户的公开信息。
type User struct {
	ID        int64     `json:"id"`
	Nickname  string    `json:"nickname"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"` // 共情时间
}

// RankMessage 是"最受共情帖子榜"的条目。匿名消息不返回作者。
type RankMessage struct {
	MessageID    int64   `json:"message_id"`
	ChannelID    int64   `json:"channel_id"`
	ChannelName  string  `json:"channel_name"`
	Content      string  `json:"content"`
	Emotion      string  `json:"emotion,omitempty"`
	IsAnonymous  bool    `json:"is_anonymous"`
	EmpathyCount int64   `json:"empathy_count"`
	Author       *Author `json:"author,omitempty"`
}

// RankUser 是牛马榜（最暖 / 最活跃）的条目。
type RankUser struct {
	UserID    int64  `json:"user_id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Level     int    `json:"level"`
	Metric    int64  `json:"metric"` // 榜单排序指标：被共情数 / 给出共情数 / 本周消息数
}

// Author 是消息作者的公开信息。
type Author struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Level     int    `json:"level"`
}
