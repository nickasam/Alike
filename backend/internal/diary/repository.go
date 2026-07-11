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

// ErrAlreadyEmpathized 表示已对该日记共情过。
var ErrAlreadyEmpathized = errors.New("already empathized")

// ErrNotEmpathized 表示尚未对该日记共情。
var ErrNotEmpathized = errors.New("not empathized")

// ErrSelfEmpathy 表示不能对自己的日记共情。
var ErrSelfEmpathy = errors.New("cannot empathize own diary")

// Repository 封装 diaries / diary_comments 表的数据库访问。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建 diary 仓储。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// diaryCols 返回查询日记时选取的列，顺序与 scanDiary 一致。
// comment_count / empathy_count 用冗余或子查询；empathized 依据 viewer 参数占位符判断当前用户是否已共情。
func diaryCols(viewerParam string) string {
	return `d.id, COALESCE(d.title, ''), d.content, COALESCE(d.mood, ''), d.is_public,
		(SELECT COUNT(*) FROM diary_comments dc WHERE dc.diary_id = d.id AND dc.deleted_at IS NULL) AS comment_count,
		d.empathy_count,
		EXISTS(SELECT 1 FROM diary_empathies e WHERE e.diary_id = d.id AND e.user_id = ` + viewerParam + `) AS empathized,
		d.user_id, u.nickname, COALESCE(u.avatar_url, ''), d.created_at`
}

const diaryFrom = ` FROM diaries d JOIN users u ON u.id = d.user_id`

// ListPublic 返回公开日记流（is_public=true），按 created_at DESC 游标分页。
// before>0 时仅返回更早的记录。viewerID 用于计算当前用户是否已共情（<=0 视为未登录）。
func (r *Repository) ListPublic(ctx context.Context, viewerID, before int64, limit int) ([]*Diary, bool, error) {
	args := []any{viewerID} // $1 = viewer（empathized 判断）
	where := `WHERE d.is_public = TRUE`
	if before > 0 {
		where += ` AND (d.created_at, d.id) < (SELECT created_at, id FROM diaries WHERE id = $2)`
		args = append(args, before)
	}
	args = append(args, limit+1) // 多取 1 条判断是否还有更多
	q := `SELECT ` + diaryCols("$1") + diaryFrom + ` ` + where +
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
	q := `SELECT ` + diaryCols("$2") + diaryFrom + ` WHERE d.id = $1 AND (d.is_public = TRUE OR d.user_id = $2)`
	d, err := scanDiary(r.db.QueryRowContext(ctx, q, diaryID, viewerID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrDiaryNotFound
	}
	if err != nil {
		return nil, err
	}
	return d, nil
}

// CreateEmpathy 事务化地为日记添加一次共情：
// INSERT diary_empathies + diaries.empathy_count+1 + 作者 empathy_received+1 + 当前用户 empathy_given+1。
// 日记不存在返回 ErrDiaryNotFound；重复共情返回 ErrAlreadyEmpathized；对自己日记共情返回 ErrSelfEmpathy。
// 返回日记最新共情计数。
func (r *Repository) CreateEmpathy(ctx context.Context, diaryID, userID int64) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck // 提交成功后回滚为 no-op

	var authorID int64
	err = tx.QueryRowContext(ctx, `SELECT user_id FROM diaries WHERE id = $1`, diaryID).Scan(&authorID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrDiaryNotFound
	}
	if err != nil {
		return 0, err
	}
	if authorID == userID {
		return 0, ErrSelfEmpathy
	}

	res, err := tx.ExecContext(ctx,
		`INSERT INTO diary_empathies (diary_id, user_id) VALUES ($1, $2) ON CONFLICT (diary_id, user_id) DO NOTHING`,
		diaryID, userID)
	if err != nil {
		return 0, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return 0, ErrAlreadyEmpathized
	}

	var count int64
	if err := tx.QueryRowContext(ctx,
		`UPDATE diaries SET empathy_count = empathy_count + 1 WHERE id = $1 RETURNING empathy_count`,
		diaryID).Scan(&count); err != nil {
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

// DeleteEmpathy 事务化地取消一次日记共情，做 CreateEmpathy 的反向操作。计数不降到 0 以下。
// 日记不存在返回 ErrDiaryNotFound；未共情返回 ErrNotEmpathized。返回日记最新共情计数。
func (r *Repository) DeleteEmpathy(ctx context.Context, diaryID, userID int64) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() //nolint:errcheck // 提交成功后回滚为 no-op

	var authorID int64
	err = tx.QueryRowContext(ctx, `SELECT user_id FROM diaries WHERE id = $1`, diaryID).Scan(&authorID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrDiaryNotFound
	}
	if err != nil {
		return 0, err
	}

	res, err := tx.ExecContext(ctx,
		`DELETE FROM diary_empathies WHERE diary_id = $1 AND user_id = $2`, diaryID, userID)
	if err != nil {
		return 0, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return 0, ErrNotEmpathized
	}

	var count int64
	if err := tx.QueryRowContext(ctx,
		`UPDATE diaries SET empathy_count = GREATEST(empathy_count - 1, 0) WHERE id = $1 RETURNING empathy_count`,
		diaryID).Scan(&count); err != nil {
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

// MyStreakRank 返回某用户在连续打卡榜中的精确名次与当前连续天数。
// 名次算法与 RankingStreak 一致（仅统计 last_day >= 昨天 的存活连续段）；
// 若该用户当前无存活连续段（未上榜）返回 rank=0。
func (r *Repository) MyStreakRank(ctx context.Context, userID int64) (rank, days int, err error) {
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
		),
		me AS (
			SELECT len FROM current WHERE user_id = $1 ORDER BY len DESC LIMIT 1
		)
		SELECT
			COALESCE((SELECT len FROM me), 0) AS my_len,
			CASE WHEN (SELECT len FROM me) IS NULL THEN 0
			ELSE (SELECT COUNT(*) FROM current c
			      WHERE c.len > (SELECT len FROM me)
			         OR (c.len = (SELECT len FROM me) AND c.user_id < $1)) + 1
			END AS my_rank`
	if err = r.db.QueryRowContext(ctx, q, userID).Scan(&days, &rank); err != nil {
		return 0, 0, err
	}
	return rank, days, nil
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
		&d.EmpathyCount, &d.Empathized,
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
