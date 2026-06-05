package middleware

import "github.com/gogf/gf/v2/net/ghttp"

// CORS 跨域资源共享中间件
func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
