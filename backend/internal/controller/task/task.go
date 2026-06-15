package task

import (
	"strconv"

	v1 "novel2script-backend/api/v1"
	"novel2script-backend/internal/middleware"
	"novel2script-backend/internal/model/entity"
	"novel2script-backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var Controller = &taskController{}

type taskController struct{}

// Create 创建任务
func (c *taskController) Create(r *ghttp.Request) {
	projectId, err := strconv.ParseInt(r.GetRouter("id").String(), 10, 64)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": "无效的项目ID"})
		return
	}

	var req v1.CreateTaskReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": err.Error()})
		return
	}
	req.ProjectId = projectId

	userId := middleware.GetUserID(r)
	ok, err := service.Project.CheckPermission(r.Context(), userId, projectId)
	if err != nil || !ok {
		r.Response.WriteJsonExit(g.Map{"code": 403, "message": "无权访问此项目"})
		return
	}

	task, err := service.Task.Create(r.Context(), userId, projectId, req.TaskType)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 500, "message": "创建任务失败"})
		return
	}

	_ = service.Audit.Log(r.Context(), userId, projectId, "task.create", "ai_task", task.Id,
		r.GetClientIp(), r.GetHeader("User-Agent"), middleware.GetRequestID(r))

	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": v1.CreateTaskRes{TaskId: task.Id, Status: task.Status},
	})
}

// Status 查询任务状态
func (c *taskController) Status(r *ghttp.Request) {
	taskId, err := strconv.ParseInt(r.GetRouter("id").String(), 10, 64)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": "无效的任务ID"})
		return
	}

	task, err := service.Task.GetStatus(r.Context(), taskId)
	if err != nil || task == nil {
		r.Response.WriteJsonExit(g.Map{"code": 404, "message": "任务不存在"})
		return
	}

	// 权限校验
	userId := middleware.GetUserID(r)
	if task.OwnerId != userId {
		r.Response.WriteJsonExit(g.Map{"code": 403, "message": "无权查看此任务"})
		return
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": toTaskStatusRes(task),
	})
}

// LatestByProject 查询项目最近一次任务，用于页面刷新后恢复工作流状态
func (c *taskController) LatestByProject(r *ghttp.Request) {
	projectId, err := strconv.ParseInt(r.GetRouter("id").String(), 10, 64)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": "无效的项目ID"})
		return
	}

	userId := middleware.GetUserID(r)
	ok, err := service.Project.CheckPermission(r.Context(), userId, projectId)
	if err != nil || !ok {
		r.Response.WriteJsonExit(g.Map{"code": 403, "message": "无权访问此项目"})
		return
	}

	task, err := service.Task.GetLatestByProject(r.Context(), projectId)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 500, "message": "获取任务记录失败"})
		return
	}
	if task == nil {
		r.Response.WriteJsonExit(g.Map{"code": 0, "data": nil})
		return
	}

	r.Response.WriteJsonExit(g.Map{"code": 0, "data": toTaskStatusRes(task)})
}

func toTaskStatusRes(task *entity.AiTask) v1.TaskStatusRes {
	return v1.TaskStatusRes{
		Id:           task.Id,
		ProjectId:    task.ProjectId,
		TaskType:     task.TaskType,
		Status:       task.Status,
		Progress:     task.Progress,
		CurrentStep:  task.CurrentStep,
		ErrorMessage: task.ErrorMessage,
		CreatedAt:    task.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    task.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
