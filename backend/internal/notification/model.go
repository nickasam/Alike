package notification

import "time"

// 通知类型常量，对齐 notifications.type 取值。
const (
	TypeMention = "mention"
	TypeEmpathy = "empathy"
	TypeReply   = "reply"
	TypeSystem  = "system"
)

// Notification 是 notifications 表的领域模型。
type Notification struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	Content   string    `json:"content,omitempty"`
	RefID     *int64    `json:"ref_id,omitempty"` // 关联的消息/日记 ID
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
