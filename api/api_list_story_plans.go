package api

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/a-novel/service-story-schematics/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type ListStoryPlansService interface {
	ListStoryPlans(ctx context.Context, request services.ListStoryPlansRequest) ([]*models.StoryPlanPreview, error)
}

func (api *API) GetStoryPlans(
	ctx context.Context, params codegen.GetStoryPlansParams,
) (codegen.GetStoryPlansRes, error) {
	storyPlans, err := api.ListStoryPlansService.ListStoryPlans(ctx, services.ListStoryPlansRequest{
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("list story plans: %w", err)
	}

	res := codegen.GetStoryPlansOKApplicationJSON(
		lo.Map(storyPlans, func(item *models.StoryPlanPreview, _ int) codegen.StoryPlanPreview {
			return codegen.StoryPlanPreview{
				ID:          codegen.StoryPlanID(item.ID),
				Slug:        codegen.Slug(item.Slug),
				Name:        item.Name,
				Description: item.Description,
				CreatedAt:   item.CreatedAt,
			}
		},
		))

	return &res, nil
}
