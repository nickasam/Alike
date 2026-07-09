// Package channel 负责频道 CRUD、成员管理、频道列表、情绪看板入口。
package channel

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
)

// Handler 承载 channel 模块的依赖，替代阶段一的无状态包级函数。
type Handler struct {
	repo *Repository
}

// NewHandler 创建 channel handler。
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// List 处理 GET /api/channels，支持 category 过滤与分页。
func (h *Handler) List(c *gin.Context) {
	category := c.Query("category")
	page, pageSize := paginate(c)

	list, total, err := h.repo.List(c.Request.Context(), category, page, pageSize)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Page(c, nonNil(list), total, page, pageSize)
}

// Create 处理 POST /api/channels，创建频道（需登录）。
func (h *Handler) Create(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized)
		return
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeValidationError, "频道名称、slug 或分类格式不正确")
		return
	}

	ch, err := h.repo.Create(c.Request.Context(), req, uid)
	if errors.Is(err, ErrSlugConflict) {
		response.Error(c, response.CodeConflict, "该 slug 已被占用")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, ch)
}

// Get 处理 GET /api/channels/:id，返回频道详情（含 member_count）。
func (h *Handler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	ch, err := h.repo.GetByID(c.Request.Context(), id)
	if errors.Is(err, ErrChannelNotFound) {
		response.Error(c, response.CodeNotFound, "频道不存在")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, ch)
}

// Join 处理 POST /api/channels/:id/join，加入频道（需登录）。
func (h *Handler) Join(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized)
		return
	}
	id, ok := parseID(c)
	if !ok {
		return
	}

	err := h.repo.Join(c.Request.Context(), id, uid)
	switch {
	case errors.Is(err, ErrChannelNotFound):
		response.Error(c, response.CodeNotFound, "频道不存在")
		return
	case errors.Is(err, ErrAlreadyMember):
		response.Error(c, response.CodeConflict, "你已加入该频道")
		return
	case err != nil:
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{"message": "已加入频道"})
}

// Leave 处理 POST /api/channels/:id/leave，离开频道（需登录）。
func (h *Handler) Leave(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized)
		return
	}
	id, ok := parseID(c)
	if !ok {
		return
	}

	err := h.repo.Leave(c.Request.Context(), id, uid)
	if errors.Is(err, ErrNotMember) {
		response.Error(c, response.CodeNotFound, "你尚未加入该频道")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{"message": "已离开频道"})
}

// Members 处理 GET /api/channels/:id/members，返回成员列表（分页）。
func (h *Handler) Members(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	page, pageSize := paginate(c)

	list, total, err := h.repo.Members(c.Request.Context(), id, page, pageSize)
	if errors.Is(err, ErrChannelNotFound) {
		response.Error(c, response.CodeNotFound, "频道不存在")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Page(c, nonNil(list), total, page, pageSize)
}

// EmotionBoard 处理 GET /api/channels/:id/emotion-board。
// 阶段五由 emotion 模块提供数据，此处暂返回空看板。
func (h *Handler) EmotionBoard(c *gin.Context) {
	if _, ok := parseID(c); !ok {
		return
	}
	response.Success(c, gin.H{"emotions": []any{}, "total": 0})
}

// parseID 解析路径参数 :id，非法时写入 404 响应并返回 false。
func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, response.CodeNotFound, "频道不存在")
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
