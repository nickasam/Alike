package handler

import (
	"strconv"

	"github.com/Alike/internal/chat"
	"github.com/Alike/internal/domain"
	"github.com/Alike/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService *chat.Service
}

func NewChatHandler(chatService *chat.Service) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) ListChats(c *gin.Context) {
	userID := c.GetString("user_id")
	
	chats, err := h.chatService.GetChats(userID)
	if err != nil {
		response.InternalError(c, "Failed to get chats")
		return
	}
	
	response.Success(c, chats)
}

func (h *ChatHandler) GetChat(c *gin.Context) {
	id := c.Param("id")
	
	chat, err := h.chatService.GetChat(id)
	if err != nil {
		response.NotFound(c, "Chat not found")
		return
	}
	
	response.Success(c, chat)
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	chatID := c.Param("id")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")
	
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	
	messages, err := h.chatService.GetMessages(chatID, page, limit)
	if err != nil {
		response.InternalError(c, "Failed to get messages")
		return
	}
	
	response.Success(c, messages)
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	chatID := c.Param("id")
	userID := c.GetString("user_id")
	
	var req struct {
		Content string `json:"content" binding:"required"`
		Type    string `json:"type"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	
	if req.Type == "" {
		req.Type = "text"
	}
	
	message := &domain.Message{
		ChatID:      chatID,
		SenderID:    userID,
		Content:     req.Content,
		MessageType: req.Type,
	}
	
	if err := h.chatService.CreateMessage(message); err != nil {
		response.InternalError(c, "Failed to send message")
		return
	}
	
	response.Success(c, message)
}
