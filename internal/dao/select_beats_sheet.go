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

var ErrSelectBeatsSheetRepository = errors.New("SelectBeatsSheetRepository.SelectBeatsSheet")

func NewErrSelectBeatsSheetRepository(err error) error {
	return errors.Join(err, ErrSelectBeatsSheetRepository)
}

type SelectBeatsSheetRepository struct{}

func NewSelectBeatsSheetRepository() *SelectBeatsSheetRepository {
	return &SelectBeatsSheetRepository{}
}

func (repository *SelectBeatsSheetRepository) SelectBeatsSheet(
	ctx context.Context, data uuid.UUID,
) (*BeatsSheetEntity, error) {
	span := sentry.StartSpan(ctx, "SelectBeatsSheetRepository.SelectBeatsSheet")
	defer span.Finish()

	span.SetData("sheet.id", data.String())

	tx, err := lib.PostgresContext(span.Context())
	if err != nil {
		span.SetData("postgres.context.error", err.Error())

		return nil, NewErrSelectBeatsSheetRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entity := &BeatsSheetEntity{}

	err = tx.NewSelect().Model(entity).Where("id = ?", data).Scan(span.Context())
	if err != nil {
		span.SetData("scan.error", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewErrSelectBeatsSheetRepository(ErrBeatsSheetNotFound)
		}

		return nil, NewErrSelectBeatsSheetRepository(fmt.Errorf("select beats sheet: %w", err))
	}

	return entity, nil
}
