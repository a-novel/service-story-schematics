package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrCreateStoryPlanService = errors.New("CreateStoryPlanService.CreateStoryPlan")

func NewErrCreateStoryPlanService(err error) error {
	return errors.Join(err, ErrCreateStoryPlanService)
}

type CreateStoryPlanSource interface {
	InsertStoryPlan(ctx context.Context, data dao.InsertStoryPlanData) (*dao.StoryPlanEntity, error)
	SelectSlugIteration(ctx context.Context, data dao.SelectSlugIterationData) (models.Slug, int, error)
}

type CreateStoryPlanRequest struct {
	Slug        models.Slug
	Name        string
	Description string
	Beats       []models.BeatDefinition
}

type CreateStoryPlanService struct {
	source CreateStoryPlanSource
}

func (service *CreateStoryPlanService) CreateStoryPlan(
	ctx context.Context, request CreateStoryPlanRequest,
) (*models.StoryPlan, error) {
	data := dao.InsertStoryPlanData{
		Plan: models.StoryPlan{
			ID:          uuid.New(),
			Slug:        request.Slug,
			Name:        request.Name,
			Description: request.Description,
			Beats:       request.Beats,
			CreatedAt:   time.Now(),
		},
	}

	resp, err := service.source.InsertStoryPlan(ctx, data)

	// If slug is taken, try to modify it by appending a version number.
	if errors.Is(err, dao.ErrStoryPlanAlreadyExists) {
		data.Plan.Slug, _, err = service.source.SelectSlugIteration(ctx, dao.SelectSlugIterationData{
			Slug:  data.Plan.Slug,
			Table: "story_plans",
			Order: []string{"created_at DESC"},
		})
		if err != nil {
			return nil, NewErrCreateStoryPlanService(err)
		}

		resp, err = service.source.InsertStoryPlan(ctx, data)
	}

	if err != nil {
		return nil, NewErrCreateStoryPlanService(err)
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

func NewCreateStoryPlanServiceSource(
	insertStoryPlanDAO *dao.InsertStoryPlanRepository,
	selectSlugIterationDAO *dao.SelectSlugIterationRepository,
) CreateStoryPlanSource {
	return &struct {
		*dao.InsertStoryPlanRepository
		*dao.SelectSlugIterationRepository
	}{
		InsertStoryPlanRepository:     insertStoryPlanDAO,
		SelectSlugIterationRepository: selectSlugIterationDAO,
	}
}

func NewCreateStoryPlanService(source CreateStoryPlanSource) *CreateStoryPlanService {
	return &CreateStoryPlanService{source: source}
}
