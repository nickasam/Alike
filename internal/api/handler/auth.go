package handler

import (
	"github.com/Alike/internal/api/middleware"
	"github.com/Alike/internal/auth"
	"github.com/Alike/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *auth.Service
}

func New(authService *auth.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// TODO: 实现注册逻辑
	response.Success(c, gin.H{"message": "Register endpoint - to be implemented"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// TODO: 实现登录逻辑
	response.Success(c, gin.H{"message": "Login endpoint - to be implemented"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	
	userID, err := h.authService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "Invalid refresh token")
		return
	}
	
	accessToken, refreshToken, err := h.authService.GenerateToken(userID)
	if err != nil {
		response.InternalError(c, "Failed to generate tokens")
		return
	}
	
	response.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	response.Success(c, gin.H{"user_id": userID, "message": "Logged out"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	response.Success(c, gin.H{"user_id": userID, "message": "User info"})
}
