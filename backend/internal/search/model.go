package search

import "time"

// SearchType 是搜索类型枚举。
type SearchType string

const (
	TypeMessage SearchType = "message" // 搜索消息内容
	TypeDiary   SearchType = "diary"   // 搜索日记标题/内容
	TypeChannel SearchType = "channel" // 搜索频道名/描述
	TypeUser    SearchType = "user"    // 搜索用户昵称/简介
)

// Author 是消息/日记搜索结果中的作者公开信息。匿名内容不返回该字段。
type Author struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

// MessageResult 是消息搜索命中项。
type MessageResult struct {
	ID          int64     `json:"id"`
	ChannelID   int64     `json:"channel_id"`
	Content     string    `json:"content"`
	Emotion     string    `json:"emotion,omitempty"`
	IsAnonymous bool      `json:"is_anonymous"`
	Author      *Author   `json:"author,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// DiaryResult 是日记搜索命中项（仅公开日记）。
type DiaryResult struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content"`
	Mood      string    `json:"mood,omitempty"`
	Author    *Author   `json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ChannelResult 是频道搜索命中项。
type ChannelResult struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	Category    string    `json:"category"`
	Icon        string    `json:"icon,omitempty"`
	MemberCount int       `json:"member_count"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserResult 是用户搜索命中项。
type UserResult struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Bio       string `json:"bio,omitempty"`
	Industry  string `json:"industry,omitempty"`
	JobTitle  string `json:"job_title,omitempty"`
	Level     int    `json:"level"`
}
