// Package router 组装 Gin 引擎、中间件与路由分组。
package router

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/Alike/backend/internal/middleware"
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

// New 构建并返回配置好的 Gin 引擎。
func New(deps *Deps) *gin.Engine {
	if deps.Cfg != nil && deps.Cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 健康检查（统一响应格式），探测 DB / Redis 连通性。
	r.GET("/api/health", healthHandler(deps))

	api := r.Group("/api")
	registerRoutes(api, deps)

	return r
}

// healthHandler 返回服务及依赖的健康状态。
func healthHandler(deps *Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := gin.H{
			"status":   "ok",
			"database": pingDB(deps.DB),
			"redis":    pingRedis(deps.Redis),
		}
		response.Success(c, status)
	}
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
