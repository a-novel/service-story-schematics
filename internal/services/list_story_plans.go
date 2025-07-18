package services

import (
	"context"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

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

func NewListStoryPlansService(source ListStoryPlansSource) *ListStoryPlansService {
	return &ListStoryPlansService{source: source}
}

func (service *ListStoryPlansService) ListStoryPlans(
	ctx context.Context, request ListStoryPlansRequest,
) ([]*models.StoryPlanPreview, error) {
	ctx, span := otel.Tracer().Start(ctx, "service.ListStoryPlans")
	defer span.End()

	span.SetAttributes(
		attribute.Int("request.limit", request.Limit),
		attribute.Int("request.offset", request.Offset),
	)

	resp, err := service.source.ListStoryPlans(ctx, dao.ListStoryPlansData{
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	span.SetAttributes(attribute.Int("dao.listStoryPlans.count", len(resp)))

	output := lo.Map(resp, func(item *dao.StoryPlanPreviewEntity, _ int) *models.StoryPlanPreview {
		return &models.StoryPlanPreview{
			ID:          item.ID,
			Slug:        item.Slug,
			Name:        item.Name,
			Description: item.Description,
			Lang:        item.Lang,
			CreatedAt:   item.CreatedAt,
		}
	})

	return otel.ReportSuccess(span, output), nil
}
