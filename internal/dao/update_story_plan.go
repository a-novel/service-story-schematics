package dao

import (
	"context"
	_ "embed"
	"fmt"

	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"
	"github.com/a-novel/golib/postgres"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed update_story_plan.sql
var updateStoryPlanQuery string

type UpdateStoryPlanData struct {
	Plan models.StoryPlan
}

type UpdateStoryPlanRepository struct {
	existsStoryPlan *ExistsStoryPlanSlugRepository
}

func NewUpdateStoryPlanRepository(existsStoryPlan *ExistsStoryPlanSlugRepository) *UpdateStoryPlanRepository {
	return &UpdateStoryPlanRepository{
		existsStoryPlan: existsStoryPlan,
	}
}

func (repository *UpdateStoryPlanRepository) UpdateStoryPlan(
	ctx context.Context, data UpdateStoryPlanData,
) (*StoryPlanEntity, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.UpdateStoryPlan")
	defer span.End()

	span.SetAttributes(
		attribute.String("storyPlan.id", data.Plan.ID.String()),
		attribute.String("storyPlan.slug", data.Plan.Slug.String()),
		attribute.String("storyPlan.name", data.Plan.Name),
		attribute.String("storyPlan.lang", data.Plan.Lang.String()),
	)

	exists, err := repository.existsStoryPlan.ExistsStoryPlanSlug(ctx, data.Plan.Slug)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("check if story plan exists: %w", err))
	}

	if !exists {
		return nil, otel.ReportError(span, ErrStoryPlanNotFound)
	}

	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get postgres client: %w", err))
	}

	entity := &StoryPlanEntity{}

	err = tx.
		NewRaw(
			updateStoryPlanQuery,
			data.Plan.ID,
			data.Plan.Slug,
			data.Plan.Name,
			data.Plan.Description,
			data.Plan.Beats,
			data.Plan.Lang,
			data.Plan.CreatedAt,
		).
		Scan(ctx, entity)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("insert story plan: %w", err))
	}

	return otel.ReportSuccess(span, entity), nil
}
