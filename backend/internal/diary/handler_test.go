package diary

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/internal/middleware"
	"github.com/Alike/backend/pkg/httputil"
	"github.com/Alike/backend/pkg/response"
)

func init() {
	gin.SetMode(gin.TestMode)
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

func withParam(key, val string) func(*gin.Context) {
	return func(c *gin.Context) { c.Params = append(c.Params, gin.Param{Key: key, Value: val}) }
}

func TestCreateRequiresAuth(t *testing.T) {
	h := NewHandler(nil)
	code, body := invoke(h.Create, http.MethodPost, "/api/diaries", `{"content":"hi"}`, nil)
	if code != http.StatusUnauthorized || body.Code != response.CodeUnauthorized {
		t.Fatalf("status=%d code=%d, want 401/%d", code, body.Code, response.CodeUnauthorized)
	}
}

func TestCreateCommentRequiresAuth(t *testing.T) {
	h := NewHandler(nil)
	code, body := invoke(h.CreateComment, http.MethodPost, "/api/diaries/1/comments", `{"content":"hi"}`, withParam("id", "1"))
	if code != http.StatusUnauthorized || body.Code != response.CodeUnauthorized {
		t.Fatalf("status=%d code=%d, want 401/%d", code, body.Code, response.CodeUnauthorized)
	}
}

func TestGetInvalidID(t *testing.T) {
	h := NewHandler(nil)
	code, body := invoke(h.Get, http.MethodGet, "/api/diaries/x", "", withParam("id", "x"))
	if code != http.StatusNotFound || body.Code != response.CodeNotFound {
		t.Fatalf("status=%d code=%d, want 404/%d", code, body.Code, response.CodeNotFound)
	}
}

func TestStreakInvalidUserID(t *testing.T) {
	h := NewHandler(nil)
	code, body := invoke(h.Streak, http.MethodGet, "/api/diaries/streak/0", "", withParam("user_id", "0"))
	if code != http.StatusNotFound || body.Code != response.CodeNotFound {
		t.Fatalf("status=%d code=%d, want 404/%d", code, body.Code, response.CodeNotFound)
	}
}

func TestCreateValidationError(t *testing.T) {
	h := NewHandler(nil)
	code, body := invoke(h.Create, http.MethodPost, "/api/diaries", `{"content":""}`, withUser(1))
	if code != http.StatusUnprocessableEntity || body.Code != response.CodeValidationError {
		t.Fatalf("status=%d code=%d, want 422/%d", code, body.Code, response.CodeValidationError)
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
		{"?limit=9999", maxLimit},
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

func TestPaginateDefaultsAndCaps(t *testing.T) {
	cases := []struct {
		query            string
		wantPage, wantPS int
	}{
		{"", httputil.DefaultPage, httputil.DefaultPageSize},
		{"?page=3&page_size=10", 3, 10},
		{"?page=0&page_size=0", httputil.DefaultPage, httputil.DefaultPageSize},
		{"?page=-1&page_size=999", httputil.DefaultPage, httputil.MaxPageSize},
	}
	for _, tc := range cases {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/x"+tc.query, nil)
		page, ps := httputil.Paginate(c)
		if page != tc.wantPage || ps != tc.wantPS {
			t.Errorf("query=%q -> page=%d ps=%d, want %d/%d", tc.query, page, ps, tc.wantPage, tc.wantPS)
		}
	}
}

func TestCommentMask(t *testing.T) {
	anon := &Comment{Content: "secret", IsAnonymous: true, Author: &Author{ID: 1, Nickname: "n"}}
	anon.mask()
	if anon.Author != nil {
		t.Errorf("anonymous comment should hide author")
	}
	if anon.Content != "secret" {
		t.Errorf("anonymous comment content should be preserved")
	}

	del := &Comment{Content: "gone", IsDeleted: true, Author: &Author{ID: 1}}
	del.mask()
	if del.Author != nil || del.Content != deletedPlaceholder {
		t.Errorf("deleted comment should mask content and author, got content=%q author=%v", del.Content, del.Author)
	}
}
