package user

import "time"

// User 是 users 表的领域模型（用户模块视角）。
// email 属隐私字段，仅在本人视角序列化；公开主页使用 PublicUser。
type User struct {
	ID               int64     `json:"id"`
	Email            string    `json:"email"`
	Nickname         string    `json:"nickname"`
	AvatarURL        string    `json:"avatar_url"`
	Bio              string    `json:"bio"`
	Industry         string    `json:"industry"`
	JobTitle         string    `json:"job_title"`
	WorkYears        int       `json:"work_years"`
	Level            int       `json:"level"`
	EmpathyReceived  int       `json:"empathy_received"`
	EmpathyGiven     int       `json:"empathy_given"`
	TotalCheckInDays int       `json:"total_check_in_days"`
	IsAnonymous      bool      `json:"is_anonymous"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// PublicUser 是公开主页视图，剔除 email 等隐私字段。
type PublicUser struct {
	ID               int64     `json:"id"`
	Nickname         string    `json:"nickname"`
	AvatarURL        string    `json:"avatar_url"`
	Bio              string    `json:"bio"`
	Industry         string    `json:"industry"`
	JobTitle         string    `json:"job_title"`
	WorkYears        int       `json:"work_years"`
	Level            int       `json:"level"`
	EmpathyReceived  int       `json:"empathy_received"`
	EmpathyGiven     int       `json:"empathy_given"`
	TotalCheckInDays int       `json:"total_check_in_days"`
	CreatedAt        time.Time `json:"created_at"`
}

// Public 将完整用户转为公开视图。
func (u *User) Public() *PublicUser {
	return &PublicUser{
		ID:               u.ID,
		Nickname:         u.Nickname,
		AvatarURL:        u.AvatarURL,
		Bio:              u.Bio,
		Industry:         u.Industry,
		JobTitle:         u.JobTitle,
		WorkYears:        u.WorkYears,
		Level:            u.Level,
		EmpathyReceived:  u.EmpathyReceived,
		EmpathyGiven:     u.EmpathyGiven,
		TotalCheckInDays: u.TotalCheckInDays,
		CreatedAt:        u.CreatedAt,
	}
}

// UpdateRequest 是更新资料请求体。所有字段可选，仅更新提供的字段。
// 使用指针区分「未提供」与「提供空值」。
type UpdateRequest struct {
	Nickname    *string `json:"nickname" binding:"omitempty,min=1,max=100"`
	AvatarURL   *string `json:"avatar_url" binding:"omitempty,max=500"`
	Bio         *string `json:"bio" binding:"omitempty,max=200"`
	Industry    *string `json:"industry" binding:"omitempty,max=100"`
	JobTitle    *string `json:"job_title" binding:"omitempty,max=100"`
	WorkYears   *int    `json:"work_years" binding:"omitempty,min=0"`
	IsAnonymous *bool   `json:"is_anonymous"`
}
