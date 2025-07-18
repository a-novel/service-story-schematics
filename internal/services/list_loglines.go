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

type ListLoglinesSource interface {
	ListLoglines(ctx context.Context, data dao.ListLoglinesData) ([]*dao.LoglinePreviewEntity, error)
}

type ListLoglinesRequest struct {
	UserID uuid.UUID
	Limit  int
	Offset int
}

type ListLoglinesService struct {
	source ListLoglinesSource
}

func NewListLoglinesService(source ListLoglinesSource) *ListLoglinesService {
	return &ListLoglinesService{source: source}
}

func (service *ListLoglinesService) ListLoglines(
	ctx context.Context, request ListLoglinesRequest,
) ([]*models.LoglinePreview, error) {
	ctx, span := otel.Tracer().Start(ctx, "service.ListLoglines")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.userID", request.UserID.String()),
		attribute.Int("request.limit", request.Limit),
		attribute.Int("request.offset", request.Offset),
	)

	resp, err := service.source.ListLoglines(ctx, dao.ListLoglinesData{
		UserID: request.UserID,
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	span.SetAttributes(attribute.Int("dao.listLoglines.count", len(resp)))

	output := lo.Map(resp, func(item *dao.LoglinePreviewEntity, _ int) *models.LoglinePreview {
		return &models.LoglinePreview{
			Slug:      item.Slug,
			Name:      item.Name,
			Content:   item.Content,
			Lang:      item.Lang,
			CreatedAt: item.CreatedAt,
		}
	})

	return otel.ReportSuccess(span, output), nil
}
