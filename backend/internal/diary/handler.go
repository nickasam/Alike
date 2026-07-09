// Package diary 负责打工日记 CRUD、评论、打卡日历（v1.5）。
// 阶段一仅提供占位 handler。
package diary

import (
	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/response"
)

func todo(c *gin.Context, name string) {
	response.Success(c, gin.H{"todo": name})
}

// ListHandler 日记广场（公开流）。TODO。
func ListHandler(c *gin.Context) { todo(c, "diary.list") }

// CreateHandler 写日记/打卡。TODO。
func CreateHandler(c *gin.Context) { todo(c, "diary.create") }

// GetHandler 日记详情。TODO。
func GetHandler(c *gin.Context) { todo(c, "diary.get") }

// CommentsHandler 日记评论列表。TODO。
func CommentsHandler(c *gin.Context) { todo(c, "diary.comments") }

// CreateCommentHandler 发表日记评论。TODO。
func CreateCommentHandler(c *gin.Context) { todo(c, "diary.create_comment") }

// StreakHandler 打卡日历（连续 + 累计天数）。TODO。
func StreakHandler(c *gin.Context) { todo(c, "diary.streak") }

// RankingStreakHandler 连续打卡牛马榜。TODO。
func RankingStreakHandler(c *gin.Context) { todo(c, "ranking.streak") }
