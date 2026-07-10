package message

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	emotionpkg "github.com/Alike/backend/internal/emotion"
)

// itoa 是 strconv.Itoa 的别名，用于拼接 SQL 占位符序号。
func itoa(n int) string { return strconv.Itoa(n) }

// ErrChannelNotFound 表示频道不存在。
var ErrChannelNotFound = errors.New("channel not found")

// ErrMessageNotFound 表示消息不存在（或已彻底不可见）。
var ErrMessageNotFound = errors.New("message not found")

// ErrNotMember 表示用户尚未加入该频道，无权发言。
var ErrNotMember = errors.New("not a channel member")

// ErrForbidden 表示当前用户无权操作该资源（如删除他人消息）。
var ErrForbidden = errors.New("forbidden")

// ErrInvalidEmotion 表示提供的情绪标签不在受支持集合内。
var ErrInvalidEmotion = errors.New("invalid emotion tag")

// Repository 封装 messages 表的数据库访问。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建 message 仓储。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// messageColumns 是查询消息时选取的列，顺序与 scanMessage 一致。
// reply_count 通过子查询统计未删除的直接回复数量；
// empathized 表示 $1（viewer）是否已对该消息共情（viewer<=0 时恒 false）。
// 约定：使用 messageColumns 的查询必须以 $1 = viewerID 作为首个参数。
const messageColumns = `m.id, m.channel_id, m.parent_id, m.content, COALESCE(m.emotion, ''),
	m.is_anonymous, m.empathy_count,
	(SELECT COUNT(*) FROM messages r WHERE r.parent_id = m.id AND r.deleted_at IS NULL) AS reply_count,
	EXISTS(SELECT 1 FROM empathies e WHERE e.message_id = m.id AND e.user_id = $1) AS empathized,
	m.user_id, u.nickname, COALESCE(u.avatar_url, ''), m.created_at, m.deleted_at`

const messageFrom = ` FROM messages m JOIN users u ON u.id = m.user_id`

// channelExists 报告频道是否存在。
func (r *Repository) channelExists(ctx context.Context, channelID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM channels WHERE id = $1)`, channelID).Scan(&exists)
	return exists, err
}

// IsMember 报告用户是否为频道成员。
func (r *Repository) IsMember(ctx context.Context, channelID, userID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM channel_members WHERE channel_id = $1 AND user_id = $2)`,
		channelID, userID).Scan(&exists)
	return exists, err
}

// EnsureMember 幂等地将用户加入频道（若尚未加入），并维护 member_count。
// 用于"发消息即自动入群"，消除前端异步加入与发送之间的竞态。
// 频道不存在返回 ErrChannelNotFound。
func (r *Repository) EnsureMember(ctx context.Context, channelID, userID int64) error {
	var exists bool
	if err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM channels WHERE id = $1)`, channelID).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return ErrChannelNotFound
	}
	// ON CONFLICT DO NOTHING 保证幂等；仅真正新增时才 +1 member_count。
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO channel_members (channel_id, user_id) VALUES ($1, $2)
		 ON CONFLICT (channel_id, user_id) DO NOTHING`,
		channelID, userID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		if _, err := r.db.ExecContext(ctx,
			`UPDATE channels SET member_count = member_count + 1 WHERE id = $1`, channelID); err != nil {
			return err
		}
	}
	return nil
}

// ListByChannel 返回频道主消息列表（parent_id IS NULL），按 created_at DESC 游标分页。
// before>0 时仅返回早于该消息的记录（用其 created_at + id 作游标，避免 OFFSET 深翻）。
// viewerID 用于判定每条消息当前用户是否已共情（<=0 表示未登录，恒 false）。
// 返回列表与 hasMore 标记。
func (r *Repository) ListByChannel(ctx context.Context, channelID, before, viewerID int64, limit int) ([]*Message, bool, error) {
	exists, err := r.channelExists(ctx, channelID)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, ErrChannelNotFound
	}

	// $1 = viewerID（messageColumns 的 empathized 子查询依赖），$2 = channelID。
	args := []any{viewerID, channelID}
	where := `WHERE m.channel_id = $2 AND m.parent_id IS NULL`
	if before > 0 {
		// 游标：取 before 消息的 (created_at, id)，返回严格更早的记录。
		where += ` AND (m.created_at, m.id) < (SELECT created_at, id FROM messages WHERE id = $3)`
		args = append(args, before)
	}
	args = append(args, limit+1) // 多取 1 条判断是否还有更多
	q := `SELECT ` + messageColumns + messageFrom + ` ` + where +
		` ORDER BY m.created_at DESC, m.id DESC LIMIT $` + itoa(len(args))

	return r.queryList(ctx, q, args, limit)
}

// ListThreads 返回某主消息下的线程回复列表（parent_id = 该消息），按 created_at ASC 游标分页。
// after>0 时仅返回晚于该消息的记录。viewerID 用于判定已共情态（<=0 恒 false）。
func (r *Repository) ListThreads(ctx context.Context, parentID, after, viewerID int64, limit int) ([]*Message, bool, error) {
	// 父消息必须存在且未删除。
	var exists bool
	if err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM messages WHERE id = $1 AND deleted_at IS NULL)`,
		parentID).Scan(&exists); err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, ErrMessageNotFound
	}

	// $1 = viewerID（empathized 子查询），$2 = parentID。
	args := []any{viewerID, parentID}
	where := `WHERE m.parent_id = $2`
	if after > 0 {
		where += ` AND (m.created_at, m.id) > (SELECT created_at, id FROM messages WHERE id = $3)`
		args = append(args, after)
	}
	args = append(args, limit+1)
	q := `SELECT ` + messageColumns + messageFrom + ` ` + where +
		` ORDER BY m.created_at ASC, m.id ASC LIMIT $` + itoa(len(args))

	return r.queryList(ctx, q, args, limit)
}

// queryList 执行分页查询，取回 limit+1 条以推导 hasMore，并对结果做匿名/软删除脱敏。
func (r *Repository) queryList(ctx context.Context, q string, args []any, limit int) ([]*Message, bool, error) {
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var list []*Message
	for rows.Next() {
		m, err := scanMessage(rows)
		if err != nil {
			return nil, false, err
		}
		list = append(list, m)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	hasMore := len(list) > limit
	if hasMore {
		list = list[:limit]
	}
	for _, m := range list {
		m.mask()
	}
	return list, hasMore, nil
}

// Create 插入一条消息（主消息或线程回复），校验频道存在与成员身份，返回完整记录。
// parentID 为 nil 时为主消息；非 nil 时为线程回复（channelID 继承自父消息）。
func (r *Repository) Create(ctx context.Context, channelID int64, parentID *int64, userID int64, req CreateRequest) (*Message, error) {
	exists, err := r.channelExists(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrChannelNotFound
	}

	member, err := r.IsMember(ctx, channelID, userID)
	if err != nil {
		return nil, err
	}
	if !member {
		return nil, ErrNotMember
	}

	var emotion any
	if req.Emotion != "" {
		if !emotionpkg.IsValid(req.Emotion) {
			return nil, ErrInvalidEmotion
		}
		emotion = req.Emotion
	}

	const q = `INSERT INTO messages (channel_id, user_id, parent_id, content, emotion, is_anonymous)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at,
			(SELECT nickname FROM users WHERE id = $2),
			(SELECT COALESCE(avatar_url, '') FROM users WHERE id = $2)`
	var id int64
	var created time.Time
	var nickname, avatarURL string
	if err := r.db.QueryRowContext(ctx, q,
		channelID, userID, parentID, req.Content, emotion, req.IsAnonymous,
	).Scan(&id, &created, &nickname, &avatarURL); err != nil {
		return nil, err
	}

	m := &Message{
		ID:          id,
		ChannelID:   channelID,
		ParentID:    parentID,
		Content:     req.Content,
		Emotion:     req.Emotion,
		IsAnonymous: req.IsAnonymous,
		Author:      &Author{ID: userID, Nickname: nickname, AvatarURL: avatarURL},
		CreatedAt:   created,
		authorID:    userID,
	}
	return m, nil
}

// CreateReply 校验父消息存在且未删除，并在其所属频道下插入一条线程回复。
func (r *Repository) CreateReply(ctx context.Context, parentID, userID int64, req CreateRequest) (*Message, error) {
	var channelID int64
	err := r.db.QueryRowContext(ctx,
		`SELECT channel_id FROM messages WHERE id = $1 AND deleted_at IS NULL`,
		parentID).Scan(&channelID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrMessageNotFound
	}
	if err != nil {
		return nil, err
	}
	return r.Create(ctx, channelID, &parentID, userID, req)
}

// AuthorAndChannel 返回消息作者 ID 与所属频道 ID，供删除鉴权使用。
// 消息不存在或已软删除返回 ErrMessageNotFound。
func (r *Repository) AuthorAndChannel(ctx context.Context, messageID int64) (authorID, channelID int64, err error) {
	err = r.db.QueryRowContext(ctx,
		`SELECT user_id, channel_id FROM messages WHERE id = $1 AND deleted_at IS NULL`,
		messageID).Scan(&authorID, &channelID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, 0, ErrMessageNotFound
	}
	return authorID, channelID, err
}

// isChannelAdmin 报告用户是否为频道管理员。
func (r *Repository) isChannelAdmin(ctx context.Context, channelID, userID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM channel_members WHERE channel_id = $1 AND user_id = $2 AND role = 'admin')`,
		channelID, userID).Scan(&exists)
	return exists, err
}

// SoftDelete 软删除消息（置 deleted_at）。仅作者或频道管理员可删。
// 消息不存在返回 ErrMessageNotFound；无权返回 ErrForbidden。
// 事务内同时回收该消息的共情计数：按其 empathies 行数扣减作者的 empathy_received，
// 否则被共情的消息删除后作者 empathy_received 会永久虚高、污染排行榜。
// 成功返回该消息所属频道 ID，供上层广播 message_deleted 事件。
func (r *Repository) SoftDelete(ctx context.Context, messageID, userID int64) (int64, error) {
	authorID, channelID, err := r.AuthorAndChannel(ctx, messageID)
	if err != nil {
		return 0, err
	}
	if authorID != userID {
		admin, err := r.isChannelAdmin(ctx, channelID, userID)
		if err != nil {
			return 0, err
		}
		if !admin {
			return 0, ErrForbidden
		}
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck // 提交成功后回滚为 no-op

	res, err := tx.ExecContext(ctx,
		`UPDATE messages SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
		messageID)
	if err != nil {
		return 0, err
	}
	// 已被并发删除则无需重复回收计数。
	if n, _ := res.RowsAffected(); n == 0 {
		return channelID, tx.Commit()
	}

	// 回收该消息的共情计数，保持聚合与 empathies 行一致：
	//   1) 作者 empathy_received 扣减该消息收到的共情数；
	//   2) 每个共情者 empathy_given 各扣 1（否则暖榜 empathy_given DESC 永久虚高）；
	//   3) 删除该消息的 empathies 行，避免残留悬挂数据。
	var empCount int64
	if err := tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM empathies WHERE message_id = $1`, messageID).Scan(&empCount); err != nil {
		return 0, err
	}
	if empCount > 0 {
		if _, err := tx.ExecContext(ctx,
			`UPDATE users SET empathy_received = GREATEST(empathy_received - $1, 0) WHERE id = $2`,
			empCount, authorID); err != nil {
			return 0, err
		}
		if _, err := tx.ExecContext(ctx,
			`UPDATE users SET empathy_given = GREATEST(empathy_given - 1, 0)
			 WHERE id IN (SELECT user_id FROM empathies WHERE message_id = $1)`,
			messageID); err != nil {
			return 0, err
		}
		if _, err := tx.ExecContext(ctx,
			`DELETE FROM empathies WHERE message_id = $1`, messageID); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return channelID, nil
}

// rowScanner 抽象 *sql.Row / *sql.Rows 的 Scan 能力。
type rowScanner interface {
	Scan(dest ...any) error
}

func scanMessage(s rowScanner) (*Message, error) {
	var (
		m         Message
		parentID  sql.NullInt64
		nickname  string
		avatarURL string
		deletedAt sql.NullTime
	)
	if err := s.Scan(
		&m.ID, &m.ChannelID, &parentID, &m.Content, &m.Emotion,
		&m.IsAnonymous, &m.EmpathyCount, &m.ReplyCount, &m.Empathized,
		&m.authorID, &nickname, &avatarURL, &m.CreatedAt, &deletedAt,
	); err != nil {
		return nil, err
	}
	if parentID.Valid {
		m.ParentID = &parentID.Int64
	}
	if deletedAt.Valid {
		m.IsDeleted = true
		m.DeletedAt = &deletedAt.Time
	}
	m.Author = &Author{ID: m.authorID, Nickname: nickname, AvatarURL: avatarURL}
	return &m, nil
}
