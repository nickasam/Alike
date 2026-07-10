// Package user 负责用户信息、牛马等级、个人主页。
package user

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/internal/middleware"
	"github.com/Alike/backend/pkg/response"
)

// Handler 承载 user 模块的依赖，替代阶段一的无状态包级函数。
type Handler struct {
	repo *Repository
}

// NewHandler 创建 user handler。
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// Get 处理 GET /api/users/:id，返回公开主页信息（不含 email/password）。
func (h *Handler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	u, err := h.repo.GetByID(c.Request.Context(), id)
	if errors.Is(err, ErrUserNotFound) {
		response.Error(c, response.CodeNotFound, "用户不存在")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, u.Public())
}

// Update 处理 PUT /api/users/:id，仅本人可改（需登录）。
func (h *Handler) Update(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized)
		return
	}
	id, ok := parseID(c)
	if !ok {
		return
	}
	if id != uid {
		response.Error(c, response.CodeForbidden, "只能修改自己的资料")
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeValidationError, "资料字段格式不正确")
		return
	}

	u, err := h.repo.Update(c.Request.Context(), id, req)
	if errors.Is(err, ErrUserNotFound) {
		response.Error(c, response.CodeNotFound, "用户不存在")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, u)
}

// Diaries 处理 GET /api/users/:id/diaries，返回该用户的公开日记（分页）。
func (h *Handler) Diaries(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	page, pageSize := paginate(c)

	list, total, err := h.repo.ListPublicDiaries(c.Request.Context(), id, page, pageSize)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Page(c, nonNilDiaries(list), total, page, pageSize)
}

// Stats 处理 GET /api/users/:id/stats，返回该用户的公开统计聚合。
func (h *Handler) Stats(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	stats, err := h.repo.Stats(c.Request.Context(), id)
	if errors.Is(err, ErrUserNotFound) {
		response.Error(c, response.CodeNotFound, "用户不存在")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, stats)
}

// parseID 解析路径参数 :id，非法时写入 404 响应并返回 false。
func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, response.CodeNotFound, "用户不存在")
		return 0, false
	}
	return id, true
}

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 50
)

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

// nonNilDiaries 保证空列表序列化为 [] 而非 null。
func nonNilDiaries(list []*PublicDiary) []*PublicDiary {
	if list == nil {
		return []*PublicDiary{}
	}
	return list
}
