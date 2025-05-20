package dao

import (
	"context"
	"errors"
	"fmt"

	pgctx "github.com/a-novel-kit/context/pgbun"

	"github.com/a-novel/service-story-schematics/models"
)

var ErrInsertStoryPlanRepository = errors.New("InsertStoryPlanRepository.InsertStoryPlan")

func NewErrInsertStoryPlanRepository(err error) error {
	return errors.Join(err, ErrInsertStoryPlanRepository)
}

type InsertStoryPlanData struct {
	Plan models.StoryPlan
}

type InsertStoryPlanRepository struct{}

func (repository *InsertStoryPlanRepository) InsertStoryPlan(
	ctx context.Context, data InsertStoryPlanData,
) (*StoryPlanEntity, error) {
	tx, err := pgctx.Context(ctx)
	if err != nil {
		return nil, NewErrInsertStoryPlanRepository(fmt.Errorf("get postgres client: %w", err))
	}

	// Make sure the slug is unique.
	exists, err := tx.NewSelect().Model(&StoryPlanEntity{}).Where("slug = ?", data.Plan.Slug).Exists(ctx)
	if err != nil {
		return nil, NewErrInsertStoryPlanRepository(fmt.Errorf("check slug uniqueness: %w", err))
	}

	if exists {
		return nil, NewErrInsertStoryPlanRepository(ErrStoryPlanAlreadyExists)
	}

	entity := &StoryPlanEntity{
		ID:          data.Plan.ID,
		Slug:        data.Plan.Slug,
		Name:        data.Plan.Name,
		Description: data.Plan.Description,
		Beats:       data.Plan.Beats,
		CreatedAt:   data.Plan.CreatedAt,
	}

	_, err = tx.NewInsert().Model(entity).Returning("*").Exec(ctx)
	if err != nil {
		return nil, NewErrInsertStoryPlanRepository(fmt.Errorf("insert story plan: %w", err))
	}

	return entity, nil
}

func NewInsertStoryPlanRepository() *InsertStoryPlanRepository {
	return &InsertStoryPlanRepository{}
}
