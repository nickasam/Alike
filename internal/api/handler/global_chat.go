package handler

import (
	"github.com/Alike/internal/domain"
	"github.com/Alike/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type GlobalChatHandler struct {
	GlobalChatRepo interface {
		CreateRoom(room *domain.GlobalChatRoom) error
		GetRoom(id string) (*domain.GlobalChatRoom, error)
		CreateMessage(msg *domain.GlobalMessage) error
		GetMessages(roomID string, limit int) ([]domain.GlobalMessage, error)
		JoinRoom(userID, roomID string) error
	}
}

// CreateGlobalRoom 创建全局聊天室
func (h *GlobalChatHandler) CreateGlobalRoom(c *gin.Context) {
	room := &domain.GlobalChatRoom{
		ID:          "global",
		Name:        "Alike大家庭",
		Description: "欢迎来到Alike大家庭！这里是所有用户的聊天空间，认识新朋友，分享生活点滴。",
		MaxMembers:  1000,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.GlobalChatRepo.CreateRoom(room); err != nil {
		response.Error(c, "CREATE_ROOM_FAILED", "创建聊天室失败", http.StatusInternalServerError)
		return
	}

	response.Success(c, room)
}

// GetGlobalRoom 获取全局聊天室信息
func (h *GlobalChatHandler) GetGlobalRoom(c *gin.Context) {
	room, err := h.GlobalChatRepo.GetRoom("global")
	if err != nil {
		response.Error(c, "ROOM_NOT_FOUND", "聊天室不存在", http.StatusNotFound)
		return
	}

	response.Success(c, room)
}

// GetGlobalMessages 获取全局聊天消息
func (h *GlobalChatHandler) GetGlobalMessages(c *gin.Context) {
	messages, err := h.GlobalChatRepo.GetMessages("global", 100)
	if err != nil {
		response.Error(c, "GET_MESSAGES_FAILED", "获取消息失败", http.StatusInternalServerError)
		return
	}

	response.Success(c, messages)
}

// SendGlobalMessage 发送全局消息
func (h *GlobalChatHandler) SendGlobalMessage(c *gin.Context) {
	userID := c.GetString("user_id")
	username := c.GetString("username")

	var req struct {
		Content string `json:"content" binding:"required,min=1,max=500"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "INVALID_CONTENT", "消息内容不能为空且最多500字", http.StatusBadRequest)
		return
	}

	message := &domain.GlobalMessage{
		ID:        time.Now().Format("20060102150405") + randomString(6),
		RoomID:    "global",
		UserID:    userID,
		Username:  username,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := h.GlobalChatRepo.CreateMessage(message); err != nil {
		response.Error(c, "SEND_FAILED", "发送消息失败", http.StatusInternalServerError)
		return
	}

	response.Success(c, message)
}

// JoinGlobalRoom 加入全局聊天室
func (h *GlobalChatHandler) JoinGlobalRoom(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.GlobalChatRepo.JoinRoom(userID, "global"); err != nil {
		response.Error(c, "JOIN_FAILED", "加入聊天室失败", http.StatusInternalServerError)
		return
	}

	response.Success(c, gin.H{
		"message": "成功加入全局聊天室",
	})
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
