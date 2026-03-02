package handler

import (
	"strconv"

	"github.com/Alike/internal/user"
	"github.com/Alike/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *user.Service
}

func NewUserHandler(userService *user.Service) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")
	
	user, err := h.userService.GetByID(userID)
	if err != nil {
		response.InternalError(c, "Failed to get user")
		return
	}
	
	response.Success(c, user)
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := c.GetString("user_id")
	
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	
	user, err := h.userService.GetByID(userID)
	if err != nil {
		response.InternalError(c, "Failed to get user")
		return
	}
	
	// Update fields (simplified)
	if nickname, ok := req["nickname"].(string); ok {
		user.Nickname = nickname
	}
	if bio, ok := req["bio"].(string); ok {
		user.Bio = bio
	}
	
	if err := h.userService.Update(user); err != nil {
		response.InternalError(c, "Failed to update user")
		return
	}
	
	response.Success(c, user)
}

func (h *UserHandler) GetNearby(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.Query("radius")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	
	lat, _ := strconv.ParseFloat(latStr, 64)
	lng, _ := strconv.ParseFloat(lngStr, 64)
	radius, _ := strconv.ParseFloat(radiusStr, 64)
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	
	if radius == 0 {
		radius = 10.0 // default 10km
	}
	
	users, err := h.userService.FindNearby(lat, lng, radius, page, limit)
	if err != nil {
		response.InternalError(c, "Failed to get nearby users")
		return
	}
	
	response.Success(c, gin.H{
		"users": users,
		"page":  page,
		"limit": limit,
	})
}
