package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"

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
	span := sentry.StartSpan(ctx, "UpdateStoryPlanRepository.UpdateStoryPlan")
	defer span.Finish()

	span.SetData("story_plan.id", data.Plan.ID.String())
	span.SetData("story_plan.slug", data.Plan.Slug)
	span.SetData("story_plan.name", data.Plan.Name)
	span.SetData("story_plan.lang", data.Plan.Lang)

	db, err := lib.PostgresContext(span.Context())
	if err != nil {
		span.SetData("postgres.context.error", err.Error())

		return nil, NewErrUpdateStoryPlanRepository(fmt.Errorf("get postgres client: %w", err))
	}

	// Make sure the slug is unique.
	exists, err := db.NewSelect().Model(&StoryPlanEntity{}).Where("slug = ?", data.Plan.Slug).Exists(span.Context())
	if err != nil {
		span.SetData("check.slug.error", err.Error())

		return nil, NewErrUpdateStoryPlanRepository(fmt.Errorf("check slug uniqueness: %w", err))
	}

	if !exists {
		span.SetData("check.slug.error", "slug not found")

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

	_, err = db.NewInsert().Model(entity).Returning("*").Exec(span.Context())
	if err != nil {
		span.SetData("insert.error", err.Error())

		return nil, NewErrUpdateStoryPlanRepository(fmt.Errorf("update story plan: %w", err))
	}

	return entity, nil
}

func NewUpdateStoryPlanRepository() *UpdateStoryPlanRepository {
	return &UpdateStoryPlanRepository{}
}
