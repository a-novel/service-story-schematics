package services

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/models"
)

type GenerateLoglinesSource interface {
	GenerateLoglines(ctx context.Context, request daoai.GenerateLoglinesRequest) ([]models.LoglineIdea, error)
}

type GenerateLoglinesRequest struct {
	Count  int
	Theme  string
	UserID uuid.UUID
	Lang   models.Lang
}

type GenerateLoglinesService struct {
	source GenerateLoglinesSource
}

func NewGenerateLoglinesService(source GenerateLoglinesSource) *GenerateLoglinesService {
	return &GenerateLoglinesService{source: source}
}

func (service *GenerateLoglinesService) GenerateLoglines(
	ctx context.Context, request GenerateLoglinesRequest,
) ([]models.LoglineIdea, error) {
	ctx, span := otel.Tracer().Start(ctx, "service.GenerateLoglines")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.theme", request.Theme),
		attribute.Int("request.count", request.Count),
		attribute.String("request.userID", request.UserID.String()),
		attribute.String("request.lang", request.Lang.String()),
	)

	resp, err := service.source.GenerateLoglines(ctx, daoai.GenerateLoglinesRequest{
		Count:  request.Count,
		Theme:  request.Theme,
		UserID: request.UserID.String(),
		Lang:   request.Lang,
	})
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	return otel.ReportSuccess(span, resp), nil
}
