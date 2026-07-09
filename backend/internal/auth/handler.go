// Package auth 负责注册、登录、JWT 签发与验证、匿名令牌。
// 阶段一仅提供占位 handler，业务实现见阶段二。
package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/response"
)

// todo 统一的占位响应。
func todo(c *gin.Context, name string) {
	response.Success(c, gin.H{"todo": name})
}

// RegisterHandler 注册。TODO: 阶段二实现。
func RegisterHandler(c *gin.Context) { todo(c, "auth.register") }

// LoginHandler 登录。TODO: 阶段二实现。
func LoginHandler(c *gin.Context) { todo(c, "auth.login") }

// RefreshHandler 刷新 JWT。TODO: 阶段二实现。
func RefreshHandler(c *gin.Context) { todo(c, "auth.refresh") }

// LogoutHandler 登出。TODO: 阶段二实现。
func LogoutHandler(c *gin.Context) { todo(c, "auth.logout") }

// MeHandler 当前用户信息。TODO: 阶段二实现。
func MeHandler(c *gin.Context) { todo(c, "auth.me") }
