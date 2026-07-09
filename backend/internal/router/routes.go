package router

import (
	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/internal/auth"
	"github.com/Alike/backend/internal/channel"
	"github.com/Alike/backend/internal/diary"
	"github.com/Alike/backend/internal/empathy"
	"github.com/Alike/backend/internal/message"
	"github.com/Alike/backend/internal/middleware"
	"github.com/Alike/backend/internal/notification"
	"github.com/Alike/backend/internal/search"
	"github.com/Alike/backend/internal/user"
	"github.com/Alike/backend/internal/ws"
)

// registerRoutes 注册各业务模块的路由分组骨架。
// 阶段一仅搭建分组与占位 handler，具体业务在后续阶段填充。
func registerRoutes(api *gin.RouterGroup, deps *Deps) {
	// 需鉴权的中间件
	authMW := middleware.Auth(deps.JWT)

	// 认证：注入 DB / JWT 依赖。
	authHandler := auth.NewHandler(auth.NewRepository(deps.DB), deps.JWT)
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.GET("/me", authMW, authHandler.Me)
	}

	// 用户：注入 DB 依赖。
	userHandler := user.NewHandler(user.NewRepository(deps.DB))
	users := api.Group("/users")
	{
		users.GET("/:id", userHandler.Get)
		users.PUT("/:id", authMW, userHandler.Update)
		users.GET("/:id/diaries", user.DiariesHandler)
		users.GET("/:id/stats", user.StatsHandler)
	}

	// 频道：注入 DB 依赖。
	channelHandler := channel.NewHandler(channel.NewRepository(deps.DB))
	channels := api.Group("/channels")
	{
		channels.GET("", channelHandler.List)
		channels.POST("", authMW, channelHandler.Create)
		channels.GET("/:id", channelHandler.Get)
		channels.POST("/:id/join", authMW, channelHandler.Join)
		channels.POST("/:id/leave", authMW, channelHandler.Leave)
		channels.GET("/:id/members", channelHandler.Members)
		channels.GET("/:id/emotion-board", channelHandler.EmotionBoard)
		// 频道内消息
		channels.GET("/:id/messages", message.ListHandler)
		channels.POST("/:id/messages", authMW, message.CreateHandler)
	}

	// 消息 / 线程 / 共情
	messages := api.Group("/messages")
	{
		messages.GET("/:id/threads", message.ThreadsHandler)
		messages.POST("/:id/replies", authMW, message.ReplyHandler)
		messages.DELETE("/:id", authMW, message.DeleteHandler)
		messages.POST("/:id/empathy", authMW, empathy.CreateHandler)
		messages.DELETE("/:id/empathy", authMW, empathy.DeleteHandler)
		messages.GET("/:id/empathy-users", empathy.UsersHandler)
	}

	// 打工日记（v1.5）
	diaries := api.Group("/diaries")
	{
		diaries.GET("", diary.ListHandler)
		diaries.POST("", authMW, diary.CreateHandler)
		diaries.GET("/:id", diary.GetHandler)
		diaries.GET("/:id/comments", diary.CommentsHandler)
		diaries.POST("/:id/comments", authMW, diary.CreateCommentHandler)
		diaries.GET("/streak/:user_id", diary.StreakHandler)
	}

	// 排行榜
	ranking := api.Group("/ranking")
	{
		ranking.GET("/empathy", empathy.RankingEmpathyHandler)
		ranking.GET("/warmest", empathy.RankingWarmestHandler)
		ranking.GET("/streak", diary.RankingStreakHandler)
		ranking.GET("/active", empathy.RankingActiveHandler)
	}

	// 通知
	notifications := api.Group("/notifications", authMW)
	{
		notifications.GET("", notification.ListHandler)
		notifications.PUT("/:id/read", notification.ReadHandler)
		notifications.PUT("/read-all", notification.ReadAllHandler)
	}

	// 搜索
	api.GET("/search", search.SearchHandler)

	// WebSocket 端点
	api.GET("/ws", ws.Handler)
}
