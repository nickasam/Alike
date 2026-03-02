package handler

import (
	"strconv"

	"github.com/Alike/internal/notification"
	"github.com/Alike/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *notification.Service
}

func NewNotificationHandler(notificationService *notification.Service) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	userID := c.GetString("user_id")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	
	notifications, err := h.notificationService.GetNotifications(userID, limit, (page-1)*limit)
	if err != nil {
		response.InternalError(c, "Failed to get notifications")
		return
	}
	
	response.Success(c, notifications)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	
	if err := h.notificationService.MarkAsRead(id); err != nil {
		response.InternalError(c, "Failed to mark as read")
		return
	}
	
	response.Success(c, gin.H{"message": "Notification marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	
	if err := h.notificationService.MarkAllAsRead(userID); err != nil {
		response.InternalError(c, "Failed to mark all as read")
		return
	}
	
	response.Success(c, gin.H{"message": "All notifications marked as read"})
}
