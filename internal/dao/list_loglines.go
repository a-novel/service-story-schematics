package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	pgctx "github.com/a-novel-kit/context/pgbun"
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
	tx, err := pgctx.Context(ctx)
	if err != nil {
		return nil, NewErrListLoglinesRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entities := make([]*LoglinePreviewEntity, 0)

	err = tx.NewSelect().
		Model(&entities).
		Where("user_id = ?", data.UserID).
		Order("created_at DESC", "name DESC", "slug DESC").
		Limit(data.Limit).
		Offset(data.Offset).
		Scan(ctx)
	if err != nil {
		return nil, NewErrListLoglinesRepository(fmt.Errorf("list loglines: %w", err))
	}

	return entities, nil
}

func NewListLoglinesRepository() *ListLoglinesRepository {
	return &ListLoglinesRepository{}
}
