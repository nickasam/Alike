package channel

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// ErrSlugConflict 表示 slug 已被占用（唯一约束冲突）。
var ErrSlugConflict = errors.New("channel slug already exists")

// ErrChannelNotFound 表示频道不存在。
var ErrChannelNotFound = errors.New("channel not found")

// ErrAlreadyMember 表示用户已加入该频道。
var ErrAlreadyMember = errors.New("already a member")

// ErrNotMember 表示用户未加入该频道。
var ErrNotMember = errors.New("not a member")

// pgUniqueViolation 是 PostgreSQL 唯一约束冲突的 SQLSTATE。
const pgUniqueViolation = "23505"

// Repository 封装 channels / channel_members 表的数据库访问。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建 channel 仓储。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// channelColumns 是查询频道时选取的完整列，顺序与 scanChannel 一致。
const channelColumns = `id, name, slug, COALESCE(description, ''), category,
	COALESCE(icon, ''), status, member_count, COALESCE(created_by, 0), created_at`

// List 返回频道列表（默认仅 status=active）。category 为空时不过滤分类。
func (r *Repository) List(ctx context.Context, category string, page, pageSize int) ([]*Channel, int64, error) {
	var (
		where = `WHERE status = 'active'`
		args  []any
	)
	if category != "" {
		where += ` AND category = $1`
		args = append(args, category)
	}

	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM channels `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limitArg := len(args) + 1
	offsetArg := len(args) + 2
	q := `SELECT ` + channelColumns + ` FROM channels ` + where +
		` ORDER BY member_count DESC, id DESC LIMIT $` + strconv.Itoa(limitArg) + ` OFFSET $` + strconv.Itoa(offsetArg)
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*Channel
	for rows.Next() {
		ch, err := scanChannel(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, ch)
	}
	return list, total, rows.Err()
}

// Create 插入新频道并将创建者登记为 admin 成员（member_count=1），返回完整记录。
// slug 唯一冲突时返回 ErrSlugConflict。整个过程在事务中完成，保证频道与成员一致。
func (r *Repository) Create(ctx context.Context, req CreateRequest, createdBy int64) (*Channel, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck // 提交成功后回滚为 no-op

	const q = `INSERT INTO channels (name, slug, description, category, icon, status, created_by, member_count)
		VALUES ($1, $2, $3, $4, $5, 'active', $6, 1)
		RETURNING ` + channelColumns

	row := tx.QueryRowContext(ctx, q, req.Name, req.Slug, req.Description, req.Category, req.Icon, createdBy)
	ch, err := scanChannel(row)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, ErrSlugConflict
		}
		return nil, err
	}

	// 创建者自动加入并成为管理员，否则无法在自己的频道发言/管理。
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO channel_members (channel_id, user_id, role) VALUES ($1, $2, 'admin')`,
		ch.ID, createdBy); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return ch, nil
}

// GetByID 按 ID 查询频道。不存在返回 ErrChannelNotFound。
func (r *Repository) GetByID(ctx context.Context, id int64) (*Channel, error) {
	const q = `SELECT ` + channelColumns + ` FROM channels WHERE id = $1`
	ch, err := scanChannel(r.db.QueryRowContext(ctx, q, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrChannelNotFound
	}
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// Join 将用户加入频道并原子递增 member_count。
// 频道不存在返回 ErrChannelNotFound；重复加入返回 ErrAlreadyMember。
func (r *Repository) Join(ctx context.Context, channelID, userID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck // 提交成功后回滚为 no-op

	var exists bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM channels WHERE id = $1)`, channelID).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return ErrChannelNotFound
	}

	res, err := tx.ExecContext(ctx,
		`INSERT INTO channel_members (channel_id, user_id) VALUES ($1, $2)
		 ON CONFLICT (channel_id, user_id) DO NOTHING`,
		channelID, userID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrAlreadyMember
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE channels SET member_count = member_count + 1 WHERE id = $1`, channelID); err != nil {
		return err
	}
	return tx.Commit()
}

// Leave 将用户移出频道并原子递减 member_count。
// 未加入返回 ErrNotMember。member_count 不会降到 0 以下。
func (r *Repository) Leave(ctx context.Context, channelID, userID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck // 提交成功后回滚为 no-op

	res, err := tx.ExecContext(ctx,
		`DELETE FROM channel_members WHERE channel_id = $1 AND user_id = $2`, channelID, userID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotMember
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE channels SET member_count = GREATEST(member_count - 1, 0) WHERE id = $1`, channelID); err != nil {
		return err
	}
	return tx.Commit()
}

// Members 返回频道成员列表（分页），按加入时间升序。
func (r *Repository) Members(ctx context.Context, channelID int64, page, pageSize int) ([]*Member, int64, error) {
	if _, err := r.GetByID(ctx, channelID); err != nil {
		return nil, 0, err
	}

	var total int64
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM channel_members WHERE channel_id = $1`, channelID).Scan(&total); err != nil {
		return nil, 0, err
	}

	const q = `SELECT cm.user_id, u.nickname, COALESCE(u.avatar_url, ''), cm.role, cm.joined_at
		FROM channel_members cm
		JOIN users u ON u.id = cm.user_id
		WHERE cm.channel_id = $1
		ORDER BY cm.joined_at ASC, cm.id ASC
		LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, q, channelID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*Member
	for rows.Next() {
		var m Member
		var joined time.Time
		if err := rows.Scan(&m.UserID, &m.Nickname, &m.AvatarURL, &m.Role, &joined); err != nil {
			return nil, 0, err
		}
		m.JoinedAt = joined
		list = append(list, &m)
	}
	return list, total, rows.Err()
}

// rowScanner 抽象 *sql.Row / *sql.Rows 的 Scan 能力，便于复用。
type rowScanner interface {
	Scan(dest ...any) error
}

func scanChannel(s rowScanner) (*Channel, error) {
	var ch Channel
	var created time.Time
	if err := s.Scan(
		&ch.ID, &ch.Name, &ch.Slug, &ch.Description, &ch.Category,
		&ch.Icon, &ch.Status, &ch.MemberCount, &ch.CreatedBy, &created,
	); err != nil {
		return nil, err
	}
	ch.CreatedAt = created
	return &ch, nil
}
