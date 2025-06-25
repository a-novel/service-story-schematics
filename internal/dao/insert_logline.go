package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/getsentry/sentry-go"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun/driver/pgdriver"

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
	span := sentry.StartSpan(ctx, "InsertLoglineRepository.InsertLogline")
	defer span.Finish()

	span.SetData("logline.id", data.ID.String())
	span.SetData("logline.user_id", data.UserID.String())
	span.SetData("logline.slug", data.Slug)
	span.SetData("logline.name", data.Name)
	span.SetData("logline.lang", data.Lang)

	tx, err := lib.PostgresContext(span.Context())
	if err != nil {
		span.SetData("postgres.context.error", err.Error())

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

	if _, err = tx.NewInsert().Model(entity).Returning("*").Exec(span.Context()); err != nil {
		span.SetData("insert.error", err.Error())

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
