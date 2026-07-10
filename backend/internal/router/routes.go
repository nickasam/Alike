package router

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/internal/auth"
	"github.com/Alike/backend/internal/channel"
	"github.com/Alike/backend/internal/diary"
	"github.com/Alike/backend/internal/emotion"
	"github.com/Alike/backend/internal/empathy"
	"github.com/Alike/backend/internal/message"
	"github.com/Alike/backend/internal/middleware"
	"github.com/Alike/backend/internal/notification"
	"github.com/Alike/backend/internal/search"
	"github.com/Alike/backend/internal/storage"
	"github.com/Alike/backend/internal/user"
	"github.com/Alike/backend/internal/ws"
)

// registerRoutes 注册各业务模块的路由分组，返回 WebSocket Hub 供优雅关闭。
func registerRoutes(api *gin.RouterGroup, deps *Deps) *ws.Hub {
	// 需鉴权的中间件
	authMW := middleware.Auth(deps.JWT)
	// 尽力而为鉴权：登录则识别本人，未登录仍放行（用于日记详情等半公开读接口）
	optionalAuthMW := middleware.OptionalAuth(deps.JWT)

	// 限流器：认证接口按 IP 严格限流（防暴力破解/注册刷号），
	// 写接口按用户限流（防刷消息/共情）。单实例内存令牌桶。
	authLimiter := middleware.NewLimiter(1, 10)  // 认证：~1 req/s，突发 10
	writeLimiter := middleware.NewLimiter(5, 20) // 写操作：~5 req/s，突发 20
	authRL := middleware.RateLimitByIP(authLimiter)
	writeRL := middleware.RateLimitByUser(writeLimiter)

	// 认证：注入 DB / JWT 依赖。
	authHandler := auth.NewHandler(auth.NewRepository(deps.DB), deps.JWT)
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authRL, authHandler.Register)
		authGroup.POST("/login", authRL, authHandler.Login)
		authGroup.POST("/refresh", authRL, authHandler.Refresh)
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.GET("/me", authMW, authHandler.Me)
	}

	// 用户：注入 DB 依赖。
	userHandler := user.NewHandler(user.NewRepository(deps.DB))
	users := api.Group("/users")
	{
		users.GET("/:id", userHandler.Get)
		users.PUT("/:id", authMW, userHandler.Update)
		users.GET("/:id/diaries", userHandler.Diaries)
		users.GET("/:id/stats", userHandler.Stats)
	}

	// 频道：注入 DB 依赖。
	channelHandler := channel.NewHandler(channel.NewRepository(deps.DB))
	channels := api.Group("/channels")

	// 情绪看板 + 共情：注入 DB 依赖。
	emotionRepo := emotion.NewRepository(deps.DB)
	emotionHandler := emotion.NewHandler(emotionRepo)
	empathyRepo := empathy.NewRepository(deps.DB)

	// 通知仓储：供共情/回复等事件写入通知（fire-and-forget）。
	notificationRepo := notification.NewRepository(deps.DB)

	// 消息 + WebSocket：Hub 依赖 message 服务，message handler 反向依赖 Hub 广播。
	msgRepo := message.NewRepository(deps.DB)
	pubsub := ws.NewPubSub(deps.Redis)
	hub := ws.NewHub(message.NewService(msgRepo), pubsub)
	// 注入情绪看板提供者，启用 emotion_update 实时推送。
	hub.SetEmotionProvider(emotion.NewService(emotionRepo))
	messageHandler := message.NewHandler(msgRepo, hub, notificationRepo)
	// 共情变更经 Hub 广播 empathy 事件，实现跨端实时同步。
	empathyHandler := empathy.NewHandler(empathyRepo, hub, notificationRepo)
	wsHandler := ws.NewHandler(hub, deps.JWT, corsOrigins(deps.Cfg)...)

	{
		channels.GET("", channelHandler.List)
		channels.POST("", authMW, writeRL, channelHandler.Create)
		channels.GET("/:id", channelHandler.Get)
		channels.POST("/:id/join", authMW, channelHandler.Join)
		channels.POST("/:id/leave", authMW, channelHandler.Leave)
		channels.GET("/:id/members", channelHandler.Members)
		channels.GET("/:id/emotion-board", emotionHandler.Board)
		// 频道内消息（OptionalAuth：登录时返回每条消息的 empathized 态）
		channels.GET("/:id/messages", optionalAuthMW, messageHandler.List)
		channels.POST("/:id/messages", authMW, writeRL, messageHandler.Create)
	}

	// 消息 / 线程 / 共情
	messages := api.Group("/messages")
	{
		messages.GET("/:id/threads", optionalAuthMW, messageHandler.Threads)
		messages.POST("/:id/replies", authMW, writeRL, messageHandler.Reply)
		messages.DELETE("/:id", authMW, messageHandler.Delete)
		messages.POST("/:id/empathy", authMW, writeRL, empathyHandler.Create)
		messages.DELETE("/:id/empathy", authMW, empathyHandler.Delete)
		messages.GET("/:id/empathy-users", empathyHandler.Users)
	}

	// 打工日记（v1.5）：注入 DB 依赖。
	diaryHandler := diary.NewHandler(diary.NewRepository(deps.DB))
	diaries := api.Group("/diaries")
	{
		diaries.GET("", diaryHandler.List)
		diaries.POST("", authMW, writeRL, diaryHandler.Create)
		diaries.GET("/:id", optionalAuthMW, diaryHandler.Get)
		diaries.GET("/:id/comments", optionalAuthMW, diaryHandler.Comments)
		diaries.POST("/:id/comments", authMW, writeRL, diaryHandler.CreateComment)
		diaries.GET("/streak/:user_id", diaryHandler.Streak)
	}

	// 排行榜
	ranking := api.Group("/ranking")
	{
		ranking.GET("/empathy", empathyHandler.RankingEmpathy)
		ranking.GET("/warmest", empathyHandler.RankingWarmest)
		ranking.GET("/streak", diaryHandler.RankingStreak)
		ranking.GET("/active", empathyHandler.RankingActive)
	}

	// 通知：复用上文构造的 notificationRepo。
	notificationHandler := notification.NewHandler(notificationRepo)
	notifications := api.Group("/notifications", authMW)
	{
		notifications.GET("", notificationHandler.List)
		notifications.PUT("/:id/read", notificationHandler.Read)
		notifications.PUT("/read-all", notificationHandler.ReadAll)
	}

	// 搜索：注入 DB 依赖。
	searchHandler := search.NewHandler(search.NewRepository(deps.DB))
	api.GET("/search", searchHandler.Search)

	// 全站今日情绪看板（首页用，公开读）。
	api.GET("/emotion/board", emotionHandler.GlobalBoard)

	// 文件上传：MinIO 不可用时跳过路由注册。
	if store, err := storage.New(deps.Cfg); err == nil {
		// 确保目标 bucket 存在（幂等），否则首次上传会因 bucket 缺失失败。
		// 失败不阻断路由注册，仅记录，上传时会返回相应错误。
		if err := store.EnsureBucket(context.Background()); err != nil {
			log.Printf("storage: ensure bucket failed: %v", err)
		}
		storageHandler := storage.NewHandler(store)
		api.POST("/upload", authMW, writeRL, storageHandler.Upload)
	}

	// WebSocket 端点
	api.GET("/ws", wsHandler.Handle)

	return hub
}
