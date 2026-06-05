package project

import (
	"context"
	"errors"

	"novel2script-backend/internal/dao"
	"novel2script-backend/internal/model/entity"
	"novel2script-backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type projectImpl struct{}

func init() {
	service.Project = &projectImpl{}
}

func (p *projectImpl) Create(ctx context.Context, userId int64, title, description, mode, visibility string) (*entity.NovelProject, error) {
	if mode == "" {
		mode = "screen_script"
	}
	if visibility == "" {
		visibility = "private"
	}

	project := &entity.NovelProject{
		OwnerId:        userId,
		Title:          title,
		Description:    description,
		AdaptationMode: mode,
		Visibility:     visibility,
	}

	id, err := dao.CreateProject(ctx, project)
	if err != nil {
		return nil, err
	}

	project.Id = id
	return project, nil
}

func (p *projectImpl) GetList(ctx context.Context, userId int64, page, pageSize int) ([]entity.NovelProject, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return dao.GetProjectList(ctx, userId, page, pageSize)
}

func (p *projectImpl) GetById(ctx context.Context, userId, projectId int64) (*entity.NovelProject, error) {
	project, err := dao.GetProjectById(ctx, projectId)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("项目不存在")
	}

	// 权限校验
	if project.OwnerId != userId {
		user, _ := dao.GetUserById(ctx, userId)
		if user == nil || user.Role != "admin" {
			return nil, errors.New("无权访问此项目")
		}
	}

	return project, nil
}

func (p *projectImpl) Update(ctx context.Context, userId int64, projectId int64, title, description, mode, visibility string) error {
	// 先校验权限
	_, err := p.GetById(ctx, userId, projectId)
	if err != nil {
		return err
	}

	data := g.Map{}
	if title != "" {
		data["title"] = title
	}
	if description != "" {
		data["description"] = description
	}
	if mode != "" {
		data["adaptation_mode"] = mode
	}
	if visibility != "" {
		data["visibility"] = visibility
	}

	if len(data) == 0 {
		return errors.New("没有需要更新的字段")
	}

	return dao.UpdateProject(ctx, projectId, data)
}

func (p *projectImpl) Delete(ctx context.Context, userId, projectId int64) error {
	// 先校验权限
	_, err := p.GetById(ctx, userId, projectId)
	if err != nil {
		return err
	}

	return dao.DeleteProject(ctx, projectId)
}

func (p *projectImpl) CheckPermission(ctx context.Context, userId, projectId int64) (bool, error) {
	project, err := dao.GetProjectById(ctx, projectId)
	if err != nil {
		return false, err
	}
	if project == nil {
		return false, nil
	}

	if project.OwnerId == userId {
		return true, nil
	}

	user, _ := dao.GetUserById(ctx, userId)
	if user != nil && user.Role == "admin" {
		return true, nil
	}

	return false, nil
}
