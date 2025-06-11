package services

import (
	"context"
	"errors"

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

type SelectStoryPlanRequest struct {
	Slug *models.Slug
	ID   *uuid.UUID
}

type SelectStoryPlanService struct {
	source SelectStoryPlanSource
}

func (service *SelectStoryPlanService) SelectStoryPlan(
	ctx context.Context, request SelectStoryPlanRequest,
) (*models.StoryPlan, error) {
	if request.Slug != nil {
		data, err := service.source.SelectStoryPlanBySlug(ctx, lo.FromPtr(request.Slug))
		if err != nil {
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

	data, err := service.source.SelectStoryPlan(ctx, lo.FromPtr(request.ID))
	if err != nil {
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

func NewSelectStoryPlanService(source SelectStoryPlanSource) *SelectStoryPlanService {
	return &SelectStoryPlanService{source: source}
}
