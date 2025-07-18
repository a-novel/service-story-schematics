package services

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/models"
)

type RegenerateBeatsSource interface {
	RegenerateBeats(ctx context.Context, request daoai.RegenerateBeatsRequest) ([]models.Beat, error)
	SelectBeatsSheet(ctx context.Context, data uuid.UUID) (*dao.BeatsSheetEntity, error)
	SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error)
	SelectStoryPlan(ctx context.Context, data uuid.UUID) (*dao.StoryPlanEntity, error)
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

type RegenerateBeatsRequest struct {
	BeatsSheetID   uuid.UUID
	UserID         uuid.UUID
	RegenerateKeys []string
}

type RegenerateBeatsService struct {
	source RegenerateBeatsSource
}

func NewRegenerateBeatsService(source RegenerateBeatsSource) *RegenerateBeatsService {
	return &RegenerateBeatsService{source: source}
}

func (service *RegenerateBeatsService) RegenerateBeats(
	ctx context.Context, request RegenerateBeatsRequest,
) ([]models.Beat, error) {
	ctx, span := otel.Tracer().Start(ctx, "service.RegenerateBeats")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.beatsSheetID", request.BeatsSheetID.String()),
		attribute.String("request.userID", request.UserID.String()),
		attribute.StringSlice("request.regenerateKeys", request.RegenerateKeys),
	)

	beatsSheet, err := service.source.SelectBeatsSheet(ctx, request.BeatsSheetID)
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	// Make sure the selected beats sheet is linked to a logline that belongs to the user.
	logline, err := service.source.SelectLogline(ctx, dao.SelectLoglineData{
		ID:     beatsSheet.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	storyPlan, err := service.source.SelectStoryPlan(ctx, beatsSheet.StoryPlanID)
	if err != nil {
		return nil, otel.ReportError(span, err)
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
		return nil, otel.ReportError(span, err)
	}

	return otel.ReportSuccess(span, regenerated), nil
}
