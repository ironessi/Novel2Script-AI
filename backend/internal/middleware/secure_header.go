package middleware

import "github.com/gogf/gf/v2/net/ghttp"

// SecureHeader 安全响应头中间件
func SecureHeader(r *ghttp.Request) {
	r.Response.Header().Set("X-Content-Type-Options", "nosniff")
	r.Response.Header().Set("X-Frame-Options", "DENY")
	r.Response.Header().Set("X-XSS-Protection", "1; mode=block")
	r.Middleware.Next()
}
