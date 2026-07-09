// Package message 负责消息发布、线程回复、历史查询。
// 阶段一仅提供占位 handler。
package message

import (
	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/response"
)

func todo(c *gin.Context, name string) {
	response.Success(c, gin.H{"todo": name})
}

// ListHandler 频道消息列表（分页）。TODO。
func ListHandler(c *gin.Context) { todo(c, "message.list") }

// CreateHandler 发布消息。TODO。
func CreateHandler(c *gin.Context) { todo(c, "message.create") }

// ThreadsHandler 线程回复列表。TODO。
func ThreadsHandler(c *gin.Context) { todo(c, "message.threads") }

// ReplyHandler 线程回复。TODO。
func ReplyHandler(c *gin.Context) { todo(c, "message.reply") }

// DeleteHandler 删除消息（软删）。TODO。
func DeleteHandler(c *gin.Context) { todo(c, "message.delete") }
