package pulse

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/response"
)

// Handler 承载 pulse 模块的路由。全部公开可读（无 authMW）。
type Handler struct {
	repo *Repository
}

// NewHandler 创建 pulse handler。
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// ListTopics 处理 GET /api/pulse/topics，返回所有活跃专题（前端 Tab 数据源）。
func (h *Handler) ListTopics(c *gin.Context) {
	topics, err := h.repo.ListActiveTopics(c.Request.Context())
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	// 保证空列表也返回 [] 而非 null（前端 v-for 依赖）
	if topics == nil {
		topics = []*Topic{}
	}
	response.Success(c, gin.H{"list": topics})
}

// GetItems 处理 GET /api/pulse/topics/:slug/items，返回某专题的 Top-N 条目。
func (h *Handler) GetItems(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.Fail(c, response.CodeNotFound)
		return
	}
	ctx := c.Request.Context()

	topic, err := h.repo.GetTopicBySlug(ctx, slug)
	switch {
	case errors.Is(err, ErrTopicNotFound):
		response.Error(c, response.CodeNotFound, "专题不存在")
		return
	case err != nil:
		response.Fail(c, response.CodeInternalError)
		return
	}

	items, err := h.repo.ListItemsByTopic(ctx, topic.ID, 0)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	if items == nil {
		items = []*Item{}
	}

	// dereference 到值切片，避免前端拿到指针数组的 null 元素
	list := make([]Item, len(items))
	for i, it := range items {
		list[i] = *it
	}

	response.Success(c, ItemsResponse{
		Topic:         topic,
		List:          list,
		LastFetchedAt: topic.LastFetchedAt,
		Stale:         isStale(topic.LastFetchedAt),
	})
}

// isStale 判定当前时间距最后成功抓取是否超过阈值（M0 数据都未抓，返回 true 属预期）。
func isStale(last *time.Time) bool {
	if last == nil {
		return true
	}
	return time.Since(*last) > staleThreshold
}
