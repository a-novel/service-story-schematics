package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"

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
	Lang        models.Lang
}

type GenerateBeatsSheetService struct {
	source GenerateBeatsSheetSource
}

func (service *GenerateBeatsSheetService) GenerateBeatsSheet(
	ctx context.Context, request GenerateBeatsSheetRequest,
) ([]models.Beat, error) {
	span := sentry.StartSpan(ctx, "GenerateBeatsSheetService.GenerateBeatsSheet")
	defer span.Finish()

	span.SetData("request.loglineID", request.LoglineID)
	span.SetData("request.storyPlanID", request.StoryPlanID)
	span.SetData("request.lang", request.Lang)
	span.SetData("request.userID", request.UserID)

	logline, err := service.source.SelectLogline(span.Context(), dao.SelectLoglineData{
		ID:     request.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		span.SetData("dao.selectLogline.err", err.Error())

		return nil, NewErrGenerateBeatsSheetService(fmt.Errorf("get logline: %w", err))
	}

	storyPlan, err := service.source.SelectStoryPlan(span.Context(), request.StoryPlanID)
	if err != nil {
		span.SetData("dao.selectStoryPlan.err", err.Error())

		return nil, NewErrGenerateBeatsSheetService(fmt.Errorf("get story plan: %w", err))
	}

	resp, err := service.source.GenerateBeatsSheet(span.Context(), daoai.GenerateBeatsSheetRequest{
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
		Lang:   request.Lang,
		UserID: request.UserID.String(),
	})
	if err != nil {
		span.SetData("daoai.generateBeatsSheet.err", err.Error())

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
