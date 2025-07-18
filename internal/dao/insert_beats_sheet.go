package dao

import (
	"context"
	_ "embed"
	"fmt"

	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"
	"github.com/a-novel/golib/postgres"

	"github.com/a-novel/service-story-schematics/models"
)

//go:embed insert_beats_sheet.sql
var insertBeatsSheetQuery string

type InsertBeatsSheetData struct {
	Sheet models.BeatsSheet
}

type InsertBeatsSheetRepository struct{}

func NewInsertBeatsSheetRepository() *InsertBeatsSheetRepository {
	return &InsertBeatsSheetRepository{}
}

func (repository *InsertBeatsSheetRepository) InsertBeatsSheet(
	ctx context.Context, data InsertBeatsSheetData,
) (*BeatsSheetEntity, error) {
	ctx, span := otel.Tracer().Start(ctx, "dao.InsertBeatsSheet")
	defer span.End()

	span.SetAttributes(
		attribute.String("sheet.id", data.Sheet.ID.String()),
		attribute.String("sheet.loglineID", data.Sheet.LoglineID.String()),
		attribute.String("sheet.storyPlanID", data.Sheet.StoryPlanID.String()),
		attribute.String("sheet.lang", data.Sheet.Lang.String()),
	)

	tx, err := postgres.GetContext(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get postgres client: %w", err))
	}

	entity := &BeatsSheetEntity{}

	err = tx.
		NewRaw(
			insertBeatsSheetQuery,
			data.Sheet.ID,
			data.Sheet.LoglineID,
			data.Sheet.StoryPlanID,
			data.Sheet.Content,
			data.Sheet.Lang,
			data.Sheet.CreatedAt,
		).
		Scan(ctx, entity)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("insert sheet: %w", err))
	}

	return otel.ReportSuccess(span, entity), nil
}
