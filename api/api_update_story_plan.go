package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/samber/lo"

	"github.com/a-novel/service-story-schematics/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type UpdateStoryPlanService interface {
	UpdateStoryPlan(ctx context.Context, request services.UpdateStoryPlanRequest) (*models.StoryPlan, error)
}

func (api *API) UpdateStoryPlan(
	ctx context.Context, req *codegen.UpdateStoryPlanForm,
) (codegen.UpdateStoryPlanRes, error) {
	span := sentry.StartSpan(ctx, "API.UpdateStoryPlan")
	defer span.Finish()

	span.SetData("request.slug", req.GetSlug())
	span.SetData("request.name", req.GetName())
	span.SetData("request.lang", req.GetLang())

	storyPlan, err := api.UpdateStoryPlanService.UpdateStoryPlan(span.Context(), services.UpdateStoryPlanRequest{
		Slug:        models.Slug(req.GetSlug()),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Beats: lo.Map(req.GetBeats(), func(item codegen.BeatDefinition, _ int) models.BeatDefinition {
			return models.BeatDefinition{
				Name:      item.GetName(),
				Key:       item.GetKey(),
				KeyPoints: item.GetKeyPoints(),
				Purpose:   item.GetPurpose(),
				MinScenes: item.GetMinScenes(),
				MaxScenes: item.GetMaxScenes(),
			}
		}),
		Lang: models.Lang(req.GetLang()),
	})

	switch {
	case errors.Is(err, dao.ErrStoryPlanNotFound):
		span.SetData("service.err", err.Error())

		return &codegen.NotFoundError{Error: err.Error()}, nil
	case err != nil:
		span.SetData("service.err", err.Error())

		return nil, fmt.Errorf("update story plan: %w", err)
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
