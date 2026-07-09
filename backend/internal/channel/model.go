package channel

import "time"

// Channel 是 channels 表的领域模型。
type Channel struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Icon        string    `json:"icon"`
	Status      string    `json:"status"`
	MemberCount int       `json:"member_count"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// Member 是 channel_members 表的领域模型，附带成员的公开信息。
type Member struct {
	UserID    int64     `json:"user_id"`
	Nickname  string    `json:"nickname"`
	AvatarURL string    `json:"avatar_url"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
}

// CreateRequest 是创建频道的请求体。
type CreateRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Slug        string `json:"slug" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=2000"`
	Category    string `json:"category" binding:"required,oneof=industry job topic custom"`
	Icon        string `json:"icon" binding:"max=50"`
}
