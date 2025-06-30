package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrCreateBeatsSheetService = errors.New("CreateBeatsSheetService.CreateBeatsSheet")

func NewErrCreateBeatsSheetService(err error) error {
	return errors.Join(err, ErrCreateBeatsSheetService)
}

type CreateBeatsSheetSource interface {
	InsertBeatsSheet(ctx context.Context, data dao.InsertBeatsSheetData) (*dao.BeatsSheetEntity, error)
	SelectStoryPlan(ctx context.Context, data uuid.UUID) (*dao.StoryPlanEntity, error)
	SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error)
}

func NewCreateBeatsSheetServiceSource(
	insertBeatSheetDAO *dao.InsertBeatsSheetRepository,
	selectStoryPlanDAO *dao.SelectStoryPlanRepository,
	selectLoglineDAO *dao.SelectLoglineRepository,
) CreateBeatsSheetSource {
	return &struct {
		*dao.InsertBeatsSheetRepository
		*dao.SelectStoryPlanRepository
		*dao.SelectLoglineRepository
	}{
		InsertBeatsSheetRepository: insertBeatSheetDAO,
		SelectStoryPlanRepository:  selectStoryPlanDAO,
		SelectLoglineRepository:    selectLoglineDAO,
	}
}

type CreateBeatsSheetRequest struct {
	LoglineID   uuid.UUID
	UserID      uuid.UUID
	StoryPlanID uuid.UUID
	Content     []models.Beat
	Lang        models.Lang
}

type CreateBeatsSheetService struct {
	source CreateBeatsSheetSource
}

func NewCreateBeatsSheetService(source CreateBeatsSheetSource) *CreateBeatsSheetService {
	return &CreateBeatsSheetService{source: source}
}

func (service *CreateBeatsSheetService) CreateBeatsSheet(
	ctx context.Context, request CreateBeatsSheetRequest,
) (*models.BeatsSheet, error) {
	span := sentry.StartSpan(ctx, "CreateBeatsSheetService.CreateBeatsSheet")
	defer span.Finish()

	span.SetData("request.loglineID", request.LoglineID)
	span.SetData("request.storyPlanID", request.StoryPlanID)
	span.SetData("request.lang", request.Lang)
	span.SetData("request.userID", request.UserID)

	_, err := service.source.SelectLogline(span.Context(), dao.SelectLoglineData{
		ID:     request.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		span.SetData("dao.selectLogline.err", err.Error())

		return nil, NewErrCreateBeatsSheetService(fmt.Errorf("check logline: %w", err))
	}

	storyPlan, err := service.source.SelectStoryPlan(span.Context(), request.StoryPlanID)
	if err != nil {
		span.SetData("dao.selectStoryPlan.err", err.Error())

		return nil, NewErrCreateBeatsSheetService(fmt.Errorf("check story plan: %w", err))
	}

	// Ensure story plan matches the beats sheet.
	err = lib.CheckStoryPlan(request.Content, storyPlan.Beats)
	if err != nil {
		span.SetData("lib.checkStoryPlan.err", err.Error())

		return nil, NewErrCreateBeatsSheetService(fmt.Errorf("check story plan: %w", err))
	}

	resp, err := service.source.InsertBeatsSheet(span.Context(), dao.InsertBeatsSheetData{
		Sheet: models.BeatsSheet{
			ID:          uuid.New(),
			LoglineID:   request.LoglineID,
			StoryPlanID: request.StoryPlanID,
			Content:     request.Content,
			Lang:        request.Lang,
			CreatedAt:   time.Now(),
		},
	})
	if err != nil {
		span.SetData("dao.insertBeatsSheet.err", err.Error())

		return nil, NewErrCreateBeatsSheetService(fmt.Errorf("insert beats sheet: %w", err))
	}

	span.SetData("dao.insertBeatsSheet.id", resp.ID)

	return &models.BeatsSheet{
		ID:          resp.ID,
		LoglineID:   resp.LoglineID,
		StoryPlanID: resp.StoryPlanID,
		Content:     resp.Content,
		Lang:        resp.Lang,
		CreatedAt:   resp.CreatedAt,
	}, nil
}
