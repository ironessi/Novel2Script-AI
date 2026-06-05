package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/grand"
)

// RequestID 请求 ID 中间件
func RequestID(r *ghttp.Request) {
	requestId := r.GetHeader("X-Request-Id")
	if requestId == "" {
		requestId = grand.S(16)
	}
	r.SetCtxVar("requestId", requestId)
	r.Response.Header().Set("X-Request-Id", requestId)
	r.Middleware.Next()
}

// GetRequestID 从上下文获取请求 ID
func GetRequestID(r *ghttp.Request) string {
	if v := r.GetCtxVar("requestId"); !v.IsNil() {
		return v.String()
	}
	return ""
}
