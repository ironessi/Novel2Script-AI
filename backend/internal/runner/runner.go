package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"novel2script-backend/internal/client"
	"novel2script-backend/internal/dao"
	myredis "novel2script-backend/internal/redis"

	"github.com/gogf/gf/v2/frame/g"
)

var aiClient *client.AIServiceClient

// Init 初始化任务执行器
func Init() {
	aiClient = client.NewAIServiceClient()
}

// RunTask 异步执行 AI 任务
func RunTask(taskId int64) {
	go func() {
		ctx := context.Background()

		// 获取任务信息
		task, err := dao.GetTaskById(ctx, taskId)
		if err != nil || task == nil {
			g.Log().Errorf(ctx, "Task %d not found: %v", taskId, err)
			return
		}

		// 更新状态为 running
		_ = updateTaskStatus(ctx, taskId, "running", 0, "开始处理")
		_ = myredis.SetTaskProgress(ctx, taskId, 0, "开始处理")

		// 获取项目章节
		chapters, err := dao.GetChapterList(ctx, task.ProjectId)
		if err != nil || len(chapters) == 0 {
			_ = updateTaskFailed(ctx, taskId, "获取章节失败或项目无章节")
			return
		}

		// 构造 AI 请求
		chapterInputs := make([]client.ChapterInput, 0, len(chapters))
		for _, ch := range chapters {
			chapterInputs = append(chapterInputs, client.ChapterInput{
				ChapterIndex: ch.ChapterIndex,
				Title:        ch.ChapterTitle,
				Content:      ch.Content,
			})
		}

		// 获取项目信息以确定改编模式
		project, err := dao.GetProjectById(ctx, task.ProjectId)
		if err != nil || project == nil {
			_ = updateTaskFailed(ctx, taskId, "获取项目信息失败")
			return
		}

		adaptMode := project.AdaptationMode
		if adaptMode == "" {
			adaptMode = "screen_script"
		}

		req := &client.AnalyzeRequest{
			ProjectId:      task.ProjectId,
			Chapters:       chapterInputs,
			AdaptationMode: adaptMode,
		}

		// 更新进度：调用 AI 服务
		_ = updateTaskStatus(ctx, taskId, "running", 10, "正在调用 AI 服务")
		_ = myredis.SetTaskProgress(ctx, taskId, 10, "正在调用 AI 服务")

		// 调用 AI 服务（同步等待）
		maxRetries := 3
		var result *client.AnalyzeResponse
		for attempt := 1; attempt <= maxRetries; attempt++ {
			result, err = aiClient.Analyze(ctx, req)
			if err == nil {
				break
			}
			g.Log().Warningf(ctx, "AI service attempt %d/%d failed: %v", attempt, maxRetries, err)
			if attempt < maxRetries {
				_ = updateTaskStatus(ctx, taskId, "running", 10+attempt*5, fmt.Sprintf("重试中 (%d/%d)", attempt, maxRetries))
				time.Sleep(time.Duration(attempt) * 5 * time.Second)
			}
		}

		if err != nil {
			_ = updateTaskFailed(ctx, taskId, fmt.Sprintf("AI 服务调用失败: %v", err))
			return
		}

		// 保存结果
		_ = updateTaskStatus(ctx, taskId, "running", 50, "正在保存人物档案")
		_ = myredis.SetTaskProgress(ctx, taskId, 50, "正在保存人物档案")

		if err := saveCharacters(ctx, task.ProjectId, result.Characters); err != nil {
			g.Log().Errorf(ctx, "Save characters failed: %v", err)
		}

		_ = updateTaskStatus(ctx, taskId, "running", 60, "正在保存剧情事件")
		_ = myredis.SetTaskProgress(ctx, taskId, 60, "正在保存剧情事件")

		if err := savePlotEvents(ctx, task.ProjectId, result.PlotEvents); err != nil {
			g.Log().Errorf(ctx, "Save plot events failed: %v", err)
		}

		_ = updateTaskStatus(ctx, taskId, "running", 70, "正在保存剧本")
		_ = myredis.SetTaskProgress(ctx, taskId, 70, "正在保存剧本")

		if err := saveScript(ctx, task.ProjectId, task.OwnerId, result); err != nil {
			g.Log().Errorf(ctx, "Save script failed: %v", err)
		}

		// 更新项目状态
		_ = dao.UpdateProjectStatus(ctx, task.ProjectId, "completed")

		// 完成
		_ = updateTaskStatus(ctx, taskId, "completed", 100, "完成")
		_ = myredis.SetTaskProgress(ctx, taskId, 100, "完成")

		// 更新完成时间
		now := time.Now()
		_, _ = dao.DB().Model("ai_task").Ctx(ctx).
			Where("id", taskId).
			Data(g.Map{"finished_at": now}).
			Update()

		g.Log().Infof(ctx, "Task %d completed successfully", taskId)
	}()
}

// updateTaskStatus 更新任务状态
func updateTaskStatus(ctx context.Context, taskId int64, status string, progress int, step string) error {
	_, err := dao.DB().Model("ai_task").Ctx(ctx).
		Where("id", taskId).
		Data(g.Map{
			"status":       status,
			"progress":     progress,
			"current_step": step,
		}).
		Update()
	return err
}

// updateTaskFailed 更新任务为失败状态
func updateTaskFailed(ctx context.Context, taskId int64, errMsg string) error {
	_, err := dao.DB().Model("ai_task").Ctx(ctx).
		Where("id", taskId).
		Data(g.Map{
			"status":        "failed",
			"error_message": errMsg,
		}).
		Update()
	_ = myredis.SetTaskProgress(ctx, taskId, 0, "失败: "+errMsg)
	return err
}

// saveCharacters 保存人物档案
func saveCharacters(ctx context.Context, projectId int64, characters []client.CharacterData) error {
	for _, char := range characters {
		aliases, _ := json.Marshal(char.Aliases)
		personality, _ := json.Marshal(char.Personality)
		relationships, _ := json.Marshal(char.Relationships)
		sourceRefs, _ := json.Marshal(char.SourceTrace)

		_, err := dao.DB().Model("character_profile").Ctx(ctx).Insert(g.Map{
			"project_id":    projectId,
			"character_key": char.Id,
			"name":          char.Name,
			"aliases":       string(aliases),
			"role_type":     char.Role,
			"description":   char.Description,
			"personality":   string(personality),
			"relationships": string(relationships),
			"source_refs":   string(sourceRefs),
			"confidence":    char.Confidence,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// savePlotEvents 保存剧情事件
func savePlotEvents(ctx context.Context, projectId int64, events []client.PlotEventData) error {
	for _, event := range events {
		sourceRefs, _ := json.Marshal(event.SourceTrace)

		_, err := dao.DB().Model("plot_event").Ctx(ctx).Insert(g.Map{
			"project_id":    projectId,
			"event_key":     event.EventKey,
			"chapter_index": event.ChapterIndex,
			"trigger_text":  event.Trigger,
			"action_text":   event.Action,
			"result_text":   event.Result,
			"importance":    event.Importance,
			"source_refs":   string(sourceRefs),
			"confidence":    event.Confidence,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// saveScript 保存剧本版本
func saveScript(ctx context.Context, projectId, userId int64, result *client.AnalyzeResponse) error {
	// 获取当前最大版本号
	var maxVersion int
	record, err := dao.DB().Model("script_version").Ctx(ctx).
		Where("project_id", projectId).
		OrderDesc("version_no").
		Fields("version_no").
		One()
	if err == nil && !record.IsEmpty() {
		maxVersion = record["version_no"].Int()
	}

	// 确定校验状态
	validationStatus := "pending"
	if result.Validation.Valid {
		validationStatus = "valid"
	} else {
		validationStatus = "invalid"
	}

	_, err = dao.DB().Model("script_version").Ctx(ctx).Insert(g.Map{
		"project_id":          projectId,
		"owner_id":            userId,
		"version_no":          maxVersion + 1,
		"yaml_content":        result.YamlContent,
		"yaml_hash":           "",
		"validation_status":   validationStatus,
		"hallucination_risk":  result.Hallucination.RiskLevel,
		"safety_risk":         result.Safety.RiskLevel,
		"created_by":          "ai",
	})
	return err
}
