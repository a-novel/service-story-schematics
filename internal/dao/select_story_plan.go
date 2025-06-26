package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"

	"github.com/a-novel/service-story-schematics/internal/lib"
)

var ErrSelectStoryPlanRepository = errors.New("SelectStoryPlanRepository.SelectStoryPlan")

func NewErrSelectStoryPlanRepository(err error) error {
	return errors.Join(err, ErrSelectStoryPlanRepository)
}

type SelectStoryPlanRepository struct{}

func (repository *SelectStoryPlanRepository) SelectStoryPlan(
	ctx context.Context, data uuid.UUID,
) (*StoryPlanEntity, error) {
	span := sentry.StartSpan(ctx, "SelectStoryPlanRepository.SelectStoryPlan")
	defer span.Finish()

	span.SetData("story_plan.id", data.String())

	tx, err := lib.PostgresContext(span.Context())
	if err != nil {
		span.SetData("postgres.context.error", err.Error())

		return nil, NewErrSelectStoryPlanRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entity := &StoryPlanEntity{}

	err = tx.NewSelect().Model(entity).Where("id = ?", data).Scan(span.Context())
	if err != nil {
		span.SetData("scan.error", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewErrSelectStoryPlanRepository(ErrStoryPlanNotFound)
		}

		return nil, NewErrSelectStoryPlanRepository(fmt.Errorf("select story plan: %w", err))
	}

	return entity, nil
}

func NewSelectStoryPlanRepository() *SelectStoryPlanRepository {
	return &SelectStoryPlanRepository{}
}
