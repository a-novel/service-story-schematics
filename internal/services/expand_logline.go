package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/models"
)

type ExpandLoglineSource interface {
	ExpandLogline(ctx context.Context, request daoai.ExpandLoglineRequest) (*models.LoglineIdea, error)
}

type ExpandLoglineRequest struct {
	Logline models.LoglineIdea
	UserID  uuid.UUID
}

type ExpandLoglineService struct {
	source ExpandLoglineSource
}

func NewExpandLoglineService(source ExpandLoglineSource) *ExpandLoglineService {
	return &ExpandLoglineService{source: source}
}

func (service *ExpandLoglineService) ExpandLogline(
	ctx context.Context, request ExpandLoglineRequest,
) (*models.LoglineIdea, error) {
	ctx, span := otel.Tracer().Start(ctx, "service.ExpandLogline")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.logline.name", request.Logline.Name),
		attribute.String("request.logline.lang", request.Logline.Lang.String()),
		attribute.String("request.userID", request.UserID.String()),
	)

	resp, err := service.source.ExpandLogline(ctx, daoai.ExpandLoglineRequest{
		Logline: request.Logline.Name + "\n\n" + request.Logline.Content,
		UserID:  request.UserID.String(),
		Lang:    request.Logline.Lang,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("expand logline: %w", err))
	}

	return otel.ReportSuccess(span, resp), nil
}
