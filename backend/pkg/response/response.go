// Package response 提供统一的 API 响应封装与业务错误码。
// 响应格式：{ "code": 0, "message": "success", "data": {...} }
package response

import "github.com/gin-gonic/gin"

// 业务错误码，对齐架构文档 5.3 错误码表。
const (
	CodeSuccess         = 0     // 成功
	CodeBadRequest      = 40000 // 请求参数错误
	CodeUnauthorized    = 40100 // 未认证 / token 失效
	CodeForbidden       = 40300 // 无权限
	CodeNotFound        = 40400 // 资源不存在
	CodeConflict        = 40900 // 冲突（如邮箱已注册、重复共情）
	CodeValidationError = 42200 // 校验失败
	CodeTooManyRequests = 42900 // 触发限流
	CodeInternalError   = 50000 // 服务端内部错误
)

// Body 是统一响应体结构。
type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData 是分页响应的 data 结构。
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// Success 返回 HTTP 200 + code 0 的成功响应。
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Body{Code: CodeSuccess, Message: "success", Data: data})
}

// Page 返回分页成功响应。
func Page(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	Success(c, PageData{List: list, Total: total, Page: page, PageSize: pageSize})
}

// Error 以给定业务错误码返回失败响应，HTTP 状态码由 code 推导。
func Error(c *gin.Context, code int, message string) {
	if message == "" {
		message = defaultMessage(code)
	}
	c.JSON(httpStatus(code), Body{Code: code, Message: message})
}

// Fail 是 Error 的语义别名，使用错误码对应的默认文案。
func Fail(c *gin.Context, code int) {
	Error(c, code, "")
}

// httpStatus 将业务错误码映射为 HTTP 状态码。
func httpStatus(code int) int {
	switch code {
	case CodeSuccess:
		return 200
	case CodeBadRequest:
		return 400
	case CodeUnauthorized:
		return 401
	case CodeForbidden:
		return 403
	case CodeNotFound:
		return 404
	case CodeConflict:
		return 409
	case CodeValidationError:
		return 422
	case CodeTooManyRequests:
		return 429
	default:
		return 500
	}
}

// defaultMessage 返回错误码的默认中文文案。
func defaultMessage(code int) string {
	switch code {
	case CodeBadRequest:
		return "请求参数错误"
	case CodeUnauthorized:
		return "未认证或登录已失效"
	case CodeForbidden:
		return "无权限"
	case CodeNotFound:
		return "资源不存在"
	case CodeConflict:
		return "资源冲突"
	case CodeValidationError:
		return "校验失败"
	case CodeTooManyRequests:
		return "请求过于频繁"
	default:
		return "服务端内部错误"
	}
}
