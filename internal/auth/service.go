package auth

import (
	"fmt"
	"time"

	"github.com/Alike/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	config *config.JWTConfig
}

func New(cfg *config.JWTConfig) *Service {
	return &Service{config: cfg}
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *Service) GenerateToken(userID string) (string, string, error) {
	accessClaims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}
	
	refreshClaims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	
	return accessToken, refreshToken, nil
}

func (s *Service) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Secret), nil
	})
	
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}
	
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}
	
	return "", fmt.Errorf("invalid token")
}

func (s *Service) ValidateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Secret), nil
	})
	
	if err != nil {
		return "", fmt.Errorf("failed to parse refresh token: %w", err)
	}
	
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}
	
	return "", fmt.Errorf("invalid refresh token")
}
