// Package search 负责消息/日记/频道/用户全文搜索。
// 阶段一仅提供占位 handler。
package search

import (
	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/response"
)

// SearchHandler 全文搜索（按 type 参数分流）。TODO：阶段九实现。
func SearchHandler(c *gin.Context) {
	response.Success(c, gin.H{"todo": "search", "type": c.Query("type"), "q": c.Query("q")})
}
