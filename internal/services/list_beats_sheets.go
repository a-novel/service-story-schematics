package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

type ListBeatsSheetsSource interface {
	ListBeatsSheets(ctx context.Context, data dao.ListBeatsSheetsData) ([]*dao.BeatsSheetPreviewEntity, error)
	SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error)
}

func NewListBeatsSheetsServiceSource(
	listBeatsSheetsDAO *dao.ListBeatsSheetsRepository,
	selectLoglineDAO *dao.SelectLoglineRepository,
) ListBeatsSheetsSource {
	return &struct {
		*dao.ListBeatsSheetsRepository
		*dao.SelectLoglineRepository
	}{
		ListBeatsSheetsRepository: listBeatsSheetsDAO,
		SelectLoglineRepository:   selectLoglineDAO,
	}
}

type ListBeatsSheetsRequest struct {
	UserID    uuid.UUID
	LoglineID uuid.UUID
	Limit     int
	Offset    int
}

type ListBeatsSheetsService struct {
	source ListBeatsSheetsSource
}

func NewListBeatsSheetsService(source ListBeatsSheetsSource) *ListBeatsSheetsService {
	return &ListBeatsSheetsService{source: source}
}

func (service *ListBeatsSheetsService) ListBeatsSheets(
	ctx context.Context, request ListBeatsSheetsRequest,
) ([]*models.BeatsSheetPreview, error) {
	ctx, span := otel.Tracer().Start(ctx, "service.ListBeatsSheets")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.userID", request.UserID.String()),
		attribute.String("request.loglineID", request.LoglineID.String()),
		attribute.Int("request.limit", request.Limit),
		attribute.Int("request.offset", request.Offset),
	)

	_, err := service.source.SelectLogline(ctx, dao.SelectLoglineData{
		ID:     request.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	data := dao.ListBeatsSheetsData{
		LoglineID: request.LoglineID,
		Limit:     request.Limit,
		Offset:    request.Offset,
	}

	resp, err := service.source.ListBeatsSheets(ctx, data)
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	span.SetAttributes(attribute.Int("dao.listBeatsSheets.count", len(resp)))

	output := lo.Map(resp, func(item *dao.BeatsSheetPreviewEntity, _ int) *models.BeatsSheetPreview {
		return &models.BeatsSheetPreview{
			ID:        item.ID,
			Lang:      item.Lang,
			CreatedAt: item.CreatedAt,
		}
	})

	return otel.ReportSuccess(span, output), nil
}
