// Package router 组装 Gin 引擎、中间件与路由分组。
package router

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/Alike/backend/internal/middleware"
	"github.com/Alike/backend/internal/pulse/scheduler"
	"github.com/Alike/backend/internal/storage"
	"github.com/Alike/backend/internal/ws"
	"github.com/Alike/backend/pkg/config"
	"github.com/Alike/backend/pkg/jwt"
	"github.com/Alike/backend/pkg/response"
)

// Deps 汇聚路由处理所需的依赖。DB / Redis 允许为 nil（依赖不可用时仍能提供 health）。
type Deps struct {
	Cfg   *config.Config
	DB    *sql.DB
	Redis *redis.Client
	JWT   *jwt.Manager
}

// New 构建并返回配置好的 Gin 引擎，以及 WebSocket Hub 与 Pulse Scheduler
// （供优雅关闭使用；无 DB 时可能为 nil）。
func New(deps *Deps) (*gin.Engine, *ws.Hub, *scheduler.Scheduler) {
	if deps.Cfg != nil && deps.Cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	// 仅信任私有网段代理（nginx 在 Docker 私有网络内）。
	// Gin 默认信任所有代理头，会采信客户端伪造的 X-Forwarded-For，
	// 使按 IP 限流可被逐请求换 XFF 绕过。声明可信代理后，ClientIP 只取
	// nginx 追加（$proxy_add_x_forwarded_for，真实 peer 置于最右）的可信来源。
	if err := r.SetTrustedProxies([]string{
		"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.1/32",
	}); err != nil {
		panic("router: set trusted proxies: " + err.Error())
	}
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(middleware.CORSOptions{
		AllowedOrigins: corsOrigins(deps.Cfg),
		// 非生产环境且未配置白名单时放开跨域，便于本地开发。
		AllowAllInDev: deps.Cfg == nil || !deps.Cfg.IsProduction(),
	}))

	// 健康检查（统一响应格式），探测 DB / Redis 连通性。
	r.GET("/api/health", healthHandler(deps))

	api := r.Group("/api")
	hub, pulseSched := registerRoutes(api, deps)

	return r, hub, pulseSched
}

// corsOrigins 从配置读取 CORS 白名单（nil 安全）。
func corsOrigins(cfg *config.Config) []string {
	if cfg == nil {
		return nil
	}
	return cfg.CORSAllowedOrigins
}

// healthHandler 返回服务及依赖的健康状态。
// 任一依赖 down 时整体 status=down 并返回 503，供编排/探针识别。
func healthHandler(deps *Deps) gin.HandlerFunc {
	// 复用一个 storage.Store 探测 MinIO（客户端惰性连接，构造无副作用）。
	var store *storage.Store
	if deps.Cfg != nil {
		store, _ = storage.New(deps.Cfg)
	}
	return func(c *gin.Context) {
		db := pingDB(deps.DB)
		rdb := pingRedis(deps.Redis)
		minioStatus := pingMinIO(c.Request.Context(), store)

		overall := "ok"
		httpCode := http.StatusOK
		// "unavailable"（未配置）不视为故障；仅 "down" 判定整体不健康。
		if db == "down" || rdb == "down" || minioStatus == "down" {
			overall = "down"
			httpCode = http.StatusServiceUnavailable
		}
		c.JSON(httpCode, response.Body{
			Code:    response.CodeSuccess,
			Message: "success",
			Data: gin.H{
				"status":   overall,
				"database": db,
				"redis":    rdb,
				"minio":    minioStatus,
			},
		})
	}
}

func pingMinIO(ctx context.Context, store *storage.Store) string {
	if store == nil {
		return "unavailable"
	}
	// 限时探测，避免 MinIO 不可达时健康检查长时间阻塞。
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := store.Ping(ctx); err != nil {
		if errors.Is(err, storage.ErrStorageDisabled) {
			return "unavailable"
		}
		return "down"
	}
	return "ok"
}

func pingDB(db *sql.DB) string {
	if db == nil {
		return "unavailable"
	}
	if err := db.Ping(); err != nil {
		return "down"
	}
	return "ok"
}

func pingRedis(rdb *redis.Client) string {
	if rdb == nil {
		return "unavailable"
	}
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return "down"
	}
	return "ok"
}
