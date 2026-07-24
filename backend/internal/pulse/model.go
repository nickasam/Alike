// Package pulse 负责"最近发生"：外部信号（GitHub Trending / HN AI 等）的
// 抓取、存储与只读展示。首发只做 2 个专题的展示闭环，不接入任何站内互动。
// 见 docs/plans/pulse-module-design.md。
package pulse

import (
	"encoding/json"
	"time"
)

// Topic 是 pulse_topics 表的领域模型（前端 Tab 数据源）。
// M0 阶段只需要 slug/name/emoji/description + last_fetched_at 供前端渲染专题条与"截至 HH:MM 更新"。
type Topic struct {
	ID                 int64           `json:"id"`
	Slug               string          `json:"slug"`
	Name               string          `json:"name"`
	Emoji              string          `json:"emoji,omitempty"`
	Description        string          `json:"description,omitempty"`
	CollectorKind      string          `json:"-"` // 前端不感知，仅后端调度用
	CollectorConfig    json.RawMessage `json:"-"`
	SortOrder          int             `json:"sort_order"`
	IsActive           bool            `json:"-"`
	RefreshIntervalMin int             `json:"refresh_interval_min,omitempty"`
	LastFetchedAt      *time.Time      `json:"last_fetched_at,omitempty"`
	LastError          string          `json:"-"` // 内部用，不外抛
}

// Item 是 pulse_items 表的领域模型（前端卡片数据源）。
type Item struct {
	ID          int64           `json:"id"`
	TopicID     int64           `json:"topic_id"`
	Source      string          `json:"source"` // 'github' / 'hackernews'
	SourceID    string          `json:"source_id"`
	Title       string          `json:"title"`
	Summary     string          `json:"summary,omitempty"`
	URL         string          `json:"url"`
	Author      string          `json:"author,omitempty"`
	Score       int             `json:"score"`
	ScoreDelta  int             `json:"score_delta,omitempty"`
	Extra       json.RawMessage `json:"extra,omitempty"`
	PublishedAt *time.Time      `json:"published_at,omitempty"`
	CapturedAt  time.Time       `json:"captured_at"`
}

// ItemsResponse 是 GET /api/pulse/topics/:slug/items 的响应体。
// stale=true 表示数据超过 6 小时未更新，前端会显示黄色警告条。
type ItemsResponse struct {
	Topic         *Topic `json:"topic"`
	List          []Item `json:"list"`
	LastFetchedAt *time.Time `json:"last_fetched_at,omitempty"`
	Stale         bool   `json:"stale"`
}

// staleThreshold 是判定数据"不新鲜"的门槛。
// 超过此值前端会显示黄色警告"数据可能不新鲜"。
const staleThreshold = 6 * time.Hour
