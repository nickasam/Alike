package router

import (
	"github.com/Alike/internal/api/handler"
	"github.com/Alike/internal/api/middleware"
	"github.com/Alike/internal/auth"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, authService *auth.Service) {
	r.Use(middleware.CORSMiddleware())
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	
	v1 := r.Group("/api/v1")
	{
		authHandler := handler.New(authService)
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
			protected.GET("/users/me", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Get current user"})
			})
			protected.GET("/users/nearby", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Get nearby users"})
			})
		}
	}
}
