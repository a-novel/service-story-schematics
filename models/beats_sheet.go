package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BeatsSheet represents a story idea as a list of beats, i.e. single phrases that summarize each beat of the
// story.
type BeatsSheet struct {
	ID        uuid.UUID `json:"id"`
	LoglineID uuid.UUID `json:"loglineID"`

	// The beats (in order) that make up the story.
	Content []Beat `bun:"content,type:jsonb" json:"content"`
	Lang    Lang   `bun:"lang"               json:"lang"`

	CreatedAt time.Time `bun:"created_at" json:"createdAt"`
}

type BeatsSheetPreview struct {
	ID   uuid.UUID `json:"id"`
	Lang Lang      `json:"lang"`

	CreatedAt time.Time `bun:"created_at" json:"createdAt"`
}

// Beat represents a single phrase that summarizes a beat of the story.
type Beat struct {
	// Key links the current Beat to a beat in the StoryPlan.
	Key string `json:"key" jsonschema_description:"The jsonKey of the given beat" yaml:"key"`

	// The human-readable title of the beat.
	Title string `json:"title" jsonschema_description:"A short title for the beat" yaml:"title"`

	// A summary of the beat.
	Content string `json:"content" jsonschema_description:"A summary of the scenes in the beat" yaml:"content"`
}

func (beat Beat) String() string {
	return fmt.Sprintf("%s (%s)\n%s", beat.Title, beat.Key, beat.Content)
}
