// Package collector 定义 pulse 抓取器的通用契约与注册表。
// 每个专题（如 GitHub Trending / HN AI）对应一个 Collector 实现，
// scheduler 按 pulse_topics.collector_kind 从 Registry 中查表调度。
// M0 阶段只提供接口 + Registry，M1/M2 补上 github_trending / hackernews_ai 实现。
package collector

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// Item 是抓取器返回的标准条目，字段对应 pulse_items 表的核心列。
// collector 负责组装 Extra JSONB（源特有字段：lang / total_stars / fork_count / hn_comments 等）。
type Item struct {
	Source      string          // 'github' / 'hackernews'
	SourceID    string          // 去重键：'owner/repo' / HN story_id
	Title       string
	Summary     string
	URL         string
	Author      string
	Score       int             // 主分：today_stars / hn_points
	Extra       json.RawMessage // JSONB 源特有字段
	PublishedAt *time.Time
}

// Collector 是单个专题的抓取器契约。实现必须无状态、并发安全（scheduler 可能并行触发）。
// Kind 是注册键，与 pulse_topics.collector_kind 对应；Fetch 抓一批 Item 或返回 error。
// 失败时 scheduler 不覆写数据库，前端仍看到上次成功的快照（见 pulse-module-design.md 第 3.6 节）。
type Collector interface {
	Kind() string
	Fetch(ctx context.Context, config json.RawMessage) ([]Item, error)
}

// Registry 是 collector 注册表。启动阶段各 collector 实现调用 Register 挂上，
// scheduler 查 pulse_topics.collector_kind 时按键取用。
var (
	registryMu sync.RWMutex
	registry   = map[string]Collector{}
)

// Register 注册一个 collector。重复注册同一 kind 时后者覆盖前者（通常发生在测试）。
func Register(c Collector) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[c.Kind()] = c
}

// Get 按 kind 查找 collector。未找到返回 nil, false。
func Get(kind string) (Collector, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	c, ok := registry[kind]
	return c, ok
}

// Kinds 返回已注册的所有 kind（顺序不保证）。用于启动时的健康检查。
func Kinds() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	out := make([]string, 0, len(registry))
	for k := range registry {
		out = append(out, k)
	}
	return out
}
