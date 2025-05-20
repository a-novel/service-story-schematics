package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	pgctx "github.com/a-novel-kit/context/pgbun"
)

var ErrListBeatsSheetsRepository = errors.New("ListBeatsSheetsRepository.ListBeatsSheets")

func NewErrListBeatsSheetsRepository(err error) error {
	return errors.Join(err, ErrListBeatsSheetsRepository)
}

type ListBeatsSheetsRepository struct{}

type ListBeatsSheetsData struct {
	LoglineID uuid.UUID
	Limit     int
	Offset    int
}

func (repository *ListBeatsSheetsRepository) ListBeatsSheets(
	ctx context.Context, data ListBeatsSheetsData,
) ([]*BeatsSheetPreviewEntity, error) {
	tx, err := pgctx.Context(ctx)
	if err != nil {
		return nil, NewErrListBeatsSheetsRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entities := make([]*BeatsSheetPreviewEntity, 0)

	err = tx.NewSelect().
		Model(&entities).
		Where("logline_id = ?", data.LoglineID).
		Order("created_at DESC").
		Limit(data.Limit).
		Offset(data.Offset).
		Scan(ctx)
	if err != nil {
		return nil, NewErrListBeatsSheetsRepository(fmt.Errorf("list beats sheet: %w", err))
	}

	return entities, nil
}

func NewListBeatsSheetsRepository() *ListBeatsSheetsRepository {
	return &ListBeatsSheetsRepository{}
}
