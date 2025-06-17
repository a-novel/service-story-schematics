package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-novel/service-story-schematics/internal/lib"
)

var ErrListStoryPlansRepository = errors.New("ListStoryPlansRepository.ListStoryPlans")

func NewErrListStoryPlansRepository(err error) error {
	return errors.Join(err, ErrListStoryPlansRepository)
}

type ListStoryPlansRepository struct{}

type ListStoryPlansData struct {
	Limit  int
	Offset int
}

func (repository *ListStoryPlansRepository) ListStoryPlans(
	ctx context.Context, data ListStoryPlansData,
) ([]*StoryPlanPreviewEntity, error) {
	tx, err := lib.PostgresContext(ctx)
	if err != nil {
		return nil, NewErrListStoryPlansRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entities := make([]*StoryPlanPreviewEntity, 0)

	err = tx.NewSelect().
		Model(&entities).
		Order("created_at DESC").
		Limit(data.Limit).
		Offset(data.Offset).
		Scan(ctx)
	if err != nil {
		return nil, NewErrListStoryPlansRepository(fmt.Errorf("list story plans: %w", err))
	}

	return entities, nil
}

func NewListStoryPlansRepository() *ListStoryPlansRepository {
	return &ListStoryPlansRepository{}
}
