package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"

	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrSelectLoglineBySlugRepository = errors.New("SelectLoglineBySlugRepository.SelectLoglineBySlug")

func NewErrSelectLoglineBySlugRepository(err error) error {
	return errors.Join(err, ErrSelectLoglineBySlugRepository)
}

type SelectLoglineBySlugData struct {
	Slug   models.Slug
	UserID uuid.UUID
}

type SelectLoglineBySlugRepository struct{}

func (repository *SelectLoglineBySlugRepository) SelectLoglineBySlug(
	ctx context.Context, data SelectLoglineBySlugData,
) (*LoglineEntity, error) {
	span := sentry.StartSpan(ctx, "SelectLoglineBySlugRepository.SelectLoglineBySlug")
	defer span.Finish()

	span.SetData("logline.slug", data.Slug)
	span.SetData("logline.user_id", data.UserID.String())

	tx, err := lib.PostgresContext(span.Context())
	if err != nil {
		span.SetData("postgres.context.error", err.Error())

		return nil, NewErrSelectLoglineBySlugRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entity := &LoglineEntity{}

	err = tx.NewSelect().Model(entity).
		Where("slug = ?", data.Slug).
		Where("user_id = ?", data.UserID).
		Scan(span.Context())
	if err != nil {
		span.SetData("scan.error", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewErrSelectLoglineBySlugRepository(ErrLoglineNotFound)
		}

		return nil, NewErrSelectLoglineBySlugRepository(fmt.Errorf("select logline: %w", err))
	}

	return entity, nil
}

func NewSelectLoglineBySlugRepository() *SelectLoglineBySlugRepository {
	return &SelectLoglineBySlugRepository{}
}
