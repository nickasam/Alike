// Package empathy 负责共情（抱团取暖）与排行榜。
// 阶段一仅提供占位 handler。
package empathy

import (
	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/response"
)

func todo(c *gin.Context, name string) {
	response.Success(c, gin.H{"todo": name})
}

// CreateHandler 抱团取暖（共情）。TODO。
func CreateHandler(c *gin.Context) { todo(c, "empathy.create") }

// DeleteHandler 取消共情。TODO。
func DeleteHandler(c *gin.Context) { todo(c, "empathy.delete") }

// UsersHandler 共情用户列表。TODO。
func UsersHandler(c *gin.Context) { todo(c, "empathy.users") }

// RankingEmpathyHandler 最受共情帖子榜。TODO。
func RankingEmpathyHandler(c *gin.Context) { todo(c, "ranking.empathy") }

// RankingWarmestHandler 最暖牛马榜。TODO。
func RankingWarmestHandler(c *gin.Context) { todo(c, "ranking.warmest") }

// RankingActiveHandler 本周最活跃牛马榜。TODO。
func RankingActiveHandler(c *gin.Context) { todo(c, "ranking.active") }
