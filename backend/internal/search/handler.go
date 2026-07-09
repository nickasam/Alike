// Package search 负责消息/日记/频道/用户模糊搜索（PostgreSQL ILIKE）。
package search

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/response"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 50
	maxQueryLen     = 100
)

// Handler 承载 search 模块的依赖。
type Handler struct {
	repo *Repository
}

// NewHandler 创建 search handler。
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// Search 处理 GET /api/search?q=xxx&type=message&channel_id=&page=&page_size=。
// type 默认 message，可选 message/diary/channel/user。返回分页结果。
func (h *Handler) Search(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		response.Error(c, response.CodeValidationError, "搜索关键词不能为空")
		return
	}
	if len(q) > maxQueryLen {
		q = q[:maxQueryLen]
	}

	page, pageSize := paginate(c)
	ctx := c.Request.Context()

	switch parseType(c.Query("type")) {
	case TypeDiary:
		list, total, err := h.repo.SearchDiaries(ctx, q, page, pageSize)
		respond(c, nonNil(list), total, page, pageSize, err)
	case TypeChannel:
		list, total, err := h.repo.SearchChannels(ctx, q, page, pageSize)
		respond(c, nonNil(list), total, page, pageSize, err)
	case TypeUser:
		list, total, err := h.repo.SearchUsers(ctx, q, page, pageSize)
		respond(c, nonNil(list), total, page, pageSize, err)
	default: // TypeMessage
		list, total, err := h.repo.SearchMessages(ctx, q, parseChannelID(c), page, pageSize)
		respond(c, nonNil(list), total, page, pageSize, err)
	}
}

// respond 统一处理错误与分页响应。
func respond[T any](c *gin.Context, list []T, total int64, page, pageSize int, err error) {
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Page(c, list, total, page, pageSize)
}

// parseType 归一化 type 参数，未知或缺省视为 message。
func parseType(t string) SearchType {
	switch SearchType(strings.ToLower(strings.TrimSpace(t))) {
	case TypeDiary:
		return TypeDiary
	case TypeChannel:
		return TypeChannel
	case TypeUser:
		return TypeUser
	default:
		return TypeMessage
	}
}

// parseChannelID 解析可选的 channel_id 过滤参数，非法或缺省返回 0（不过滤）。
func parseChannelID(c *gin.Context) int64 {
	id, _ := strconv.ParseInt(c.Query("channel_id"), 10, 64)
	if id < 0 {
		return 0
	}
	return id
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
