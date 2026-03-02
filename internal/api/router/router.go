package router

import (
	"github.com/Alike/internal/api/handler"
	"github.com/Alike/internal/api/middleware"
	"github.com/Alike/internal/auth"
	"github.com/Alike/internal/chat"
	"github.com/Alike/internal/match"
	"github.com/Alike/internal/notification"
	"github.com/Alike/internal/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, authService *auth.Service, db *gorm.DB) {
	r.Use(middleware.CORSMiddleware())
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	
	v1 := r.Group("/api/v1")
	{
		authHandler := handler.New(authService)
		userService := user.NewService(db, authService)
		userHandler := handler.NewUserHandler(userService)
		matchService := match.NewService(db)
		matchHandler := handler.NewMatchHandler(matchService)
		chatService := chat.NewService(db)
		chatHandler := handler.NewChatHandler(chatService)
		notificationService := notification.NewService(db)
		notificationHandler := handler.NewNotificationHandler(notificationService)
		
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", middleware.AuthMiddleware(authService), authHandler.Me)
		}
		
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetMe)
				users.PUT("/me", userHandler.UpdateMe)
				users.GET("/nearby", userHandler.GetNearby)
			}
			
			matches := protected.Group("/matches")
			{
				matches.GET("", matchHandler.ListMatches)
				matches.GET("/:id", matchHandler.GetMatch)
				matches.POST("/:id/like", matchHandler.LikeUser)
			}
			
			chats := protected.Group("/chats")
			{
				chats.GET("", chatHandler.ListChats)
				chats.GET("/:id", chatHandler.GetChat)
				chats.GET("/:id/messages", chatHandler.GetMessages)
				chats.POST("/:id/messages", chatHandler.SendMessage)
			}
			
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", notificationHandler.ListNotifications)
				notifications.POST("/:id/read", notificationHandler.MarkAsRead)
				notifications.POST("/read-all", notificationHandler.MarkAllAsRead)
			}
		}
	}
}
