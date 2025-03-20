package services

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/a-novel/story-schematics/internal/dao"
	"github.com/a-novel/story-schematics/internal/daoai"
	"github.com/a-novel/story-schematics/models"
)

var ErrExpandBeatService = errors.New("ExpandBeatService.ExpandBeat")

func NewErrExpandBeatService(err error) error {
	return errors.Join(err, ErrExpandBeatService)
}

type ExpandBeatSource interface {
	ExpandBeat(ctx context.Context, request daoai.ExpandBeatRequest) (*models.Beat, error)
	SelectBeatsSheet(ctx context.Context, data uuid.UUID) (*dao.BeatsSheetEntity, error)
	SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error)
	SelectStoryPlan(ctx context.Context, data uuid.UUID) (*dao.StoryPlanEntity, error)
}

type ExpandBeatRequest struct {
	BeatsSheetID uuid.UUID
	TargetKey    string
	UserID       uuid.UUID
}

type ExpandBeatService struct {
	source ExpandBeatSource
}

func (service *ExpandBeatService) ExpandBeat(
	ctx context.Context, request ExpandBeatRequest,
) (*models.Beat, error) {
	beatsSheet, err := service.source.SelectBeatsSheet(ctx, request.BeatsSheetID)
	if err != nil {
		return nil, NewErrExpandBeatService(err)
	}

	// Make sure the selected beats sheet is linked to a logline that belongs to the user.
	logline, err := service.source.SelectLogline(ctx, dao.SelectLoglineData{
		ID:     beatsSheet.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, NewErrExpandBeatService(err)
	}

	storyPlan, err := service.source.SelectStoryPlan(ctx, beatsSheet.StoryPlanID)
	if err != nil {
		return nil, NewErrExpandBeatService(err)
	}

	expanded, err := service.source.ExpandBeat(ctx, daoai.ExpandBeatRequest{
		Logline: logline.Name + "\n\n" + logline.Content,
		Beats:   beatsSheet.Content,
		Plan: models.StoryPlan{
			ID:          storyPlan.ID,
			Slug:        storyPlan.Slug,
			Name:        storyPlan.Name,
			Description: storyPlan.Description,
			Beats:       storyPlan.Beats,
			CreatedAt:   storyPlan.CreatedAt,
		},
		TargetKey: request.TargetKey,
		UserID:    request.UserID.String(),
	})
	if err != nil {
		return nil, NewErrExpandBeatService(err)
	}

	return expanded, nil
}

func NewExpandBeatServiceSource(
	expandBeatDAO *daoai.ExpandBeatRepository,
	selectBeatsSheetDAO *dao.SelectBeatsSheetRepository,
	selectLoglineDAO *dao.SelectLoglineRepository,
	selectStoryPlanDAO *dao.SelectStoryPlanRepository,
) ExpandBeatSource {
	return &struct {
		*daoai.ExpandBeatRepository
		*dao.SelectBeatsSheetRepository
		*dao.SelectLoglineRepository
		*dao.SelectStoryPlanRepository
	}{
		ExpandBeatRepository:       expandBeatDAO,
		SelectBeatsSheetRepository: selectBeatsSheetDAO,
		SelectLoglineRepository:    selectLoglineDAO,
		SelectStoryPlanRepository:  selectStoryPlanDAO,
	}
}

func NewExpandBeatService(source ExpandBeatSource) *ExpandBeatService {
	return &ExpandBeatService{source: source}
}
