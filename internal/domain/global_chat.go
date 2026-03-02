package domain

import "time"

// GlobalChatRoom 全局聊天室
type GlobalChatRoom struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100"`
	Description string    `json:"description" gorm:"size:500"`
	MaxMembers  int       `json:"max_members" gorm:"default:1000"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GlobalMessage 全局聊天消息
type GlobalMessage struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	RoomID    string    `json:"room_id" gorm:"index"`
	UserID    string    `json:"user_id" gorm:"index"`
	Username  string    `json:"username" gorm:"size:100"`
	Content   string    `json:"content" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
}
