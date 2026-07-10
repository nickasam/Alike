package diary

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// ErrDiaryNotFound 表示日记不存在。
var ErrDiaryNotFound = errors.New("diary not found")

// ErrUserNotFound 表示用户不存在。
var ErrUserNotFound = errors.New("user not found")

// Repository 封装 diaries / diary_comments 表的数据库访问。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建 diary 仓储。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// diaryColumns 是查询日记时选取的列，顺序与 scanDiary 一致。
// comment_count 通过子查询统计未删除评论数。
const diaryColumns = `d.id, COALESCE(d.title, ''), d.content, COALESCE(d.mood, ''), d.is_public,
	(SELECT COUNT(*) FROM diary_comments dc WHERE dc.diary_id = d.id AND dc.deleted_at IS NULL) AS comment_count,
	d.user_id, u.nickname, COALESCE(u.avatar_url, ''), d.created_at`

const diaryFrom = ` FROM diaries d JOIN users u ON u.id = d.user_id`

// ListPublic 返回公开日记流（is_public=true），按 created_at DESC 游标分页。
// before>0 时仅返回更早的记录。返回列表与 hasMore 标记。
func (r *Repository) ListPublic(ctx context.Context, before int64, limit int) ([]*Diary, bool, error) {
	args := []any{}
	where := `WHERE d.is_public = TRUE`
	if before > 0 {
		where += ` AND (d.created_at, d.id) < (SELECT created_at, id FROM diaries WHERE id = $1)`
		args = append(args, before)
	}
	args = append(args, limit+1) // 多取 1 条判断是否还有更多
	q := `SELECT ` + diaryColumns + diaryFrom + ` ` + where +
		` ORDER BY d.created_at DESC, d.id DESC LIMIT $` + itoa(len(args))

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var list []*Diary
	for rows.Next() {
		d, err := scanDiary(rows)
		if err != nil {
			return nil, false, err
		}
		list = append(list, d)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	hasMore := len(list) > limit
	if hasMore {
		list = list[:limit]
	}
	return list, hasMore, nil
}

// Get 返回单篇日记详情。私密日记仅作者本人（viewerID==作者）可见，
// 其他访问者（含未登录，viewerID<=0）对私密日记视为不存在，避免泄露。
// 日记不存在或无权查看返回 ErrDiaryNotFound。
func (r *Repository) Get(ctx context.Context, diaryID, viewerID int64) (*Diary, error) {
	q := `SELECT ` + diaryColumns + diaryFrom + ` WHERE d.id = $1 AND (d.is_public = TRUE OR d.user_id = $2)`
	d, err := scanDiary(r.db.QueryRowContext(ctx, q, diaryID, viewerID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrDiaryNotFound
	}
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Create 插入一篇日记并同步累计打卡天数，返回完整记录（含作者信息）。
func (r *Repository) Create(ctx context.Context, userID int64, req CreateRequest) (*Diary, error) {
	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck // 提交成功后回滚为 no-op

	var title, mood any
	if req.Title != "" {
		title = req.Title
	}
	if req.Mood != "" {
		mood = req.Mood
	}

	const insert = `INSERT INTO diaries (user_id, title, content, mood, is_public)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`
	var id int64
	var created time.Time
	if err := tx.QueryRowContext(ctx, insert,
		userID, title, req.Content, mood, isPublic,
	).Scan(&id, &created); err != nil {
		return nil, err
	}

	// users.total_check_in_days 记录累计打卡（去重日期）天数。
	// 仅当当天首篇日记时 +1，避免同一天多写重复计数。
	const bump = `UPDATE users SET total_check_in_days = (
		SELECT COUNT(DISTINCT DATE(created_at)) FROM diaries WHERE user_id = $1
	) WHERE id = $1`
	if _, err := tx.ExecContext(ctx, bump, userID); err != nil {
		return nil, err
	}

	var nickname, avatar string
	if err := tx.QueryRowContext(ctx,
		`SELECT nickname, COALESCE(avatar_url, '') FROM users WHERE id = $1`, userID,
	).Scan(&nickname, &avatar); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Diary{
		ID:        id,
		Title:     req.Title,
		Content:   req.Content,
		Mood:      req.Mood,
		IsPublic:  isPublic,
		Author:    &Author{ID: userID, Nickname: nickname, AvatarURL: avatar},
		CreatedAt: created,
	}, nil
}

// ListComments 返回某日记下的评论（含软删除占位），按 created_at ASC 分页。
// 私密日记仅作者本人（viewerID==作者）可查看评论；其他访问者视为不存在。
// 日记不存在或无权查看返回 ErrDiaryNotFound。
func (r *Repository) ListComments(ctx context.Context, diaryID, viewerID int64, page, pageSize int) ([]*Comment, int64, error) {
	var visible bool
	if err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM diaries WHERE id = $1 AND (is_public = TRUE OR user_id = $2))`,
		diaryID, viewerID).Scan(&visible); err != nil {
		return nil, 0, err
	}
	if !visible {
		return nil, 0, ErrDiaryNotFound
	}

	var total int64
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM diary_comments WHERE diary_id = $1`, diaryID).Scan(&total); err != nil {
		return nil, 0, err
	}

	const q = `SELECT dc.id, dc.diary_id, dc.content, dc.is_anonymous,
		dc.user_id, u.nickname, COALESCE(u.avatar_url, ''), dc.created_at, dc.deleted_at
		FROM diary_comments dc JOIN users u ON u.id = dc.user_id
		WHERE dc.diary_id = $1
		ORDER BY dc.created_at ASC, dc.id ASC
		LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, q, diaryID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*Comment
	for rows.Next() {
		cm, err := scanComment(rows)
		if err != nil {
			return nil, 0, err
		}
		cm.mask()
		list = append(list, cm)
	}
	return list, total, rows.Err()
}

// CreateComment 校验日记对该用户可见后插入一条评论，返回完整记录（脱敏后）。
// 私密日记仅作者本人可评论；不可见或不存在返回 ErrDiaryNotFound。
func (r *Repository) CreateComment(ctx context.Context, diaryID, userID int64, req CommentRequest) (*Comment, error) {
	var visible bool
	if err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM diaries WHERE id = $1 AND (is_public = TRUE OR user_id = $2))`,
		diaryID, userID).Scan(&visible); err != nil {
		return nil, err
	}
	if !visible {
		return nil, ErrDiaryNotFound
	}

	const insert = `INSERT INTO diary_comments (diary_id, user_id, content, is_anonymous)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	var id int64
	var created time.Time
	if err := r.db.QueryRowContext(ctx, insert,
		diaryID, userID, req.Content, req.IsAnonymous,
	).Scan(&id, &created); err != nil {
		return nil, err
	}

	var nickname, avatar string
	if err := r.db.QueryRowContext(ctx,
		`SELECT nickname, COALESCE(avatar_url, '') FROM users WHERE id = $1`, userID,
	).Scan(&nickname, &avatar); err != nil {
		return nil, err
	}

	cm := &Comment{
		ID:          id,
		DiaryID:     diaryID,
		Content:     req.Content,
		IsAnonymous: req.IsAnonymous,
		Author:      &Author{ID: userID, Nickname: nickname, AvatarURL: avatar},
		CreatedAt:   created,
	}
	cm.mask()
	return cm, nil
}

// GetStreak 返回某用户的连续打卡天数与累计打卡天数。
// 连续天数为包含今天或昨天的最新连续区间长度；若最近一次打卡早于昨天则为 0。
// 用户不存在返回 ErrUserNotFound。
func (r *Repository) GetStreak(ctx context.Context, userID int64) (*Streak, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, userID).Scan(&exists); err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(DISTINCT DATE(created_at)) FROM diaries WHERE user_id = $1`,
		userID).Scan(&total); err != nil {
		return nil, err
	}

	// 缺口与孤岛法（gaps and islands）：将连续日期归为同一区间，
	// 取包含今天/昨天的区间长度作为当前连续天数。
	const q = `
		WITH days AS (
			SELECT DATE(created_at) AS d FROM diaries WHERE user_id = $1 GROUP BY DATE(created_at)
		),
		islands AS (
			SELECT d, d - (ROW_NUMBER() OVER (ORDER BY d))::int AS grp FROM days
		),
		streaks AS (
			SELECT COUNT(*) AS len, MAX(d) AS last_day FROM islands GROUP BY grp
		)
		SELECT COALESCE(MAX(len), 0) FROM streaks WHERE last_day >= CURRENT_DATE - 1`
	var current int
	if err := r.db.QueryRowContext(ctx, q, userID).Scan(&current); err != nil {
		return nil, err
	}

	return &Streak{UserID: userID, CurrentDays: current, TotalDays: total}, nil
}

// RankingStreak 返回连续打卡牛马榜（当前连续天数 DESC）。
func (r *Repository) RankingStreak(ctx context.Context, limit int) ([]*RankStreak, error) {
	const q = `
		WITH days AS (
			SELECT user_id, DATE(created_at) AS d FROM diaries GROUP BY user_id, DATE(created_at)
		),
		islands AS (
			SELECT user_id, d,
				d - (ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY d))::int AS grp
			FROM days
		),
		streaks AS (
			SELECT user_id, COUNT(*) AS len, MAX(d) AS last_day
			FROM islands GROUP BY user_id, grp
		),
		current AS (
			SELECT user_id, len FROM streaks WHERE last_day >= CURRENT_DATE - 1
		)
		SELECT u.id, u.nickname, COALESCE(u.avatar_url, ''), u.level, c.len
		FROM current c JOIN users u ON u.id = c.user_id
		ORDER BY c.len DESC, u.id ASC
		LIMIT $1`
	rows, err := r.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*RankStreak
	for rows.Next() {
		var rs RankStreak
		if err := rows.Scan(&rs.UserID, &rs.Nickname, &rs.AvatarURL, &rs.Level, &rs.Days); err != nil {
			return nil, err
		}
		list = append(list, &rs)
	}
	return list, rows.Err()
}

// rowScanner 抽象 *sql.Row / *sql.Rows 的 Scan 能力。
type rowScanner interface {
	Scan(dest ...any) error
}

func scanDiary(s rowScanner) (*Diary, error) {
	var (
		d         Diary
		authorID  int64
		nickname  string
		avatarURL string
	)
	if err := s.Scan(
		&d.ID, &d.Title, &d.Content, &d.Mood, &d.IsPublic, &d.CommentCount,
		&authorID, &nickname, &avatarURL, &d.CreatedAt,
	); err != nil {
		return nil, err
	}
	d.Author = &Author{ID: authorID, Nickname: nickname, AvatarURL: avatarURL}
	return &d, nil
}

func scanComment(s rowScanner) (*Comment, error) {
	var (
		cm        Comment
		authorID  int64
		nickname  string
		avatarURL string
		deletedAt sql.NullTime
	)
	if err := s.Scan(
		&cm.ID, &cm.DiaryID, &cm.Content, &cm.IsAnonymous,
		&authorID, &nickname, &avatarURL, &cm.CreatedAt, &deletedAt,
	); err != nil {
		return nil, err
	}
	if deletedAt.Valid {
		cm.IsDeleted = true
		cm.DeletedAt = &deletedAt.Time
	}
	cm.Author = &Author{ID: authorID, Nickname: nickname, AvatarURL: avatarURL}
	return &cm, nil
}
