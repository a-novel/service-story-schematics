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

type ListStoryPlansService interface {
	ListStoryPlans(ctx context.Context, request services.ListStoryPlansRequest) ([]*models.StoryPlanPreview, error)
}

func (api *API) GetStoryPlans(
	ctx context.Context, params apimodels.GetStoryPlansParams,
) (apimodels.GetStoryPlansRes, error) {
	ctx, span := otel.Tracer().Start(ctx, "api.GetStoryPlans")
	defer span.End()

	storyPlans, err := api.ListStoryPlansService.ListStoryPlans(ctx, services.ListStoryPlansRequest{
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("list story plans: %w", err))
	}

	res := apimodels.GetStoryPlansOKApplicationJSON(
		lo.Map(storyPlans, func(item *models.StoryPlanPreview, _ int) apimodels.StoryPlanPreview {
			return apimodels.StoryPlanPreview{
				ID:          apimodels.StoryPlanID(item.ID),
				Slug:        apimodels.Slug(item.Slug),
				Name:        item.Name,
				Description: item.Description,
				Lang:        apimodels.Lang(item.Lang),
				CreatedAt:   item.CreatedAt,
			}
		}),
	)

	return otel.ReportSuccess(span, &res), nil
}
