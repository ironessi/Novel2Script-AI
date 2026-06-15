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
		"username":      user.Username,
		"email":         user.Email,
		"password_hash": user.PasswordHash,
		"role":          user.Role,
		"status":        user.Status,
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
	err = m.Fields("novel_project.*", "(SELECT COUNT(*) FROM novel_chapter WHERE novel_chapter.project_id = novel_project.id) AS chapter_count").
		Page(page, pageSize).
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

// UpdateProjectStatus 更新项目状态
func UpdateProjectStatus(ctx context.Context, projectId int64, status string) error {
	_, err := db.Model("novel_project").Ctx(ctx).
		Where("id", projectId).
		Where("deleted_at IS NULL").
		Data(g.Map{"status": status}).
		Update()
	return err
}

// ==================== Source File DAO ====================

// CreateSourceFile 创建源文件记录
func CreateSourceFile(ctx context.Context, file *entity.NovelSourceFile) (int64, error) {
	result, err := db.Model("novel_source_file").Ctx(ctx).Insert(g.Map{
		"project_id":        file.ProjectId,
		"owner_id":          file.OwnerId,
		"original_filename": file.OriginalFilename,
		"storage_path":      file.StoragePath,
		"file_hash":         file.FileHash,
		"file_size":         file.FileSize,
		"mime_type":         file.MimeType,
		"scan_status":       "clean",
	})
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// GetSourceFilesByProject 获取项目的所有源文件
func GetSourceFilesByProject(ctx context.Context, projectId int64) ([]entity.NovelSourceFile, error) {
	var files []entity.NovelSourceFile
	err := db.Model("novel_source_file").Ctx(ctx).
		Where("project_id", projectId).
		OrderDesc("created_at").
		Scan(&files)
	return files, err
}

// ==================== Chapter DAO ====================

// CreateChapter 创建章节
func CreateChapter(ctx context.Context, chapter *entity.NovelChapter) (int64, error) {
	result, err := db.Model("novel_chapter").Ctx(ctx).Insert(g.Map{
		"project_id":    chapter.ProjectId,
		"chapter_index": chapter.ChapterIndex,
		"chapter_title": chapter.ChapterTitle,
		"content":       chapter.Content,
		"content_hash":  chapter.ContentHash,
	})
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// BatchCreateChapters 批量创建章节
func BatchCreateChapters(ctx context.Context, chapters []entity.NovelChapter) error {
	if len(chapters) == 0 {
		return nil
	}
	items := make(g.List, 0, len(chapters))
	for _, ch := range chapters {
		items = append(items, g.Map{
			"project_id":    ch.ProjectId,
			"chapter_index": ch.ChapterIndex,
			"chapter_title": ch.ChapterTitle,
			"content":       ch.Content,
			"content_hash":  ch.ContentHash,
		})
	}
	_, err := db.Model("novel_chapter").Ctx(ctx).Insert(items)
	return err
}

// DeleteChaptersByProject 删除项目的所有章节（重新上传时用）
func DeleteChaptersByProject(ctx context.Context, projectId int64) error {
	_, err := db.Model("novel_chapter").Ctx(ctx).
		Where("project_id", projectId).
		Delete()
	return err
}

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

// GetLatestTaskByProject 获取项目最近一次 AI 任务
func GetLatestTaskByProject(ctx context.Context, projectId int64) (*entity.AiTask, error) {
	record, err := db.Model("ai_task").Ctx(ctx).
		Where("project_id", projectId).
		OrderDesc("created_at").
		OrderDesc("id").
		Limit(1).
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

// GetScriptVersions 获取项目的剧本版本历史
func GetScriptVersions(ctx context.Context, projectId int64) ([]entity.ScriptVersion, error) {
	var versions []entity.ScriptVersion
	err := db.Model("script_version").Ctx(ctx).
		Where("project_id", projectId).
		OrderDesc("version_no").
		Scan(&versions)
	return versions, err
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

// CreateValidationIssue 创建校验问题
func CreateValidationIssue(ctx context.Context, issue *entity.ValidationIssue) error {
	_, err := db.Model("validation_issue").Ctx(ctx).Insert(g.Map{
		"project_id":        issue.ProjectId,
		"script_version_id": issue.ScriptVersionId,
		"issue_type":        issue.IssueType,
		"severity":          issue.Severity,
		"message":           issue.Message,
		"location_path":     issue.LocationPath,
		"suggestion":        issue.Suggestion,
	})
	return err
}

// BatchCreateValidationIssues 批量创建校验问题
func BatchCreateValidationIssues(ctx context.Context, issues []entity.ValidationIssue) error {
	if len(issues) == 0 {
		return nil
	}
	items := make(g.List, 0, len(issues))
	for _, issue := range issues {
		items = append(items, g.Map{
			"project_id":        issue.ProjectId,
			"script_version_id": issue.ScriptVersionId,
			"issue_type":        issue.IssueType,
			"severity":          issue.Severity,
			"message":           issue.Message,
			"location_path":     issue.LocationPath,
			"suggestion":        issue.Suggestion,
		})
	}
	_, err := db.Model("validation_issue").Ctx(ctx).Insert(items)
	return err
}

// UpdateScriptValidation 更新剧本校验状态
func UpdateScriptValidation(ctx context.Context, scriptId int64, validationStatus, hallucinationRisk, safetyRisk string) error {
	data := g.Map{}
	if validationStatus != "" {
		data["validation_status"] = validationStatus
	}
	if hallucinationRisk != "" {
		data["hallucination_risk"] = hallucinationRisk
	}
	if safetyRisk != "" {
		data["safety_risk"] = safetyRisk
	}
	if len(data) == 0 {
		return nil
	}
	_, err := db.Model("script_version").Ctx(ctx).
		Where("id", scriptId).
		Data(data).
		Update()
	return err
}

// GetCharactersByProject 获取项目的人物档案
func GetCharactersByProject(ctx context.Context, projectId int64) ([]entity.CharacterProfile, error) {
	var chars []entity.CharacterProfile
	err := db.Model("character_profile").Ctx(ctx).
		Where("project_id", projectId).
		Scan(&chars)
	return chars, err
}

// GetPlotEventsByProject 获取项目的剧情事件链
func GetPlotEventsByProject(ctx context.Context, projectId int64) ([]entity.PlotEvent, error) {
	var events []entity.PlotEvent
	err := db.Model("plot_event").Ctx(ctx).
		Where("project_id", projectId).
		OrderAsc("chapter_index, id").
		Scan(&events)
	return events, err
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
