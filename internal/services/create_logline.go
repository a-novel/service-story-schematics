package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

type CreateLoglineSource interface {
	InsertLogline(ctx context.Context, data dao.InsertLoglineData) (*dao.LoglineEntity, error)
	SelectSlugIteration(ctx context.Context, data dao.SelectSlugIterationData) (models.Slug, int, error)
}

func NewCreateLoglineServiceSource(
	insertLoglineDAO *dao.InsertLoglineRepository,
	selectSlugIterationDAO *dao.SelectSlugIterationRepository,
) CreateLoglineSource {
	return &struct {
		*dao.InsertLoglineRepository
		*dao.SelectSlugIterationRepository
	}{
		InsertLoglineRepository:       insertLoglineDAO,
		SelectSlugIterationRepository: selectSlugIterationDAO,
	}
}

type CreateLoglineRequest struct {
	UserID  uuid.UUID
	Slug    models.Slug
	Name    string
	Content string
	Lang    models.Lang
}

type CreateLoglineService struct {
	source CreateLoglineSource
}

func NewCreateLoglineService(source CreateLoglineSource) *CreateLoglineService {
	return &CreateLoglineService{source: source}
}

func (service *CreateLoglineService) CreateLogline(
	ctx context.Context, request CreateLoglineRequest,
) (*models.Logline, error) {
	ctx, span := otel.Tracer().Start(ctx, "service.CreateLogline")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.userID", request.UserID.String()),
		attribute.String("request.slug", request.Slug.String()),
		attribute.String("request.name", request.Name),
		attribute.String("request.lang", request.Lang.String()),
		attribute.Bool("slug.taken", false),
	)

	data := dao.InsertLoglineData{
		ID:      uuid.New(),
		UserID:  request.UserID,
		Slug:    request.Slug,
		Name:    request.Name,
		Content: request.Content,
		Lang:    request.Lang,
		Now:     time.Now(),
	}

	resp, err := service.source.InsertLogline(ctx, data)

	// If slug is taken, try to modify it by appending a version number.
	if errors.Is(err, dao.ErrLoglineAlreadyExists) {
		span.SetAttributes(attribute.Bool("slug.taken", true))

		data.Slug, _, err = service.source.SelectSlugIteration(ctx, dao.SelectSlugIterationData{
			Slug:   data.Slug,
			Target: dao.SlugIterationTargetLogline,
			Args:   []any{data.UserID},
		})
		if err != nil {
			return nil, otel.ReportError(span, fmt.Errorf("check slug uniqueness: %w", err))
		}

		resp, err = service.source.InsertLogline(ctx, data)
	}

	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("insert logline: %w", err))
	}

	span.SetAttributes(attribute.String("dao.insertLogline.id", resp.ID.String()))

	return otel.ReportSuccess(span, &models.Logline{
		ID:        resp.ID,
		UserID:    resp.UserID,
		Slug:      resp.Slug,
		Name:      resp.Name,
		Content:   resp.Content,
		Lang:      resp.Lang,
		CreatedAt: resp.CreatedAt,
	}), nil
}
