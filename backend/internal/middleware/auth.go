package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/jwt"
	"github.com/Alike/backend/pkg/response"
)

// ContextUserID 是存放在 gin.Context 中的当前用户 ID 键名。
const ContextUserID = "user_id"

// Auth 返回一个校验 access token 的 JWT 鉴权中间件。
// 期望请求头 Authorization: Bearer <token>，校验通过后将 user_id 写入上下文。
func Auth(mgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Fail(c, response.CodeUnauthorized)
			c.Abort()
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Fail(c, response.CodeUnauthorized)
			c.Abort()
			return
		}

		claims, err := mgr.Parse(parts[1])
		if err != nil || claims.Type != jwt.AccessToken {
			response.Fail(c, response.CodeUnauthorized)
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Next()
	}
}

// CurrentUserID 从上下文取出当前用户 ID，不存在时返回 0, false。
func CurrentUserID(c *gin.Context) (int64, bool) {
	v, ok := c.Get(ContextUserID)
	if !ok {
		return 0, false
	}
	id, ok := v.(int64)
	return id, ok
}
