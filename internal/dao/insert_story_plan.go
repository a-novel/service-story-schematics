package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/getsentry/sentry-go"

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
	span := sentry.StartSpan(ctx, "InsertStoryPlanRepository.InsertStoryPlan")
	defer span.Finish()

	span.SetData("story_plan.id", data.Plan.ID.String())
	span.SetData("story_plan.slug", data.Plan.Slug)
	span.SetData("story_plan.name", data.Plan.Name)
	span.SetData("story_plan.lang", data.Plan.Lang)

	tx, err := lib.PostgresContext(span.Context())
	if err != nil {
		span.SetData("postgres.context.error", err.Error())

		return nil, NewErrInsertStoryPlanRepository(fmt.Errorf("get postgres client: %w", err))
	}

	// Make sure the slug is unique.
	exists, err := tx.NewSelect().Model(&StoryPlanEntity{}).Where("slug = ?", data.Plan.Slug).Exists(span.Context())
	if err != nil {
		span.SetData("check.slug.uniqueness.error", err.Error())

		return nil, NewErrInsertStoryPlanRepository(fmt.Errorf("check slug uniqueness: %w", err))
	}

	if exists {
		span.SetData("check.slug.uniqueness.exists", true)

		return nil, NewErrInsertStoryPlanRepository(ErrStoryPlanAlreadyExists)
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

	_, err = tx.NewInsert().Model(entity).Returning("*").Exec(span.Context())
	if err != nil {
		span.SetData("insert.error", err.Error())

		return nil, NewErrInsertStoryPlanRepository(fmt.Errorf("insert story plan: %w", err))
	}

	return entity, nil
}

func NewInsertStoryPlanRepository() *InsertStoryPlanRepository {
	return &InsertStoryPlanRepository{}
}
