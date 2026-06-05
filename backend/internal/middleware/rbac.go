package middleware

import (
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// RBAC 基于角色的权限控制中间件
// 需要在 Auth 中间件之后使用
func RBAC(allowedRoles ...string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		role := GetRole(r)
		if role == "" {
			r.Response.WriteJsonExit(g.Map{
				"code":    http.StatusForbidden,
				"message": "权限不足",
			})
			return
		}

		// admin 拥有所有权限
		if role == "admin" {
			r.Middleware.Next()
			return
		}

		// 检查角色是否在允许列表中
		for _, allowed := range allowedRoles {
			if role == allowed {
				r.Middleware.Next()
				return
			}
		}

		r.Response.WriteJsonExit(g.Map{
			"code":    http.StatusForbidden,
			"message": "权限不足",
		})
	}
}
