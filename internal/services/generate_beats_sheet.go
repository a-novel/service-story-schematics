package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/models"
)

type GenerateBeatsSheetSource interface {
	GenerateBeatsSheet(ctx context.Context, request daoai.GenerateBeatsSheetRequest) ([]models.Beat, error)
	SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error)
	SelectStoryPlan(ctx context.Context, data uuid.UUID) (*dao.StoryPlanEntity, error)
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

type GenerateBeatsSheetRequest struct {
	LoglineID   uuid.UUID
	StoryPlanID uuid.UUID
	UserID      uuid.UUID
	Lang        models.Lang
}

type GenerateBeatsSheetService struct {
	source GenerateBeatsSheetSource
}

func NewGenerateBeatsSheetService(source GenerateBeatsSheetSource) *GenerateBeatsSheetService {
	return &GenerateBeatsSheetService{source: source}
}

func (service *GenerateBeatsSheetService) GenerateBeatsSheet(
	ctx context.Context, request GenerateBeatsSheetRequest,
) ([]models.Beat, error) {
	ctx, span := otel.Tracer().Start(ctx, "service.GenerateBeatsSheet")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.loglineID", request.LoglineID.String()),
		attribute.String("request.storyPlanID", request.StoryPlanID.String()),
		attribute.String("request.lang", request.Lang.String()),
		attribute.String("request.userID", request.UserID.String()),
	)

	logline, err := service.source.SelectLogline(ctx, dao.SelectLoglineData{
		ID:     request.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get logline: %w", err))
	}

	storyPlan, err := service.source.SelectStoryPlan(ctx, request.StoryPlanID)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get story plan: %w", err))
	}

	resp, err := service.source.GenerateBeatsSheet(ctx, daoai.GenerateBeatsSheetRequest{
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
		return nil, otel.ReportError(span, err)
	}

	return otel.ReportSuccess(span, resp), nil
}
