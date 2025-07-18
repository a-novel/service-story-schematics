package dao

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"
	"github.com/a-novel/golib/postgres"
)

//go:embed list_story_plans.sql
var listStoryPlansQuery string

type ListStoryPlansData struct {
	Limit  int
	Offset int
}

type ListStoryPlansRepository struct{}

func NewListStoryPlansRepository() *ListStoryPlansRepository {
	return &ListStoryPlansRepository{}
}

func (repository *ListStoryPlansRepository) ListStoryPlans(
	ctx context.Context, data ListStoryPlansData,
) ([]*StoryPlanPreviewEntity, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.ListStoryPlans")
	defer span.End()

	span.SetAttributes(
		attribute.Int("limit", data.Limit),
		attribute.Int("offset", data.Offset),
	)

	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get postgres client: %w", err))
	}

	entities := make([]*StoryPlanPreviewEntity, 0)

	err = tx.NewRaw(listStoryPlansQuery, bun.NullZero(data.Limit), data.Offset).Scan(ctx, &entities)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("list story plans: %w", err))
	}

	return otel.ReportSuccess(span, entities), nil
}
