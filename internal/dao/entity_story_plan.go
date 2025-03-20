package dao

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/a-novel/story-schematics/models"
)

var (
	ErrStoryPlanNotFound      = errors.New("story plan not found")
	ErrStoryPlanAlreadyExists = errors.New("story plan already exists")
)

type StoryPlanEntity struct {
	bun.BaseModel `bun:"table:story_plans,alias:sp,select:story_plans_active_view"`

	ID   uuid.UUID   `bun:"id,pk,type:uuid"`
	Slug models.Slug `bun:"slug"`

	Name        string `bun:"name"`
	Description string `bun:"description"`

	Beats []models.BeatDefinition `bun:"beats,type:jsonb"`

	CreatedAt time.Time `bun:"created_at"`
}

type StoryPlanPreviewEntity struct {
	bun.BaseModel `bun:"table:story_plans,alias:sp,select:story_plans_active_view"`

	ID   uuid.UUID   `bun:"id,pk,type:uuid"`
	Slug models.Slug `bun:"slug"`

	Name        string `bun:"name"`
	Description string `bun:"description"`

	CreatedAt time.Time `bun:"created_at"`
}
