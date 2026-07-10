// Package emotion 负责情绪标签、情绪看板、趋势统计。
package emotion

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Tag 是支持的情绪标签，与 PRD/架构文档 2.4 一致（共 8 种）。
type Tag string

const (
	TagTired   Tag = "tired"   // 😮‍💨 疲惫
	TagAngry   Tag = "angry"   // 😡 愤怒
	TagWronged Tag = "wronged" // 😢 委屈
	TagBreak   Tag = "break"   // 🤯 崩溃
	TagNumb    Tag = "numb"    // 😴 麻木
	TagQuit    Tag = "quit"    // 🔥 想润
	TagAnxious Tag = "anxious" // 😰 焦虑
	TagCheer   Tag = "cheer"   // 💪 加油
)

// AllTags 返回全部合法情绪标签。
func AllTags() []Tag {
	return []Tag{TagTired, TagAngry, TagWronged, TagBreak, TagNumb, TagQuit, TagAnxious, TagCheer}
}

// IsValid 报告 tag 是否为受支持的情绪标签。
func IsValid(tag string) bool {
	switch Tag(tag) {
	case TagTired, TagAngry, TagWronged, TagBreak, TagNumb, TagQuit, TagAnxious, TagCheer:
		return true
	default:
		return false
	}
}

// ErrChannelNotFound 表示频道不存在。
var ErrChannelNotFound = errors.New("channel not found")

// boardLocation 是情绪看板"今日"判定使用的时区（面向国内打工人，用东八区）。
// 加载失败（如容器缺 tzdata）时回退到固定 +8 偏移，保证行为稳定。
var boardLocation = func() *time.Location {
	if loc, err := time.LoadLocation("Asia/Shanghai"); err == nil {
		return loc
	}
	return time.FixedZone("CST", 8*3600)
}()

// todayStart 返回 boardLocation 时区下当天 00:00 的时间点。
func todayStart() time.Time {
	now := time.Now().In(boardLocation)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, boardLocation)
}

// Count 是单个情绪标签的计数。
type Count struct {
	Emotion string `json:"emotion"`
	Count   int64  `json:"count"`
}

// Board 是频道情绪看板的聚合结果：各情绪计数 + 总计 + 统计范围。
// Emotions 始终按 AllTags 顺序返回全部 8 种（缺省计数为 0），便于前端稳定渲染热力图。
// Scope 为 "today" 或 "all"，Dominant 为占比最高的情绪 key（无数据时为空）。
type Board struct {
	Scope    string  `json:"scope"`
	Emotions []Count `json:"emotions"`
	Total    int64   `json:"total"`
	Dominant string  `json:"dominant,omitempty"`
}

// Repository 封装情绪看板相关的数据库访问。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建 emotion 仓储。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// channelExists 报告频道是否存在。
func (r *Repository) channelExists(ctx context.Context, channelID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM channels WHERE id = $1)`, channelID).Scan(&exists)
	return exists, err
}

// BoardByChannel 聚合频道内未删除消息的情绪计数。
// todayOnly 为 true 时仅统计 boardLocation 时区当天的消息。
// 频道不存在返回 ErrChannelNotFound。返回结果覆盖全部 8 种情绪（缺省为 0）。
func (r *Repository) BoardByChannel(ctx context.Context, channelID int64, todayOnly bool) (*Board, error) {
	exists, err := r.channelExists(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrChannelNotFound
	}

	q := `SELECT emotion, COUNT(*) FROM messages
		WHERE channel_id = $1 AND emotion IS NOT NULL AND emotion <> '' AND deleted_at IS NULL`
	args := []any{channelID}
	scope := "all"
	if todayOnly {
		q += ` AND created_at >= $2`
		args = append(args, todayStart())
		scope = "today"
	}
	q += ` GROUP BY emotion`

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int64)
	for rows.Next() {
		var emotion string
		var n int64
		if err := rows.Scan(&emotion, &n); err != nil {
			return nil, err
		}
		counts[emotion] = n
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	board := &Board{Scope: scope, Emotions: make([]Count, 0, len(AllTags()))}
	var maxCount int64
	for _, tag := range AllTags() {
		n := counts[string(tag)]
		board.Emotions = append(board.Emotions, Count{Emotion: string(tag), Count: n})
		board.Total += n
		if n > maxCount {
			maxCount = n
			board.Dominant = string(tag)
		}
	}
	return board, nil
}
