package domain

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Phone        string     `json:"phone" gorm:"type:varchar(20);uniqueIndex;not null"`
	Nickname     string     `json:"nickname" gorm:"type:varchar(50);not null"`
	AvatarURL    string     `json:"avatar_url" gorm:"type:varchar(500)"`
	BirthDate    *time.Time `json:"birth_date" gorm:"type:date"`
	Height       *int       `json:"height" gorm:"type:smallint"`
	Weight       *int       `json:"weight" gorm:"type:smallint"`
	Role         string     `json:"role" gorm:"type:varchar(10)"`
	Bio          string     `json:"bio" gorm:"type:text"`
	LocationLat  *float64   `json:"location_lat" gorm:"type:decimal(10,8)"`
	LocationLng  *float64   `json:"location_lng" gorm:"type:decimal(11,8)"`
	LocationName string     `json:"location_name" gorm:"type:varchar(200)"`
	IsVerified   bool       `json:"is_verified" gorm:"default:false"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	LastOnlineAt *time.Time `json:"last_online_at" gorm:"index"`
	Timestamps
}

func (User) TableName() string {
	return "users"
}

type UserTag struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;index"`
	Tag       string    `json:"tag" gorm:"type:varchar(50);not null;index"`
	Timestamps
}

func (UserTag) TableName() string {
	return "user_tags"
}

type UserImage struct {
	ID         string `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     string `json:"user_id" gorm:"type:uuid;not null;index"`
	ImageURL   string `json:"image_url" gorm:"type:varchar(500);not null"`
	OrderIndex int    `json:"order_index" gorm:"type:smallint;default:0"`
	Timestamps
}

func (UserImage) TableName() string {
	return "user_images"
}
