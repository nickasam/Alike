package message

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
	return NewHandler(nil, nil, nil, nil)
}

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
	code, body := invoke(h.Create, http.MethodPost, "/api/channels/1/messages", `{}`, withParamID("1"))
	if code != http.StatusUnauthorized || body.Code != response.CodeUnauthorized {
		t.Fatalf("status=%d code=%d, want 401/%d", code, body.Code, response.CodeUnauthorized)
	}
}

func TestReplyRequiresAuth(t *testing.T) {
	h := newHandler()
	code, body := invoke(h.Reply, http.MethodPost, "/api/messages/1/replies", `{}`, withParamID("1"))
	if code != http.StatusUnauthorized || body.Code != response.CodeUnauthorized {
		t.Fatalf("status=%d code=%d, want 401/%d", code, body.Code, response.CodeUnauthorized)
	}
}

func TestDeleteRequiresAuth(t *testing.T) {
	h := newHandler()
	code, body := invoke(h.Delete, http.MethodDelete, "/api/messages/1", ``, withParamID("1"))
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
		{"blank content", `{"content":""}`},
		{"malformed json", `{`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, body := invoke(h.Create, http.MethodPost, "/api/channels/1/messages", tc.body,
				func(c *gin.Context) { withUser(1)(c); withParamID("1")(c) })
			if code != http.StatusUnprocessableEntity || body.Code != response.CodeValidationError {
				t.Fatalf("status=%d code=%d, want 422/%d", code, body.Code, response.CodeValidationError)
			}
		})
	}
}

func TestCreateInvalidChannelID(t *testing.T) {
	h := newHandler()
	code, body := invoke(h.Create, http.MethodPost, "/api/channels/x/messages", `{"content":"hi"}`,
		func(c *gin.Context) { withUser(1)(c); withParamID("x")(c) })
	if code != http.StatusNotFound || body.Code != response.CodeNotFound {
		t.Fatalf("status=%d code=%d, want 404/%d", code, body.Code, response.CodeNotFound)
	}
}

func TestListInvalidChannelID(t *testing.T) {
	h := newHandler()
	code, body := invoke(h.List, http.MethodGet, "/api/channels/0/messages", ``, withParamID("0"))
	if code != http.StatusNotFound || body.Code != response.CodeNotFound {
		t.Fatalf("status=%d code=%d, want 404/%d", code, body.Code, response.CodeNotFound)
	}
}

func TestParseLimitDefaultsAndCaps(t *testing.T) {
	cases := []struct {
		query string
		want  int
	}{
		{"", defaultLimit},
		{"?limit=10", 10},
		{"?limit=0", defaultLimit},
		{"?limit=-5", defaultLimit},
		{"?limit=999", maxLimit},
		{"?limit=abc", defaultLimit},
	}
	for _, tc := range cases {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/x"+tc.query, nil)
		if got := parseLimit(c); got != tc.want {
			t.Errorf("query=%q -> %d, want %d", tc.query, got, tc.want)
		}
	}
}

func TestParseCursor(t *testing.T) {
	cases := []struct {
		query string
		want  int64
	}{
		{"", 0},
		{"?before=42", 42},
		{"?before=-1", 0},
		{"?before=abc", 0},
	}
	for _, tc := range cases {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/x"+tc.query, nil)
		if got := parseCursor(c, "before"); got != tc.want {
			t.Errorf("query=%q -> %d, want %d", tc.query, got, tc.want)
		}
	}
}

func TestMaskAnonymousHidesAuthor(t *testing.T) {
	m := &Message{IsAnonymous: true, Author: &Author{ID: 7, Nickname: "牛马"}}
	m.mask()
	if m.Author != nil {
		t.Fatalf("anonymous message must not expose author, got %+v", m.Author)
	}
}

func TestMaskDeletedHidesContent(t *testing.T) {
	m := &Message{IsDeleted: true, Content: "秘密", Emotion: "累", Author: &Author{ID: 1}}
	m.mask()
	if m.Content != deletedPlaceholder {
		t.Errorf("content=%q, want placeholder", m.Content)
	}
	if m.Author != nil || m.Emotion != "" {
		t.Errorf("deleted message must hide author/emotion, got author=%v emotion=%q", m.Author, m.Emotion)
	}
}

func TestListDataNextCursor(t *testing.T) {
	list := []*Message{{ID: 5}, {ID: 3}}
	data := listData(list, true)
	if data["next_cursor"].(int64) != 3 {
		t.Errorf("next_cursor=%v, want 3", data["next_cursor"])
	}
	if data["has_more"].(bool) != true {
		t.Errorf("has_more=%v, want true", data["has_more"])
	}
	// 无更多时 next_cursor 为 0。
	if listData(list, false)["next_cursor"].(int64) != 0 {
		t.Error("next_cursor should be 0 when no more pages")
	}
}
