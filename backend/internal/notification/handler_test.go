package notification

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

func withParamID(id string) func(*gin.Context) {
	return func(c *gin.Context) { c.Params = gin.Params{{Key: "id", Value: id}} }
}

func TestListRequiresAuth(t *testing.T) {
	h := NewHandler(nil)
	code, body := invoke(h.List, http.MethodGet, "/api/notifications", nil)
	if code != http.StatusUnauthorized || body.Code != response.CodeUnauthorized {
		t.Fatalf("status=%d code=%d, want 401/%d", code, body.Code, response.CodeUnauthorized)
	}
}

func TestReadRequiresAuth(t *testing.T) {
	h := NewHandler(nil)
	code, body := invoke(h.Read, http.MethodPut, "/api/notifications/1/read", withParamID("1"))
	if code != http.StatusUnauthorized || body.Code != response.CodeUnauthorized {
		t.Fatalf("status=%d code=%d, want 401/%d", code, body.Code, response.CodeUnauthorized)
	}
}

func TestReadAllRequiresAuth(t *testing.T) {
	h := NewHandler(nil)
	code, body := invoke(h.ReadAll, http.MethodPut, "/api/notifications/read-all", nil)
	if code != http.StatusUnauthorized || body.Code != response.CodeUnauthorized {
		t.Fatalf("status=%d code=%d, want 401/%d", code, body.Code, response.CodeUnauthorized)
	}
}

func TestReadInvalidID(t *testing.T) {
	h := NewHandler(nil)
	code, body := invoke(h.Read, http.MethodPut, "/api/notifications/x/read",
		func(c *gin.Context) { c.Set(middleware.ContextUserID, int64(1)); withParamID("x")(c) })
	if code != http.StatusNotFound || body.Code != response.CodeNotFound {
		t.Fatalf("status=%d code=%d, want 404/%d", code, body.Code, response.CodeNotFound)
	}
}

func TestPaginateDefaultsAndCaps(t *testing.T) {
	cases := []struct {
		query            string
		wantPage, wantPS int
	}{
		{"", defaultPage, defaultPageSize},
		{"?page=2&page_size=10", 2, 10},
		{"?page=0&page_size=0", defaultPage, defaultPageSize},
		{"?page=-1&page_size=999", defaultPage, maxPageSize},
	}
	for _, tc := range cases {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/x"+tc.query, nil)
		page, ps := paginate(c)
		if page != tc.wantPage || ps != tc.wantPS {
			t.Errorf("query=%q -> page=%d ps=%d, want %d/%d", tc.query, page, ps, tc.wantPage, tc.wantPS)
		}
	}
}
