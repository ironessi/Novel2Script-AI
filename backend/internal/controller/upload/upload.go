package upload

import (
	"strconv"

	v1 "novel2script-backend/api/v1"
	"novel2script-backend/internal/middleware"
	"novel2script-backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var Controller = &uploadController{}

type uploadController struct{}

// Upload 上传小说文件
func (c *uploadController) Upload(r *ghttp.Request) {
	projectId, err := strconv.ParseInt(r.GetRouter("id").String(), 10, 64)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": "无效的项目ID"})
		return
	}

	// 权限校验
	userId := middleware.GetUserID(r)
	ok, err := service.Project.CheckPermission(r.Context(), userId, projectId)
	if err != nil || !ok {
		r.Response.WriteJsonExit(g.Map{"code": 403, "message": "无权访问此项目"})
		return
	}

	// 获取上传文件
	file := r.GetUploadFile("file")
	if file == nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": "请选择要上传的文件"})
		return
	}

	// 打开文件流
	f, err := file.Open()
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 500, "message": "读取文件失败"})
		return
	}
	defer f.Close()

	// 调用上传服务
	sourceFile, chapters, err := service.Upload.Upload(
		r.Context(), userId, projectId,
		file.Filename, file.Size, file.Header.Get("Content-Type"), f,
	)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 400, "message": err.Error()})
		return
	}

	// 记录审计日志
	_ = service.Audit.Log(r.Context(), userId, projectId, "file.upload", "novel_source_file", sourceFile.Id,
		r.GetClientIp(), r.GetHeader("User-Agent"), middleware.GetRequestID(r))

	// 构造响应
	chapterItems := make([]v1.ChapterItem, 0, len(chapters))
	for _, ch := range chapters {
		chapterItems = append(chapterItems, v1.ChapterItem{
			Id:           ch.Id,
			ChapterIndex: ch.ChapterIndex,
			ChapterTitle: ch.ChapterTitle,
			ContentHash:  ch.ContentHash,
			CreatedAt:    ch.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": v1.UploadRes{
			FileId:           sourceFile.Id,
			OriginalFilename: sourceFile.OriginalFilename,
			FileSize:         sourceFile.FileSize,
			ChapterCount:     len(chapters),
			Chapters:         chapterItems,
		},
	})
}
