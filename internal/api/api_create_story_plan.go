package api

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

type CreateStoryPlanService interface {
	CreateStoryPlan(ctx context.Context, request services.CreateStoryPlanRequest) (*models.StoryPlan, error)
}

func (api *API) CreateStoryPlan(
	ctx context.Context, req *apimodels.CreateStoryPlanForm,
) (apimodels.CreateStoryPlanRes, error) {
	ctx, span := otel.Tracer().Start(ctx, "api.CreateStoryPlan")
	defer span.End()

	storyPlan, err := api.CreateStoryPlanService.CreateStoryPlan(ctx, services.CreateStoryPlanRequest{
		Slug:        models.Slug(req.GetSlug()),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Beats: lo.Map(req.GetBeats(), func(item apimodels.BeatDefinition, _ int) models.BeatDefinition {
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
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("create story plan: %w", err))
	}

	return otel.ReportSuccess(span, &apimodels.StoryPlan{
		ID:          apimodels.StoryPlanID(storyPlan.ID),
		Slug:        apimodels.Slug(storyPlan.Slug),
		Name:        storyPlan.Name,
		Description: storyPlan.Description,
		Beats: lo.Map(storyPlan.Beats, func(item models.BeatDefinition, _ int) apimodels.BeatDefinition {
			return apimodels.BeatDefinition{
				Name:      item.Name,
				Key:       item.Key,
				KeyPoints: item.KeyPoints,
				Purpose:   item.Purpose,
				MinScenes: item.MinScenes,
				MaxScenes: item.MaxScenes,
			}
		}),
		Lang:      apimodels.Lang(storyPlan.Lang),
		CreatedAt: storyPlan.CreatedAt,
	}), nil
}
