package services

import (
	"context"
	"errors"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrSelectStoryPlanService = errors.New("SelectStoryPlanService.SelectStoryPlan")

func NewErrSelectStoryPlanService(err error) error {
	return errors.Join(err, ErrSelectStoryPlanService)
}

type SelectStoryPlanSource interface {
	SelectStoryPlan(ctx context.Context, data uuid.UUID) (*dao.StoryPlanEntity, error)
	SelectStoryPlanBySlug(ctx context.Context, data models.Slug) (*dao.StoryPlanEntity, error)
}

func NewSelectStoryPlanServiceSource(
	selectDAO *dao.SelectStoryPlanRepository,
	selectBySlugRepository *dao.SelectStoryPlanBySlugRepository,
) SelectStoryPlanSource {
	return &struct {
		*dao.SelectStoryPlanRepository
		*dao.SelectStoryPlanBySlugRepository
	}{
		SelectStoryPlanRepository:       selectDAO,
		SelectStoryPlanBySlugRepository: selectBySlugRepository,
	}
}

type SelectStoryPlanRequest struct {
	Slug *models.Slug
	ID   *uuid.UUID
}

type SelectStoryPlanService struct {
	source SelectStoryPlanSource
}

func NewSelectStoryPlanService(source SelectStoryPlanSource) *SelectStoryPlanService {
	return &SelectStoryPlanService{source: source}
}

func (service *SelectStoryPlanService) SelectStoryPlan(
	ctx context.Context, request SelectStoryPlanRequest,
) (*models.StoryPlan, error) {
	span := sentry.StartSpan(ctx, "SelectStoryPlanService.SelectStoryPlan")
	defer span.Finish()

	span.SetData("request.slug", request.Slug)
	span.SetData("request.id", request.ID)

	if request.Slug != nil {
		data, err := service.source.SelectStoryPlanBySlug(span.Context(), lo.FromPtr(request.Slug))
		if err != nil {
			span.SetData("dao.selectStoryPlanBySlug.err", err.Error())

			return nil, NewErrSelectStoryPlanService(err)
		}

		return &models.StoryPlan{
			ID:          data.ID,
			Slug:        data.Slug,
			Name:        data.Name,
			Description: data.Description,
			Beats:       data.Beats,
			Lang:        data.Lang,
			CreatedAt:   data.CreatedAt,
		}, nil
	}

	data, err := service.source.SelectStoryPlan(span.Context(), lo.FromPtr(request.ID))
	if err != nil {
		span.SetData("dao.selectStoryPlan.err", err.Error())

		return nil, NewErrSelectStoryPlanService(err)
	}

	return &models.StoryPlan{
		ID:          data.ID,
		Slug:        data.Slug,
		Name:        data.Name,
		Description: data.Description,
		Beats:       data.Beats,
		Lang:        data.Lang,
		CreatedAt:   data.CreatedAt,
	}, nil
}
