package empathy

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

func invoke(h gin.HandlerFunc, method, target string, prep func(*gin.Context)) (int, response.Body) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, target, strings.NewReader(""))
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
	h := NewHandler(nil, nil, nil)
	code, body := invoke(h.Create, http.MethodPost, "/api/messages/1/empathy", withParamID("1"))
	if code != http.StatusUnauthorized || body.Code != response.CodeUnauthorized {
		t.Fatalf("status=%d code=%d, want 401/%d", code, body.Code, response.CodeUnauthorized)
	}
}

func TestDeleteRequiresAuth(t *testing.T) {
	h := NewHandler(nil, nil, nil)
	code, body := invoke(h.Delete, http.MethodDelete, "/api/messages/1/empathy", withParamID("1"))
	if code != http.StatusUnauthorized || body.Code != response.CodeUnauthorized {
		t.Fatalf("status=%d code=%d, want 401/%d", code, body.Code, response.CodeUnauthorized)
	}
}

func TestCreateInvalidMessageID(t *testing.T) {
	h := NewHandler(nil, nil, nil)
	code, body := invoke(h.Create, http.MethodPost, "/api/messages/x/empathy",
		func(c *gin.Context) { withUser(1)(c); withParamID("x")(c) })
	if code != http.StatusNotFound || body.Code != response.CodeNotFound {
		t.Fatalf("status=%d code=%d, want 404/%d", code, body.Code, response.CodeNotFound)
	}
}

func TestUsersInvalidMessageID(t *testing.T) {
	h := NewHandler(nil, nil, nil)
	code, body := invoke(h.Users, http.MethodGet, "/api/messages/0/empathy-users", withParamID("0"))
	if code != http.StatusNotFound || body.Code != response.CodeNotFound {
		t.Fatalf("status=%d code=%d, want 404/%d", code, body.Code, response.CodeNotFound)
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

func TestRankLimitDefaultsAndCaps(t *testing.T) {
	cases := []struct {
		query string
		want  int
	}{
		{"", defaultRankSize},
		{"?limit=10", 10},
		{"?limit=0", defaultRankSize},
		{"?limit=9999", maxRankSize},
		{"?limit=abc", defaultRankSize},
	}
	for _, tc := range cases {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/x"+tc.query, nil)
		if got := rankLimit(c); got != tc.want {
			t.Errorf("query=%q -> %d, want %d", tc.query, got, tc.want)
		}
	}
}
