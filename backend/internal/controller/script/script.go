package script

import (
	"strconv"

	v1 "novel2script-backend/api/v1"
	"novel2script-backend/internal/middleware"
	"novel2script-backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var Controller = &scriptController{}

type scriptController struct{}

// Get 获取剧本
func (c *scriptController) Get(r *ghttp.Request) {
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

	script, err := service.Script.GetLatest(r.Context(), projectId)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 500, "message": "获取剧本失败"})
		return
	}
	if script == nil {
		r.Response.WriteJsonExit(g.Map{"code": 404, "message": "暂无剧本"})
		return
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": v1.ScriptGetRes{
			Id:                script.Id,
			VersionNo:         script.VersionNo,
			YamlContent:       script.YamlContent,
			ValidationStatus:  script.ValidationStatus,
			HallucinationRisk: script.HallucinationRisk,
			SafetyRisk:        script.SafetyRisk,
			CreatedBy:         script.CreatedBy,
			CreatedAt:         script.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

// Update 修改剧本
func (c *scriptController) Update(r *ghttp.Request) {
	projectId, err := strconv.ParseInt(r.GetRouter("id").String(), 10, 64)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": "无效的项目ID"})
		return
	}

	var req v1.ScriptUpdateReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": err.Error()})
		return
	}

	userId := middleware.GetUserID(r)
	ok, err := service.Project.CheckPermission(r.Context(), userId, projectId)
	if err != nil || !ok {
		r.Response.WriteJsonExit(g.Map{"code": 403, "message": "无权访问此项目"})
		return
	}

	if err := service.Script.Update(r.Context(), userId, projectId, req.YamlContent); err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 500, "message": "保存剧本失败"})
		return
	}

	_ = service.Audit.Log(r.Context(), userId, projectId, "script.edit", "script_version", 0,
		r.GetClientIp(), r.GetHeader("User-Agent"), middleware.GetRequestID(r))

	r.Response.WriteJsonExit(g.Map{"code": 0, "message": "保存成功"})
}

// Validate 校验剧本
func (c *scriptController) Validate(r *ghttp.Request) {
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

	valid, issues, err := service.Script.Validate(r.Context(), projectId)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 500, "message": err.Error()})
		return
	}

	vIssues := make([]v1.ValidationIssue, 0, len(issues))
	for _, issue := range issues {
		vIssues = append(vIssues, v1.ValidationIssue{
			IssueType:    issue.IssueType,
			Severity:     issue.Severity,
			Message:      issue.Message,
			LocationPath: issue.LocationPath,
			Suggestion:   issue.Suggestion,
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": v1.ScriptValidateRes{Valid: valid, Issues: vIssues},
	})
}

// Export 导出剧本
func (c *scriptController) Export(r *ghttp.Request) {
	projectId, err := strconv.ParseInt(r.GetRouter("id").String(), 10, 64)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": "无效的项目ID"})
		return
	}

	format := r.Get("format", "yaml").String()

	userId := middleware.GetUserID(r)
	ok, err := service.Project.CheckPermission(r.Context(), userId, projectId)
	if err != nil || !ok {
		r.Response.WriteJsonExit(g.Map{"code": 403, "message": "无权访问此项目"})
		return
	}

	content, err := service.Script.Export(r.Context(), projectId, format)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 500, "message": err.Error()})
		return
	}

	_ = service.Audit.Log(r.Context(), userId, projectId, "script.export", "script_version", 0,
		r.GetClientIp(), r.GetHeader("User-Agent"), middleware.GetRequestID(r))

	if format == "markdown" {
		r.Response.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	} else {
		r.Response.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	}
	r.Response.Header().Set("Content-Disposition", "attachment; filename=script."+format)
	r.Response.Write(content)
}
