package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"

	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrInsertBeatsSheetRepository = errors.New("InsertBeatsSheetRepository.InsertBeatsSheet")

func NewErrInsertBeatsSheetRepository(err error) error {
	return errors.Join(err, ErrInsertBeatsSheetRepository)
}

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
	span := sentry.StartSpan(ctx, "InsertBeatsSheetRepository.InsertBeatsSheet")
	defer span.Finish()

	span.SetData("sheet.id", data.Sheet.ID.String())
	span.SetData("sheet.logline_id", data.Sheet.LoglineID.String())
	span.SetData("sheet.story_plan_id", data.Sheet.StoryPlanID.String())
	span.SetData("sheet.lang", data.Sheet.Lang)

	tx, err := lib.PostgresContext(span.Context())
	if err != nil {
		span.SetData("postgres.context.error", err.Error())

		return nil, NewErrInsertBeatsSheetRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entity := &BeatsSheetEntity{
		ID:          data.Sheet.ID,
		LoglineID:   data.Sheet.LoglineID,
		StoryPlanID: data.Sheet.StoryPlanID,
		Content:     make([]models.Beat, len(data.Sheet.Content)),
		Lang:        data.Sheet.Lang,
		CreatedAt:   data.Sheet.CreatedAt,
	}

	for i, beat := range data.Sheet.Content {
		entity.Content[i] = models.Beat{
			Key:     beat.Key,
			Title:   beat.Title,
			Content: beat.Content,
		}
	}

	_, err = tx.NewInsert().Model(entity).Returning("*").Exec(span.Context())
	if err != nil {
		span.SetData("insert.error", err.Error())

		return nil, NewErrInsertBeatsSheetRepository(fmt.Errorf("insert sheet: %w", err))
	}

	return entity, nil
}
