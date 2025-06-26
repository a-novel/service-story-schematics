package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"

	"github.com/a-novel/service-story-schematics/internal/lib"
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
	span := sentry.StartSpan(ctx, "ListBeatsSheetsRepository.ListBeatsSheets")
	defer span.Finish()

	span.SetData("logline.id", data.LoglineID.String())
	span.SetData("limit", data.Limit)
	span.SetData("offset", data.Offset)

	tx, err := lib.PostgresContext(span.Context())
	if err != nil {
		span.SetData("postgres.context.error", err.Error())

		return nil, NewErrListBeatsSheetsRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entities := make([]*BeatsSheetPreviewEntity, 0)

	err = tx.NewSelect().
		Model(&entities).
		Where("logline_id = ?", data.LoglineID).
		Order("created_at DESC").
		Limit(data.Limit).
		Offset(data.Offset).
		Scan(span.Context())
	if err != nil {
		span.SetData("scan.error", err.Error())

		return nil, NewErrListBeatsSheetsRepository(fmt.Errorf("list beats sheet: %w", err))
	}

	return entities, nil
}

func NewListBeatsSheetsRepository() *ListBeatsSheetsRepository {
	return &ListBeatsSheetsRepository{}
}
