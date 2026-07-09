// Package user 负责用户信息、牛马等级、个人主页。
// 阶段一仅提供占位 handler。
package user

import (
	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/response"
)

func todo(c *gin.Context, name string) {
	response.Success(c, gin.H{"todo": name})
}

// GetHandler 用户主页。TODO。
func GetHandler(c *gin.Context) { todo(c, "user.get") }

// UpdateHandler 更新资料。TODO。
func UpdateHandler(c *gin.Context) { todo(c, "user.update") }

// DiariesHandler 用户日记列表。TODO。
func DiariesHandler(c *gin.Context) { todo(c, "user.diaries") }

// StatsHandler 用户统计。TODO。
func StatsHandler(c *gin.Context) { todo(c, "user.stats") }
