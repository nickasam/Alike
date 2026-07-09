package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/internal/middleware"
	"github.com/Alike/backend/pkg/jwt"
	"github.com/Alike/backend/pkg/response"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func testManager() *jwt.Manager {
	return jwt.NewManager("test-secret", time.Hour, 24*time.Hour)
}

// newHandler 构造仅依赖 JWT 的 handler；repo 为 nil，仅用于在
// 校验/令牌类型检查阶段就返回、不触达数据库的用例。
func newHandler() *Handler {
	return NewHandler(nil, testManager())
}

// doJSON 发起一次 JSON 请求并返回解析后的响应体与状态码。
func doJSON(h gin.HandlerFunc, method, path, body string) (int, response.Body) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	h(c)

	var parsed response.Body
	_ = json.Unmarshal(w.Body.Bytes(), &parsed)
	return w.Code, parsed
}

func TestRegisterValidation(t *testing.T) {
	h := newHandler()
	cases := []struct {
		name string
		body string
	}{
		{"missing email", `{"password":"secret6","nickname":"牛马"}`},
		{"invalid email", `{"email":"not-an-email","password":"secret6","nickname":"牛马"}`},
		{"short password", `{"email":"a@b.com","password":"123","nickname":"牛马"}`},
		{"missing nickname", `{"email":"a@b.com","password":"secret6"}`},
		{"malformed json", `{`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, body := doJSON(h.Register, http.MethodPost, "/api/auth/register", tc.body)
			if code != http.StatusUnprocessableEntity {
				t.Fatalf("status = %d, want 422", code)
			}
			if body.Code != response.CodeValidationError {
				t.Errorf("code = %d, want %d", body.Code, response.CodeValidationError)
			}
		})
	}
}

func TestLoginValidation(t *testing.T) {
	h := newHandler()
	code, body := doJSON(h.Login, http.MethodPost, "/api/auth/login", `{"email":"bad","password":"123"}`)
	if code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", code)
	}
	if body.Code != response.CodeValidationError {
		t.Errorf("code = %d, want %d", body.Code, response.CodeValidationError)
	}
}

func TestRefreshMissingToken(t *testing.T) {
	h := newHandler()
	code, body := doJSON(h.Refresh, http.MethodPost, "/api/auth/refresh", `{}`)
	if code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", code)
	}
	if body.Code != response.CodeValidationError {
		t.Errorf("code = %d, want %d", body.Code, response.CodeValidationError)
	}
}

func TestRefreshRejectsAccessToken(t *testing.T) {
	h := newHandler()
	// 用 access token 冒充 refresh token，应被类型校验拒绝（早于 DB 访问）。
	access, err := h.jwt.GenerateAccess(42)
	if err != nil {
		t.Fatalf("generate access: %v", err)
	}
	code, body := doJSON(h.Refresh, http.MethodPost, "/api/auth/refresh",
		`{"refresh_token":"`+access+`"}`)
	if code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", code)
	}
	if body.Code != response.CodeUnauthorized {
		t.Errorf("code = %d, want %d", body.Code, response.CodeUnauthorized)
	}
}

func TestRefreshRejectsGarbageToken(t *testing.T) {
	h := newHandler()
	code, body := doJSON(h.Refresh, http.MethodPost, "/api/auth/refresh",
		`{"refresh_token":"not.a.jwt"}`)
	if code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", code)
	}
	if body.Code != response.CodeUnauthorized {
		t.Errorf("code = %d, want %d", body.Code, response.CodeUnauthorized)
	}
}

func TestLogout(t *testing.T) {
	h := newHandler()
	code, body := doJSON(h.Logout, http.MethodPost, "/api/auth/logout", ``)
	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	if body.Code != response.CodeSuccess {
		t.Errorf("code = %d, want 0", body.Code)
	}
}

func TestMeWithoutContextUser(t *testing.T) {
	h := newHandler()
	// 未经中间件设置 user_id，应返回未认证。
	code, body := doJSON(h.Me, http.MethodGet, "/api/auth/me", ``)
	if code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", code)
	}
	if body.Code != response.CodeUnauthorized {
		t.Errorf("code = %d, want %d", body.Code, response.CodeUnauthorized)
	}
}

func TestIssueTokensRoundTrip(t *testing.T) {
	h := newHandler()
	tokens, err := h.issueTokens(7)
	if err != nil {
		t.Fatalf("issueTokens: %v", err)
	}

	access, err := h.jwt.Parse(tokens.AccessToken)
	if err != nil {
		t.Fatalf("parse access: %v", err)
	}
	if access.Type != jwt.AccessToken || access.UserID != 7 {
		t.Errorf("access claims = %+v", access)
	}

	refresh, err := h.jwt.Parse(tokens.RefreshToken)
	if err != nil {
		t.Fatalf("parse refresh: %v", err)
	}
	if refresh.Type != jwt.RefreshToken || refresh.UserID != 7 {
		t.Errorf("refresh claims = %+v", refresh)
	}
}

func TestMeReadsContextUserID(t *testing.T) {
	// 直接验证中间件与 handler 之间的上下文契约。
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.ContextUserID, int64(99))
	if id, ok := middleware.CurrentUserID(c); !ok || id != 99 {
		t.Fatalf("CurrentUserID = %d, %v; want 99, true", id, ok)
	}
}
