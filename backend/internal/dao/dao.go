package dao

import (
	"context"
	"fmt"
	"novel2script-backend/internal/config"
	"novel2script-backend/internal/model/entity"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// db 全局数据库实例
var db gdb.DB

// Init 初始化数据库连接
func Init() error {
	c := config.C.MySQL
	link := fmt.Sprintf("mysql:%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Database)

	node := gdb.ConfigNode{
		Link: link,
	}
	instance, err := gdb.New(node)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	db = instance

	// 测试连接
	ctx := context.Background()
	if _, err := db.GetOne(ctx, "SELECT 1"); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// DB 获取数据库实例
func DB() gdb.DB {
	return db
}

// ==================== User DAO ====================

// CreateUser 创建用户
func CreateUser(ctx context.Context, user *entity.SysUser) (int64, error) {
	result, err := db.Model("sys_user").Ctx(ctx).Insert(g.Map{
		"username":     user.Username,
		"email":        user.Email,
		"password_hash": user.PasswordHash,
		"role":         user.Role,
		"status":       user.Status,
	})
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// GetUserByUsername 根据用户名查询用户
func GetUserByUsername(ctx context.Context, username string) (*entity.SysUser, error) {
	record, err := db.Model("sys_user").Ctx(ctx).
		Where("username", username).
		Where("status", "active").
		One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, nil
	}
	var user entity.SysUser
	if err := record.Struct(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserById 根据 ID 查询用户
func GetUserById(ctx context.Context, id int64) (*entity.SysUser, error) {
	record, err := db.Model("sys_user").Ctx(ctx).
		Where("id", id).
		Where("status", "active").
		One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, nil
	}
	var user entity.SysUser
	if err := record.Struct(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateLastLogin 更新最后登录时间
func UpdateLastLogin(ctx context.Context, userId int64) error {
	_, err := db.Model("sys_user").Ctx(ctx).
		Where("id", userId).
		Data(g.Map{"last_login_at": time.Now()}).
		Update()
	return err
}

// ==================== Project DAO ====================

// CreateProject 创建项目
func CreateProject(ctx context.Context, project *entity.NovelProject) (int64, error) {
	result, err := db.Model("novel_project").Ctx(ctx).Insert(g.Map{
		"owner_id":        project.OwnerId,
		"title":           project.Title,
		"description":     project.Description,
		"adaptation_mode": project.AdaptationMode,
		"visibility":      project.Visibility,
		"status":          "created",
	})
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// GetProjectList 获取项目列表
func GetProjectList(ctx context.Context, userId int64, page, pageSize int) ([]entity.NovelProject, int, error) {
	m := db.Model("novel_project").Ctx(ctx).
		Where("owner_id", userId).
		Where("deleted_at IS NULL")

	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	var projects []entity.NovelProject
	err = m.Page(page, pageSize).
		OrderDesc("created_at").
		Scan(&projects)
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// GetProjectById 根据 ID 获取项目
func GetProjectById(ctx context.Context, projectId int64) (*entity.NovelProject, error) {
	record, err := db.Model("novel_project").Ctx(ctx).
		Where("id", projectId).
		Where("deleted_at IS NULL").
		One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, nil
	}
	var project entity.NovelProject
	if err := record.Struct(&project); err != nil {
		return nil, err
	}
	return &project, nil
}

// UpdateProject 更新项目
func UpdateProject(ctx context.Context, projectId int64, data g.Map) error {
	_, err := db.Model("novel_project").Ctx(ctx).
		Where("id", projectId).
		Where("deleted_at IS NULL").
		Data(data).
		Update()
	return err
}

// DeleteProject 软删除项目
func DeleteProject(ctx context.Context, projectId int64) error {
	_, err := db.Model("novel_project").Ctx(ctx).
		Where("id", projectId).
		Data(g.Map{"deleted_at": time.Now()}).
		Update()
	return err
}

// ==================== Chapter DAO ====================

// GetChapterList 获取章节列表
func GetChapterList(ctx context.Context, projectId int64) ([]entity.NovelChapter, error) {
	var chapters []entity.NovelChapter
	err := db.Model("novel_chapter").Ctx(ctx).
		Where("project_id", projectId).
		OrderAsc("chapter_index").
		Scan(&chapters)
	return chapters, err
}

// GetChapterById 根据 ID 获取章节
func GetChapterById(ctx context.Context, projectId, chapterId int64) (*entity.NovelChapter, error) {
	record, err := db.Model("novel_chapter").Ctx(ctx).
		Where("id", chapterId).
		Where("project_id", projectId).
		One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, nil
	}
	var chapter entity.NovelChapter
	if err := record.Struct(&chapter); err != nil {
		return nil, err
	}
	return &chapter, nil
}

// ==================== Task DAO ====================

// CreateTask 创建任务
func CreateTask(ctx context.Context, task *entity.AiTask) (int64, error) {
	result, err := db.Model("ai_task").Ctx(ctx).Insert(g.Map{
		"project_id": task.ProjectId,
		"owner_id":   task.OwnerId,
		"task_type":  task.TaskType,
		"status":     "pending",
	})
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// GetTaskById 根据 ID 获取任务
func GetTaskById(ctx context.Context, taskId int64) (*entity.AiTask, error) {
	record, err := db.Model("ai_task").Ctx(ctx).
		Where("id", taskId).
		One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, nil
	}
	var task entity.AiTask
	if err := record.Struct(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

// ==================== Script Version DAO ====================

// GetLatestScript 获取最新剧本版本
func GetLatestScript(ctx context.Context, projectId int64) (*entity.ScriptVersion, error) {
	record, err := db.Model("script_version").Ctx(ctx).
		Where("project_id", projectId).
		OrderDesc("version_no").
		Limit(1).
		One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, nil
	}
	var script entity.ScriptVersion
	if err := record.Struct(&script); err != nil {
		return nil, err
	}
	return &script, nil
}

// CreateScriptVersion 创建剧本版本
func CreateScriptVersion(ctx context.Context, script *entity.ScriptVersion) (int64, error) {
	result, err := db.Model("script_version").Ctx(ctx).Insert(g.Map{
		"project_id":   script.ProjectId,
		"owner_id":     script.OwnerId,
		"version_no":   script.VersionNo,
		"yaml_content": script.YamlContent,
		"yaml_hash":    "",
		"created_by":   script.CreatedBy,
	})
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// GetValidationIssues 获取校验问题列表
func GetValidationIssues(ctx context.Context, scriptVersionId int64) ([]entity.ValidationIssue, error) {
	var issues []entity.ValidationIssue
	err := db.Model("validation_issue").Ctx(ctx).
		Where("script_version_id", scriptVersionId).
		Where("resolved", 0).
		Scan(&issues)
	return issues, err
}

// ==================== Audit Log DAO ====================

// CreateAuditLog 创建审计日志
func CreateAuditLog(ctx context.Context, log *entity.AuditLog) error {
	_, err := db.Model("audit_log").Ctx(ctx).Insert(g.Map{
		"user_id":       log.UserId,
		"project_id":    log.ProjectId,
		"action":        log.Action,
		"resource_type": log.ResourceType,
		"resource_id":   log.ResourceId,
		"ip_address":    log.IpAddress,
		"user_agent":    log.UserAgent,
		"request_id":    log.RequestId,
	})
	return err
}

// GetAuditLogList 获取审计日志列表
func GetAuditLogList(ctx context.Context, projectId int64, page, pageSize int) ([]entity.AuditLog, int, error) {
	m := db.Model("audit_log").Ctx(ctx).
		Where("project_id", projectId)

	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}

	var logs []entity.AuditLog
	err = m.Page(page, pageSize).
		OrderDesc("created_at").
		Scan(&logs)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
