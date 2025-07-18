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

//go:embed exists_story_plan_slug.sql
var existsStoryPlanSlugQuery string

type ExistsStoryPlanSlugRepository struct{}

func NewExistsStoryPlanSlugRepository() *ExistsStoryPlanSlugRepository {
	return &ExistsStoryPlanSlugRepository{}
}

// ExistsStoryPlanSlug returns whether a story plan with the given slug exists in the database.
func (repository *ExistsStoryPlanSlugRepository) ExistsStoryPlanSlug(
	ctx context.Context, data models.Slug,
) (bool, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.ExistsStoryPlanSlug")
	defer span.End()

	span.SetAttributes(attribute.String("slug", data.String()))

	// Retrieve a connection to postgres from the context.
	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return false, otel.ReportError(span, fmt.Errorf("get postgres client: %w", err))
	}

	res, err := tx.NewRaw(existsStoryPlanSlugQuery, data.String()).Exec(ctx)
	if err != nil {
		return false, otel.ReportError(span, fmt.Errorf("check database: %w", err))
	}

	n, err := res.RowsAffected()
	if err != nil {
		return false, otel.ReportError(span, fmt.Errorf("get rows affected: %w", err))
	}

	return otel.ReportSuccess(span, n > 0), nil
}
