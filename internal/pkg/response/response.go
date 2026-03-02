package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func Error(c *gin.Context, code, message string, statusCode int) {
	c.JSON(statusCode, gin.H{"success": false, "error": gin.H{"code": code, "message": message}})
}

func ValidationError(c *gin.Context, message string) {
	Error(c, "VALIDATION_ERROR", message, http.StatusBadRequest)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, "UNAUTHORIZED", message, http.StatusUnauthorized)
}

func InternalError(c *gin.Context, message string) {
	Error(c, "INTERNAL_ERROR", message, http.StatusInternalServerError)
}
