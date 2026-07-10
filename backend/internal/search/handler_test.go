package search

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

func invoke(h gin.HandlerFunc, target string) (int, response.Body) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, target, strings.NewReader(""))
	h(c)

	var parsed response.Body
	_ = json.Unmarshal(w.Body.Bytes(), &parsed)
	return w.Code, parsed
}

func TestSearchEmptyQuery(t *testing.T) {
	h := NewHandler(nil)
	code, body := invoke(h.Search, "/api/search?q=")
	if body.Code != response.CodeValidationError {
		t.Fatalf("code=%d, want %d (empty query must be rejected)", body.Code, response.CodeValidationError)
	}
	if code != http.StatusUnprocessableEntity {
		t.Fatalf("status=%d, want 422", code)
	}
}

func TestSearchWhitespaceQuery(t *testing.T) {
	h := NewHandler(nil)
	_, body := invoke(h.Search, "/api/search?q=%20%20%20")
	if body.Code != response.CodeValidationError {
		t.Fatalf("code=%d, want %d (whitespace-only query must be rejected)", body.Code, response.CodeValidationError)
	}
}

func TestLikePatternEscapesWildcards(t *testing.T) {
	cases := map[string]string{
		"abc":   `%abc%`,
		"50%":   `%50\%%`,
		"a_b":   `%a\_b%`,
		`a\b`:   `%a\\b%`,
		"%_\\":  `%\%\_\\%`,
	}
	for in, want := range cases {
		if got := likePattern(in); got != want {
			t.Errorf("likePattern(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParseTypeNormalization(t *testing.T) {
	cases := map[string]SearchType{
		"":         TypeMessage,
		"message":  TypeMessage,
		"MESSAGE":  TypeMessage,
		" diary ":  TypeDiary,
		"channel":  TypeChannel,
		"user":     TypeUser,
		"bogus":    TypeMessage,
	}
	for in, want := range cases {
		if got := parseType(in); got != want {
			t.Errorf("parseType(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPaginateDefaultsAndClamps(t *testing.T) {
	cases := []struct {
		query            string
		wantPage, wantPS int
	}{
		{"/s?q=x", defaultPage, defaultPageSize},
		{"/s?q=x&page=0&page_size=0", defaultPage, defaultPageSize},
		{"/s?q=x&page=-3&page_size=-1", defaultPage, defaultPageSize},
		{"/s?q=x&page=2&page_size=10", 2, 10},
		{"/s?q=x&page_size=999", defaultPage, maxPageSize},
	}
	for _, tc := range cases {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, tc.query, nil)
		page, ps := paginate(c)
		if page != tc.wantPage || ps != tc.wantPS {
			t.Errorf("paginate(%q) = (%d,%d), want (%d,%d)", tc.query, page, ps, tc.wantPage, tc.wantPS)
		}
	}
}

func TestParseChannelID(t *testing.T) {
	cases := map[string]int64{
		"/s?q=x":               0,
		"/s?q=x&channel_id=5":  5,
		"/s?q=x&channel_id=-3": 0,
		"/s?q=x&channel_id=ab": 0,
	}
	for target, want := range cases {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, target, nil)
		if got := parseChannelID(c); got != want {
			t.Errorf("parseChannelID(%q) = %d, want %d", target, got, want)
		}
	}
}

func TestNonNilReturnsEmptySlice(t *testing.T) {
	var nilSlice []int
	if got := nonNil(nilSlice); got == nil || len(got) != 0 {
		t.Fatalf("nonNil(nil) = %v, want empty non-nil slice", got)
	}
	full := []int{1, 2}
	if got := nonNil(full); len(got) != 2 {
		t.Fatalf("nonNil(%v) mutated slice: %v", full, got)
	}
}
