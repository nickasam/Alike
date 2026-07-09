// Package jwt 提供 JWT access/refresh token 的签发与验证工具。
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType 区分 access 与 refresh token。
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// ErrInvalidToken 表示 token 无效或已过期。
var ErrInvalidToken = errors.New("invalid or expired token")

// Claims 是 Alike 的自定义 JWT 载荷。
type Claims struct {
	UserID int64     `json:"user_id"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

// Manager 封装签名密钥与过期配置。
type Manager struct {
	secret        []byte
	accessExpire  time.Duration
	refreshExpire time.Duration
}

// NewManager 创建 JWT 管理器。
func NewManager(secret string, accessExpire, refreshExpire time.Duration) *Manager {
	return &Manager{
		secret:        []byte(secret),
		accessExpire:  accessExpire,
		refreshExpire: refreshExpire,
	}
}

// Generate 为指定用户签发某类型 token。
func (m *Manager) Generate(userID int64, t TokenType) (string, error) {
	expire := m.accessExpire
	if t == RefreshToken {
		expire = m.refreshExpire
	}
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Type:   t,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expire)),
			Issuer:    "alike",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// GenerateAccess 签发 access token。
func (m *Manager) GenerateAccess(userID int64) (string, error) {
	return m.Generate(userID, AccessToken)
}

// GenerateRefresh 签发 refresh token。
func (m *Manager) GenerateRefresh(userID int64) (string, error) {
	return m.Generate(userID, RefreshToken)
}

// Parse 校验并解析 token，返回 Claims。
func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
