package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun/driver/pgdriver"

	pgctx "github.com/a-novel-kit/context/pgbun"

	"github.com/a-novel/service-story-schematics/models"
)

var ErrInsertLoglineRepository = errors.New("InsertLoglineRepository.InsertLogline")

func NewErrInsertLoglineRepository(err error) error {
	return errors.Join(err, ErrInsertLoglineRepository)
}

type InsertLoglineData struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Slug   models.Slug

	Name    string
	Content string
	Lang    models.Lang

	Now time.Time
}

type InsertLoglineRepository struct{}

func (repository *InsertLoglineRepository) InsertLogline(
	ctx context.Context, data InsertLoglineData,
) (*LoglineEntity, error) {
	tx, err := pgctx.Context(ctx)
	if err != nil {
		return nil, NewErrInsertLoglineRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entity := &LoglineEntity{
		ID:        data.ID,
		UserID:    data.UserID,
		Slug:      data.Slug,
		Name:      data.Name,
		Content:   data.Content,
		Lang:      data.Lang,
		CreatedAt: data.Now,
	}

	if _, err = tx.NewInsert().Model(entity).Returning("*").Exec(ctx); err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.Field('C') == "23505" {
			return nil, NewErrInsertLoglineRepository(errors.Join(err, ErrLoglineAlreadyExists))
		}

		return nil, NewErrInsertLoglineRepository(fmt.Errorf("insert logline: %w", err))
	}

	return entity, nil
}

func NewInsertLoglineRepository() *InsertLoglineRepository {
	return &InsertLoglineRepository{}
}
