package dao

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"
	"github.com/a-novel/golib/postgres"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed select_story_plan_by_slug.sql
var selectStoryPlanBySlugQuery string

type SelectStoryPlanBySlugRepository struct{}

func NewSelectStoryPlanBySlugRepository() *SelectStoryPlanBySlugRepository {
	return &SelectStoryPlanBySlugRepository{}
}

func (repository *SelectStoryPlanBySlugRepository) SelectStoryPlanBySlug(
	ctx context.Context, data models.Slug,
) (*StoryPlanEntity, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.SelectStoryPlanBySlug")
	defer span.End()

	span.SetAttributes(attribute.String("storyPlan.slug", data.String()))

	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get postgres client: %w", err))
	}

	entity := &StoryPlanEntity{}

	err = tx.NewRaw(selectStoryPlanBySlugQuery, data).Scan(ctx, entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, otel.ReportError(span, ErrStoryPlanNotFound)
		}

		return nil, otel.ReportError(span, fmt.Errorf("select story plan: %w", err))
	}

	return otel.ReportSuccess(span, entity), nil
}
