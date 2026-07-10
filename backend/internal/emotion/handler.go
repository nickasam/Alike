package emotion

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/response"
)

// Handler 承载 emotion 模块的依赖。
type Handler struct {
	repo *Repository
}

// NewHandler 创建 emotion handler。
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// Board 处理 GET /api/channels/:id/emotion-board，返回频道情绪看板。
// query scope=today|all，默认 today（"今日情绪看板"）。
func (h *Handler) Board(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	todayOnly := c.Query("scope") != "all" // 默认今日

	board, err := h.repo.BoardByChannel(c.Request.Context(), id, todayOnly)
	if errors.Is(err, ErrChannelNotFound) {
		response.Error(c, response.CodeNotFound, "频道不存在")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, board)
}

// GlobalBoard 处理 GET /api/emotion/board，返回全站情绪看板（供首页今日情绪看板）。
// query scope=today|all，默认 today。
func (h *Handler) GlobalBoard(c *gin.Context) {
	todayOnly := c.Query("scope") != "all"
	board, err := h.repo.BoardGlobal(c.Request.Context(), todayOnly)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, board)
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
