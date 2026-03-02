package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alike/internal/api/router"
	"github.com/Alike/internal/auth"
	"github.com/Alike/pkg/config"
	"github.com/Alike/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// Map server mode to gin mode
	switch cfg.Server.Mode {
	case "development", "dev":
		gin.SetMode(gin.DebugMode)
	case "production", "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
	
	db, err := database.PostgreSQL(&cfg.Database)
	if err != nil {
		log.Printf("Warning: Database connection failed: %v", err)
		log.Println("Continuing without database...")
	} else {
		log.Println("✅ Database connected")
		defer func() {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		}()
	}
	
	authService := auth.New(&cfg.JWT)
	
	r := gin.Default()
	router.Setup(r, authService, db)
	
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
	
	go func() {
		log.Printf("🚀 Server starting on http://localhost:%d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	srv.Shutdown(ctx)
	log.Println("Server shutdown complete")
}
