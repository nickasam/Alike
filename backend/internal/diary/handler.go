// Package diary 负责打工日记 CRUD、评论、打卡日历（v1.5）。
package diary

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/internal/middleware"
	"github.com/Alike/backend/pkg/httputil"
	"github.com/Alike/backend/pkg/response"
)

const (
	defaultLimit    = 20
	maxLimit        = 50
	defaultRankSize = 20
	maxRankSize     = 100
)

// itoa 是 strconv.Itoa 的别名，用于拼接 SQL 占位符序号。
func itoa(n int) string { return strconv.Itoa(n) }

// Handler 承载 diary 模块的依赖。
type Handler struct {
	repo *Repository
}

// NewHandler 创建 diary handler。
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// List 处理 GET /api/diaries，游标分页返回公开日记流。
// query: before=<日记ID>（返回更早的），limit（默认 20，上限 50）。
func (h *Handler) List(c *gin.Context) {
	before := parseCursor(c, "before")
	limit := parseLimit(c)

	list, hasMore, err := h.repo.ListPublic(c.Request.Context(), before, limit)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, listData(list, hasMore))
}

// Create 处理 POST /api/diaries，写日记（需登录）。
func (h *Handler) Create(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized)
		return
	}
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeValidationError, "日记内容格式不正确")
		return
	}

	d, err := h.repo.Create(c.Request.Context(), uid, req)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, d)
}

// Get 处理 GET /api/diaries/:id，日记详情。
// 需挂 OptionalAuth：私密日记仅作者本人可见，其余访问者返回 404。
func (h *Handler) Get(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	viewerID, _ := middleware.CurrentUserID(c)
	d, err := h.repo.Get(c.Request.Context(), id, viewerID)
	if errors.Is(err, ErrDiaryNotFound) {
		response.Error(c, response.CodeNotFound, "日记不存在")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, d)
}

// Comments 处理 GET /api/diaries/:id/comments，评论列表（分页）。
// 需挂 OptionalAuth：私密日记仅作者本人可查看其评论。
func (h *Handler) Comments(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	viewerID, _ := middleware.CurrentUserID(c)
	page, pageSize := httputil.Paginate(c)

	list, total, err := h.repo.ListComments(c.Request.Context(), id, viewerID, page, pageSize)
	if errors.Is(err, ErrDiaryNotFound) {
		response.Error(c, response.CodeNotFound, "日记不存在")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Page(c, httputil.NonNil(list), total, page, pageSize)
}

// CreateComment 处理 POST /api/diaries/:id/comments，发表评论（需登录，支持匿名）。
func (h *Handler) CreateComment(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized)
		return
	}
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeValidationError, "评论内容格式不正确")
		return
	}

	cm, err := h.repo.CreateComment(c.Request.Context(), id, uid, req)
	if errors.Is(err, ErrDiaryNotFound) {
		response.Error(c, response.CodeNotFound, "日记不存在")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, cm)
}

// Streak 处理 GET /api/diaries/streak/:user_id，连续 + 累计打卡天数。
func (h *Handler) Streak(c *gin.Context) {
	uid, ok := parseID(c, "user_id")
	if !ok {
		return
	}
	s, err := h.repo.GetStreak(c.Request.Context(), uid)
	if errors.Is(err, ErrUserNotFound) {
		response.Error(c, response.CodeNotFound, "用户不存在")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, s)
}

// RankingStreak 处理 GET /api/ranking/streak，连续打卡牛马榜。
func (h *Handler) RankingStreak(c *gin.Context) {
	list, err := h.repo.RankingStreak(c.Request.Context(), rankLimit(c))
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, gin.H{"list": httputil.NonNil(list)})
}

// listData 组装游标分页的响应结构。next_cursor 指向本页最后一条日记 ID。
func listData(list []*Diary, hasMore bool) gin.H {
	items := httputil.NonNil(list)
	var next int64
	if hasMore && len(items) > 0 {
		next = items[len(items)-1].ID
	}
	return gin.H{"list": items, "has_more": hasMore, "next_cursor": next}
}

// parseID 解析路径参数（:id / :user_id），非法时写入 404 响应并返回 false。
func parseID(c *gin.Context, key string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(key), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, response.CodeNotFound, "资源不存在")
		return 0, false
	}
	return id, true
}

// parseCursor 解析游标 query（before），非法或缺省返回 0。
func parseCursor(c *gin.Context, key string) int64 {
	v, _ := strconv.ParseInt(c.Query(key), 10, 64)
	if v < 0 {
		return 0
	}
	return v
}

// parseLimit 解析 limit query，应用默认值与上限。
func parseLimit(c *gin.Context) int {
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit < 1 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
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
