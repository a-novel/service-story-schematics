package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	pgctx "github.com/a-novel-kit/context/pgbun"

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
	tx, err := pgctx.Context(ctx)
	if err != nil {
		return nil, NewErrSelectLoglineBySlugRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entity := &LoglineEntity{}

	err = tx.NewSelect().Model(entity).
		Where("slug = ?", data.Slug).
		Where("user_id = ?", data.UserID).
		Scan(ctx)
	if err != nil {
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
