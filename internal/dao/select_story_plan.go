package dao

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"
	"github.com/a-novel/golib/postgres"
)

//go:embed select_story_plan.sql
var selectStoryPlanQuery string

type SelectStoryPlanRepository struct{}

func NewSelectStoryPlanRepository() *SelectStoryPlanRepository {
	return &SelectStoryPlanRepository{}
}

func (repository *SelectStoryPlanRepository) SelectStoryPlan(
	ctx context.Context, data uuid.UUID,
) (*StoryPlanEntity, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.SelectStoryPlan")
	defer span.End()

	span.SetAttributes(attribute.String("storyPlan.id", data.String()))

	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get postgres client: %w", err))
	}

	entity := &StoryPlanEntity{}

	err = tx.NewRaw(selectStoryPlanQuery, data).Scan(ctx, entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, otel.ReportError(span, ErrStoryPlanNotFound)
		}

		return nil, otel.ReportError(span, fmt.Errorf("select story plan: %w", err))
	}

	return otel.ReportSuccess(span, entity), nil
}
