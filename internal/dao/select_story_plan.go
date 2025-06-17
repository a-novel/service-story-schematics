package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/a-novel/service-story-schematics/internal/lib"

	"github.com/google/uuid"
)

var ErrSelectStoryPlanRepository = errors.New("SelectStoryPlanRepository.SelectStoryPlan")

func NewErrSelectStoryPlanRepository(err error) error {
	return errors.Join(err, ErrSelectStoryPlanRepository)
}

type SelectStoryPlanRepository struct{}

func (repository *SelectStoryPlanRepository) SelectStoryPlan(
	ctx context.Context, data uuid.UUID,
) (*StoryPlanEntity, error) {
	tx, err := lib.PostgresContext(ctx)
	if err != nil {
		return nil, NewErrSelectStoryPlanRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entity := &StoryPlanEntity{}

	err = tx.NewSelect().Model(entity).Where("id = ?", data).Scan(ctx)
	if err != nil {
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
