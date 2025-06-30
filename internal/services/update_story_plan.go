package services

import (
	"context"
	"errors"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrUpdateStoryPlanService = errors.New("UpdateStoryPlanService.UpdateStoryPlan")

func NewErrUpdateStoryPlanService(err error) error {
	return errors.Join(err, ErrUpdateStoryPlanService)
}

type UpdateStoryPlanSource interface {
	UpdateStoryPlan(ctx context.Context, data dao.UpdateStoryPlanData) (*dao.StoryPlanEntity, error)
}

type UpdateStoryPlanRequest struct {
	Slug        models.Slug
	Name        string
	Description string
	Beats       []models.BeatDefinition
	Lang        models.Lang
}

type UpdateStoryPlanService struct {
	source UpdateStoryPlanSource
}

func NewUpdateStoryPlanService(source UpdateStoryPlanSource) *UpdateStoryPlanService {
	return &UpdateStoryPlanService{source: source}
}

func (service *UpdateStoryPlanService) UpdateStoryPlan(
	ctx context.Context, request UpdateStoryPlanRequest,
) (*models.StoryPlan, error) {
	span := sentry.StartSpan(ctx, "UpdateStoryPlanService.UpdateStoryPlan")
	defer span.Finish()

	span.SetData("request.slug", request.Slug)
	span.SetData("request.name", request.Name)
	span.SetData("request.lang", request.Lang)

	resp, err := service.source.UpdateStoryPlan(span.Context(), dao.UpdateStoryPlanData{
		Plan: models.StoryPlan{
			ID:          uuid.New(),
			Slug:        request.Slug,
			Name:        request.Name,
			Description: request.Description,
			Beats:       request.Beats,
			Lang:        request.Lang,
			CreatedAt:   time.Now(),
		},
	})
	if err != nil {
		span.SetData("dao.updateStoryPlan.err", err.Error())

		return nil, NewErrUpdateStoryPlanService(err)
	}

	return &models.StoryPlan{
		ID:          resp.ID,
		Slug:        resp.Slug,
		Name:        resp.Name,
		Description: resp.Description,
		Beats:       resp.Beats,
		Lang:        resp.Lang,
		CreatedAt:   resp.CreatedAt,
	}, nil
}
