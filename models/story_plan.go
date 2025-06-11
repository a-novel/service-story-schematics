package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// StoryPlan represents the structure of a story, based on a list of predefined beats.
type StoryPlan struct {
	ID uuid.UUID `json:"id"`

	// Slug used to uniquely retrieve the story plan.
	Slug Slug `json:"slug"`

	// Human-readable name of the story plan.
	Name string `json:"name"`
	// Description of the story plan.
	Description string `json:"description"`

	// A list of beats (ordered) that make up the story.
	Beats []BeatDefinition `json:"beats"`

	Lang Lang `json:"lang"`

	CreatedAt time.Time `json:"createdAt"`
}

type StoryPlanPreview struct {
	ID uuid.UUID `json:"id"`

	// Slug used to uniquely retrieve the story plan.
	Slug Slug `json:"slug"`

	// Human-readable name of the story plan.
	Name string `json:"name"`
	// Description of the story plan.
	Description string `json:"description"`

	Lang Lang `json:"lang"`

	CreatedAt time.Time `json:"createdAt"`
}

// BeatDefinition represents a single part of the story plan. It describes the expected content of that part
// of the story.
type BeatDefinition struct {
	// Human-readable name of the beat.
	Name string `json:"name"`
	// Key used to uniquely identify the beat.
	Key string `json:"key"`

	// The important highlights of the beat.
	KeyPoints []string `json:"goals"`
	// Summarize the purpose of the current beat.
	Purpose string `json:"purpose"`

	MinScenes int `json:"minScenes"`
	MaxScenes int `json:"maxScenes"`
}

func storyPlanStaticSceneCount(scenes int) string {
	if scenes == 0 {
		return ""
	}

	if scenes == 1 {
		return "1 scene"
	}

	return fmt.Sprintf("%d scenes", scenes)
}

func (beat BeatDefinition) GetScenesCount() string {
	if beat.MinScenes == 0 && beat.MaxScenes == 0 {
		return ""
	}

	if beat.MinScenes == 0 {
		return fmt.Sprintf("(%s)", storyPlanStaticSceneCount(beat.MaxScenes))
	}

	if beat.MaxScenes == 0 {
		return fmt.Sprintf("(%s)", storyPlanStaticSceneCount(beat.MinScenes))
	}

	if beat.MinScenes == beat.MaxScenes {
		return fmt.Sprintf("(%s)", storyPlanStaticSceneCount(beat.MinScenes))
	}

	return fmt.Sprintf("(%d to %d scenes)", beat.MinScenes, beat.MaxScenes)
}
