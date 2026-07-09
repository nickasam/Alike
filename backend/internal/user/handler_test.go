package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/internal/middleware"
	"github.com/Alike/backend/pkg/response"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// doReq 发起请求，可选注入上下文 user_id（模拟鉴权中间件）与路径参数 :id。
func doReq(h gin.HandlerFunc, method, body string, uid *int64, idParam string) (int, response.Body) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, "/api/users/"+idParam, strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: idParam}}
	if uid != nil {
		c.Set(middleware.ContextUserID, *uid)
	}
	h(c)

	var parsed response.Body
	_ = json.Unmarshal(w.Body.Bytes(), &parsed)
	return w.Code, parsed
}

func ptr[T any](v T) *T { return &v }

func TestUpdateRequiresAuth(t *testing.T) {
	h := NewHandler(nil)
	code, body := doReq(h.Update, http.MethodPut, `{"nickname":"牛马"}`, nil, "1")
	if code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", code)
	}
	if body.Code != response.CodeUnauthorized {
		t.Errorf("code = %d, want %d", body.Code, response.CodeUnauthorized)
	}
}

func TestUpdateForbidsOtherUser(t *testing.T) {
	h := NewHandler(nil)
	// 当前用户 1，尝试修改用户 2 的资料，应被拒绝（早于 DB 访问）。
	code, body := doReq(h.Update, http.MethodPut, `{"nickname":"牛马"}`, ptr(int64(1)), "2")
	if code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", code)
	}
	if body.Code != response.CodeForbidden {
		t.Errorf("code = %d, want %d", body.Code, response.CodeForbidden)
	}
}

func TestGetInvalidID(t *testing.T) {
	h := NewHandler(nil)
	code, body := doReq(h.Get, http.MethodGet, ``, nil, "abc")
	if code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", code)
	}
	if body.Code != response.CodeNotFound {
		t.Errorf("code = %d, want %d", body.Code, response.CodeNotFound)
	}
}

func TestUpdateValidation(t *testing.T) {
	h := NewHandler(nil)
	// work_years 为负、自身可改，但校验应在 DB 访问前失败。
	code, body := doReq(h.Update, http.MethodPut, `{"work_years":-1}`, ptr(int64(1)), "1")
	if code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", code)
	}
	if body.Code != response.CodeValidationError {
		t.Errorf("code = %d, want %d", body.Code, response.CodeValidationError)
	}
}

func TestPublicOmitsEmail(t *testing.T) {
	u := &User{
		ID:        7,
		Email:     "secret@example.com",
		Nickname:  "牛马",
		WorkYears: 3,
		CreatedAt: time.Now(),
	}
	pub := u.Public()
	if pub.ID != u.ID || pub.Nickname != u.Nickname || pub.WorkYears != u.WorkYears {
		t.Errorf("public fields mismatch: %+v", pub)
	}
	// 序列化后不应包含 email 字段。
	raw, err := json.Marshal(pub)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(raw), "email") || strings.Contains(string(raw), "secret@example.com") {
		t.Errorf("public JSON leaked email: %s", raw)
	}
}
