// Package notification 负责 @提及、共情、回复通知。
// 阶段一仅提供占位 handler。
package notification

import (
	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/response"
)

func todo(c *gin.Context, name string) {
	response.Success(c, gin.H{"todo": name})
}

// ListHandler 通知列表。TODO。
func ListHandler(c *gin.Context) { todo(c, "notification.list") }

// ReadHandler 标记单条已读。TODO。
func ReadHandler(c *gin.Context) { todo(c, "notification.read") }

// ReadAllHandler 全部已读。TODO。
func ReadAllHandler(c *gin.Context) { todo(c, "notification.read_all") }
