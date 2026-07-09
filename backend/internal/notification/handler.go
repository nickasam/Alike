// Package notification 负责 @提及、共情、回复通知。
package notification

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/internal/middleware"
	"github.com/Alike/backend/pkg/response"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 50
)

// Handler 承载 notification 模块的依赖。
type Handler struct {
	repo *Repository
}

// NewHandler 创建 notification handler。
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// List 处理 GET /api/notifications，通知列表（需登录，分页）。
func (h *Handler) List(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized)
		return
	}
	page, pageSize := paginate(c)

	list, total, unread, err := h.repo.List(c.Request.Context(), uid, page, pageSize)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{
		"list":      nonNil(list),
		"total":     total,
		"unread":    unread,
		"page":      page,
		"page_size": pageSize,
	})
}

// Read 处理 PUT /api/notifications/:id/read，标记单条已读（需登录）。
func (h *Handler) Read(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized)
		return
	}
	id, ok := parseID(c)
	if !ok {
		return
	}

	found, err := h.repo.MarkRead(c.Request.Context(), id, uid)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	if !found {
		response.Error(c, response.CodeNotFound, "通知不存在")
		return
	}
	response.Success(c, gin.H{"message": "已读"})
}

// ReadAll 处理 PUT /api/notifications/read-all，全部标记已读（需登录）。
func (h *Handler) ReadAll(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized)
		return
	}
	n, err := h.repo.MarkAllRead(c.Request.Context(), uid)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{"marked": n})
}

// parseID 解析路径参数 :id，非法时写入 404 响应并返回 false。
func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, response.CodeNotFound, "通知不存在")
		return 0, false
	}
	return id, true
}

// paginate 从 query 解析分页参数，应用默认值与上限。
func paginate(c *gin.Context) (page, pageSize int) {
	page, _ = strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = defaultPage
	}
	pageSize, _ = strconv.Atoi(c.Query("page_size"))
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}

// nonNil 保证空列表序列化为 [] 而非 null。
func nonNil[T any](list []T) []T {
	if list == nil {
		return []T{}
	}
	return list
}
