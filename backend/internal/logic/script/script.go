package script

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"novel2script-backend/internal/client"
	"novel2script-backend/internal/dao"
	"novel2script-backend/internal/model/entity"
	"novel2script-backend/internal/service"
)

var aiClient *client.AIServiceClient

type scriptImpl struct{}

func init() {
	service.Script = &scriptImpl{}
}

func getAIClient() *client.AIServiceClient {
	if aiClient == nil {
		aiClient = client.NewAIServiceClient()
	}
	return aiClient
}

func (s *scriptImpl) GetLatest(ctx context.Context, projectId int64) (*entity.ScriptVersion, error) {
	return dao.GetLatestScript(ctx, projectId)
}

func (s *scriptImpl) Update(ctx context.Context, userId, projectId int64, yamlContent string) error {
	latest, err := dao.GetLatestScript(ctx, projectId)
	if err != nil {
		return err
	}

	versionNo := 1
	if latest != nil {
		versionNo = latest.VersionNo + 1
	}

	_, err = dao.CreateScriptVersion(ctx, &entity.ScriptVersion{
		ProjectId:   projectId,
		OwnerId:     userId,
		VersionNo:   versionNo,
		YamlContent: yamlContent,
		CreatedBy:   "user",
	})
	return err
}

func (s *scriptImpl) Validate(ctx context.Context, projectId int64) (bool, []entity.ValidationIssue, error) {
	script, err := dao.GetLatestScript(ctx, projectId)
	if err != nil {
		return false, nil, err
	}
	if script == nil {
		return false, nil, errors.New("没有可校验的剧本")
	}

	// 调用 AI 服务校验
	resp, err := getAIClient().ValidateYAML(ctx, script.YamlContent)
	if err != nil {
		return false, nil, fmt.Errorf("校验失败: %w", err)
	}

	// 转换为 entity.ValidationIssue
	issues := make([]entity.ValidationIssue, 0, len(resp.Issues))
	for _, issue := range resp.Issues {
		issues = append(issues, entity.ValidationIssue{
			ProjectId:       projectId,
			ScriptVersionId: script.Id,
			IssueType:       issue.Type,
			Severity:        issue.Severity,
			Message:         issue.Message,
			LocationPath:    issue.LocationPath,
			Suggestion:      issue.Suggestion,
		})
	}

	// 保存校验问题到数据库
	if len(issues) > 0 {
		_ = dao.BatchCreateValidationIssues(ctx, issues)
	}

	// 更新剧本校验状态
	status := "valid"
	if !resp.Valid {
		status = "invalid"
	}
	_ = dao.UpdateScriptValidation(ctx, script.Id, status, "", "")

	return resp.Valid, issues, nil
}

func (s *scriptImpl) Export(ctx context.Context, projectId int64, format string) (string, error) {
	script, err := dao.GetLatestScript(ctx, projectId)
	if err != nil {
		return "", err
	}
	if script == nil {
		return "", errors.New("没有可导出的剧本")
	}

	if format == "yaml" || format == "" {
		return script.YamlContent, nil
	}

	// Markdown 导出
	return convertToMarkdown(script.YamlContent), nil
}

func (s *scriptImpl) CheckHallucination(ctx context.Context, projectId int64) (*entity.ScriptVersion, error) {
	script, err := dao.GetLatestScript(ctx, projectId)
	if err != nil {
		return nil, err
	}
	if script == nil {
		return nil, errors.New("没有可检测的剧本")
	}

	// 获取章节
	chapters, err := dao.GetChapterList(ctx, projectId)
	if err != nil {
		return nil, err
	}

	// 获取人物档案
	characters, err := dao.GetCharactersByProject(ctx, projectId)
	if err != nil {
		return nil, err
	}

	// 转换为 AI 客户端格式
	chapterInputs := make([]client.ChapterInput, 0, len(chapters))
	for _, ch := range chapters {
		chapterInputs = append(chapterInputs, client.ChapterInput{
			ChapterIndex: ch.ChapterIndex,
			Title:        ch.ChapterTitle,
			Content:      ch.Content,
		})
	}

	charDatas := make([]client.CharacterData, 0, len(characters))
	for _, c := range characters {
		var aliases, personality []string
		_ = json.Unmarshal(c.Aliases, &aliases)
		_ = json.Unmarshal(c.Personality, &personality)
		charDatas = append(charDatas, client.CharacterData{
			Id:          c.CharacterKey,
			Name:        c.Name,
			Aliases:     aliases,
			Role:        c.RoleType,
			Description: c.Description,
			Personality: personality,
		})
	}

	// 调用 AI 服务
	resp, err := getAIClient().CheckHallucination(ctx, script.YamlContent, chapterInputs, charDatas)
	if err != nil {
		return nil, fmt.Errorf("幻觉检测失败: %w", err)
	}

	// 更新剧本风险等级
	_ = dao.UpdateScriptValidation(ctx, script.Id, "", resp.RiskLevel, "")

	script.HallucinationRisk = resp.RiskLevel
	return script, nil
}

func (s *scriptImpl) CheckSafety(ctx context.Context, projectId int64) (*entity.ScriptVersion, error) {
	script, err := dao.GetLatestScript(ctx, projectId)
	if err != nil {
		return nil, err
	}
	if script == nil {
		return nil, errors.New("没有可检测的剧本")
	}

	// 调用 AI 服务
	resp, err := getAIClient().CheckSafety(ctx, script.YamlContent)
	if err != nil {
		return nil, fmt.Errorf("安全审查失败: %w", err)
	}

	// 更新剧本风险等级
	_ = dao.UpdateScriptValidation(ctx, script.Id, "", "", resp.RiskLevel)

	script.SafetyRisk = resp.RiskLevel
	return script, nil
}

// convertToMarkdown 将 YAML 剧本转换为 Markdown 格式
func convertToMarkdown(yamlContent string) string {
	// 简化实现：直接返回 YAML 内容
	// 完整实现应该解析 YAML 并格式化为 Markdown
	return yamlContent
}
