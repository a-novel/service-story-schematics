package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel/codes"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

type UpdateStoryPlanService interface {
	UpdateStoryPlan(ctx context.Context, request services.UpdateStoryPlanRequest) (*models.StoryPlan, error)
}

func (api *API) UpdateStoryPlan(
	ctx context.Context, req *apimodels.UpdateStoryPlanForm,
) (apimodels.UpdateStoryPlanRes, error) {
	ctx, span := otel.Tracer().Start(ctx, "api.UpdateStoryPlan")
	defer span.End()

	storyPlan, err := api.UpdateStoryPlanService.UpdateStoryPlan(ctx, services.UpdateStoryPlanRequest{
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

	switch {
	case errors.Is(err, dao.ErrStoryPlanNotFound):
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return &apimodels.NotFoundError{Error: err.Error()}, nil
	case err != nil:
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return nil, fmt.Errorf("update story plan: %w", err)
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
