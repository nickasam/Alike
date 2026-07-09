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
