// Package message 负责消息发布、线程回复、历史查询与软删除。
package message

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/internal/middleware"
	"github.com/Alike/backend/pkg/response"
)

const (
	defaultLimit = 20
	maxLimit     = 50
)

// Broadcaster 抽象将消息事件推送到 WebSocket 频道的能力。
// 由 ws.Hub 实现，经依赖注入解耦，避免 message 直接依赖 ws。
type Broadcaster interface {
	BroadcastNewMessage(channelID int64, payload any)
	BroadcastThreadReply(channelID, parentID int64, payload any)
	// BroadcastEmotionUpdate 触发频道情绪看板重聚合并广播（消息带情绪时调用）。
	BroadcastEmotionUpdate(channelID int64)
	// BroadcastMessageDeleted 广播消息被软删除，供其他客户端就地置为已删除。
	BroadcastMessageDeleted(channelID, messageID int64)
}

// Handler 承载 message 模块的依赖。
type Handler struct {
	repo *Repository
	bc   Broadcaster
}

// NewHandler 创建 message handler。bc 可为 nil（WebSocket 不可用时仅落库不广播）。
func NewHandler(repo *Repository, bc Broadcaster) *Handler {
	return &Handler{repo: repo, bc: bc}
}

// List 处理 GET /api/channels/:id/messages，游标分页返回主消息。
// query: before=<消息ID>（返回更早的），limit（默认 20，上限 50）。
func (h *Handler) List(c *gin.Context) {
	channelID, ok := parseID(c)
	if !ok {
		return
	}
	before := parseCursor(c, "before")
	limit := parseLimit(c)
	viewerID, _ := middleware.CurrentUserID(c) // 未登录为 0，empathized 恒 false

	list, hasMore, err := h.repo.ListByChannel(c.Request.Context(), channelID, before, viewerID, limit)
	if errors.Is(err, ErrChannelNotFound) {
		response.Error(c, response.CodeNotFound, "频道不存在")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, listData(list, hasMore))
}

// Create 处理 POST /api/channels/:id/messages，发布主消息（需登录 + 频道成员）。
func (h *Handler) Create(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized)
		return
	}
	channelID, ok := parseID(c)
	if !ok {
		return
	}
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeValidationError, "消息内容格式不正确")
		return
	}

	msg, err := h.repo.Create(c.Request.Context(), channelID, nil, uid, req)
	if !h.writeCreateError(c, err) {
		return
	}
	msg.mask()
	if h.bc != nil {
		h.bc.BroadcastNewMessage(channelID, msg)
		// 带情绪的消息触发情绪看板重聚合与推送。
		if req.Emotion != "" {
			h.bc.BroadcastEmotionUpdate(channelID)
		}
	}
	response.Success(c, msg)
}

// Threads 处理 GET /api/messages/:id/threads，返回线程回复（游标分页）。
// query: after=<消息ID>（返回更晚的），limit。
func (h *Handler) Threads(c *gin.Context) {
	parentID, ok := parseID(c)
	if !ok {
		return
	}
	after := parseCursor(c, "after")
	limit := parseLimit(c)
	viewerID, _ := middleware.CurrentUserID(c)

	list, hasMore, err := h.repo.ListThreads(c.Request.Context(), parentID, after, viewerID, limit)
	if errors.Is(err, ErrMessageNotFound) {
		response.Error(c, response.CodeNotFound, "消息不存在")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, listData(list, hasMore))
}

// Reply 处理 POST /api/messages/:id/replies，线程回复（需登录 + 频道成员）。
func (h *Handler) Reply(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized)
		return
	}
	parentID, ok := parseID(c)
	if !ok {
		return
	}
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeValidationError, "回复内容格式不正确")
		return
	}

	msg, err := h.repo.CreateReply(c.Request.Context(), parentID, uid, req)
	if !h.writeCreateError(c, err) {
		return
	}
	msg.mask()
	if h.bc != nil {
		h.bc.BroadcastThreadReply(msg.ChannelID, parentID, msg)
	}
	response.Success(c, msg)
}

// Delete 处理 DELETE /api/messages/:id，软删除（需登录 + 作者或频道管理员）。
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

	channelID, err := h.repo.SoftDelete(c.Request.Context(), id, uid)
	switch {
	case errors.Is(err, ErrMessageNotFound):
		response.Error(c, response.CodeNotFound, "消息不存在")
		return
	case errors.Is(err, ErrForbidden):
		response.Error(c, response.CodeForbidden, "只能删除自己的消息")
		return
	case err != nil:
		response.Fail(c, response.CodeInternalError)
		return
	}
	if h.bc != nil {
		h.bc.BroadcastMessageDeleted(channelID, id)
	}
	response.Success(c, gin.H{"message": "已删除"})
}

// writeCreateError 处理创建类操作的错误，返回 true 表示无错误可继续。
func (h *Handler) writeCreateError(c *gin.Context, err error) bool {
	switch {
	case errors.Is(err, ErrChannelNotFound):
		response.Error(c, response.CodeNotFound, "频道不存在")
	case errors.Is(err, ErrMessageNotFound):
		response.Error(c, response.CodeNotFound, "消息不存在")
	case errors.Is(err, ErrNotMember):
		response.Error(c, response.CodeForbidden, "请先加入频道再发言")
	case errors.Is(err, ErrInvalidEmotion):
		response.Error(c, response.CodeValidationError, "情绪标签无效")
	case err != nil:
		response.Fail(c, response.CodeInternalError)
	default:
		return true
	}
	return false
}

// listData 组装游标分页的响应结构。next_cursor 指向本页最后一条消息 ID。
func listData(list []*Message, hasMore bool) gin.H {
	items := nonNil(list)
	var next int64
	if hasMore && len(items) > 0 {
		next = items[len(items)-1].ID
	}
	return gin.H{"list": items, "has_more": hasMore, "next_cursor": next}
}

// parseID 解析路径参数 :id，非法时写入 404 响应并返回 false。
func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, response.CodeNotFound, "资源不存在")
		return 0, false
	}
	return id, true
}

// parseCursor 解析游标 query（before/after），非法或缺省返回 0。
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

// nonNil 保证空列表序列化为 [] 而非 null。
func nonNil[T any](list []T) []T {
	if list == nil {
		return []T{}
	}
	return list
}
