package services

import (
	"context"
	"errors"

	"github.com/samber/lo"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrListStoryPlansService = errors.New("ListStoryPlansService.ListStoryPlans")

func NewErrListStoryPlansService(err error) error {
	return errors.Join(err, ErrListStoryPlansService)
}

type ListStoryPlansSource interface {
	ListStoryPlans(ctx context.Context, data dao.ListStoryPlansData) ([]*dao.StoryPlanPreviewEntity, error)
}

type ListStoryPlansRequest struct {
	Limit  int
	Offset int
}

type ListStoryPlansService struct {
	source ListStoryPlansSource
}

func (service *ListStoryPlansService) ListStoryPlans(
	ctx context.Context, request ListStoryPlansRequest,
) ([]*models.StoryPlanPreview, error) {
	resp, err := service.source.ListStoryPlans(ctx, dao.ListStoryPlansData{
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		return nil, NewErrListStoryPlansService(err)
	}

	return lo.Map(resp, func(item *dao.StoryPlanPreviewEntity, _ int) *models.StoryPlanPreview {
		return &models.StoryPlanPreview{
			ID:          item.ID,
			Slug:        item.Slug,
			Name:        item.Name,
			Description: item.Description,
			Lang:        item.Lang,
			CreatedAt:   item.CreatedAt,
		}
	}), nil
}

func NewListStoryPlansService(source ListStoryPlansSource) *ListStoryPlansService {
	return &ListStoryPlansService{source: source}
}
