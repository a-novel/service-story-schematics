package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-novel/service-story-schematics/internal/lib"

	"github.com/a-novel/service-story-schematics/models"
)

var ErrUpdateStoryPlanRepository = errors.New("UpdateStoryPlanRepository.UpdateStoryPlan")

func NewErrUpdateStoryPlanRepository(err error) error {
	return errors.Join(err, ErrUpdateStoryPlanRepository)
}

type UpdateStoryPlanData struct {
	Plan models.StoryPlan
}

type UpdateStoryPlanRepository struct{}

func (repository *UpdateStoryPlanRepository) UpdateStoryPlan(
	ctx context.Context, data UpdateStoryPlanData,
) (*StoryPlanEntity, error) {
	db, err := lib.PostgresContext(ctx)
	if err != nil {
		return nil, NewErrUpdateStoryPlanRepository(fmt.Errorf("get postgres client: %w", err))
	}

	// Make sure the slug is unique.
	exists, err := db.NewSelect().Model(&StoryPlanEntity{}).Where("slug = ?", data.Plan.Slug).Exists(ctx)
	if err != nil {
		return nil, NewErrUpdateStoryPlanRepository(fmt.Errorf("check slug uniqueness: %w", err))
	}

	if !exists {
		return nil, NewErrUpdateStoryPlanRepository(ErrStoryPlanNotFound)
	}

	entity := &StoryPlanEntity{
		ID:          data.Plan.ID,
		Slug:        data.Plan.Slug,
		Name:        data.Plan.Name,
		Description: data.Plan.Description,
		Beats:       data.Plan.Beats,
		Lang:        data.Plan.Lang,
		CreatedAt:   data.Plan.CreatedAt,
	}

	_, err = db.NewInsert().Model(entity).Returning("*").Exec(ctx)
	if err != nil {
		return nil, NewErrUpdateStoryPlanRepository(fmt.Errorf("update story plan: %w", err))
	}

	return entity, nil
}

func NewUpdateStoryPlanRepository() *UpdateStoryPlanRepository {
	return &UpdateStoryPlanRepository{}
}
