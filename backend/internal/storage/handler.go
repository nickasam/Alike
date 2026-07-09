package storage

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/internal/middleware"
	"github.com/Alike/backend/pkg/response"
)

// Handler 承载 storage 模块的上传接口依赖。
type Handler struct {
	store *Store
}

// NewHandler 创建 storage handler。
func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

// Upload 处理 POST /api/upload（需登录）。
// multipart/form-data，字段名 file。服务端做 MIME 校验与大小限制。
func (h *Handler) Upload(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Fail(c, response.CodeUnauthorized)
		return
	}
	if h.store == nil || h.store.client == nil {
		response.Error(c, response.CodeInternalError, "对象存储未配置")
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, response.CodeValidationError, "缺少上传文件（字段名 file）")
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}
	defer f.Close()

	result, err := h.store.Upload(c.Request.Context(), f, fileHeader.Size, strconv.FormatInt(uid, 10))
	switch {
	case errors.Is(err, ErrFileTooLarge):
		response.Error(c, response.CodeValidationError, "文件超过大小限制（图片≤5MB，文档≤10MB）")
		return
	case errors.Is(err, ErrUnsupportedType):
		response.Error(c, response.CodeValidationError, "不支持的文件类型（仅 jpg/png/gif/webp/pdf/zip/txt）")
		return
	case errors.Is(err, ErrEmptyFile):
		response.Error(c, response.CodeValidationError, "文件为空")
		return
	case errors.Is(err, ErrStorageDisabled):
		response.Error(c, response.CodeInternalError, "对象存储未配置")
		return
	case err != nil:
		response.Fail(c, response.CodeInternalError)
		return
	}

	response.Success(c, result)
}
