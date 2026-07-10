package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestLimiterAllowsBurstThenBlocks(t *testing.T) {
	l := NewLimiter(1, 3) // 突发 3
	now := time.Now()
	// 前 3 次应放行（满桶）。
	for i := 0; i < 3; i++ {
		if !l.allow("k", now) {
			t.Fatalf("请求 %d 应放行", i+1)
		}
	}
	// 第 4 次应被拒（桶空）。
	if l.allow("k", now) {
		t.Fatal("超出突发后应被限流")
	}
}

func TestLimiterRefillsOverTime(t *testing.T) {
	l := NewLimiter(2, 2) // 每秒补 2
	now := time.Now()
	l.allow("k", now)
	l.allow("k", now)
	if l.allow("k", now) {
		t.Fatal("桶空时应被拒")
	}
	// 过 1 秒补 2 个令牌。
	later := now.Add(time.Second)
	if !l.allow("k", later) {
		t.Fatal("补充后应放行")
	}
}

func TestLimiterKeysAreIndependent(t *testing.T) {
	l := NewLimiter(1, 1)
	now := time.Now()
	if !l.allow("a", now) || !l.allow("b", now) {
		t.Fatal("不同 key 应各自独立计数")
	}
	if l.allow("a", now) {
		t.Fatal("key a 桶已空应被拒")
	}
}

func TestRateLimitByIPReturns429(t *testing.T) {
	gin.SetMode(gin.TestMode)
	l := NewLimiter(1, 1)
	mw := RateLimitByIP(l)

	run := func() int {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		c.Request.RemoteAddr = "203.0.113.9:1234"
		chain := gin.HandlersChain{mw, func(c *gin.Context) { c.Status(200) }}
		for _, h := range chain {
			if c.IsAborted() {
				break
			}
			h(c)
		}
		return w.Code
	}

	if code := run(); code != 200 {
		t.Fatalf("首次请求应放行, got %d", code)
	}
	if code := run(); code != http.StatusTooManyRequests {
		t.Fatalf("第二次应限流 429, got %d", code)
	}
}
