package services

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/models"
	storyplanmodel "github.com/a-novel/service-story-schematics/models/story_plan"
)

var ErrStoryPlanNotFound = errors.New("story plan not found")

type SelectStoryPlanRequest struct {
	Lang models.Lang
}

type SelectStoryPlanService struct{}

func NewSelectStoryPlanService() *SelectStoryPlanService {
	return &SelectStoryPlanService{}
}

func (service *SelectStoryPlanService) SelectStoryPlan(
	ctx context.Context, request SelectStoryPlanRequest,
) (*storyplanmodel.Plan, error) {
	_, span := otel.Tracer().Start(ctx, "service.SelectStoryPlan")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.lang", request.Lang.String()),
	)

	plan, ok := storyplanmodel.SaveTheCat[request.Lang]
	if !ok {
		return nil, ErrStoryPlanNotFound
	}

	return plan, nil
}
