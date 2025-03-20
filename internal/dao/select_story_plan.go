package dao

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/a-novel-kit/context"
	pgctx "github.com/a-novel-kit/context/pgbun"
)

var ErrSelectStoryPlanRepository = errors.New("SelectStoryPlanRepository.SelectStoryPlan")

func NewErrSelectStoryPlanRepository(err error) error {
	return errors.Join(err, ErrSelectStoryPlanRepository)
}

type SelectStoryPlanRepository struct{}

func (repository *SelectStoryPlanRepository) SelectStoryPlan(
	ctx context.Context, data uuid.UUID,
) (*StoryPlanEntity, error) {
	tx, err := pgctx.Context(ctx)
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
