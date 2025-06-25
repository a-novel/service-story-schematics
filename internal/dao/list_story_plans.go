package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/getsentry/sentry-go"
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
	span := sentry.StartSpan(ctx, "ListStoryPlansRepository.ListStoryPlans")
	defer span.Finish()

	span.SetData("limit", data.Limit)
	span.SetData("offset", data.Offset)

	tx, err := lib.PostgresContext(span.Context())
	if err != nil {
		span.SetData("postgres.context.error", err.Error())

		return nil, NewErrListStoryPlansRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entities := make([]*StoryPlanPreviewEntity, 0)

	err = tx.NewSelect().
		Model(&entities).
		Order("created_at DESC").
		Limit(data.Limit).
		Offset(data.Offset).
		Scan(span.Context())
	if err != nil {
		span.SetData("scan.error", err.Error())

		return nil, NewErrListStoryPlansRepository(fmt.Errorf("list story plans: %w", err))
	}

	return entities, nil
}

func NewListStoryPlansRepository() *ListStoryPlansRepository {
	return &ListStoryPlansRepository{}
}
