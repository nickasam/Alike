package pulse

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// ErrTopicNotFound 表示按 slug 找不到专题。
var ErrTopicNotFound = errors.New("pulse topic not found")

// Repository 封装 pulse_topics / pulse_items 表的数据库访问。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建 pulse 仓储。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// topicCols 是 SELECT pulse_topics 的列清单，顺序与 scanTopic 一致。
const topicCols = `id, slug, name, COALESCE(emoji, ''), COALESCE(description, ''),
	collector_kind, COALESCE(collector_config, '{}'::jsonb), sort_order, is_active,
	refresh_interval_min, last_fetched_at, COALESCE(last_error, '')`

func scanTopic(row interface {
	Scan(dest ...any) error
}) (*Topic, error) {
	t := &Topic{}
	err := row.Scan(&t.ID, &t.Slug, &t.Name, &t.Emoji, &t.Description,
		&t.CollectorKind, &t.CollectorConfig, &t.SortOrder, &t.IsActive,
		&t.RefreshIntervalMin, &t.LastFetchedAt, &t.LastError)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// ListActiveTopics 返回所有活跃专题，按 sort_order 升序。
// 前端 Tab / 侧栏专题项数据来源。
func (r *Repository) ListActiveTopics(ctx context.Context) ([]*Topic, error) {
	q := `SELECT ` + topicCols + ` FROM pulse_topics WHERE is_active = TRUE ORDER BY sort_order ASC, id ASC`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Topic
	for rows.Next() {
		t, err := scanTopic(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

// GetTopicBySlug 按 slug 查找单个专题。找不到返回 ErrTopicNotFound。
func (r *Repository) GetTopicBySlug(ctx context.Context, slug string) (*Topic, error) {
	q := `SELECT ` + topicCols + ` FROM pulse_topics WHERE slug = $1`
	t, err := scanTopic(r.db.QueryRowContext(ctx, q, slug))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTopicNotFound
	}
	return t, err
}

// itemCols 是 SELECT pulse_items 的列清单。
const itemCols = `id, topic_id, source, source_id, title, COALESCE(summary, ''), url,
	COALESCE(author, ''), score, score_delta, COALESCE(extra, '{}'::jsonb),
	published_at, captured_at`

func scanItem(row interface {
	Scan(dest ...any) error
}) (*Item, error) {
	it := &Item{}
	err := row.Scan(&it.ID, &it.TopicID, &it.Source, &it.SourceID, &it.Title, &it.Summary, &it.URL,
		&it.Author, &it.Score, &it.ScoreDelta, &it.Extra, &it.PublishedAt, &it.CapturedAt)
	if err != nil {
		return nil, err
	}
	return it, nil
}

// ListItemsByTopic 返回某专题的 Top N 条目，按 score DESC, captured_at DESC。
// limit<=0 时用默认值 25。
func (r *Repository) ListItemsByTopic(ctx context.Context, topicID int64, limit int) ([]*Item, error) {
	if limit <= 0 {
		limit = 25
	}
	q := `SELECT ` + itemCols + ` FROM pulse_items WHERE topic_id = $1
		ORDER BY score DESC, captured_at DESC LIMIT $2`
	rows, err := r.db.QueryContext(ctx, q, topicID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Item
	for rows.Next() {
		it, err := scanItem(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, it)
	}
	return list, rows.Err()
}

// UpsertItem 幂等写入一条 pulse_item（M0 尚不调用，M1 collector 使用）。
// 用 ON CONFLICT (topic_id, source, source_id) 保证同一外部事件重复抓取时更新而非新增。
// score_delta 记录相较上次抓取的增量。
func (r *Repository) UpsertItem(ctx context.Context, it *Item) error {
	q := `INSERT INTO pulse_items
		(topic_id, source, source_id, title, summary, url, author, score, extra, published_at, captured_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW())
		ON CONFLICT (topic_id, source, source_id) DO UPDATE SET
			title = EXCLUDED.title,
			summary = EXCLUDED.summary,
			url = EXCLUDED.url,
			author = EXCLUDED.author,
			score_delta = EXCLUDED.score - pulse_items.score,
			score = EXCLUDED.score,
			extra = EXCLUDED.extra,
			captured_at = NOW()`
	_, err := r.db.ExecContext(ctx, q,
		it.TopicID, it.Source, it.SourceID, it.Title, it.Summary, it.URL, it.Author,
		it.Score, it.Extra, it.PublishedAt)
	return err
}

// MarkTopicFetched 更新专题的成功抓取时间戳（M0 尚不调用，M1 scheduler 用）。
func (r *Repository) MarkTopicFetched(ctx context.Context, topicID int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE pulse_topics SET last_fetched_at = NOW(), last_error = NULL, updated_at = NOW() WHERE id = $1`,
		topicID)
	return err
}

// MarkTopicError 记录抓取失败（M0 尚不调用，M1 scheduler 用）。
// 不覆写 last_fetched_at —— 前端仍显示上次成功时间。
func (r *Repository) MarkTopicError(ctx context.Context, topicID int64, msg string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE pulse_topics SET last_error = $2, updated_at = NOW() WHERE id = $1`,
		topicID, msg)
	return err
}

// CleanupOldItems 删除 topicID 超过 window 时间的旧数据（M0 尚不调用，M1 scheduler 每日跑）。
// 返回删除行数。
func (r *Repository) CleanupOldItems(ctx context.Context, window time.Duration) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM pulse_items WHERE captured_at < NOW() - $1::interval`,
		window.String())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
