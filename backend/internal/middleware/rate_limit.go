package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// rateLimiter 简单的内存限流器
type rateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

var loginLimiter = &rateLimiter{
	requests: make(map[string][]time.Time),
	limit:    10,
	window:   time.Minute,
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// 清理过期记录
	valid := make([]time.Time, 0)
	for _, t := range rl.requests[key] {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.limit {
		return false
	}

	rl.requests[key] = append(valid, now)
	return true
}

// RateLimitLogin 登录接口限流中间件
func RateLimitLogin(r *ghttp.Request) {
	ip := r.GetClientIp()
	if !loginLimiter.allow(ip) {
		r.Response.WriteJsonExit(g.Map{
			"code":    http.StatusTooManyRequests,
			"message": "请求过于频繁，请稍后再试",
		})
		return
	}
	r.Middleware.Next()
}
