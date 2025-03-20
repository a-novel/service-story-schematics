package models

import (
	"time"

	"github.com/google/uuid"
)

type Logline struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"userID"`
	Slug   Slug      `json:"slug"`

	Name    string `json:"name"`
	Content string `json:"content"`

	CreatedAt time.Time `json:"createdAt"`
}

type LoglinePreview struct {
	Slug Slug `json:"slug"`

	Name    string `json:"name"`
	Content string `json:"content"`

	CreatedAt time.Time `json:"createdAt"`
}

type LoglineIdea struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}
