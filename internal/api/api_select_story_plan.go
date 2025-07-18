package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/codes"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

type SelectStoryPlanService interface {
	SelectStoryPlan(ctx context.Context, request services.SelectStoryPlanRequest) (*models.StoryPlan, error)
}

func (api *API) GetStoryPlan(
	ctx context.Context, params apimodels.GetStoryPlanParams) (apimodels.GetStoryPlanRes, error,
) {
	ctx, span := otel.Tracer().Start(ctx, "api.GetStoryPlan")
	defer span.End()

	storyPlan, err := api.SelectStoryPlanService.SelectStoryPlan(ctx, services.SelectStoryPlanRequest{
		Slug: lo.Ternary(params.Slug.IsSet(), lo.ToPtr(models.Slug(params.Slug.Value)), nil),
		ID:   lo.Ternary(params.ID.IsSet(), lo.ToPtr(uuid.UUID(params.ID.Value)), nil),
	})

	switch {
	case errors.Is(err, dao.ErrStoryPlanNotFound):
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return &apimodels.NotFoundError{Error: err.Error()}, nil
	case err != nil:
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return nil, fmt.Errorf("get story plan: %w", err)
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
