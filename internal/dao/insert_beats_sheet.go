package dao

import (
	"context"
	"errors"
	"fmt"

	pgctx "github.com/a-novel-kit/context/pgbun"

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

func (repository *InsertBeatsSheetRepository) InsertBeatsSheet(
	ctx context.Context, data InsertBeatsSheetData,
) (*BeatsSheetEntity, error) {
	tx, err := pgctx.Context(ctx)
	if err != nil {
		return nil, NewErrInsertBeatsSheetRepository(fmt.Errorf("get postgres client: %w", err))
	}

	entity := &BeatsSheetEntity{
		ID:          data.Sheet.ID,
		LoglineID:   data.Sheet.LoglineID,
		StoryPlanID: data.Sheet.StoryPlanID,
		Content:     make([]models.Beat, len(data.Sheet.Content)),
		CreatedAt:   data.Sheet.CreatedAt,
	}

	for i, beat := range data.Sheet.Content {
		entity.Content[i] = models.Beat{
			Key:     beat.Key,
			Title:   beat.Title,
			Content: beat.Content,
		}
	}

	_, err = tx.NewInsert().Model(entity).Returning("*").Exec(ctx)
	if err != nil {
		return nil, NewErrInsertBeatsSheetRepository(fmt.Errorf("insert sheet: %w", err))
	}

	return entity, nil
}

func NewInsertBeatsSheetRepository() *InsertBeatsSheetRepository {
	return &InsertBeatsSheetRepository{}
}
