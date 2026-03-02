package domain

import "time"

type Block struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BlockerID string    `json:"blocker_id" gorm:"type:uuid;not null;index"`
	BlockedID string    `json:"blocked_id" gorm:"type:uuid;not null;index"`
	Reason    string    `json:"reason" gorm:"type:varchar(200)"`
	Timestamps
}

func (Block) TableName() string {
	return "blocks"
}

type View struct {
	ID       string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ViewerID string    `json:"viewer_id" gorm:"type:uuid;not null;index"`
	ViewedID string    `json:"viewed_id" gorm:"type:uuid;not null;index"`
	Timestamps
}

func (View) TableName() string {
	return "views"
}

type Notification struct {
	ID      string                 `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID  string                 `json:"user_id" gorm:"type:uuid;not null;index"`
	Type    string                 `json:"type" gorm:"type:varchar(50);not null;index"`
	Title   string                 `json:"title" gorm:"type:varchar(200)"`
	Content string                 `json:"content" gorm:"type:text"`
	Data    map[string]interface{} `json:"data" gorm:"type:jsonb"`
	IsRead  bool                   `json:"is_read" gorm:"default:false;index"`
	ReadAt  *time.Time             `json:"read_at"`
	Timestamps
}

func (Notification) TableName() string {
	return "notifications"
}

type NotificationSettings struct {
	UserID              string    `json:"user_id" gorm:"type:uuid;primary_key"`
	LikeNotification    bool      `json:"like_notification" gorm:"default:true"`
	MatchNotification   bool      `json:"match_notification" gorm:"default:true"`
	MessageNotification bool      `json:"message_notification" gorm:"default:true"`
	ViewNotification    bool      `json:"view_notification" gorm:"default:false"`
	EmailNotification   bool      `json:"email_notification" gorm:"default:false"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (NotificationSettings) TableName() string {
	return "notification_settings"
}

type RefreshToken struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;index"`
	Token     string    `json:"token" gorm:"type:varchar(500);uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"index;not null"`
	Timestamps
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
