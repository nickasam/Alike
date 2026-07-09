package middleware

import (
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/response"
)

// Recovery 捕获 panic，记录堆栈并返回统一的 500 响应，避免进程崩溃。
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v\n%s", err, debug.Stack())
				response.Error(c, response.CodeInternalError, "")
				c.Abort()
			}
		}()
		c.Next()
	}
}
