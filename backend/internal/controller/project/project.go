package project

import (
	"strconv"

	v1 "novel2script-backend/api/v1"
	"novel2script-backend/internal/middleware"
	"novel2script-backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var Controller = &projectController{}

type projectController struct{}

// Create 创建项目
func (c *projectController) Create(r *ghttp.Request) {
	var req v1.CreateProjectReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": err.Error()})
		return
	}

	userId := middleware.GetUserID(r)
	project, err := service.Project.Create(r.Context(), userId, req.Title, req.Description, req.AdaptationMode, req.Visibility)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": err.Error()})
		return
	}

	_ = service.Audit.Log(r.Context(), userId, project.Id, "project.create", "novel_project", project.Id,
		r.GetClientIp(), r.GetHeader("User-Agent"), middleware.GetRequestID(r))

	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": v1.CreateProjectRes{Id: project.Id, Title: project.Title},
	})
}

// List 获取项目列表
func (c *projectController) List(r *ghttp.Request) {
	page, _ := strconv.Atoi(r.Get("page", "1").String())
	pageSize, _ := strconv.Atoi(r.Get("page_size", "10").String())

	userId := middleware.GetUserID(r)
	projects, total, err := service.Project.GetList(r.Context(), userId, page, pageSize)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 500, "message": "获取项目列表失败"})
		return
	}

	items := make([]v1.ProjectItem, 0, len(projects))
	for _, p := range projects {
		items = append(items, v1.ProjectItem{
			Id:             p.Id,
			Title:          p.Title,
			Description:    p.Description,
			AdaptationMode: p.AdaptationMode,
			Status:         p.Status,
			CreatedAt:      p.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": v1.ProjectListRes{Total: total, Projects: items},
	})
}

// Detail 获取项目详情
func (c *projectController) Detail(r *ghttp.Request) {
	projectId, err := strconv.ParseInt(r.GetRouter("id").String(), 10, 64)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": "无效的项目ID"})
		return
	}

	userId := middleware.GetUserID(r)
	project, err := service.Project.GetById(r.Context(), userId, projectId)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 403, "message": err.Error()})
		return
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": v1.ProjectDetailRes{
			Id:             project.Id,
			Title:          project.Title,
			Description:    project.Description,
			AdaptationMode: project.AdaptationMode,
			Visibility:     project.Visibility,
			Status:         project.Status,
			OwnerId:        project.OwnerId,
			CreatedAt:      project.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      project.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// Update 更新项目
func (c *projectController) Update(r *ghttp.Request) {
	projectId, err := strconv.ParseInt(r.GetRouter("id").String(), 10, 64)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": "无效的项目ID"})
		return
	}

	var req v1.UpdateProjectReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": err.Error()})
		return
	}

	userId := middleware.GetUserID(r)
	if err := service.Project.Update(r.Context(), userId, projectId, req.Title, req.Description, req.AdaptationMode, req.Visibility); err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": err.Error()})
		return
	}

	_ = service.Audit.Log(r.Context(), userId, projectId, "project.update", "novel_project", projectId,
		r.GetClientIp(), r.GetHeader("User-Agent"), middleware.GetRequestID(r))

	r.Response.WriteJsonExit(g.Map{"code": 0, "message": "更新成功"})
}

// Delete 删除项目
func (c *projectController) Delete(r *ghttp.Request) {
	projectId, err := strconv.ParseInt(r.GetRouter("id").String(), 10, 64)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": "无效的项目ID"})
		return
	}

	userId := middleware.GetUserID(r)
	if err := service.Project.Delete(r.Context(), userId, projectId); err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": err.Error()})
		return
	}

	_ = service.Audit.Log(r.Context(), userId, projectId, "project.delete", "novel_project", projectId,
		r.GetClientIp(), r.GetHeader("User-Agent"), middleware.GetRequestID(r))

	r.Response.WriteJsonExit(g.Map{"code": 0, "message": "删除成功"})
}
