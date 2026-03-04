package router

import (
	"github.com/Alike/internal/api/handler"
	"github.com/Alike/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func Setup(authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	matchHandler *handler.MatchHandler,
	chatHandler *handler.ChatHandler,
	notificationHandler *handler.NotificationHandler,
	globalChatHandler *handler.GlobalChatHandler) *gin.Engine {

	r := gin.Default()

	// CORS 中间件 - 必须在路由之前设置
	r.Use(middleware.CORSMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 认证
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", middleware.Auth(), authHandler.Me)
		}

		// 用户
		users := v1.Group("/users")
		users.Use(middleware.Auth())
		{
			users.GET("/me", userHandler.GetMe)
			users.PUT("/me", userHandler.UpdateMe)
			users.GET("/nearby", userHandler.GetNearby)
		}

		// 匹配
		matches := v1.Group("/matches")
		matches.Use(middleware.Auth())
		{
			matches.GET("", matchHandler.ListMatches)
			matches.GET("/:id", matchHandler.GetMatch)
			matches.POST("/:id/like", matchHandler.LikeUser)
		}

		// 聊天
		chats := v1.Group("/chats")
		chats.Use(middleware.Auth())
		{
			chats.GET("", chatHandler.ListChats)
			chats.GET("/:id", chatHandler.GetChat)
			chats.GET("/:id/messages", chatHandler.GetMessages)
			chats.POST("/:id/messages", chatHandler.SendMessage)
		}

		// 通知
		notifications := v1.Group("/notifications")
		notifications.Use(middleware.Auth())
		{
			notifications.GET("", notificationHandler.ListNotifications)
			notifications.POST("/:id/read", notificationHandler.MarkAsRead)
			notifications.POST("/read-all", notificationHandler.MarkAllAsRead)
		}

		// 全局聊天室
		global := v1.Group("/global")
		global.Use(middleware.Auth())
		{
			global.GET("/room", globalChatHandler.GetGlobalRoom)
			global.GET("/messages", globalChatHandler.GetGlobalMessages)
			global.POST("/messages", globalChatHandler.SendGlobalMessage)
			global.POST("/join", globalChatHandler.JoinGlobalRoom)
		}
	}

	return r
}
