package service

import (
	"context"
	"io"
	"novel2script-backend/internal/model/entity"
)

// IAuth 认证服务接口
type IAuth interface {
	Register(ctx context.Context, username, email, password string) (*entity.SysUser, error)
	Login(ctx context.Context, username, password string) (*entity.SysUser, string, error)
	GetUserById(ctx context.Context, id int64) (*entity.SysUser, error)
}

// IProject 项目服务接口
type IProject interface {
	Create(ctx context.Context, userId int64, title, description, mode, visibility string) (*entity.NovelProject, error)
	GetList(ctx context.Context, userId int64, page, pageSize int) ([]entity.NovelProject, int, error)
	GetById(ctx context.Context, userId, projectId int64) (*entity.NovelProject, error)
	Update(ctx context.Context, userId int64, projectId int64, title, description, mode, visibility string) error
	Delete(ctx context.Context, userId, projectId int64) error
	CheckPermission(ctx context.Context, userId, projectId int64) (bool, error)
}

// IChapter 章节服务接口
type IChapter interface {
	GetList(ctx context.Context, projectId int64) ([]entity.NovelChapter, error)
	GetById(ctx context.Context, projectId, chapterId int64) (*entity.NovelChapter, error)
}

// ITask 任务服务接口
type ITask interface {
	Create(ctx context.Context, userId, projectId int64, taskType string) (*entity.AiTask, error)
	GetStatus(ctx context.Context, taskId int64) (*entity.AiTask, error)
}

// IScript 剧本服务接口
type IScript interface {
	GetLatest(ctx context.Context, projectId int64) (*entity.ScriptVersion, error)
	GetVersions(ctx context.Context, projectId int64) ([]entity.ScriptVersion, error)
	GetCharacters(ctx context.Context, projectId int64) ([]entity.CharacterProfile, error)
	GetPlotEvents(ctx context.Context, projectId int64) ([]entity.PlotEvent, error)
	GetValidationIssues(ctx context.Context, projectId int64) ([]entity.ValidationIssue, error)
	Update(ctx context.Context, userId, projectId int64, yamlContent string) error
	Validate(ctx context.Context, projectId int64) (bool, []entity.ValidationIssue, error)
	Export(ctx context.Context, projectId int64, format string) (string, error)
	CheckHallucination(ctx context.Context, projectId int64) (*entity.ScriptVersion, error)
	CheckSafety(ctx context.Context, projectId int64) (*entity.ScriptVersion, error)
}

// IUpload 文件上传服务接口
type IUpload interface {
	Upload(ctx context.Context, userId, projectId int64, filename string, fileSize int64, mimeType string, reader io.Reader) (*entity.NovelSourceFile, []entity.NovelChapter, error)
}

// IAudit 审计日志服务接口
type IAudit interface {
	Log(ctx context.Context, userId, projectId int64, action, resourceType string, resourceId int64, ip, userAgent, requestId string) error
	GetList(ctx context.Context, projectId int64, page, pageSize int) ([]entity.AuditLog, int, error)
}

// Service 服务单例
var (
	Auth    IAuth
	Project IProject
	Chapter IChapter
	Task    ITask
	Script  IScript
	Upload  IUpload
	Audit   IAudit
)
