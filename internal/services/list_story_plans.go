package services

import (
	"context"
	"errors"

	"github.com/getsentry/sentry-go"
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
	span := sentry.StartSpan(ctx, "ListStoryPlansService.ListStoryPlans")
	defer span.Finish()

	span.SetData("request.limit", request.Limit)
	span.SetData("request.offset", request.Offset)

	resp, err := service.source.ListStoryPlans(span.Context(), dao.ListStoryPlansData{
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		span.SetData("dao.listStoryPlans.err", err.Error())

		return nil, NewErrListStoryPlansService(err)
	}

	span.SetData("dao.listStoryPlans.count", len(resp))

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
