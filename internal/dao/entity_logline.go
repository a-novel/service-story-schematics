package dao

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/a-novel/story-schematics/models"
)

var (
	ErrLoglineNotFound      = errors.New("logline not found")
	ErrLoglineAlreadyExists = errors.New("logline already exists")
)

type LoglineEntity struct {
	bun.BaseModel `bun:"table:loglines"`

	ID     uuid.UUID   `bun:"id,pk,type:uuid"`
	UserID uuid.UUID   `bun:"user_id,type:uuid"`
	Slug   models.Slug `bun:"slug"`

	Name    string `bun:"name"`
	Content string `bun:"content"`

	CreatedAt time.Time `bun:"created_at"`
}

type LoglinePreviewEntity struct {
	bun.BaseModel `bun:"table:loglines"`

	Slug models.Slug `bun:"slug"`

	Name    string `bun:"name"`
	Content string `bun:"content"`

	CreatedAt time.Time `bun:"created_at"`
}
