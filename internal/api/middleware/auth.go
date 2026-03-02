package middleware

import (
	"net/http"
	"strings"

	"github.com/Alike/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{"code": "UNAUTHORIZED", "message": "Missing authorization header"},
			})
			c.Abort()
			return
		}
		
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{"code": "UNAUTHORIZED", "message": "Invalid authorization format"},
			})
			c.Abort()
			return
		}
		
		userID, err := authService.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{"code": "UNAUTHORIZED", "message": "Invalid token"},
			})
			c.Abort()
			return
		}
		
		c.Set("user_id", userID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}
