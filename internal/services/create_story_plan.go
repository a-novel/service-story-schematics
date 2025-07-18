package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

type CreateStoryPlanSource interface {
	InsertStoryPlan(ctx context.Context, data dao.InsertStoryPlanData) (*dao.StoryPlanEntity, error)
	SelectSlugIteration(ctx context.Context, data dao.SelectSlugIterationData) (models.Slug, int, error)
}

func NewCreateStoryPlanServiceSource(
	insertStoryPlanDAO *dao.InsertStoryPlanRepository,
	selectSlugIterationDAO *dao.SelectSlugIterationRepository,
) CreateStoryPlanSource {
	return &struct {
		*dao.InsertStoryPlanRepository
		*dao.SelectSlugIterationRepository
	}{
		InsertStoryPlanRepository:     insertStoryPlanDAO,
		SelectSlugIterationRepository: selectSlugIterationDAO,
	}
}

type CreateStoryPlanRequest struct {
	Slug        models.Slug
	Name        string
	Description string
	Beats       []models.BeatDefinition
	Lang        models.Lang
}

type CreateStoryPlanService struct {
	source CreateStoryPlanSource
}

func NewCreateStoryPlanService(source CreateStoryPlanSource) *CreateStoryPlanService {
	return &CreateStoryPlanService{source: source}
}

func (service *CreateStoryPlanService) CreateStoryPlan(
	ctx context.Context, request CreateStoryPlanRequest,
) (*models.StoryPlan, error) {
	ctx, span := otel.Tracer().Start(ctx, "service.CreateStoryPlan")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.slug", request.Slug.String()),
		attribute.String("request.name", request.Name),
		attribute.String("request.lang", request.Lang.String()),
		attribute.Bool("slug.taken", false),
	)

	data := dao.InsertStoryPlanData{
		Plan: models.StoryPlan{
			ID:          uuid.New(),
			Slug:        request.Slug,
			Name:        request.Name,
			Description: request.Description,
			Beats:       request.Beats,
			Lang:        request.Lang,
			CreatedAt:   time.Now(),
		},
	}

	resp, err := service.source.InsertStoryPlan(ctx, data)

	// If slug is taken, try to modify it by appending a version number.
	if errors.Is(err, dao.ErrStoryPlanAlreadyExists) {
		span.SetAttributes(attribute.Bool("slug.taken", true))

		data.Plan.Slug, _, err = service.source.SelectSlugIteration(ctx, dao.SelectSlugIterationData{
			Slug:   data.Plan.Slug,
			Target: dao.SlugIterationTargetStoryPlan,
		})
		if err != nil {
			return nil, otel.ReportError(span, fmt.Errorf("check slug uniqueness: %w", err))
		}

		resp, err = service.source.InsertStoryPlan(ctx, data)
	}

	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("insert story plan: %w", err))
	}

	span.SetAttributes(attribute.String("dao.insertStoryPlan.id", resp.ID.String()))

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
