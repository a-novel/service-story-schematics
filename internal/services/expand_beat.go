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
	ctx, span := otel.Tracer().Start(ctx, "service.ExpandBeat")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.beatsSheetID", request.BeatsSheetID.String()),
		attribute.String("request.targetKey", request.TargetKey),
		attribute.String("request.userID", request.UserID.String()),
	)

	beatsSheet, err := service.source.SelectBeatsSheet(ctx, request.BeatsSheetID)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("select beats sheet: %w", err))
	}

	// Make sure the selected beats sheet is linked to a logline that belongs to the user.
	logline, err := service.source.SelectLogline(ctx, dao.SelectLoglineData{
		ID:     beatsSheet.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("select logline: %w", err))
	}

	storyPlan, err := service.source.SelectStoryPlan(ctx, beatsSheet.StoryPlanID)
	if err != nil {
		return nil, otel.ReportError(span, err)
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
			Lang:        beatsSheet.Lang,
		},
		Lang:      beatsSheet.Lang,
		TargetKey: request.TargetKey,
		UserID:    request.UserID.String(),
	})
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	return otel.ReportSuccess(span, expanded), nil
}
