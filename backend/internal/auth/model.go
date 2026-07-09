package auth

import "time"

// User 是 users 表的领域模型。password 字段仅用于内部校验，
// 通过 json:"-" 保证任何序列化响应都不会泄露密码 hash。
type User struct {
	ID               int64     `json:"id"`
	Email            string    `json:"email"`
	Password         string    `json:"-"`
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

// RegisterRequest 是注册请求体。industry/job_title/work_years 为选填。
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Nickname  string `json:"nickname" binding:"required,min=1,max=100"`
	Industry  string `json:"industry" binding:"max=100"`
	JobTitle  string `json:"job_title" binding:"max=100"`
	WorkYears int    `json:"work_years" binding:"min=0,max=60"`
}

// LoginRequest 是登录请求体。
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RefreshRequest 是刷新 token 请求体。
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// TokenPair 是签发给客户端的 access/refresh token 对。
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AuthResponse 是登录/注册成功后的响应，包含用户信息与 token。
// data: { "tokens": {...}, "user": {...} }
type AuthResponse struct {
	User   *User      `json:"user"`
	Tokens *TokenPair `json:"tokens"`
}

// TokensResponse 是仅返回 token 的响应（如刷新 token）。
// data: { "tokens": { "access_token": "...", "refresh_token": "..." } }
type TokensResponse struct {
	Tokens *TokenPair `json:"tokens"`
}
