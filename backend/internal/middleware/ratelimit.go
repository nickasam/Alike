package middleware

import (
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/response"
)

// ratelimit.go —— 轻量级内存令牌桶限流中间件（无第三方依赖）。
//
// 说明：单实例内存限流，多实例部署下各实例独立计数（对暴力破解/刷接口已足够；
// 如需全局精确限流可后续换 Redis 令牌桶）。按 key（IP 或 user_id）分桶。

// tokenBucket 是单个 key 的令牌桶。
type tokenBucket struct {
	tokens   float64
	lastSeen time.Time
}

// Limiter 是一个按 key 限流的令牌桶集合，并发安全。
type Limiter struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	rate     float64 // 每秒补充令牌数
	burst    float64 // 桶容量（突发上限）
	lastGC   time.Time
}

// NewLimiter 创建限流器：每秒补充 rate 个令牌，桶容量 burst。
func NewLimiter(ratePerSec, burst float64) *Limiter {
	return &Limiter{
		buckets: make(map[string]*tokenBucket),
		rate:    ratePerSec,
		burst:   burst,
		lastGC:  time.Now(),
	}
}

// allow 消费 key 的一个令牌，成功返回 true。
func (l *Limiter) allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.gc(now)

	b, ok := l.buckets[key]
	if !ok {
		// 新 key 满桶，消费 1 个。
		l.buckets[key] = &tokenBucket{tokens: l.burst - 1, lastSeen: now}
		return true
	}

	// 依据流逝时间补充令牌，不超过 burst。
	elapsed := now.Sub(b.lastSeen).Seconds()
	b.tokens = min(l.burst, b.tokens+elapsed*l.rate)
	b.lastSeen = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// gc 定期清理长时间未活跃的桶，避免内存无界增长。调用方须持有锁。
func (l *Limiter) gc(now time.Time) {
	if now.Sub(l.lastGC) < time.Minute {
		return
	}
	l.lastGC = now
	for k, b := range l.buckets {
		// 空闲超过 10 分钟的桶已必然回满，删除不影响限流语义。
		if now.Sub(b.lastSeen) > 10*time.Minute {
			delete(l.buckets, k)
		}
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// RateLimitByIP 返回按客户端 IP 限流的中间件，触发时返回 429。
func RateLimitByIP(l *Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !l.allow("ip:"+c.ClientIP(), time.Now()) {
			response.Fail(c, response.CodeTooManyRequests)
			c.Abort()
			return
		}
		c.Next()
	}
}

// RateLimitByUser 返回按当前登录用户限流的中间件（须在 Auth 之后使用）。
// 未识别到用户时回退到按 IP 限流，避免绕过。
func RateLimitByUser(l *Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "ip:" + c.ClientIP()
		if uid, ok := CurrentUserID(c); ok {
			key = "user:" + strconv.FormatInt(uid, 10)
		}
		if !l.allow(key, time.Now()) {
			response.Fail(c, response.CodeTooManyRequests)
			c.Abort()
			return
		}
		c.Next()
	}
}
