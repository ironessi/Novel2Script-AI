package script

import (
	"context"
	"errors"

	"novel2script-backend/internal/dao"
	"novel2script-backend/internal/model/entity"
	"novel2script-backend/internal/service"
)

type scriptImpl struct{}

func init() {
	service.Script = &scriptImpl{}
}

func (s *scriptImpl) GetLatest(ctx context.Context, projectId int64) (*entity.ScriptVersion, error) {
	return dao.GetLatestScript(ctx, projectId)
}

func (s *scriptImpl) Update(ctx context.Context, userId, projectId int64, yamlContent string) error {
	// 获取当前最新版本
	latest, err := dao.GetLatestScript(ctx, projectId)
	if err != nil {
		return err
	}

	versionNo := 1
	if latest != nil {
		versionNo = latest.VersionNo + 1
	}

	// 创建新版本
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

	// 调用 AI 服务校验（后续实现）
	issues, err := dao.GetValidationIssues(ctx, script.Id)
	if err != nil {
		return false, nil, err
	}

	valid := len(issues) == 0
	return valid, issues, nil
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

	// Markdown 导出（后续实现）
	return script.YamlContent, nil
}
