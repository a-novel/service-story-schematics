package dao

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/a-novel-kit/context"
	pgctx "github.com/a-novel-kit/context/pgbun"

	"github.com/a-novel/service-story-schematics/models"
)

var ErrSelectStoryPlanBySlugRepository = errors.New("SelectStoryPlanBySlugRepository.SelectStoryPlanBySlug")

func NewErrSelectStoryPlanBySlugRepository(err error) error {
	return errors.Join(err, ErrSelectStoryPlanBySlugRepository)
}

type SelectStoryPlanBySlugRepository struct{}

func (repository *SelectStoryPlanBySlugRepository) SelectStoryPlanBySlug(
	ctx context.Context, data models.Slug,
) (*StoryPlanEntity, error) {
	db, err := pgctx.Context(ctx)
	if err != nil {
		return nil, NewErrSelectStoryPlanBySlugRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entity := &StoryPlanEntity{}

	err = db.NewSelect().Model(entity).Where("slug = ?", data).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewErrSelectStoryPlanBySlugRepository(ErrStoryPlanNotFound)
		}

		return nil, NewErrSelectStoryPlanBySlugRepository(fmt.Errorf("select story plan: %w", err))
	}

	return entity, nil
}

func NewSelectStoryPlanBySlugRepository() *SelectStoryPlanBySlugRepository {
	return &SelectStoryPlanBySlugRepository{}
}
