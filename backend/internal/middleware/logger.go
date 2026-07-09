package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 记录每个请求的方法、路径、状态码与耗时。
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		if raw != "" {
			path = path + "?" + raw
		}
		log.Printf("[GIN] %3d | %13v | %-7s %s",
			c.Writer.Status(), latency, c.Request.Method, path)
	}
}
