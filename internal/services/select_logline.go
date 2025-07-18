package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

type SelectLoglineSource interface {
	SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error)
	SelectLoglineBySlug(ctx context.Context, data dao.SelectLoglineBySlugData) (*dao.LoglineEntity, error)
}

func NewSelectLoglineServiceSource(
	selectDAO *dao.SelectLoglineRepository,
	selectBySlugRepository *dao.SelectLoglineBySlugRepository,
) SelectLoglineSource {
	return &struct {
		*dao.SelectLoglineRepository
		*dao.SelectLoglineBySlugRepository
	}{
		SelectLoglineRepository:       selectDAO,
		SelectLoglineBySlugRepository: selectBySlugRepository,
	}
}

type SelectLoglineRequest struct {
	UserID uuid.UUID
	Slug   *models.Slug
	ID     *uuid.UUID
}

type SelectLoglineService struct {
	source SelectLoglineSource
}

func NewSelectLoglineService(source SelectLoglineSource) *SelectLoglineService {
	return &SelectLoglineService{source: source}
}

func (service *SelectLoglineService) SelectLogline(
	ctx context.Context, request SelectLoglineRequest,
) (*models.Logline, error) {
	ctx, span := otel.Tracer().Start(ctx, "service.SelectLogline")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.userID", request.UserID.String()),
		attribute.String("request.slug", lo.FromPtr(request.Slug).String()),
		attribute.String("request.id", lo.FromPtr(request.ID).String()),
	)

	if request.Slug != nil {
		data, err := service.source.SelectLoglineBySlug(ctx, dao.SelectLoglineBySlugData{
			Slug:   lo.FromPtr(request.Slug),
			UserID: request.UserID,
		})
		if err != nil {
			return nil, otel.ReportError(span, fmt.Errorf("select logline by slug: %w", err))
		}

		return otel.ReportSuccess(span, &models.Logline{
			ID:        data.ID,
			UserID:    data.UserID,
			Slug:      data.Slug,
			Name:      data.Name,
			Content:   data.Content,
			Lang:      data.Lang,
			CreatedAt: data.CreatedAt,
		}), nil
	}

	data, err := service.source.SelectLogline(ctx, dao.SelectLoglineData{
		ID:     lo.FromPtr(request.ID),
		UserID: request.UserID,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("select logline: %w", err))
	}

	return otel.ReportSuccess(span, &models.Logline{
		ID:        data.ID,
		UserID:    data.UserID,
		Slug:      data.Slug,
		Name:      data.Name,
		Content:   data.Content,
		Lang:      data.Lang,
		CreatedAt: data.CreatedAt,
	}), nil
}
