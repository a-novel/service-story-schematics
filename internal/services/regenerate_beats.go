package services

import (
	"context"
	"errors"
	"github.com/getsentry/sentry-go"

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
	span := sentry.StartSpan(ctx, "RegenerateBeatsService.RegenerateBeats")
	defer span.Finish()

	span.SetData("request.beatsSheetID", request.BeatsSheetID)
	span.SetData("request.userID", request.UserID)
	span.SetData("request.regenerateKeys", request.RegenerateKeys)

	beatsSheet, err := service.source.SelectBeatsSheet(span.Context(), request.BeatsSheetID)
	if err != nil {
		span.SetData("dao.selectBeatsSheet.err", err.Error())

		return nil, NewErrRegenerateBeatsService(err)
	}

	// Make sure the selected beats sheet is linked to a logline that belongs to the user.
	logline, err := service.source.SelectLogline(span.Context(), dao.SelectLoglineData{
		ID:     beatsSheet.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		span.SetData("dao.selectLogline.err", err.Error())

		return nil, NewErrRegenerateBeatsService(err)
	}

	storyPlan, err := service.source.SelectStoryPlan(span.Context(), beatsSheet.StoryPlanID)
	if err != nil {
		span.SetData("dao.selectStoryPlan.err", err.Error())

		return nil, NewErrRegenerateBeatsService(err)
	}

	regenerated, err := service.source.RegenerateBeats(span.Context(), daoai.RegenerateBeatsRequest{
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
		span.SetData("daoai.regenerateBeats.err", err.Error())

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
