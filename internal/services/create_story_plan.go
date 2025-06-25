package services

import (
	"context"
	"errors"
	"github.com/getsentry/sentry-go"
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
	Lang        models.Lang
}

type CreateStoryPlanService struct {
	source CreateStoryPlanSource
}

func (service *CreateStoryPlanService) CreateStoryPlan(
	ctx context.Context, request CreateStoryPlanRequest,
) (*models.StoryPlan, error) {
	span := sentry.StartSpan(ctx, "CreateStoryPlanService.CreateStoryPlan")
	defer span.Finish()

	span.SetData("request.slug", request.Slug)
	span.SetData("request.name", request.Name)
	span.SetData("request.lang", request.Lang)

	data := dao.InsertStoryPlanData{
		Plan: models.StoryPlan{
			ID:          uuid.New(),
			Slug:        request.Slug,
			Name:        request.Name,
			Description: request.Description,
			Beats:       request.Beats,
			Lang:        request.Lang,
			CreatedAt:   time.Now(),
		},
	}

	resp, err := service.source.InsertStoryPlan(span.Context(), data)

	// If slug is taken, try to modify it by appending a version number.
	if errors.Is(err, dao.ErrStoryPlanAlreadyExists) {
		span.SetData("dao.insertStoryPlan.slug.taken", err.Error())

		data.Plan.Slug, _, err = service.source.SelectSlugIteration(span.Context(), dao.SelectSlugIterationData{
			Slug:  data.Plan.Slug,
			Table: "story_plans",
			Order: []string{"created_at DESC"},
		})
		if err != nil {
			span.SetData("dao.selectSlugIteration.err", err.Error())

			return nil, NewErrCreateStoryPlanService(err)
		}

		resp, err = service.source.InsertStoryPlan(span.Context(), data)
	}

	if err != nil {
		span.SetData("dao.insertStoryPlan.err", err.Error())

		return nil, NewErrCreateStoryPlanService(err)
	}

	span.SetData("dao.insertStoryPlan.id", resp.ID)

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
