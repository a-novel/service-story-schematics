package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrGenerateBeatsSheetService = errors.New("GenerateBeatsSheetService.GenerateBeatsSheet")

func NewErrGenerateBeatsSheetService(err error) error {
	return errors.Join(err, ErrGenerateBeatsSheetService)
}

type GenerateBeatsSheetSource interface {
	GenerateBeatsSheet(ctx context.Context, request daoai.GenerateBeatsSheetRequest) ([]models.Beat, error)
	SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error)
	SelectStoryPlan(ctx context.Context, data uuid.UUID) (*dao.StoryPlanEntity, error)
}

type GenerateBeatsSheetRequest struct {
	LoglineID   uuid.UUID
	StoryPlanID uuid.UUID
	UserID      uuid.UUID
}

type GenerateBeatsSheetService struct {
	source GenerateBeatsSheetSource
}

func (service *GenerateBeatsSheetService) GenerateBeatsSheet(
	ctx context.Context, request GenerateBeatsSheetRequest,
) ([]models.Beat, error) {
	logline, err := service.source.SelectLogline(ctx, dao.SelectLoglineData{
		ID:     request.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, NewErrGenerateBeatsSheetService(fmt.Errorf("get logline: %w", err))
	}

	storyPlan, err := service.source.SelectStoryPlan(ctx, request.StoryPlanID)
	if err != nil {
		return nil, NewErrGenerateBeatsSheetService(fmt.Errorf("get story plan: %w", err))
	}

	resp, err := service.source.GenerateBeatsSheet(ctx, daoai.GenerateBeatsSheetRequest{
		Logline: logline.Name + "\n\n" + logline.Content,
		Plan: models.StoryPlan{
			ID:          storyPlan.ID,
			Slug:        storyPlan.Slug,
			Name:        storyPlan.Name,
			Description: storyPlan.Description,
			Beats:       storyPlan.Beats,
			CreatedAt:   storyPlan.CreatedAt,
		},
		UserID: request.UserID.String(),
	})
	if err != nil {
		return nil, NewErrGenerateBeatsSheetService(err)
	}

	return resp, nil
}

func NewGenerateBeatsSheetServiceSource(
	generateDAO *daoai.GenerateBeatsSheetRepository,
	selectLoglineDAO *dao.SelectLoglineRepository,
	selectStoryPlanDAO *dao.SelectStoryPlanRepository,
) GenerateBeatsSheetSource {
	return &struct {
		*daoai.GenerateBeatsSheetRepository
		*dao.SelectLoglineRepository
		*dao.SelectStoryPlanRepository
	}{
		GenerateBeatsSheetRepository: generateDAO,
		SelectLoglineRepository:      selectLoglineDAO,
		SelectStoryPlanRepository:    selectStoryPlanDAO,
	}
}

func NewGenerateBeatsSheetService(source GenerateBeatsSheetSource) *GenerateBeatsSheetService {
	return &GenerateBeatsSheetService{source: source}
}
