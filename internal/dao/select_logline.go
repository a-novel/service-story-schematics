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
	span := sentry.StartSpan(ctx, "SelectLoglineRepository.SelectLogline")
	defer span.Finish()

	span.SetData("logline.id", data.ID.String())
	span.SetData("logline.user_id", data.UserID.String())

	tx, err := lib.PostgresContext(span.Context())
	if err != nil {
		span.SetData("postgres.context.error", err.Error())

		return nil, NewErrSelectLoglineRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entity := &LoglineEntity{}

	err = tx.NewSelect().
		Model(entity).
		Where("id = ?", data.ID).
		Where("user_id = ?", data.UserID).
		Scan(span.Context())
	if err != nil {
		span.SetData("scan.error", err.Error())

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
