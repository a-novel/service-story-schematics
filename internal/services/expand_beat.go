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

type ExpandBeatRequest struct {
	BeatsSheetID uuid.UUID
	TargetKey    string
	UserID       uuid.UUID
}

type ExpandBeatService struct {
	source ExpandBeatSource
}

func NewExpandBeatService(source ExpandBeatSource) *ExpandBeatService {
	return &ExpandBeatService{source: source}
}

func (service *ExpandBeatService) ExpandBeat(
	ctx context.Context, request ExpandBeatRequest,
) (*models.Beat, error) {
	span := sentry.StartSpan(ctx, "ExpandBeatService.ExpandBeat")
	defer span.Finish()

	span.SetData("request.beatsSheetID", request.BeatsSheetID)
	span.SetData("request.targetKey", request.TargetKey)
	span.SetData("request.userID", request.UserID)

	beatsSheet, err := service.source.SelectBeatsSheet(span.Context(), request.BeatsSheetID)
	if err != nil {
		span.SetData("dao.selectBeatsSheet.err", err.Error())

		return nil, NewErrExpandBeatService(err)
	}

	// Make sure the selected beats sheet is linked to a logline that belongs to the user.
	logline, err := service.source.SelectLogline(span.Context(), dao.SelectLoglineData{
		ID:     beatsSheet.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		span.SetData("dao.selectLogline.err", err.Error())

		return nil, NewErrExpandBeatService(err)
	}

	storyPlan, err := service.source.SelectStoryPlan(span.Context(), beatsSheet.StoryPlanID)
	if err != nil {
		span.SetData("dao.selectStoryPlan.err", err.Error())

		return nil, NewErrExpandBeatService(err)
	}

	expanded, err := service.source.ExpandBeat(span.Context(), daoai.ExpandBeatRequest{
		Logline: logline.Name + "\n\n" + logline.Content,
		Beats:   beatsSheet.Content,
		Plan: models.StoryPlan{
			ID:          storyPlan.ID,
			Slug:        storyPlan.Slug,
			Name:        storyPlan.Name,
			Description: storyPlan.Description,
			Beats:       storyPlan.Beats,
			CreatedAt:   storyPlan.CreatedAt,
			Lang:        beatsSheet.Lang,
		},
		Lang:      beatsSheet.Lang,
		TargetKey: request.TargetKey,
		UserID:    request.UserID.String(),
	})
	if err != nil {
		span.SetData("daoai.expandBeat.err", err.Error())

		return nil, NewErrExpandBeatService(err)
	}

	return expanded, nil
}
