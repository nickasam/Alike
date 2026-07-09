// Package auth 负责注册、登录、JWT 签发与验证。
package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/Alike/backend/internal/middleware"
	"github.com/Alike/backend/pkg/jwt"
	"github.com/Alike/backend/pkg/response"
)

// Handler 承载 auth 模块的依赖，替代阶段一的无状态包级函数。
type Handler struct {
	repo *Repository
	jwt  *jwt.Manager
}

// NewHandler 创建 auth handler。
func NewHandler(repo *Repository, mgr *jwt.Manager) *Handler {
	return &Handler{repo: repo, jwt: mgr}
}

// Register 处理 POST /api/auth/register。
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeValidationError, "邮箱、密码（≥6位）或昵称格式不正确")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}

	u, err := h.repo.Create(c.Request.Context(), req.Email, string(hash), req.Nickname)
	if errors.Is(err, ErrEmailConflict) {
		response.Error(c, response.CodeConflict, "该邮箱已被注册")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}

	tokens, err := h.issueTokens(u.ID)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, AuthResponse{User: u, Tokens: tokens})
}

// Login 处理 POST /api/auth/login。
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeValidationError, "邮箱或密码格式不正确")
		return
	}

	u, err := h.repo.GetByEmail(c.Request.Context(), req.Email)
	if errors.Is(err, ErrUserNotFound) {
		// 不区分"用户不存在"与"密码错误"，避免邮箱枚举。
		response.Error(c, response.CodeUnauthorized, "邮箱或密码错误")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) != nil {
		response.Error(c, response.CodeUnauthorized, "邮箱或密码错误")
		return
	}

	tokens, err := h.issueTokens(u.ID)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, AuthResponse{User: u, Tokens: tokens})
}

// Refresh 处理 POST /api/auth/refresh，用 refresh token 换取新的 token 对。
func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeValidationError, "缺少 refresh_token")
		return
	}

	claims, err := h.jwt.Parse(req.RefreshToken)
	if err != nil || claims.Type != jwt.RefreshToken {
		response.Error(c, response.CodeUnauthorized, "refresh token 无效或已过期")
		return
	}

	// 确认用户仍存在，避免为已删除用户续签。
	if _, err := h.repo.GetByID(c.Request.Context(), claims.UserID); err != nil {
		response.Error(c, response.CodeUnauthorized, "用户不存在")
		return
	}

	tokens, err := h.issueTokens(claims.UserID)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, tokens)
}

// Logout 处理 POST /api/auth/logout。JWT 无状态，服务端无操作，前端清除本地 token 即可。
func (h *Handler) Logout(c *gin.Context) {
	response.Success(c, gin.H{"message": "已登出"})
}

// Me 处理 GET /api/auth/me，返回当前登录用户信息。需经 Auth 中间件保护。
func (h *Handler) Me(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized)
		return
	}

	u, err := h.repo.GetByID(c.Request.Context(), uid)
	if errors.Is(err, ErrUserNotFound) {
		response.Error(c, response.CodeUnauthorized, "用户不存在")
		return
	}
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	response.Success(c, u)
}

// issueTokens 为指定用户签发 access + refresh token 对。
func (h *Handler) issueTokens(userID int64) (*TokenPair, error) {
	access, err := h.jwt.GenerateAccess(userID)
	if err != nil {
		return nil, err
	}
	refresh, err := h.jwt.GenerateRefresh(userID)
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}
