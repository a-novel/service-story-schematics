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
	storyplanmodel "github.com/a-novel/service-story-schematics/models/story_plan"
)

type GenerateBeatsSheetSource interface {
	GenerateBeatsSheet(ctx context.Context, request daoai.GenerateBeatsSheetRequest) ([]models.Beat, error)
	SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error)
	SelectStoryPlan(ctx context.Context, request SelectStoryPlanRequest) (*storyplanmodel.Plan, error)
}

func NewGenerateBeatsSheetServiceSource(
	generateDAO *daoai.GenerateBeatsSheetRepository,
	selectLoglineDAO *dao.SelectLoglineRepository,
	selectStoryPlan *SelectStoryPlanService,
) GenerateBeatsSheetSource {
	return &struct {
		*daoai.GenerateBeatsSheetRepository
		*dao.SelectLoglineRepository
		*SelectStoryPlanService
	}{
		GenerateBeatsSheetRepository: generateDAO,
		SelectLoglineRepository:      selectLoglineDAO,
		SelectStoryPlanService:       selectStoryPlan,
	}
}

type GenerateBeatsSheetRequest struct {
	LoglineID uuid.UUID
	UserID    uuid.UUID
	Lang      models.Lang
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

	storyPlan, err := service.source.SelectStoryPlan(ctx, SelectStoryPlanRequest{
		Lang: request.Lang,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get story plan: %w", err))
	}

	resp, err := service.source.GenerateBeatsSheet(ctx, daoai.GenerateBeatsSheetRequest{
		Logline: logline.Name + "\n\n" + logline.Content,
		Plan:    storyPlan,
		Lang:    request.Lang,
		UserID:  request.UserID.String(),
	})
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	return otel.ReportSuccess(span, resp), nil
}
