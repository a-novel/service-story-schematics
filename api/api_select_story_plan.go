package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/a-novel/service-story-schematics/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type SelectStoryPlanService interface {
	SelectStoryPlan(ctx context.Context, request services.SelectStoryPlanRequest) (*models.StoryPlan, error)
}

func (api *API) GetStoryPlan(ctx context.Context, params codegen.GetStoryPlanParams) (codegen.GetStoryPlanRes, error) {
	span := sentry.StartSpan(ctx, "API.GetStoryPlan")
	defer span.Finish()

	span.SetData("request.slug", params.Slug)
	span.SetData("request.id", params.ID)

	storyPlan, err := api.SelectStoryPlanService.SelectStoryPlan(span.Context(), services.SelectStoryPlanRequest{
		Slug: lo.Ternary(params.Slug.IsSet(), lo.ToPtr(models.Slug(params.Slug.Value)), nil),
		ID:   lo.Ternary(params.ID.IsSet(), lo.ToPtr(uuid.UUID(params.ID.Value)), nil),
	})

	switch {
	case errors.Is(err, dao.ErrStoryPlanNotFound):
		span.SetData("service.err", err.Error())

		return &codegen.NotFoundError{Error: err.Error()}, nil
	case err != nil:
		span.SetData("service.err", err.Error())

		return nil, fmt.Errorf("get story plan: %w", err)
	}

	return &codegen.StoryPlan{
		ID:          codegen.StoryPlanID(storyPlan.ID),
		Slug:        codegen.Slug(storyPlan.Slug),
		Name:        storyPlan.Name,
		Description: storyPlan.Description,
		Beats: lo.Map(storyPlan.Beats, func(item models.BeatDefinition, _ int) codegen.BeatDefinition {
			return codegen.BeatDefinition{
				Name:      item.Name,
				Key:       item.Key,
				KeyPoints: item.KeyPoints,
				Purpose:   item.Purpose,
				MinScenes: item.MinScenes,
				MaxScenes: item.MaxScenes,
			}
		}),
		Lang:      codegen.Lang(storyPlan.Lang),
		CreatedAt: storyPlan.CreatedAt,
	}, nil
}
