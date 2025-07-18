package dao

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"regexp"

	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"
	"github.com/a-novel/golib/postgres"

	"github.com/a-novel/service-story-schematics/models"
)

var (
	//go:embed select_slug_iteration.loglines.sql
	selectSlugIterationLoglinesQuery string
	//go:embed select_slug_iteration.story_plans.sql
	selectSlugIterationStoryPlansQuery string
)

var TargetsQueries = map[SlugIterationTarget]string{
	SlugIterationTargetLogline:   selectSlugIterationLoglinesQuery,
	SlugIterationTargetStoryPlan: selectSlugIterationStoryPlansQuery,
}

type SelectSlugIterationData struct {
	Slug models.Slug

	Target SlugIterationTarget
	Args   []any
}

type SelectSlugIterationRepository struct{}

func NewSelectSlugIterationRepository() *SelectSlugIterationRepository {
	return &SelectSlugIterationRepository{}
}

func (repository *SelectSlugIterationRepository) SelectSlugIteration(
	ctx context.Context, data SelectSlugIterationData,
) (models.Slug, int, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.SelectSlugIteration")
	defer span.End()

	span.SetAttributes(
		attribute.String("data.slug", data.Slug.String()),
		attribute.String("data.target", data.Target.String()),
	)

	output := new(struct {
		Slug models.Slug `bun:"slug"`
	})

	reg, err := regexp.CompilePOSIX(`^` + regexp.QuoteMeta(data.Slug.String()) + `-([0-9]+)$`)
	if err != nil {
		return "", 0, otel.ReportError(span, fmt.Errorf("compile regex: %w", err))
	}

	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return "", 0, otel.ReportError(span, fmt.Errorf("get postgres client: %w", err))
	}

	args := append([]any{reg.String()}, data.Args...)

	err = tx.NewRaw(TargetsQueries[data.Target], args...).Scan(ctx, output)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return data.Slug + "-1", 1, nil
		}

		return "", 0, otel.ReportError(span, fmt.Errorf("select slug iteration: %w", err))
	}

	// Capture the index of the last iteration.
	index := 1

	_, err = fmt.Sscanf(output.Slug.String(), data.Slug.String()+"-%d", &index)
	if err != nil {
		return "", 0, otel.ReportError(span, fmt.Errorf("parse slug iteration: %w", err))
	}

	return otel.ReportSuccess(span, models.Slug(fmt.Sprintf("%s-%d", data.Slug, index+1))), index + 1, nil
}
