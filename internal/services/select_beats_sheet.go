package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

type SelectBeatsSheetSource interface {
	SelectBeatsSheet(ctx context.Context, data uuid.UUID) (*dao.BeatsSheetEntity, error)
	SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error)
}

func NewSelectBeatsSheetServiceSource(
	selectBeatsSheetDAO *dao.SelectBeatsSheetRepository,
	selectLoglineDAO *dao.SelectLoglineRepository,
) SelectBeatsSheetSource {
	return &struct {
		*dao.SelectBeatsSheetRepository
		*dao.SelectLoglineRepository
	}{
		SelectBeatsSheetRepository: selectBeatsSheetDAO,
		SelectLoglineRepository:    selectLoglineDAO,
	}
}

type SelectBeatsSheetRequest struct {
	BeatsSheetID uuid.UUID
	UserID       uuid.UUID
}

type SelectBeatsSheetService struct {
	source SelectBeatsSheetSource
}

func NewSelectBeatsSheetService(source SelectBeatsSheetSource) *SelectBeatsSheetService {
	return &SelectBeatsSheetService{source: source}
}

func (service *SelectBeatsSheetService) SelectBeatsSheet(
	ctx context.Context, request SelectBeatsSheetRequest,
) (*models.BeatsSheet, error) {
	ctx, span := otel.Tracer().Start(ctx, "service.SelectBeatsSheet")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.beatsSheetID", request.BeatsSheetID.String()),
		attribute.String("request.userID", request.UserID.String()),
	)

	data, err := service.source.SelectBeatsSheet(ctx, request.BeatsSheetID)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("select beats sheet: %w", err))
	}

	// Make sure the selected beats sheet is linked to a logline that belongs to the user.
	_, err = service.source.SelectLogline(ctx, dao.SelectLoglineData{
		ID:     data.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("check logline: %w", err))
	}

	return otel.ReportSuccess(span, &models.BeatsSheet{
		ID:          data.ID,
		LoglineID:   data.LoglineID,
		StoryPlanID: data.StoryPlanID,
		Content:     data.Content,
		Lang:        data.Lang,
		CreatedAt:   data.CreatedAt,
	}), nil
}
