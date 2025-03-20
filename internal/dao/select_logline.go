package dao

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/a-novel-kit/context"
	pgctx "github.com/a-novel-kit/context/pgbun"
)

var ErrSelectLoglineRepository = errors.New("SelectLoglineRepository.SelectLogline")

func NewErrSelectLoglineRepository(err error) error {
	return errors.Join(err, ErrSelectLoglineRepository)
}

type SelectLoglineRepository struct{}

type SelectLoglineData struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (repository *SelectLoglineRepository) SelectLogline(
	ctx context.Context, data SelectLoglineData,
) (*LoglineEntity, error) {
	tx, err := pgctx.Context(ctx)
	if err != nil {
		return nil, NewErrSelectLoglineRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entity := &LoglineEntity{}

	err = tx.NewSelect().
		Model(entity).
		Where("id = ?", data.ID).
		Where("user_id = ?", data.UserID).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewErrSelectLoglineRepository(ErrLoglineNotFound)
		}

		return nil, NewErrSelectLoglineRepository(fmt.Errorf("select logline: %w", err))
	}

	return entity, nil
}

func NewSelectLoglineRepository() *SelectLoglineRepository {
	return &SelectLoglineRepository{}
}
