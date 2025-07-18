package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

type UpdateStoryPlanSource interface {
	UpdateStoryPlan(ctx context.Context, data dao.UpdateStoryPlanData) (*dao.StoryPlanEntity, error)
}

type UpdateStoryPlanRequest struct {
	Slug        models.Slug
	Name        string
	Description string
	Beats       []models.BeatDefinition
	Lang        models.Lang
}

type UpdateStoryPlanService struct {
	source UpdateStoryPlanSource
}

func NewUpdateStoryPlanService(source UpdateStoryPlanSource) *UpdateStoryPlanService {
	return &UpdateStoryPlanService{source: source}
}

func (service *UpdateStoryPlanService) UpdateStoryPlan(
	ctx context.Context, request UpdateStoryPlanRequest,
) (*models.StoryPlan, error) {
	ctx, span := otel.Tracer().Start(ctx, "service.UpdateStoryPlan")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.slug", request.Slug.String()),
		attribute.String("request.name", request.Name),
		attribute.String("request.lang", request.Lang.String()),
	)

	resp, err := service.source.UpdateStoryPlan(ctx, dao.UpdateStoryPlanData{
		Plan: models.StoryPlan{
			ID:          uuid.New(),
			Slug:        request.Slug,
			Name:        request.Name,
			Description: request.Description,
			Beats:       request.Beats,
			Lang:        request.Lang,
			CreatedAt:   time.Now(),
		},
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("update story plan: %w", err))
	}

	return otel.ReportSuccess(span, &models.StoryPlan{
		ID:          resp.ID,
		Slug:        resp.Slug,
		Name:        resp.Name,
		Description: resp.Description,
		Beats:       resp.Beats,
		Lang:        resp.Lang,
		CreatedAt:   resp.CreatedAt,
	}), nil
}
