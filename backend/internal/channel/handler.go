// Package channel 负责频道 CRUD、成员管理、频道列表、情绪看板入口。
// 阶段一仅提供占位 handler。
package channel

import (
	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/response"
)

func todo(c *gin.Context, name string) {
	response.Success(c, gin.H{"todo": name})
}

// ListHandler 频道列表。TODO。
func ListHandler(c *gin.Context) { todo(c, "channel.list") }

// CreateHandler 申请创建频道。TODO。
func CreateHandler(c *gin.Context) { todo(c, "channel.create") }

// GetHandler 频道详情。TODO。
func GetHandler(c *gin.Context) { todo(c, "channel.get") }

// JoinHandler 加入频道。TODO。
func JoinHandler(c *gin.Context) { todo(c, "channel.join") }

// LeaveHandler 离开频道。TODO。
func LeaveHandler(c *gin.Context) { todo(c, "channel.leave") }

// MembersHandler 频道成员列表。TODO。
func MembersHandler(c *gin.Context) { todo(c, "channel.members") }

// EmotionBoardHandler 频道情绪看板。TODO（emotion 模块提供数据）。
func EmotionBoardHandler(c *gin.Context) { todo(c, "channel.emotion_board") }
