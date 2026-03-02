package handler

import (
	"github.com/Alike/internal/match"
	"github.com/Alike/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	matchService *match.Service
}

func NewMatchHandler(matchService *match.Service) *MatchHandler {
	return &MatchHandler{matchService: matchService}
}

func (h *MatchHandler) ListMatches(c *gin.Context) {
	userID := c.GetString("user_id")
	
	matches, err := h.matchService.GetMatches(userID)
	if err != nil {
		response.InternalError(c, "Failed to get matches")
		return
	}
	
	response.Success(c, matches)
}

func (h *MatchHandler) GetMatch(c *gin.Context) {
	id := c.Param("id")
	
	match, err := h.matchService.GetMatch(id)
	if err != nil {
		response.NotFound(c, "Match not found")
		return
	}
	
	response.Success(c, match)
}

func (h *MatchHandler) LikeUser(c *gin.Context) {
	userID := c.GetString("user_id")
	likedID := c.Param("id")
	
	if err := h.matchService.CreateLike(userID, likedID); err != nil {
		response.InternalError(c, "Failed to like user")
		return
	}
	
	response.Success(c, gin.H{"message": "User liked"})
}
