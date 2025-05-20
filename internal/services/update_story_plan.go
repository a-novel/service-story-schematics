package services

import (
	"context"
	"errors"
	"time"

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
}

type UpdateStoryPlanService struct {
	source UpdateStoryPlanSource
}

func (service *UpdateStoryPlanService) UpdateStoryPlan(
	ctx context.Context, request UpdateStoryPlanRequest,
) (*models.StoryPlan, error) {
	resp, err := service.source.UpdateStoryPlan(ctx, dao.UpdateStoryPlanData{
		Plan: models.StoryPlan{
			ID:          uuid.New(),
			Slug:        request.Slug,
			Name:        request.Name,
			Description: request.Description,
			Beats:       request.Beats,
			CreatedAt:   time.Now(),
		},
	})
	if err != nil {
		return nil, NewErrUpdateStoryPlanService(err)
	}

	return &models.StoryPlan{
		ID:          resp.ID,
		Slug:        resp.Slug,
		Name:        resp.Name,
		Description: resp.Description,
		Beats:       resp.Beats,
		CreatedAt:   resp.CreatedAt,
	}, nil
}

func NewUpdateStoryPlanService(source UpdateStoryPlanSource) *UpdateStoryPlanService {
	return &UpdateStoryPlanService{source: source}
}
