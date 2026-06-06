package auth

import (
	v1 "novel2script-backend/api/v1"
	"novel2script-backend/internal/middleware"
	"novel2script-backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var Controller = &authController{}

type authController struct{}

// Register 用户注册
func (c *authController) Register(r *ghttp.Request) {
	var req v1.RegisterReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	user, err := service.Auth.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 记录审计日志
	_ = service.Audit.Log(r.Context(), user.Id, 0, "user.register", "sys_user", user.Id,
		r.GetClientIp(), r.GetHeader("User-Agent"), middleware.GetRequestID(r))

	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": v1.RegisterRes{
			Id:       user.Id,
			Username: user.Username,
		},
	})
}

// Login 用户登录
func (c *authController) Login(r *ghttp.Request) {
	var req v1.LoginReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	user, token, err := service.Auth.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	// 记录审计日志
	_ = service.Audit.Log(r.Context(), user.Id, 0, "user.login", "sys_user", user.Id,
		r.GetClientIp(), r.GetHeader("User-Agent"), middleware.GetRequestID(r))

	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": v1.LoginRes{
			Token: token,
			User: v1.UserInfoRes{
				Id:       user.Id,
				Username: user.Username,
				Email:    user.Email,
				Role:     user.Role,
			},
		},
	})
}

// Me 获取当前用户信息
func (c *authController) Me(r *ghttp.Request) {
	userId := middleware.GetUserID(r)
	if userId == 0 {
		r.Response.WriteJsonExit(g.Map{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	user, err := service.Auth.GetUserById(r.Context(), userId)
	if err != nil || user == nil {
		r.Response.WriteJsonExit(g.Map{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": v1.UserInfoRes{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	})
}
