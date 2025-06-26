package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"

	"github.com/a-novel/service-story-schematics/internal/lib"
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
	span := sentry.StartSpan(ctx, "SelectStoryPlanBySlugRepository.SelectStoryPlanBySlug")
	defer span.Finish()

	span.SetData("story_plan.slug", data)

	db, err := lib.PostgresContext(span.Context())
	if err != nil {
		span.SetData("postgres.context.error", err.Error())

		return nil, NewErrSelectStoryPlanBySlugRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entity := &StoryPlanEntity{}

	err = db.NewSelect().Model(entity).Where("slug = ?", data).Scan(span.Context())
	if err != nil {
		span.SetData("scan.error", err.Error())

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
