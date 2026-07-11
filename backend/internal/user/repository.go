package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrUserNotFound 表示用户不存在。
var ErrUserNotFound = errors.New("user not found")

// Repository 封装 users 表的数据库访问（用户模块）。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建 user 仓储。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// userColumns 是查询用户时选取的列，顺序与 scanUser 一致（不含 password）。
const userColumns = `id, email, nickname,
	COALESCE(avatar_url, ''), COALESCE(bio, ''), COALESCE(industry, ''), COALESCE(job_title, ''),
	work_years, level, empathy_received, empathy_given, total_check_in_days, is_anonymous,
	created_at, updated_at`

// GetByID 按 ID 查询用户。不存在返回 ErrUserNotFound。
func (r *Repository) GetByID(ctx context.Context, id int64) (*User, error) {
	const q = `SELECT ` + userColumns + ` FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, q, id)
	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// IDsByNicknames 返回昵称精确匹配 names 中任一项的用户 ID（去重）。
// 用于 @提及解析：昵称非唯一，同名多人会各自返回一次。names 为空时返回空切片。
func (r *Repository) IDsByNicknames(ctx context.Context, names []string) ([]int64, error) {
	if len(names) == 0 {
		return nil, nil
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id FROM users WHERE nickname = ANY($1)`, names)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// Update 按请求中提供的字段更新用户资料并返回更新后的记录。
// 未提供的字段保持不变；总是刷新 updated_at。不存在返回 ErrUserNotFound。
func (r *Repository) Update(ctx context.Context, id int64, req UpdateRequest) (*User, error) {
	set := []string{"updated_at = NOW()"}
	args := []any{}
	n := 0
	add := func(col string, val any) {
		n++
		set = append(set, fmt.Sprintf("%s = $%d", col, n))
		args = append(args, val)
	}

	if req.Nickname != nil {
		add("nickname", *req.Nickname)
	}
	if req.AvatarURL != nil {
		add("avatar_url", *req.AvatarURL)
	}
	if req.Bio != nil {
		add("bio", *req.Bio)
	}
	if req.Industry != nil {
		add("industry", *req.Industry)
	}
	if req.JobTitle != nil {
		add("job_title", *req.JobTitle)
	}
	if req.WorkYears != nil {
		add("work_years", *req.WorkYears)
	}
	if req.IsAnonymous != nil {
		add("is_anonymous", *req.IsAnonymous)
	}

	n++
	q := fmt.Sprintf(`UPDATE users SET %s WHERE id = $%d RETURNING %s`,
		strings.Join(set, ", "), n, userColumns)
	args = append(args, id)

	row := r.db.QueryRowContext(ctx, q, args...)
	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// rowScanner 抽象 *sql.Row / *sql.Rows 的 Scan 能力。
type rowScanner interface {
	Scan(dest ...any) error
}

func scanUser(s rowScanner) (*User, error) {
	var u User
	var created, updated time.Time
	if err := s.Scan(
		&u.ID, &u.Email, &u.Nickname,
		&u.AvatarURL, &u.Bio, &u.Industry, &u.JobTitle,
		&u.WorkYears, &u.Level, &u.EmpathyReceived, &u.EmpathyGiven, &u.TotalCheckInDays, &u.IsAnonymous,
		&created, &updated,
	); err != nil {
		return nil, err
	}
	u.CreatedAt = created
	u.UpdatedAt = updated
	return &u, nil
}

// ListPublicDiaries 返回某用户的公开日记（按 created_at DESC 分页）。
// 用户不存在或无公开日记时返回空列表（total=0）。
func (r *Repository) ListPublicDiaries(ctx context.Context, userID int64, page, pageSize int) ([]*PublicDiary, int64, error) {
	var total int64
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM diaries WHERE user_id = $1 AND is_public = TRUE`,
		userID).Scan(&total); err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	const q = `SELECT id, COALESCE(title, ''), content, COALESCE(mood, ''), created_at
		FROM diaries
		WHERE user_id = $1 AND is_public = TRUE
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, q, userID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*PublicDiary
	for rows.Next() {
		var d PublicDiary
		if err := rows.Scan(&d.ID, &d.Title, &d.Content, &d.Mood, &d.CreatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, &d)
	}
	return list, total, rows.Err()
}

// Stats 返回某用户的公开统计聚合。用户不存在返回 ErrUserNotFound。
func (r *Repository) Stats(ctx context.Context, userID int64) (*UserStats, error) {
	s := &UserStats{UserID: userID}
	err := r.db.QueryRowContext(ctx,
		`SELECT empathy_received, empathy_given, total_check_in_days FROM users WHERE id = $1`,
		userID).Scan(&s.EmpathyReceived, &s.EmpathyGiven, &s.TotalCheckInDays)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM messages WHERE user_id = $1 AND deleted_at IS NULL`,
		userID).Scan(&s.MessageCount); err != nil {
		return nil, err
	}
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM diaries WHERE user_id = $1 AND is_public = TRUE`,
		userID).Scan(&s.PublicDiaryCount); err != nil {
		return nil, err
	}
	return s, nil
}
