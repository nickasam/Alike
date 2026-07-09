// Command server 是 Alike 后端 HTTP 服务入口。
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alike/backend/internal/router"
	"github.com/Alike/backend/pkg/config"
	"github.com/Alike/backend/pkg/database"
	"github.com/Alike/backend/pkg/jwt"
	"github.com/Alike/backend/pkg/redis"
)

func main() {
	cfg := config.Load()

	// DB / Redis 连接失败不阻断启动：记录错误，健康检查会反映为 down/unavailable。
	db, err := database.New(cfg)
	if err != nil {
		log.Printf("[WARN] 无法连接 PostgreSQL: %v（服务仍启动，health 将报告 down）", err)
	}
	if db != nil {
		defer db.Close()
	}

	rdb, err := redis.New(cfg)
	if err != nil {
		log.Printf("[WARN] 无法连接 Redis: %v（服务仍启动，health 将报告 down）", err)
	}
	if rdb != nil {
		defer rdb.Close()
	}

	jwtMgr := jwt.NewManager(cfg.JWTSecret, cfg.JWTAccessExpire, cfg.JWTRefreshExpire)

	engine := router.New(&router.Deps{
		Cfg:   cfg,
		DB:    db,
		Redis: rdb,
		JWT:   jwtMgr,
	})

	srv := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           engine,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// 后台启动 HTTP 服务。
	go func() {
		log.Printf("[INFO] Alike 后端启动，监听 :%s (env=%s)", cfg.ServerPort, cfg.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] HTTP 服务异常退出: %v", err)
		}
	}()

	// 等待中断信号，优雅关闭。
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[INFO] 收到关闭信号，正在优雅关闭…")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[WARN] 优雅关闭超时: %v", err)
	}
	log.Println("[INFO] 服务已退出")
}
