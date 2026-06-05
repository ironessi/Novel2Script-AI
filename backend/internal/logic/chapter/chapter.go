package chapter

import (
	"context"

	"novel2script-backend/internal/dao"
	"novel2script-backend/internal/model/entity"
	"novel2script-backend/internal/service"
)

type chapterImpl struct{}

func init() {
	service.Chapter = &chapterImpl{}
}

func (c *chapterImpl) GetList(ctx context.Context, projectId int64) ([]entity.NovelChapter, error) {
	return dao.GetChapterList(ctx, projectId)
}

func (c *chapterImpl) GetById(ctx context.Context, projectId, chapterId int64) (*entity.NovelChapter, error) {
	return dao.GetChapterById(ctx, projectId, chapterId)
}
