package channel

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/internal/middleware"
	"github.com/Alike/backend/pkg/response"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// newHandler 构造 repo 为 nil 的 handler，仅用于在触达数据库前
// 就返回的用例（鉴权、参数校验、路径参数解析）。
func newHandler() *Handler {
	return NewHandler(nil)
}

// invoke 以给定的上下文准备函数发起请求并返回状态码与响应体。
func invoke(h gin.HandlerFunc, method, target, body string, prep func(*gin.Context)) (int, response.Body) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, target, strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	if prep != nil {
		prep(c)
	}
	h(c)

	var parsed response.Body
	_ = json.Unmarshal(w.Body.Bytes(), &parsed)
	return w.Code, parsed
}

func withUser(id int64) func(*gin.Context) {
	return func(c *gin.Context) { c.Set(middleware.ContextUserID, id) }
}

func withParamID(id string) func(*gin.Context) {
	return func(c *gin.Context) { c.Params = gin.Params{{Key: "id", Value: id}} }
}

func TestCreateRequiresAuth(t *testing.T) {
	h := newHandler()
	code, body := invoke(h.Create, http.MethodPost, "/api/channels", `{}`, nil)
	if code != http.StatusUnauthorized || body.Code != response.CodeUnauthorized {
		t.Fatalf("status=%d code=%d, want 401/%d", code, body.Code, response.CodeUnauthorized)
	}
}

func TestCreateValidation(t *testing.T) {
	h := newHandler()
	cases := []struct {
		name string
		body string
	}{
		{"empty", `{}`},
		{"missing slug", `{"name":"互联网","category":"industry"}`},
		{"bad category", `{"name":"互联网","slug":"it","category":"nonsense"}`},
		{"malformed json", `{`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, body := invoke(h.Create, http.MethodPost, "/api/channels", tc.body, withUser(1))
			if code != http.StatusUnprocessableEntity || body.Code != response.CodeValidationError {
				t.Fatalf("status=%d code=%d, want 422/%d", code, body.Code, response.CodeValidationError)
			}
		})
	}
}

func TestGetInvalidID(t *testing.T) {
	h := newHandler()
	for _, id := range []string{"abc", "0", "-3"} {
		code, body := invoke(h.Get, http.MethodGet, "/api/channels/"+id, ``, withParamID(id))
		if code != http.StatusNotFound || body.Code != response.CodeNotFound {
			t.Fatalf("id=%q status=%d code=%d, want 404/%d", id, code, body.Code, response.CodeNotFound)
		}
	}
}

func TestJoinRequiresAuth(t *testing.T) {
	h := newHandler()
	code, body := invoke(h.Join, http.MethodPost, "/api/channels/1/join", ``, withParamID("1"))
	if code != http.StatusUnauthorized || body.Code != response.CodeUnauthorized {
		t.Fatalf("status=%d code=%d, want 401/%d", code, body.Code, response.CodeUnauthorized)
	}
}

func TestLeaveRequiresAuth(t *testing.T) {
	h := newHandler()
	code, body := invoke(h.Leave, http.MethodPost, "/api/channels/1/leave", ``, withParamID("1"))
	if code != http.StatusUnauthorized || body.Code != response.CodeUnauthorized {
		t.Fatalf("status=%d code=%d, want 401/%d", code, body.Code, response.CodeUnauthorized)
	}
}

func TestJoinInvalidID(t *testing.T) {
	h := newHandler()
	code, body := invoke(h.Join, http.MethodPost, "/api/channels/x/join", ``,
		func(c *gin.Context) { withUser(1)(c); withParamID("x")(c) })
	if code != http.StatusNotFound || body.Code != response.CodeNotFound {
		t.Fatalf("status=%d code=%d, want 404/%d", code, body.Code, response.CodeNotFound)
	}
}

func TestEmotionBoardReturnsEmpty(t *testing.T) {
	h := newHandler()
	code, body := invoke(h.EmotionBoard, http.MethodGet, "/api/channels/1/emotion-board", ``, withParamID("1"))
	if code != http.StatusOK || body.Code != response.CodeSuccess {
		t.Fatalf("status=%d code=%d, want 200/0", code, body.Code)
	}
	data, ok := body.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("data type=%T, want object", body.Data)
	}
	if _, ok := data["emotions"]; !ok {
		t.Errorf("missing emotions key in %v", data)
	}
}

func TestPaginateDefaultsAndCaps(t *testing.T) {
	cases := []struct {
		query            string
		wantPage, wantPS int
	}{
		{"", defaultPage, defaultPageSize},
		{"?page=3&page_size=10", 3, 10},
		{"?page=0&page_size=0", defaultPage, defaultPageSize},
		{"?page=-1&page_size=999", defaultPage, maxPageSize},
		{"?page=abc&page_size=xyz", defaultPage, defaultPageSize},
	}
	for _, tc := range cases {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/api/channels"+tc.query, nil)
		page, ps := paginate(c)
		if page != tc.wantPage || ps != tc.wantPS {
			t.Errorf("query=%q -> page=%d,ps=%d; want %d,%d", tc.query, page, ps, tc.wantPage, tc.wantPS)
		}
	}
}

func TestNonNil(t *testing.T) {
	if got := nonNil[*Channel](nil); got == nil || len(got) != 0 {
		t.Fatalf("nonNil(nil) = %v, want empty non-nil slice", got)
	}
}
