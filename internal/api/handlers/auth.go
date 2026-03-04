package handlers

import (
	"net/http"
	"time"

	"github.com/Alike/internal/domain"
	"github.com/Alike/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userRepo  *repository.UserRepository
	jwtSecret []byte
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(userRepo *repository.UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "请求参数错误",
			},
		})
		return
	}

	// 检查用户是否已存在
	existingUser, err := h.userRepo.GetByPhone(req.Phone)
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "DATABASE_ERROR",
				"message": "数据库查询失败",
			},
		})
		return
	}

	if existingUser != nil && existingUser.ID != "" {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "USER_EXISTS",
				"message": "用户已存在",
			},
		})
		return
	}

	// 创建新用户
	now := time.Now()
	user := &domain.User{
		ID:       uuid.New().String(),
		Phone:    req.Phone,
		Nickname: req.Nickname,
		IsActive: true,
		Timestamps: domain.Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "CREATE_FAILED",
				"message": "创建用户失败",
			},
		})
		return
	}

	// 生成 Token
	accessToken, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "TOKEN_ERROR",
				"message": "生成 Token 失败",
			},
		})
		return
	}

	refreshToken, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "TOKEN_ERROR",
				"message": "生成 Token 失败",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": gin.H{
			"user": gin.H{
				"id":       user.ID,
				"phone":    user.Phone,
				"nickname": user.Nickname,
			},
			"tokens": gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		},
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "请求参数错误",
			},
		})
		return
	}

	user, err := h.userRepo.GetByPhone(req.Phone)
	if err != nil || user == nil || user.ID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "USER_NOT_FOUND",
				"message": "用户不存在",
			},
		})
		return
	}

	// TODO: 验证密码（目前代码中没有密码字段）
	// if !user.VerifyPassword(req.Password) {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"success": false,
	// 		"error": gin.H{
	// 			"code":    "INVALID_PASSWORD",
	// 			"message": "密码错误",
	// 		},
	// 	})
	// 	return
	// }

	accessToken, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "TOKEN_ERROR",
				"message": "生成 Token 失败",
			},
		})
		return
	}

	refreshToken, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "TOKEN_ERROR",
				"message": "生成 Token 失败",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"user": gin.H{
				"id":       user.ID,
				"phone":    user.Phone,
				"nickname": user.Nickname,
			},
			"tokens": gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		},
	})
}

// GetCurrentUser 获取当前用户信息
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "USER_NOT_FOUND",
				"message": "用户不存在",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// ValidateToken 验证JWT Token（公开方法，供中间件使用）
func (h *AuthHandler) ValidateToken(tokenString string) (string, error) {
	return h.validateToken(tokenString)
}

// generateToken 生成JWT Token
func (h *AuthHandler) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.jwtSecret)
}

// validateToken 验证JWT Token
func (h *AuthHandler) validateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return h.jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)
		return userID, nil
	}

	return "", err
}
