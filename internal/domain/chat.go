package domain

import "time"

type Chat struct {
	ID            string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MatchID       string     `json:"match_id" gorm:"type:uuid;not null;index"`
	User1ID       string     `json:"user1_id" gorm:"type:uuid;not null;index"`
	User2ID       string     `json:"user2_id" gorm:"type:uuid;not null;index"`
	LastMessageAt *time.Time `json:"last_message_at" gorm:"index"`
	Timestamps
}

func (Chat) TableName() string {
	return "chats"
}

type Message struct {
	ID          string                 `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ChatID      string                 `json:"chat_id" gorm:"type:uuid;not null;index"`
	SenderID    string                 `json:"sender_id" gorm:"type:uuid;not null;index"`
	Content     string                 `json:"content" gorm:"type:text"`
	MessageType string                 `json:"message_type" gorm:"type:varchar(20);default:text"`
	Metadata    map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	IsRead      bool                   `json:"is_read" gorm:"default:false;index"`
	ReadAt      *time.Time             `json:"read_at"`
	Timestamps
}

func (Message) TableName() string {
	return "messages"
}

const (
	MessageTypeText     = "text"
	MessageTypeImage    = "image"
	MessageTypeLocation = "location"
	MessageTypeVoice    = "voice"
)
