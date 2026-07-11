package diary

import "time"

// deletedPlaceholder 是软删除评论对外返回的占位内容。
const deletedPlaceholder = "该评论已删除"

// Author 是日记/评论作者的公开信息。匿名评论不返回该字段。
type Author struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

// Diary 是 diaries 表的领域模型（含作者信息）。
type Diary struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title,omitempty"`
	Content      string    `json:"content"`
	Mood         string    `json:"mood,omitempty"`
	IsPublic     bool      `json:"is_public"`
	CommentCount int       `json:"comment_count"`
	EmpathyCount int       `json:"empathy_count"`
	Empathized   bool      `json:"empathized"` // 请求方是否已对本日记共情
	Author       *Author   `json:"author,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// Comment 是 diary_comments 表的领域模型（含作者信息与匿名/软删除标记）。
type Comment struct {
	ID          int64      `json:"id"`
	DiaryID     int64      `json:"diary_id"`
	Content     string     `json:"content"`
	IsAnonymous bool       `json:"is_anonymous"`
	IsDeleted   bool       `json:"is_deleted"`
	Author      *Author    `json:"author,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// mask 依据匿名与软删除标记，隐藏对外不应暴露的字段。
func (c *Comment) mask() {
	if c.IsDeleted {
		c.Content = deletedPlaceholder
		c.Author = nil
		return
	}
	if c.IsAnonymous {
		c.Author = nil
	}
}

// Streak 是某用户的打卡统计：连续天数与累计天数。
type Streak struct {
	UserID      int64 `json:"user_id"`
	CurrentDays int   `json:"current_days"` // 连续打卡天数（截至今天）
	TotalDays   int   `json:"total_days"`   // 累计打卡天数（去重日期数）
}

// RankStreak 是连续打卡牛马榜的条目。
type RankStreak struct {
	UserID    int64  `json:"user_id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Level     int    `json:"level"`
	Days      int    `json:"days"` // 连续打卡天数
}

// CreateRequest 是写日记的请求体。
type CreateRequest struct {
	Title    string `json:"title" binding:"max=200"`
	Content  string `json:"content" binding:"required,min=1,max=10000"`
	Mood     string `json:"mood" binding:"max=50"`
	IsPublic *bool  `json:"is_public"` // 缺省视为公开
}

// CommentRequest 是发表日记评论的请求体。
type CommentRequest struct {
	Content     string `json:"content" binding:"required,min=1,max=2000"`
	IsAnonymous bool   `json:"is_anonymous"`
}
