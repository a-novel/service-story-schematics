package dao

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"
	"github.com/a-novel/golib/postgres"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed insert_logline.sql
var insertLoglineQuery string

type InsertLoglineData struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Slug   models.Slug

	Name    string
	Content string
	Lang    models.Lang

	Now time.Time
}

type InsertLoglineRepository struct{}

func NewInsertLoglineRepository() *InsertLoglineRepository {
	return &InsertLoglineRepository{}
}

func (repository *InsertLoglineRepository) InsertLogline(
	ctx context.Context, data InsertLoglineData,
) (*LoglineEntity, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.InsertLogline")
	defer span.End()

	span.SetAttributes(
		attribute.String("logline.id", data.ID.String()),
		attribute.String("logline.userID", data.UserID.String()),
		attribute.String("logline.slug", data.Slug.String()),
		attribute.String("logline.name", data.Name),
		attribute.String("logline.lang", data.Lang.String()),
	)

	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get postgres client: %w", err))
	}

	entity := &LoglineEntity{}

	err = tx.
		NewRaw(
			insertLoglineQuery,
			data.ID,
			data.UserID,
			data.Slug,
			data.Name,
			data.Content,
			data.Lang,
			data.Now,
		).
		Scan(ctx, entity)
	if err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) && pgErr.Field('C') == "23505" {
			return nil, otel.ReportError(span, errors.Join(err, ErrLoglineAlreadyExists))
		}

		return nil, otel.ReportError(span, fmt.Errorf("insert logline: %w", err))
	}

	return otel.ReportSuccess(span, entity), nil
}
