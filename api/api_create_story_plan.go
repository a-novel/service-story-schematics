package api

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/a-novel/service-story-schematics/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type CreateStoryPlanService interface {
	CreateStoryPlan(ctx context.Context, request services.CreateStoryPlanRequest) (*models.StoryPlan, error)
}

func (api *API) CreateStoryPlan(
	ctx context.Context, req *codegen.CreateStoryPlanForm,
) (codegen.CreateStoryPlanRes, error) {
	storyPlan, err := api.CreateStoryPlanService.CreateStoryPlan(ctx, services.CreateStoryPlanRequest{
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
	if err != nil {
		return nil, fmt.Errorf("create story plan: %w", err)
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
