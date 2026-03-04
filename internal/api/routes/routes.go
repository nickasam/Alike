package routes

import (
	"time"

	"github.com/Alike/internal/api/handlers"
	"github.com/Alike/internal/api/middleware"
	"github.com/Alike/internal/domain"
	"github.com/Alike/internal/repository"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine, userRepo *repository.UserRepository, globalChatRepo *repository.GlobalChatRepository, jwtSecret string) {
	// 初始化 handlers
	authHandler := handlers.NewAuthHandler(userRepo, jwtSecret)

	// 应用全局中间件
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RecoveryMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "alike-api",
			"version": "1.0.0",
		})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 认证路由（不需要认证）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(authHandler), authHandler.GetCurrentUser)
		}

		// 用户路由（需要认证）
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(authHandler))
		{
			users.GET("/nearby", GetNearbyUsers(userRepo))
		}

		// 全局聊天路由（需要认证）
		global := v1.Group("/global")
		global.Use(middleware.AuthMiddleware(authHandler))
		{
			global.GET("/room", GetGlobalRoom(globalChatRepo))
			global.GET("/messages", GetGlobalMessages(globalChatRepo))
			global.POST("/messages", SendGlobalMessage(userRepo, globalChatRepo))
		}
	}
}

// GetNearbyUsers 获取附近用户
func GetNearbyUsers(userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := userRepo.List(100, 0)
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "DATABASE_ERROR",
					"message": "获取用户列表失败",
				},
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"data":    users,
		})
	}
}

// GetGlobalRoom 获取全局聊天室
func GetGlobalRoom(globalChatRepo *repository.GlobalChatRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		room, err := globalChatRepo.GetRoom("global")
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "DATABASE_ERROR",
					"message": "获取聊天室失败",
				},
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"data":    room,
		})
	}
}

// GetGlobalMessages 获取全局消息
func GetGlobalMessages(globalChatRepo *repository.GlobalChatRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		messages, err := globalChatRepo.GetMessages("global", 100)
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "DATABASE_ERROR",
					"message": "获取消息失败",
				},
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"data":    messages,
		})
	}
}

// SendGlobalMessage 发送全局消息
func SendGlobalMessage(userRepo *repository.UserRepository, globalChatRepo *repository.GlobalChatRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		user, err := userRepo.GetByID(userID)
		if err != nil {
			c.JSON(404, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "USER_NOT_FOUND",
					"message": "用户不存在",
				},
			})
			return
		}

		var req struct {
			Content string `json:"content" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_REQUEST",
					"message": "请求参数错误",
				},
			})
			return
		}

		message := &domain.GlobalMessage{
			ID:        time.Now().Format("20060102150405") + "000",
			RoomID:    "global",
			UserID:    userID,
			Username:  user.Nickname,
			Content:   req.Content,
			CreatedAt: time.Now(),
		}

		if err := globalChatRepo.CreateMessage(message); err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "CREATE_FAILED",
					"message": "发送消息失败",
				},
			})
			return
		}

		c.JSON(201, gin.H{
			"success": true,
			"data":    message,
		})
	}
}
