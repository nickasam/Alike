package pulse

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

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

func withParam(key, val string) func(*gin.Context) {
	return func(c *gin.Context) { c.Params = append(c.Params, gin.Param{Key: key, Value: val}) }
}

// GetItems 传空 slug 时按 404 处理（gin 路由通常拦不到，但 handler 也要防御）。
func TestGetItemsEmptySlug(t *testing.T) {
	h := NewHandler(nil)
	code, body := invoke(h.GetItems, http.MethodGet, "/api/pulse/topics//items", withParam("slug", ""))
	if code != http.StatusNotFound || body.Code != response.CodeNotFound {
		t.Fatalf("status=%d code=%d, want 404/%d", code, body.Code, response.CodeNotFound)
	}
}

// isStale：nil / 超时 / 近期，三种情况的判定。
func TestIsStale(t *testing.T) {
	cases := []struct {
		name string
		last *time.Time
		want bool
	}{
		{"nil", nil, true},
		{"far past", timePtr(time.Now().Add(-7 * time.Hour)), true},
		{"recent", timePtr(time.Now().Add(-2 * time.Minute)), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isStale(tc.last); got != tc.want {
				t.Fatalf("isStale=%v, want %v", got, tc.want)
			}
		})
	}
}

func timePtr(t time.Time) *time.Time { return &t }
