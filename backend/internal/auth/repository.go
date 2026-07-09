package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// ErrEmailConflict 表示邮箱已被注册（唯一约束冲突）。
var ErrEmailConflict = errors.New("email already registered")

// ErrUserNotFound 表示用户不存在。
var ErrUserNotFound = errors.New("user not found")

// pgUniqueViolation 是 PostgreSQL 唯一约束冲突的 SQLSTATE。
const pgUniqueViolation = "23505"

// Repository 封装 users 表的数据库访问。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建 auth 仓储。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// userColumns 是查询用户时选取的完整列，顺序与 scanUser 一致。
const userColumns = `id, email, password, nickname,
	COALESCE(avatar_url, ''), COALESCE(bio, ''), COALESCE(industry, ''), COALESCE(job_title, ''),
	work_years, level, empathy_received, empathy_given, total_check_in_days, is_anonymous,
	created_at, updated_at`

// Create 插入新用户并返回完整记录。邮箱唯一冲突时返回 ErrEmailConflict。
func (r *Repository) Create(ctx context.Context, email, passwordHash, nickname string) (*User, error) {
	const q = `INSERT INTO users (email, password, nickname)
		VALUES ($1, $2, $3)
		RETURNING ` + userColumns

	row := r.db.QueryRowContext(ctx, q, email, passwordHash, nickname)
	u, err := scanUser(row)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, ErrEmailConflict
		}
		return nil, err
	}
	return u, nil
}

// GetByEmail 按邮箱查询用户（CITEXT 大小写不敏感）。不存在返回 ErrUserNotFound。
func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	const q = `SELECT ` + userColumns + ` FROM users WHERE email = $1`
	return r.queryOne(ctx, q, email)
}

// GetByID 按 ID 查询用户。不存在返回 ErrUserNotFound。
func (r *Repository) GetByID(ctx context.Context, id int64) (*User, error) {
	const q = `SELECT ` + userColumns + ` FROM users WHERE id = $1`
	return r.queryOne(ctx, q, id)
}

func (r *Repository) queryOne(ctx context.Context, query string, args ...any) (*User, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// rowScanner 抽象 *sql.Row / *sql.Rows 的 Scan 能力，便于复用。
type rowScanner interface {
	Scan(dest ...any) error
}

func scanUser(s rowScanner) (*User, error) {
	var u User
	var created, updated time.Time
	if err := s.Scan(
		&u.ID, &u.Email, &u.Password, &u.Nickname,
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
