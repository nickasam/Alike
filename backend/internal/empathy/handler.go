// Package empathy 负责共情（抱团取暖）与排行榜。
package empathy

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/internal/middleware"
	"github.com/Alike/backend/pkg/response"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 50
	defaultRankSize = 20
	maxRankSize     = 100
)

// Handler 承载 empathy 模块的依赖。
type Handler struct {
	repo *Repository
}

// NewHandler 创建 empathy handler。
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// Create 处理 POST /api/messages/:id/empathy，抱团取暖（需登录）。
func (h *Handler) Create(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized)
		return
	}
	id, ok := parseID(c)
	if !ok {
		return
	}

	count, err := h.repo.Create(c.Request.Context(), id, uid)
	switch {
	case errors.Is(err, ErrMessageNotFound):
		response.Error(c, response.CodeNotFound, "消息不存在")
		return
	case errors.Is(err, ErrAlreadyEmpathized):
		response.Error(c, response.CodeConflict, "你已经抱过啦")
		return
	case err != nil:
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{"empathy_count": count, "empathized": true})
}

// Delete 处理 DELETE /api/messages/:id/empathy，取消共情（需登录）。
func (h *Handler) Delete(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized)
		return
	}
	id, ok := parseID(c)
	if !ok {
		return
	}

	count, err := h.repo.Delete(c.Request.Context(), id, uid)
	switch {
	case errors.Is(err, ErrMessageNotFound):
		response.Error(c, response.CodeNotFound, "消息不存在")
		return
	case errors.Is(err, ErrNotEmpathized):
		response.Error(c, response.CodeNotFound, "你还没有抱过这条消息")
		return
	case err != nil:
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{"empathy_count": count, "empathized": false})
}

// Users 处理 GET /api/messages/:id/empathy-users，返回共情用户列表（分页）。
func (h *Handler) Users(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	page, pageSize := paginate(c)

	list, total, err := h.repo.ListUsers(c.Request.Context(), id, page, pageSize)
	if errors.Is(err, ErrMessageNotFound) {
		response.Error(c, response.CodeNotFound, "消息不存在")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Page(c, nonNil(list), total, page, pageSize)
}

// RankingEmpathy 处理 GET /api/ranking/empathy，最受共情帖子榜。
func (h *Handler) RankingEmpathy(c *gin.Context) {
	list, err := h.repo.RankingEmpathy(c.Request.Context(), rankLimit(c))
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{"list": nonNil(list)})
}

// RankingWarmest 处理 GET /api/ranking/warmest，最暖牛马榜。
func (h *Handler) RankingWarmest(c *gin.Context) {
	list, err := h.repo.RankingWarmest(c.Request.Context(), rankLimit(c))
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{"list": nonNil(list)})
}

// RankingActive 处理 GET /api/ranking/active，本周最活跃牛马榜。
func (h *Handler) RankingActive(c *gin.Context) {
	list, err := h.repo.RankingActive(c.Request.Context(), rankLimit(c))
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{"list": nonNil(list)})
}

// parseID 解析路径参数 :id，非法时写入 404 响应并返回 false。
func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, response.CodeNotFound, "消息不存在")
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

// rankLimit 解析榜单 limit query，应用默认值与上限。
func rankLimit(c *gin.Context) int {
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit < 1 {
		return defaultRankSize
	}
	if limit > maxRankSize {
		return maxRankSize
	}
	return limit
}

// nonNil 保证空列表序列化为 [] 而非 null。
func nonNil[T any](list []T) []T {
	if list == nil {
		return []T{}
	}
	return list
}
