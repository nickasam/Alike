package empathy

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// pgUniqueViolation 是 PostgreSQL 唯一约束冲突的 SQLSTATE。
const pgUniqueViolation = "23505"

// ErrMessageNotFound 表示消息不存在（或已软删除）。
var ErrMessageNotFound = errors.New("message not found")

// ErrAlreadyEmpathized 表示当前用户已对该消息共情过（唯一约束冲突）。
var ErrAlreadyEmpathized = errors.New("already empathized")

// ErrNotEmpathized 表示当前用户尚未对该消息共情，无法取消。
var ErrNotEmpathized = errors.New("not empathized")

// Repository 封装 empathies 表及关联计数的数据库访问。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建 empathy 仓储。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// authorOf 返回未删除消息的作者 ID。消息不存在返回 ErrMessageNotFound。
func authorOf(ctx context.Context, tx *sql.Tx, messageID int64) (int64, error) {
	var authorID int64
	err := tx.QueryRowContext(ctx,
		`SELECT user_id FROM messages WHERE id = $1 AND deleted_at IS NULL`,
		messageID).Scan(&authorID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrMessageNotFound
	}
	return authorID, err
}

// Create 事务化地为消息添加一次共情：
// INSERT empathies + messages.empathy_count+1 + 作者 empathy_received+1 + 当前用户 empathy_given+1。
// 消息不存在返回 ErrMessageNotFound；重复共情返回 ErrAlreadyEmpathized。
// 返回消息的最新共情计数。
func (r *Repository) Create(ctx context.Context, messageID, userID int64) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck // 提交成功后回滚为 no-op

	authorID, err := authorOf(ctx, tx, messageID)
	if err != nil {
		return 0, err
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO empathies (message_id, user_id) VALUES ($1, $2)`,
		messageID, userID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return 0, ErrAlreadyEmpathized
		}
		return 0, err
	}

	var count int64
	if err := tx.QueryRowContext(ctx,
		`UPDATE messages SET empathy_count = empathy_count + 1 WHERE id = $1 RETURNING empathy_count`,
		messageID).Scan(&count); err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE users SET empathy_received = empathy_received + 1 WHERE id = $1`, authorID); err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE users SET empathy_given = empathy_given + 1 WHERE id = $1`, userID); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return count, nil
}

// Delete 事务化地取消一次共情，做 Create 的反向操作。计数不会降到 0 以下。
// 消息不存在返回 ErrMessageNotFound；未共情返回 ErrNotEmpathized。
// 返回消息的最新共情计数。
func (r *Repository) Delete(ctx context.Context, messageID, userID int64) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck // 提交成功后回滚为 no-op

	authorID, err := authorOf(ctx, tx, messageID)
	if err != nil {
		return 0, err
	}

	res, err := tx.ExecContext(ctx,
		`DELETE FROM empathies WHERE message_id = $1 AND user_id = $2`, messageID, userID)
	if err != nil {
		return 0, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return 0, ErrNotEmpathized
	}

	var count int64
	if err := tx.QueryRowContext(ctx,
		`UPDATE messages SET empathy_count = GREATEST(empathy_count - 1, 0) WHERE id = $1 RETURNING empathy_count`,
		messageID).Scan(&count); err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE users SET empathy_received = GREATEST(empathy_received - 1, 0) WHERE id = $1`, authorID); err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE users SET empathy_given = GREATEST(empathy_given - 1, 0) WHERE id = $1`, userID); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return count, nil
}

// ListUsers 返回对某消息共情的用户列表（按共情时间倒序，分页）。
// 消息不存在返回 ErrMessageNotFound。
func (r *Repository) ListUsers(ctx context.Context, messageID int64, page, pageSize int) ([]*User, int64, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM messages WHERE id = $1 AND deleted_at IS NULL)`,
		messageID).Scan(&exists); err != nil {
		return nil, 0, err
	}
	if !exists {
		return nil, 0, ErrMessageNotFound
	}

	var total int64
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM empathies WHERE message_id = $1`, messageID).Scan(&total); err != nil {
		return nil, 0, err
	}

	const q = `SELECT u.id, u.nickname, COALESCE(u.avatar_url, ''), e.created_at
		FROM empathies e JOIN users u ON u.id = e.user_id
		WHERE e.message_id = $1
		ORDER BY e.created_at DESC, e.id DESC
		LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, q, messageID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Nickname, &u.AvatarURL, &u.CreatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &u)
	}
	return list, total, rows.Err()
}

// RankingEmpathy 返回最受共情的帖子榜（按 empathy_count DESC），仅统计未删除主消息与回复。
func (r *Repository) RankingEmpathy(ctx context.Context, limit int) ([]*RankMessage, error) {
	const q = `SELECT m.id, m.channel_id, m.content, COALESCE(m.emotion, ''), m.is_anonymous, m.empathy_count,
		m.user_id, u.nickname, COALESCE(u.avatar_url, '')
		FROM messages m JOIN users u ON u.id = m.user_id
		WHERE m.deleted_at IS NULL AND m.empathy_count > 0
		ORDER BY m.empathy_count DESC, m.created_at DESC
		LIMIT $1`
	rows, err := r.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*RankMessage
	for rows.Next() {
		var m RankMessage
		var authorID int64
		var nickname, avatar string
		if err := rows.Scan(&m.MessageID, &m.ChannelID, &m.Content, &m.Emotion, &m.IsAnonymous, &m.EmpathyCount,
			&authorID, &nickname, &avatar); err != nil {
			return nil, err
		}
		if !m.IsAnonymous {
			m.Author = &Author{ID: authorID, Nickname: nickname, AvatarURL: avatar}
		}
		list = append(list, &m)
	}
	return list, rows.Err()
}

// RankingWarmest 返回最暖牛马榜（按 empathy_given DESC）。
func (r *Repository) RankingWarmest(ctx context.Context, limit int) ([]*RankUser, error) {
	const q = `SELECT id, nickname, COALESCE(avatar_url, ''), level, empathy_given
		FROM users
		WHERE empathy_given > 0
		ORDER BY empathy_given DESC, id ASC
		LIMIT $1`
	return r.rankUsers(ctx, q, limit)
}

// RankingActive 返回本周（近 7 天）最活跃牛马榜（按未删除消息数 DESC）。
func (r *Repository) RankingActive(ctx context.Context, limit int) ([]*RankUser, error) {
	const q = `SELECT u.id, u.nickname, COALESCE(u.avatar_url, ''), u.level, COUNT(m.id) AS msg_count
		FROM users u JOIN messages m ON m.user_id = u.id
		WHERE m.deleted_at IS NULL AND m.created_at >= NOW() - INTERVAL '7 days'
		GROUP BY u.id, u.nickname, u.avatar_url, u.level
		ORDER BY msg_count DESC, u.id ASC
		LIMIT $1`
	return r.rankUsers(ctx, q, limit)
}

// rankUsers 执行返回 (id, nickname, avatar, level, metric) 五列的榜单查询。
func (r *Repository) rankUsers(ctx context.Context, q string, limit int) ([]*RankUser, error) {
	rows, err := r.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*RankUser
	for rows.Next() {
		var u RankUser
		if err := rows.Scan(&u.UserID, &u.Nickname, &u.AvatarURL, &u.Level, &u.Metric); err != nil {
			return nil, err
		}
		list = append(list, &u)
	}
	return list, rows.Err()
}
