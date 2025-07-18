package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/models"
)

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
	ctx, span := otel.Tracer().Start(ctx, "service.CreateBeatsSheet")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.loglineID", request.LoglineID.String()),
		attribute.String("request.storyPlanID", request.StoryPlanID.String()),
		attribute.String("request.lang", request.Lang.String()),
		attribute.String("request.userID", request.UserID.String()),
	)

	_, err := service.source.SelectLogline(ctx, dao.SelectLoglineData{
		ID:     request.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("check logline: %w", err))
	}

	storyPlan, err := service.source.SelectStoryPlan(ctx, request.StoryPlanID)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("check story plan: %w", err))
	}

	// Ensure story plan matches the beats sheet.
	err = lib.CheckStoryPlan(request.Content, storyPlan.Beats)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("check story plan: %w", err))
	}

	resp, err := service.source.InsertBeatsSheet(ctx, dao.InsertBeatsSheetData{
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
		return nil, otel.ReportError(span, fmt.Errorf("insert beats sheet: %w", err))
	}

	span.SetAttributes(attribute.String("dao.insertBeatsSheet.id", resp.ID.String()))

	return otel.ReportSuccess(span, &models.BeatsSheet{
		ID:          resp.ID,
		LoglineID:   resp.LoglineID,
		StoryPlanID: resp.StoryPlanID,
		Content:     resp.Content,
		Lang:        resp.Lang,
		CreatedAt:   resp.CreatedAt,
	}), nil
}
