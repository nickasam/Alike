package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func runCORS(mw gin.HandlerFunc, method, origin string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, "/api/x", nil)
	if origin != "" {
		c.Request.Header.Set("Origin", origin)
	}
	chain := gin.HandlersChain{mw, func(c *gin.Context) { c.Status(200) }}
	for _, h := range chain {
		if c.IsAborted() {
			break
		}
		h(c)
	}
	return w
}

func TestCORSWhitelistedOriginEchoed(t *testing.T) {
	mw := CORS(CORSOptions{AllowedOrigins: []string{"https://alike.example"}})
	w := runCORS(mw, http.MethodGet, "https://alike.example")
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://alike.example" {
		t.Fatalf("Allow-Origin=%q, want 回显白名单 Origin", got)
	}
	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Fatal("命中白名单应允许凭证")
	}
}

func TestCORSNonWhitelistedOriginDenied(t *testing.T) {
	mw := CORS(CORSOptions{AllowedOrigins: []string{"https://alike.example"}})
	w := runCORS(mw, http.MethodGet, "https://evil.example")
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("非白名单 Origin 不应下发 Allow-Origin，got %q", got)
	}
}

func TestCORSPreflightNonWhitelistedForbidden(t *testing.T) {
	mw := CORS(CORSOptions{AllowedOrigins: []string{"https://alike.example"}})
	w := runCORS(mw, http.MethodOptions, "https://evil.example")
	if w.Code != http.StatusForbidden {
		t.Fatalf("非白名单预检应 403，got %d", w.Code)
	}
}

func TestCORSPreflightWhitelisted204(t *testing.T) {
	mw := CORS(CORSOptions{AllowedOrigins: []string{"https://alike.example"}})
	w := runCORS(mw, http.MethodOptions, "https://alike.example")
	if w.Code != http.StatusNoContent {
		t.Fatalf("白名单预检应 204，got %d", w.Code)
	}
}

func TestCORSDevModeEchoesWhenEmptyWhitelist(t *testing.T) {
	mw := CORS(CORSOptions{AllowAllInDev: true})
	w := runCORS(mw, http.MethodGet, "http://localhost:3000")
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("开发模式空白名单应回显 Origin，got %q", got)
	}
}

func TestCORSProdEmptyWhitelistDenies(t *testing.T) {
	// 生产语义：AllowAllInDev=false + 空白名单 → 拒绝所有跨域。
	mw := CORS(CORSOptions{AllowAllInDev: false})
	w := runCORS(mw, http.MethodGet, "https://any.example")
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("生产空白名单不应放行任何 Origin，got %q", got)
	}
}
