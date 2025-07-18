package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

type SelectStoryPlanSource interface {
	SelectStoryPlan(ctx context.Context, data uuid.UUID) (*dao.StoryPlanEntity, error)
	SelectStoryPlanBySlug(ctx context.Context, data models.Slug) (*dao.StoryPlanEntity, error)
}

func NewSelectStoryPlanServiceSource(
	selectDAO *dao.SelectStoryPlanRepository,
	selectBySlugRepository *dao.SelectStoryPlanBySlugRepository,
) SelectStoryPlanSource {
	return &struct {
		*dao.SelectStoryPlanRepository
		*dao.SelectStoryPlanBySlugRepository
	}{
		SelectStoryPlanRepository:       selectDAO,
		SelectStoryPlanBySlugRepository: selectBySlugRepository,
	}
}

type SelectStoryPlanRequest struct {
	Slug *models.Slug
	ID   *uuid.UUID
}

type SelectStoryPlanService struct {
	source SelectStoryPlanSource
}

func NewSelectStoryPlanService(source SelectStoryPlanSource) *SelectStoryPlanService {
	return &SelectStoryPlanService{source: source}
}

func (service *SelectStoryPlanService) SelectStoryPlan(
	ctx context.Context, request SelectStoryPlanRequest,
) (*models.StoryPlan, error) {
	ctx, span := otel.Tracer().Start(ctx, "service.SelectStoryPlan")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.slug", lo.FromPtr(request.Slug).String()),
		attribute.String("request.id", lo.FromPtr(request.ID).String()),
	)

	if request.Slug != nil {
		data, err := service.source.SelectStoryPlanBySlug(ctx, lo.FromPtr(request.Slug))
		if err != nil {
			return nil, otel.ReportError(span, fmt.Errorf("select story plan by slug: %w", err))
		}

		return otel.ReportSuccess(span, &models.StoryPlan{
			ID:          data.ID,
			Slug:        data.Slug,
			Name:        data.Name,
			Description: data.Description,
			Beats:       data.Beats,
			Lang:        data.Lang,
			CreatedAt:   data.CreatedAt,
		}), nil
	}

	data, err := service.source.SelectStoryPlan(ctx, lo.FromPtr(request.ID))
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("select story plan: %w", err))
	}

	return otel.ReportSuccess(span, &models.StoryPlan{
		ID:          data.ID,
		Slug:        data.Slug,
		Name:        data.Name,
		Description: data.Description,
		Beats:       data.Beats,
		Lang:        data.Lang,
		CreatedAt:   data.CreatedAt,
	}), nil
}
