package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/getsentry/sentry-go"

	"github.com/google/uuid"
)

var ErrListLoglinesRepository = errors.New("ListLoglinesRepository.ListLoglines")

func NewErrListLoglinesRepository(err error) error {
	return errors.Join(err, ErrListLoglinesRepository)
}

type ListLoglinesRepository struct{}

type ListLoglinesData struct {
	UserID uuid.UUID
	Limit  int
	Offset int
}

func (repository *ListLoglinesRepository) ListLoglines(
	ctx context.Context, data ListLoglinesData,
) ([]*LoglinePreviewEntity, error) {
	span := sentry.StartSpan(ctx, "ListLoglinesRepository.ListLoglines")
	defer span.Finish()

	span.SetData("user.id", data.UserID.String())
	span.SetData("limit", data.Limit)
	span.SetData("offset", data.Offset)

	tx, err := lib.PostgresContext(span.Context())
	if err != nil {
		span.SetData("postgres.context.error", err.Error())

		return nil, NewErrListLoglinesRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entities := make([]*LoglinePreviewEntity, 0)

	err = tx.NewSelect().
		Model(&entities).
		Where("user_id = ?", data.UserID).
		Order("created_at DESC", "name DESC", "slug DESC").
		Limit(data.Limit).
		Offset(data.Offset).
		Scan(span.Context())
	if err != nil {
		span.SetData("scan.error", err.Error())

		return nil, NewErrListLoglinesRepository(fmt.Errorf("list loglines: %w", err))
	}

	return entities, nil
}

func NewListLoglinesRepository() *ListLoglinesRepository {
	return &ListLoglinesRepository{}
}
