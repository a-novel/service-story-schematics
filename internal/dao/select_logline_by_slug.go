package dao

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"
	"github.com/a-novel/golib/postgres"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed select_logline_by_slug.sql
var selectLoglineBySlugQuery string

type SelectLoglineBySlugData struct {
	Slug   models.Slug
	UserID uuid.UUID
}

type SelectLoglineBySlugRepository struct{}

func NewSelectLoglineBySlugRepository() *SelectLoglineBySlugRepository {
	return &SelectLoglineBySlugRepository{}
}

func (repository *SelectLoglineBySlugRepository) SelectLoglineBySlug(
	ctx context.Context, data SelectLoglineBySlugData,
) (*LoglineEntity, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.SelectLoglineBySlug")
	defer span.End()

	span.SetAttributes(
		attribute.String("logline.slug", data.Slug.String()),
		attribute.String("logline.userID", data.UserID.String()),
	)

	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get postgres client: %w", err))
	}

	entity := &LoglineEntity{}

	err = tx.NewRaw(selectLoglineBySlugQuery, data.Slug, data.UserID).Scan(ctx, entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, otel.ReportError(span, ErrLoglineNotFound)
		}

		return nil, otel.ReportError(span, fmt.Errorf("select logline: %w", err))
	}

	return otel.ReportSuccess(span, entity), nil
}
