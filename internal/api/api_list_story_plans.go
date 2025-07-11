package api

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/samber/lo"

	"github.com/a-novel/service-story-schematics/internal/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type ListStoryPlansService interface {
	ListStoryPlans(ctx context.Context, request services.ListStoryPlansRequest) ([]*models.StoryPlanPreview, error)
}

func (api *API) GetStoryPlans(
	ctx context.Context, params codegen.GetStoryPlansParams,
) (codegen.GetStoryPlansRes, error) {
	span := sentry.StartSpan(ctx, "API.GetStoryPlans")
	defer span.Finish()

	span.SetData("request.limit", params.Limit)
	span.SetData("request.offset", params.Offset)

	storyPlans, err := api.ListStoryPlansService.ListStoryPlans(span.Context(), services.ListStoryPlansRequest{
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	})
	if err != nil {
		span.SetData("service.err", err.Error())

		return nil, fmt.Errorf("list story plans: %w", err)
	}

	res := codegen.GetStoryPlansOKApplicationJSON(
		lo.Map(storyPlans, func(item *models.StoryPlanPreview, _ int) codegen.StoryPlanPreview {
			return codegen.StoryPlanPreview{
				ID:          codegen.StoryPlanID(item.ID),
				Slug:        codegen.Slug(item.Slug),
				Name:        item.Name,
				Description: item.Description,
				Lang:        codegen.Lang(item.Lang),
				CreatedAt:   item.CreatedAt,
			}
		}),
	)

	span.SetData("storyPlans.count", len(storyPlans))

	return &res, nil
}
