package middleware

import (
	"net/http"
	"strings"

	"novel2script-backend/utility/jwt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Auth JWT 认证中间件
func Auth(r *ghttp.Request) {
	authorization := r.GetHeader("Authorization")
	if authorization == "" {
		r.Response.WriteJsonExit(g.Map{
			"code":    http.StatusUnauthorized,
			"message": "未登录",
		})
		return
	}

	tokenStr := strings.TrimPrefix(authorization, "Bearer ")
	if tokenStr == authorization {
		r.Response.WriteJsonExit(g.Map{
			"code":    http.StatusUnauthorized,
			"message": "Token 格式错误",
		})
		return
	}

	claims, err := jwt.ParseToken(tokenStr)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code":    http.StatusUnauthorized,
			"message": "Token 无效或已过期",
		})
		return
	}

	// 将用户信息存入上下文
	r.SetCtxVar("userId", claims.UserID)
	r.SetCtxVar("username", claims.Username)
	r.SetCtxVar("role", claims.Role)

	r.Middleware.Next()
}

// GetUserID 从上下文获取用户 ID
func GetUserID(r *ghttp.Request) int64 {
	if v := r.GetCtxVar("userId"); !v.IsNil() {
		return v.Int64()
	}
	return 0
}

// GetUsername 从上下文获取用户名
func GetUsername(r *ghttp.Request) string {
	if v := r.GetCtxVar("username"); !v.IsNil() {
		return v.String()
	}
	return ""
}

// GetRole 从上下文获取角色
func GetRole(r *ghttp.Request) string {
	if v := r.GetCtxVar("role"); !v.IsNil() {
		return v.String()
	}
	return ""
}
