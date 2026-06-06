package audit

import (
	"strconv"

	v1 "novel2script-backend/api/v1"
	"novel2script-backend/internal/middleware"
	"novel2script-backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var Controller = &auditController{}

type auditController struct{}

// List 获取审计日志列表
func (c *auditController) List(r *ghttp.Request) {
	projectId, err := strconv.ParseInt(r.GetRouter("id").String(), 10, 64)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": "无效的项目ID"})
		return
	}

	page, _ := strconv.Atoi(r.Get("page", "1").String())
	pageSize, _ := strconv.Atoi(r.Get("page_size", "20").String())

	userId := middleware.GetUserID(r)
	ok, err := service.Project.CheckPermission(r.Context(), userId, projectId)
	if err != nil || !ok {
		r.Response.WriteJsonExit(g.Map{"code": 403, "message": "无权访问此项目"})
		return
	}

	logs, total, err := service.Audit.GetList(r.Context(), projectId, page, pageSize)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 500, "message": "获取审计日志失败"})
		return
	}

	items := make([]v1.AuditLogItem, 0, len(logs))
	for _, log := range logs {
		items = append(items, v1.AuditLogItem{
			Id:           log.Id,
			UserId:       log.UserId,
			Action:       log.Action,
			ResourceType: log.ResourceType,
			ResourceId:   log.ResourceId,
			IpAddress:    log.IpAddress,
			RequestId:    log.RequestId,
			CreatedAt:    log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": v1.AuditLogListRes{Total: total, Logs: items},
	})
}
