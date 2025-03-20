package models

import (
	"time"

	"github.com/google/uuid"
)

// BeatsSheet represents a story idea as a list of beats, i.e. single phrases that summarize each beat of the
// story.
type BeatsSheet struct {
	ID          uuid.UUID `json:"id"`
	LoglineID   uuid.UUID `json:"loglineID"`
	StoryPlanID uuid.UUID `json:"storyPlanID"`

	// The beats (in order) that make up the story.
	Content   []Beat    `bun:"content,type:jsonb" json:"content"`
	CreatedAt time.Time `bun:"created_at"         json:"createdAt"`
}

type BeatsSheetPreview struct {
	ID uuid.UUID `json:"id"`

	CreatedAt time.Time `bun:"created_at" json:"createdAt"`
}

// Beat represents a single phrase that summarizes a beat of the story.
type Beat struct {
	// Key links the current Beat to a beat in the StoryPlan.
	Key string `json:"key"`

	// The human-readable title of the beat.
	Title string `json:"title"`

	// A summary of the beat.
	Content string `json:"content"`
}
