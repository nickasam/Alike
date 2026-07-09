package emotion

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func withParamID(id string) func(*gin.Context) {
	return func(c *gin.Context) { c.Params = gin.Params{{Key: "id", Value: id}} }
}

func TestBoardInvalidChannelID(t *testing.T) {
	h := NewHandler(nil)
	code, body := invoke(h.Board, http.MethodGet, "/api/channels/0/emotion-board", withParamID("0"))
	if code != http.StatusNotFound || body.Code != response.CodeNotFound {
		t.Fatalf("status=%d code=%d, want 404/%d", code, body.Code, response.CodeNotFound)
	}
}

func TestAllTagsCount(t *testing.T) {
	if got := len(AllTags()); got != 8 {
		t.Fatalf("AllTags len=%d, want 8", got)
	}
}

func TestIsValid(t *testing.T) {
	cases := map[string]bool{
		"tired":   true,
		"cheer":   true,
		"quit":    true,
		"unknown": false,
		"":        false,
	}
	for tag, want := range cases {
		if got := IsValid(tag); got != want {
			t.Errorf("IsValid(%q)=%v, want %v", tag, got, want)
		}
	}
}
