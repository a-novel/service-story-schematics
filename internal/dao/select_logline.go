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
)

//go:embed select_logline.sql
var selectLoglineQuery string

type SelectLoglineData struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

type SelectLoglineRepository struct{}

func NewSelectLoglineRepository() *SelectLoglineRepository {
	return &SelectLoglineRepository{}
}

func (repository *SelectLoglineRepository) SelectLogline(
	ctx context.Context, data SelectLoglineData,
) (*LoglineEntity, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.SelectLogline")
	defer span.End()

	span.SetAttributes(
		attribute.String("logline.id", data.ID.String()),
		attribute.String("logline.userID", data.UserID.String()),
	)

	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get postgres client: %w", err))
	}

	entity := &LoglineEntity{}

	err = tx.NewRaw(selectLoglineQuery, data.ID, data.UserID).Scan(ctx, entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, otel.ReportError(span, ErrLoglineNotFound)
		}

		return nil, otel.ReportError(span, fmt.Errorf("select logline: %w", err))
	}

	return otel.ReportSuccess(span, entity), nil
}
