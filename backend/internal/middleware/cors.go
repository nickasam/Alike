// Package middleware 提供 Gin 中间件：CORS、请求日志、JWT 鉴权、Recovery。
package middleware

import "github.com/gin-gonic/gin"

// CORSOptions 配置跨域策略。
type CORSOptions struct {
	// AllowedOrigins 精确匹配的 Origin 白名单。
	AllowedOrigins []string
	// AllowAllInDev 为 true 且白名单为空时，回显请求 Origin（开发便利）。
	// 生产环境应置 false，使空白名单等价于拒绝所有跨域。
	AllowAllInDev bool
}

// CORS 返回按白名单校验的跨域中间件。
//   - 请求 Origin 命中白名单 → 回显该 Origin 并允许带凭证；
//   - 白名单为空且 AllowAllInDev → 回显请求 Origin（仅开发）；
//   - 其余情况 → 不下发 Allow-Origin 头（浏览器据此拒绝跨域），避免反射任意来源 + 凭证。
func CORS(opts CORSOptions) gin.HandlerFunc {
	allow := make(map[string]struct{}, len(opts.AllowedOrigins))
	for _, o := range opts.AllowedOrigins {
		allow[o] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := ""
		if origin != "" {
			if _, ok := allow[origin]; ok {
				allowed = origin
			} else if len(allow) == 0 && opts.AllowAllInDev {
				allowed = origin
			}
		}

		if allowed != "" {
			c.Header("Access-Control-Allow-Origin", allowed)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept")
			c.Header("Access-Control-Expose-Headers", "Content-Length")
		}

		if c.Request.Method == "OPTIONS" {
			// 命中白名单的预检返回 204；否则 403 明确拒绝。
			if allowed != "" {
				c.AbortWithStatus(204)
			} else {
				c.AbortWithStatus(403)
			}
			return
		}
		c.Next()
	}
}
