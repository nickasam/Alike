package handler

import (
	"github.com/Alike/internal/domain"
	"github.com/Alike/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterWithAutoJoin 注册并自动加入全局聊天室
func (h *AuthHandler) RegisterWithAutoJoin(c *gin.Context) {
	var req struct {
		Phone           string `json:"phone" binding:"required"`
		VerificationCode string `json:"verification_code" binding:"required"`
		Nickname        string `json:"nickname" binding:"required,min=1,max=50"`
		Password        string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("参数错误："+err.Error()))
		return
	}

	// 检查手机号是否已注册
	exists, err := h.UserRepo.ExistsByPhone(req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("检查用户失败"))
		return
	}
	if exists {
		c.JSON(http.StatusConflict, response.Error("该手机号已注册"))
		return
	}

	// 创建用户
	user := &domain.User{
		ID:        generateUserID(),
		Phone:     req.Phone,
		Nickname:  req.Nickname,
		Password:  hashPassword(req.Password),
		Avatar:    "",
		Bio:       "",
		Latitude:  0,
		Longitude: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.UserRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("创建用户失败"))
		return
	}

	// 生成Token
	accessToken, err := h.authService.GenerateToken(user.ID, "access")
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("生成Token失败"))
		return
	}

	refreshToken, err := h.authService.GenerateToken(user.ID, "refresh")
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("生成Token失败"))
		return
	}

	// 自动加入全局聊天室
	go func() {
		// 异步加入，不阻塞注册响应
		if h.globalChatRepo != nil {
			h.globalChatRepo.JoinRoom(user.ID, "global")
		}
	}()

	// 发送欢迎消息到全局聊天室
	go func() {
		if h.globalChatRepo != nil {
			welcomeMsg := &domain.GlobalMessage{
				ID:        generateMessageID(),
				RoomID:    "global",
				UserID:    "system",
				Username:  "系统消息",
				Content:   "欢迎 " + user.Nickname + " 加入Alike大家庭！🎉",
				CreatedAt: time.Now(),
			}
			h.globalChatRepo.CreateMessage(welcomeMsg)
		}
	}()

	c.JSON(http.StatusCreated, response.Success(gin.H{
		"user": gin.H{
			"id":       user.ID,
			"phone":    user.Phone,
			"nickname": user.Nickname,
		},
		"tokens": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	}))
}
