package dao

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/a-novel/service-story-schematics/models"
)

var ErrBeatsSheetNotFound = errors.New("beats schematic not found")

type BeatsSheetEntity struct {
	bun.BaseModel `bun:"table:beats_sheets"`

	ID          uuid.UUID `bun:"id,pk,type:uuid"`
	LoglineID   uuid.UUID `bun:"logline_id,type:uuid"`
	StoryPlanID uuid.UUID `bun:"story_plan_id,type:uuid"`

	Content []models.Beat `bun:"content,type:jsonb"`
	Lang    models.Lang   `bun:"lang"`

	CreatedAt time.Time `bun:"created_at"`
}

type BeatsSheetPreviewEntity struct {
	bun.BaseModel `bun:"table:beats_sheets"`

	ID   uuid.UUID   `bun:"id,pk,type:uuid"`
	Lang models.Lang `bun:"lang"`

	CreatedAt time.Time `bun:"created_at"`
}
