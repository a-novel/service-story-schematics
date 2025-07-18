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

//go:embed select_beats_sheet.sql
var selectBeatsSheetQuery string

type SelectBeatsSheetRepository struct{}

func NewSelectBeatsSheetRepository() *SelectBeatsSheetRepository {
	return &SelectBeatsSheetRepository{}
}

func (repository *SelectBeatsSheetRepository) SelectBeatsSheet(
	ctx context.Context, data uuid.UUID,
) (*BeatsSheetEntity, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.SelectBeatsSheet")
	defer span.End()

	span.SetAttributes(attribute.String("sheet.id", data.String()))

	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get postgres client: %w", err))
	}

	entity := &BeatsSheetEntity{}

	err = tx.NewRaw(selectBeatsSheetQuery, data).Scan(ctx, entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, otel.ReportError(span, ErrBeatsSheetNotFound)
		}

		return nil, otel.ReportError(span, fmt.Errorf("select beats sheet: %w", err))
	}

	return otel.ReportSuccess(span, entity), nil
}
