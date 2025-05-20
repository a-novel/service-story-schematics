package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	pgctx "github.com/a-novel-kit/context/pgbun"
)

var ErrSelectBeatsSheetRepository = errors.New("SelectBeatsSheetRepository.SelectBeatsSheet")

func NewErrSelectBeatsSheetRepository(err error) error {
	return errors.Join(err, ErrSelectBeatsSheetRepository)
}

type SelectBeatsSheetRepository struct{}

func (repository *SelectBeatsSheetRepository) SelectBeatsSheet(
	ctx context.Context, data uuid.UUID,
) (*BeatsSheetEntity, error) {
	tx, err := pgctx.Context(ctx)
	if err != nil {
		return nil, NewErrSelectBeatsSheetRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entity := &BeatsSheetEntity{}

	err = tx.NewSelect().Model(entity).Where("id = ?", data).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewErrSelectBeatsSheetRepository(ErrBeatsSheetNotFound)
		}

		return nil, NewErrSelectBeatsSheetRepository(fmt.Errorf("select beats sheet: %w", err))
	}

	return entity, nil
}

func NewSelectBeatsSheetRepository() *SelectBeatsSheetRepository {
	return &SelectBeatsSheetRepository{}
}
