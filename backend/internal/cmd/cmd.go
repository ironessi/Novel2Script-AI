package cmd

import (
	"context"
	"strconv"

	"novel2script-backend/internal/config"
	"novel2script-backend/internal/controller/audit"
	"novel2script-backend/internal/controller/auth"
	"novel2script-backend/internal/controller/chapter"
	"novel2script-backend/internal/controller/project"
	"novel2script-backend/internal/controller/script"
	"novel2script-backend/internal/controller/task"
	uploadCtrl "novel2script-backend/internal/controller/upload"
	"novel2script-backend/internal/dao"
	"novel2script-backend/internal/middleware"
	myredis "novel2script-backend/internal/redis"
	"novel2script-backend/internal/runner"

	// 注册 MySQL 驱动
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	// 注册 logic 层（init 自动注册 service 实现）
	_ "novel2script-backend/internal/logic/audit"
	_ "novel2script-backend/internal/logic/auth"
	_ "novel2script-backend/internal/logic/chapter"
	_ "novel2script-backend/internal/logic/project"
	_ "novel2script-backend/internal/logic/script"
	_ "novel2script-backend/internal/logic/task"
	_ "novel2script-backend/internal/logic/upload"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var Main = gcmd.Command{
	Name:  "main",
	Usage: "main",
	Brief: "Novel2Script-AI Backend Server",
	Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
		// 初始化配置
		config.Init()

		// 初始化数据库
		if err := dao.Init(); err != nil {
			g.Log().Fatalf(ctx, "数据库初始化失败: %v", err)
		}
		g.Log().Info(ctx, "数据库连接成功")

		// 初始化 Redis
		if err := myredis.Init(); err != nil {
			g.Log().Fatalf(ctx, "Redis 初始化失败: %v", err)
		}
		g.Log().Info(ctx, "Redis 连接成功")

		// 初始化任务执行器
		runner.Init()
		g.Log().Info(ctx, "任务执行器初始化完成")

		// 启动 HTTP 服务器
		s := g.Server()

		// 全局中间件
		s.Use(middleware.RequestID)
		s.Use(middleware.SecureHeader)
		s.Use(middleware.CORS)

		// API 路由
		s.Group("/api", func(group *ghttp.RouterGroup) {
			// 认证相关（无需登录）
			group.Group("/auth", func(authGroup *ghttp.RouterGroup) {
				authGroup.Middleware(middleware.RateLimitLogin)
				authGroup.POST("/register", auth.Controller.Register)
				authGroup.POST("/login", auth.Controller.Login)
			})

			// 需要登录的接口
			group.Middleware(middleware.Auth)

			// 用户信息
			group.GET("/auth/me", auth.Controller.Me)

			// 项目管理
			group.Group("/projects", func(projGroup *ghttp.RouterGroup) {
				projGroup.POST("/", project.Controller.Create)
				projGroup.GET("/", project.Controller.List)
				projGroup.GET("/:id", project.Controller.Detail)
				projGroup.PUT("/:id", project.Controller.Update)
				projGroup.DELETE("/:id", project.Controller.Delete)

				// 上传
				projGroup.POST("/:id/upload", uploadCtrl.Controller.Upload)

				// 章节
				projGroup.GET("/:id/chapters", chapter.Controller.List)
				projGroup.GET("/:id/chapters/:chapterId", chapter.Controller.Detail)

				// 任务
				projGroup.POST("/:id/generate", task.Controller.Create)

				// 剧本
				projGroup.GET("/:id/script", script.Controller.Get)
				projGroup.PUT("/:id/script", script.Controller.Update)
				projGroup.POST("/:id/validate", script.Controller.Validate)
				projGroup.GET("/:id/export", script.Controller.Export)

				// 审计日志
				projGroup.GET("/:id/audit", audit.Controller.List)
			})

			// 任务状态查询
			group.GET("/tasks/:id/status", task.Controller.Status)
		})

		// 健康检查
		s.BindHandler("GET:/health", func(r *ghttp.Request) {
			r.Response.WriteJsonExit(g.Map{"status": "ok"})
		})

		port, _ := strconv.Atoi(config.C.App.Port)
		s.SetPort(port)
		s.Run()
		return nil
	},
}

