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

// Board 处理 GET /api/channels/:id/emotion-board，返回频道实时情绪看板。
func (h *Handler) Board(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	board, err := h.repo.BoardByChannel(c.Request.Context(), id)
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

// parseID 解析路径参数 :id，非法时写入 404 响应并返回 false。
func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, response.CodeNotFound, "频道不存在")
		return 0, false
	}
	return id, true
}
