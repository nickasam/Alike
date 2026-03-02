package domain

import "time"

type Like struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LikerID   string    `json:"liker_id" gorm:"type:uuid;not null;index"`
	LikedID   string    `json:"liked_id" gorm:"type:uuid;not null;index"`
	Timestamps
}

func (Like) TableName() string {
	return "likes"
}

type Match struct {
	ID            string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	User1ID       string     `json:"user1_id" gorm:"type:uuid;not null;uniqueIndex:idx_match_users"`
	User2ID       string     `json:"user2_id" gorm:"type:uuid;not null;uniqueIndex:idx_match_users"`
	IsActive      bool       `json:"is_active" gorm:"default:true;index"`
	LastMessageAt *time.Time `json:"last_message_at" gorm:"index"`
	UnreadCount1  int        `json:"unread_count_1" gorm:"default:0"`
	UnreadCount2  int        `json:"unread_count_2" gorm:"default:0"`
	Timestamps
}

func (Match) TableName() string {
	return "matches"
}
