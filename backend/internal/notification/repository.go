package notification

import (
	"context"
	"database/sql"
)

// Repository 封装 notifications 表的数据库访问。
type Repository struct {
	db *sql.DB
}

// NewRepository 创建 notification 仓储。
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// List 返回某用户的通知列表，按 created_at DESC 分页，并附带未读总数。
func (r *Repository) List(ctx context.Context, userID int64, page, pageSize int) ([]*Notification, int64, int64, error) {
	var total int64
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM notifications WHERE user_id = $1`, userID).Scan(&total); err != nil {
		return nil, 0, 0, err
	}

	var unread int64
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`,
		userID).Scan(&unread); err != nil {
		return nil, 0, 0, err
	}

	const q = `SELECT id, type, COALESCE(content, ''), ref_id, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, q, userID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var list []*Notification
	for rows.Next() {
		var (
			n     Notification
			refID sql.NullInt64
		)
		if err := rows.Scan(&n.ID, &n.Type, &n.Content, &refID, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, 0, 0, err
		}
		if refID.Valid {
			n.RefID = &refID.Int64
		}
		list = append(list, &n)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, 0, err
	}
	return list, total, unread, nil
}

// MarkRead 将某用户名下的单条通知标记为已读，返回是否命中记录。
// 仅能标记自己的通知（user_id 约束），避免越权。
func (r *Repository) MarkRead(ctx context.Context, notificationID, userID int64) (bool, error) {
	res, err := r.db.ExecContext(ctx,
		`UPDATE notifications SET is_read = TRUE WHERE id = $1 AND user_id = $2 AND is_read = FALSE`,
		notificationID, userID)
	if err != nil {
		return false, err
	}
	// 记录存在但已读也视为成功（幂等），仅当记录不属于该用户时才算未命中。
	if n, _ := res.RowsAffected(); n > 0 {
		return true, nil
	}
	var exists bool
	if err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM notifications WHERE id = $1 AND user_id = $2)`,
		notificationID, userID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// MarkAllRead 将某用户的全部未读通知标记为已读，返回受影响条数。
func (r *Repository) MarkAllRead(ctx context.Context, userID int64) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`UPDATE notifications SET is_read = TRUE WHERE user_id = $1 AND is_read = FALSE`, userID)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}
