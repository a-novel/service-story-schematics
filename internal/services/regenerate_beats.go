package services

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrRegenerateBeatsService = errors.New("RegenerateBeatsService.RegenerateBeatsSheet")

func NewErrRegenerateBeatsService(err error) error {
	return errors.Join(err, ErrRegenerateBeatsService)
}

type RegenerateBeatsSource interface {
	RegenerateBeats(ctx context.Context, request daoai.RegenerateBeatsRequest) ([]models.Beat, error)
	SelectBeatsSheet(ctx context.Context, data uuid.UUID) (*dao.BeatsSheetEntity, error)
	SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error)
	SelectStoryPlan(ctx context.Context, data uuid.UUID) (*dao.StoryPlanEntity, error)
}

type RegenerateBeatsRequest struct {
	BeatsSheetID   uuid.UUID
	UserID         uuid.UUID
	RegenerateKeys []string
}

type RegenerateBeatsService struct {
	source RegenerateBeatsSource
}

func (service *RegenerateBeatsService) RegenerateBeats(
	ctx context.Context, request RegenerateBeatsRequest,
) ([]models.Beat, error) {
	beatsSheet, err := service.source.SelectBeatsSheet(ctx, request.BeatsSheetID)
	if err != nil {
		return nil, NewErrRegenerateBeatsService(err)
	}

	// Make sure the selected beats sheet is linked to a logline that belongs to the user.
	logline, err := service.source.SelectLogline(ctx, dao.SelectLoglineData{
		ID:     beatsSheet.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, NewErrRegenerateBeatsService(err)
	}

	storyPlan, err := service.source.SelectStoryPlan(ctx, beatsSheet.StoryPlanID)
	if err != nil {
		return nil, NewErrRegenerateBeatsService(err)
	}

	regenerated, err := service.source.RegenerateBeats(ctx, daoai.RegenerateBeatsRequest{
		Logline: logline.Name + "\n\n" + logline.Content,
		Plan: models.StoryPlan{
			ID:          storyPlan.ID,
			Slug:        storyPlan.Slug,
			Name:        storyPlan.Name,
			Description: storyPlan.Description,
			Beats:       storyPlan.Beats,
			Lang:        storyPlan.Lang,
			CreatedAt:   storyPlan.CreatedAt,
		},
		UserID:         request.UserID.String(),
		Lang:           beatsSheet.Lang,
		Beats:          beatsSheet.Content,
		RegenerateKeys: request.RegenerateKeys,
	})
	if err != nil {
		return nil, NewErrRegenerateBeatsService(err)
	}

	return regenerated, nil
}

func NewRegenerateBeatsServiceSource(
	regenerateBeatsDAO *daoai.RegenerateBeatsRepository,
	selectBeatsSheetDAO *dao.SelectBeatsSheetRepository,
	selectLoglineDAO *dao.SelectLoglineRepository,
	selectStoryPlanDAO *dao.SelectStoryPlanRepository,
) RegenerateBeatsSource {
	return &struct {
		*daoai.RegenerateBeatsRepository
		*dao.SelectBeatsSheetRepository
		*dao.SelectLoglineRepository
		*dao.SelectStoryPlanRepository
	}{
		RegenerateBeatsRepository:  regenerateBeatsDAO,
		SelectBeatsSheetRepository: selectBeatsSheetDAO,
		SelectLoglineRepository:    selectLoglineDAO,
		SelectStoryPlanRepository:  selectStoryPlanDAO,
	}
}

func NewRegenerateBeatsService(source RegenerateBeatsSource) *RegenerateBeatsService {
	return &RegenerateBeatsService{source: source}
}
