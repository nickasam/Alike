package search

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
)

// Repository 封装跨表的模糊搜索数据库访问。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建 search 仓储。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// likePattern 将用户关键词转义为安全的 ILIKE 模式，%_\ 均转义避免注入通配符，
// 再包裹 % 实现子串匹配。
func likePattern(q string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return "%" + r.Replace(q) + "%"
}

// SearchMessages 按 content ILIKE 搜索未软删除的消息，可选按频道过滤。
// 匿名消息不返回作者信息。
func (r *Repository) SearchMessages(ctx context.Context, q string, channelID int64, page, pageSize int) ([]*MessageResult, int64, error) {
	where := `WHERE m.deleted_at IS NULL AND m.content ILIKE $1 ESCAPE '\'`
	args := []any{likePattern(q)}
	if channelID > 0 {
		args = append(args, channelID)
		where += ` AND m.channel_id = $` + strconv.Itoa(len(args))
	}

	var total int64
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM messages m `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	limArg := len(args) + 1
	offArg := len(args) + 2
	args = append(args, pageSize, (page-1)*pageSize)
	q2 := `SELECT m.id, m.channel_id, m.content, COALESCE(m.emotion, ''), m.is_anonymous,
		m.user_id, u.nickname, COALESCE(u.avatar_url, ''), m.created_at
		FROM messages m JOIN users u ON u.id = m.user_id ` + where +
		` ORDER BY m.created_at DESC, m.id DESC LIMIT $` + strconv.Itoa(limArg) + ` OFFSET $` + strconv.Itoa(offArg)

	rows, err := r.db.QueryContext(ctx, q2, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*MessageResult
	for rows.Next() {
		var (
			m         MessageResult
			authorID  int64
			nickname  string
			avatarURL string
		)
		if err := rows.Scan(&m.ID, &m.ChannelID, &m.Content, &m.Emotion, &m.IsAnonymous,
			&authorID, &nickname, &avatarURL, &m.CreatedAt); err != nil {
			return nil, 0, err
		}
		if !m.IsAnonymous {
			m.Author = &Author{ID: authorID, Nickname: nickname, AvatarURL: avatarURL}
		}
		list = append(list, &m)
	}
	return list, total, rows.Err()
}

// SearchDiaries 按 title/content ILIKE 搜索公开日记。
func (r *Repository) SearchDiaries(ctx context.Context, q string, page, pageSize int) ([]*DiaryResult, int64, error) {
	pattern := likePattern(q)
	const where = `WHERE d.is_public = TRUE AND (d.content ILIKE $1 ESCAPE '\' OR COALESCE(d.title, '') ILIKE $1 ESCAPE '\')`

	var total int64
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM diaries d `+where, pattern).Scan(&total); err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	q2 := `SELECT d.id, COALESCE(d.title, ''), d.content, COALESCE(d.mood, ''),
		d.user_id, u.nickname, COALESCE(u.avatar_url, ''), d.created_at
		FROM diaries d JOIN users u ON u.id = d.user_id ` + where +
		` ORDER BY d.created_at DESC, d.id DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, q2, pattern, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*DiaryResult
	for rows.Next() {
		var (
			d         DiaryResult
			authorID  int64
			nickname  string
			avatarURL string
		)
		if err := rows.Scan(&d.ID, &d.Title, &d.Content, &d.Mood,
			&authorID, &nickname, &avatarURL, &d.CreatedAt); err != nil {
			return nil, 0, err
		}
		d.Author = &Author{ID: authorID, Nickname: nickname, AvatarURL: avatarURL}
		list = append(list, &d)
	}
	return list, total, rows.Err()
}

// SearchChannels 按 name/slug/description ILIKE 搜索非归档频道。
func (r *Repository) SearchChannels(ctx context.Context, q string, page, pageSize int) ([]*ChannelResult, int64, error) {
	pattern := likePattern(q)
	const where = `WHERE c.status <> 'archived' AND (c.name ILIKE $1 ESCAPE '\' OR c.slug ILIKE $1 ESCAPE '\' OR COALESCE(c.description, '') ILIKE $1 ESCAPE '\')`

	var total int64
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM channels c `+where, pattern).Scan(&total); err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	q2 := `SELECT c.id, c.name, c.slug, COALESCE(c.description, ''), c.category,
		COALESCE(c.icon, ''), c.member_count, c.created_at
		FROM channels c ` + where +
		` ORDER BY c.member_count DESC, c.id DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, q2, pattern, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*ChannelResult
	for rows.Next() {
		var ch ChannelResult
		if err := rows.Scan(&ch.ID, &ch.Name, &ch.Slug, &ch.Description, &ch.Category,
			&ch.Icon, &ch.MemberCount, &ch.CreatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &ch)
	}
	return list, total, rows.Err()
}

// SearchUsers 按 nickname/bio ILIKE 搜索用户。
func (r *Repository) SearchUsers(ctx context.Context, q string, page, pageSize int) ([]*UserResult, int64, error) {
	pattern := likePattern(q)
	const where = `WHERE u.nickname ILIKE $1 ESCAPE '\' OR COALESCE(u.bio, '') ILIKE $1 ESCAPE '\'`

	var total int64
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM users u `+where, pattern).Scan(&total); err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	q2 := `SELECT u.id, u.nickname, COALESCE(u.avatar_url, ''), COALESCE(u.bio, ''),
		COALESCE(u.industry, ''), COALESCE(u.job_title, ''), u.level
		FROM users u ` + where +
		` ORDER BY u.empathy_received DESC, u.id DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, q2, pattern, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*UserResult
	for rows.Next() {
		var u UserResult
		if err := rows.Scan(&u.ID, &u.Nickname, &u.AvatarURL, &u.Bio,
			&u.Industry, &u.JobTitle, &u.Level); err != nil {
			return nil, 0, err
		}
		list = append(list, &u)
	}
	return list, total, rows.Err()
}
