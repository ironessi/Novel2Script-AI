package task

import (
	"context"

	"novel2script-backend/internal/dao"
	"novel2script-backend/internal/model/entity"
	"novel2script-backend/internal/runner"
	"novel2script-backend/internal/service"
)

type taskImpl struct{}

func init() {
	service.Task = &taskImpl{}
}

func (t *taskImpl) Create(ctx context.Context, userId, projectId int64, taskType string) (*entity.AiTask, error) {
	task := &entity.AiTask{
		ProjectId: projectId,
		OwnerId:   userId,
		TaskType:  taskType,
	}

	id, err := dao.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	task.Id = id
	task.Status = "pending"

	// 更新项目状态
	_ = dao.UpdateProjectStatus(ctx, projectId, "processing")

	// 异步启动 AI 任务
	runner.RunTask(id)

	return task, nil
}

func (t *taskImpl) GetStatus(ctx context.Context, taskId int64) (*entity.AiTask, error) {
	return dao.GetTaskById(ctx, taskId)
}

func (t *taskImpl) GetLatestByProject(ctx context.Context, projectId int64) (*entity.AiTask, error) {
	return dao.GetLatestTaskByProject(ctx, projectId)
}
