package chapter

import (
	"strconv"

	v1 "novel2script-backend/api/v1"
	"novel2script-backend/internal/middleware"
	"novel2script-backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var Controller = &chapterController{}

type chapterController struct{}

// List 获取章节列表
func (c *chapterController) List(r *ghttp.Request) {
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

	chapters, err := service.Chapter.GetList(r.Context(), projectId)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 500, "message": "获取章节列表失败"})
		return
	}

	items := make([]v1.ChapterItem, 0, len(chapters))
	for _, ch := range chapters {
		items = append(items, v1.ChapterItem{
			Id:           ch.Id,
			ChapterIndex: ch.ChapterIndex,
			ChapterTitle: ch.ChapterTitle,
			ContentHash:  ch.ContentHash,
			CreatedAt:    ch.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": v1.ChapterListRes{Chapters: items},
	})
}

// Detail 获取章节详情
func (c *chapterController) Detail(r *ghttp.Request) {
	projectId, err := strconv.ParseInt(r.GetRouter("id").String(), 10, 64)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": "无效的项目ID"})
		return
	}

	chapterId, err := strconv.ParseInt(r.GetRouter("chapterId").String(), 10, 64)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": "无效的章节ID"})
		return
	}

	userId := middleware.GetUserID(r)
	ok, err := service.Project.CheckPermission(r.Context(), userId, projectId)
	if err != nil || !ok {
		r.Response.WriteJsonExit(g.Map{"code": 403, "message": "无权访问此项目"})
		return
	}

	chapter, err := service.Chapter.GetById(r.Context(), projectId, chapterId)
	if err != nil || chapter == nil {
		r.Response.WriteJsonExit(g.Map{"code": 404, "message": "章节不存在"})
		return
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": v1.ChapterDetailRes{
			Id:           chapter.Id,
			ChapterIndex: chapter.ChapterIndex,
			ChapterTitle: chapter.ChapterTitle,
			Content:      chapter.Content,
		},
	})
}
