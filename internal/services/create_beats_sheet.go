package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/a-novel-kit/context"

	"github.com/a-novel/story-schematics/internal/dao"
	"github.com/a-novel/story-schematics/internal/lib"
	"github.com/a-novel/story-schematics/models"
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

type CreateBeatsSheetRequest struct {
	LoglineID   uuid.UUID
	UserID      uuid.UUID
	StoryPlanID uuid.UUID
	Content     []models.Beat
}

type CreateBeatsSheetService struct {
	source CreateBeatsSheetSource
}

func (service *CreateBeatsSheetService) CreateBeatsSheet(
	ctx context.Context, request CreateBeatsSheetRequest,
) (*models.BeatsSheet, error) {
	_, err := service.source.SelectLogline(ctx, dao.SelectLoglineData{
		ID:     request.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, NewErrCreateBeatsSheetService(fmt.Errorf("check logline: %w", err))
	}

	storyPlan, err := service.source.SelectStoryPlan(ctx, request.StoryPlanID)
	if err != nil {
		return nil, NewErrCreateBeatsSheetService(fmt.Errorf("check story plan: %w", err))
	}

	// Ensure story plan matches the beats sheet.
	if err = lib.CheckStoryPlan(request.Content, storyPlan.Beats); err != nil {
		return nil, NewErrCreateBeatsSheetService(fmt.Errorf("check story plan: %w", err))
	}

	resp, err := service.source.InsertBeatsSheet(ctx, dao.InsertBeatsSheetData{
		Sheet: models.BeatsSheet{
			ID:          uuid.New(),
			LoglineID:   request.LoglineID,
			StoryPlanID: request.StoryPlanID,
			Content:     request.Content,
			CreatedAt:   time.Now(),
		},
	})
	if err != nil {
		return nil, NewErrCreateBeatsSheetService(fmt.Errorf("insert beats sheet: %w", err))
	}

	return &models.BeatsSheet{
		ID:          resp.ID,
		LoglineID:   resp.LoglineID,
		StoryPlanID: resp.StoryPlanID,
		Content:     resp.Content,
		CreatedAt:   resp.CreatedAt,
	}, nil
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

func NewCreateBeatsSheetService(source CreateBeatsSheetSource) *CreateBeatsSheetService {
	return &CreateBeatsSheetService{source: source}
}
