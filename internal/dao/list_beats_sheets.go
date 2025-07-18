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

//go:embed list_beats_sheets.sql
var listBeatsSheetsQuery string

type ListBeatsSheetsData struct {
	LoglineID uuid.UUID
	Limit     int
	Offset    int
}

type ListBeatsSheetsRepository struct{}

func NewListBeatsSheetsRepository() *ListBeatsSheetsRepository {
	return &ListBeatsSheetsRepository{}
}

func (repository *ListBeatsSheetsRepository) ListBeatsSheets(
	ctx context.Context, data ListBeatsSheetsData,
) ([]*BeatsSheetPreviewEntity, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.ListBeatsSheets")
	defer span.End()

	span.SetAttributes(
		attribute.String("logline.id", data.LoglineID.String()),
		attribute.Int("limit", data.Limit),
		attribute.Int("offset", data.Offset),
	)

	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get postgres client: %w", err))
	}

	entities := make([]*BeatsSheetPreviewEntity, 0)

	err = tx.NewRaw(listBeatsSheetsQuery, data.LoglineID, bun.NullZero(data.Limit), data.Offset).Scan(ctx, &entities)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("list beats sheet: %w", err))
	}

	return otel.ReportSuccess(span, entities), nil
}
