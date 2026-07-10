// Package httputil 提供跨业务模块复用的 HTTP 辅助函数，
// 消除各 handler 中重复的分页解析与空切片规整逻辑。
package httputil

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// 分页默认值与上限（各模块统一）。
const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 50
)

// Paginate 从 query 解析 page/page_size，应用默认值与上限。
// page<1 归一为 1；page_size<1 归一为 20，超过 50 截断为 50。
func Paginate(c *gin.Context) (page, pageSize int) {
	page, _ = strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = DefaultPage
	}
	pageSize, _ = strconv.Atoi(c.Query("page_size"))
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return page, pageSize
}

// NonNil 保证空切片序列化为 [] 而非 null。
func NonNil[T any](list []T) []T {
	if list == nil {
		return []T{}
	}
	return list
}
