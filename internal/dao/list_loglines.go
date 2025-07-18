package dao

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"
	"github.com/a-novel/golib/postgres"
)

//go:embed list_loglines.sql
var listLoglinesQuery string

type ListLoglinesData struct {
	UserID uuid.UUID
	Limit  int
	Offset int
}

type ListLoglinesRepository struct{}

func NewListLoglinesRepository() *ListLoglinesRepository {
	return &ListLoglinesRepository{}
}

func (repository *ListLoglinesRepository) ListLoglines(
	ctx context.Context, data ListLoglinesData,
) ([]*LoglinePreviewEntity, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.ListLoglines")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.id", data.UserID.String()),
		attribute.Int("limit", data.Limit),
		attribute.Int("offset", data.Offset),
	)

	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get postgres client: %w", err))
	}

	entities := make([]*LoglinePreviewEntity, 0)

	err = tx.NewRaw(listLoglinesQuery, data.UserID, bun.NullZero(data.Limit), data.Offset).Scan(ctx, &entities)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("list loglines: %w", err))
	}

	return otel.ReportSuccess(span, entities), nil
}
