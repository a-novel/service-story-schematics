package storyplanmodel

import (
	"errors"
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/a-novel/service-story-schematics/models"
)

var (
	ErrInvalidPlan   = errors.New("invalid plan")
	ErrMissingBeat   = fmt.Errorf("%w: missing beat", ErrInvalidPlan)
	ErrMisplacedBeat = fmt.Errorf("%w: misplaced beat", ErrInvalidPlan)
	ErrExtraBeat     = fmt.Errorf("%w: extra beat", ErrInvalidPlan)
)

type Plan struct {
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	Beats    []Beat   `json:"beats"    yaml:"beats"`
}

func (plan Plan) Pick(beats ...string) *Plan {
	pickedBeats := lo.Filter(plan.Beats, func(beat Beat, _ int) bool {
		return lo.Contains(beats, beat.Key)
	})

	return &Plan{
		Metadata: plan.Metadata,
		Beats:    pickedBeats,
	}
}

func (plan Plan) GetBeat(key string) (*Beat, error) {
	beat, ok := lo.Find(plan.Beats, func(b Beat) bool { return b.Key == key })
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrMissingBeat, key)
	}

	return &beat, nil
}

func (plan Plan) Validate(beats []models.Beat) error {
	var errs []error

	for planIndex, planBeat := range plan.Beats {
		_, beatIndex, ok := lo.FindIndexOf(beats, func(b models.Beat) bool {
			return b.Key == planBeat.Key
		})

		if !ok {
			errs = append(
				errs,
				fmt.Errorf("%w: %s at index %d", ErrMissingBeat, planBeat.Key, planIndex),
			)

			continue
		}

		if planIndex != beatIndex {
			errs = append(
				errs,
				fmt.Errorf(
					"%w: expected beat %s at index %d, found at index %d",
					ErrMisplacedBeat, planBeat.Key, planIndex, beatIndex,
				),
			)
		}
	}

	for _, beat := range beats {
		if !lo.ContainsBy(plan.Beats, func(b Beat) bool { return b.Key == beat.Key }) {
			errs = append(errs, fmt.Errorf("%w: %s", ErrExtraBeat, beat.Key))
		}
	}

	return errors.Join(errs...)
}

func (plan Plan) OutputSchema() any {
	return map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required":             []string{"beats"},
		"properties": map[string]any{
			"beats": map[string]any{
				"type": "array",
				"description": "The beats that compose the story. " +
					"A beat is a unit of story structure that represents a specific moment or event in the narrative.",
				"prefixItems": lo.Map(plan.Beats, func(item Beat, index int) any {
					return item.OutputSchema()
				}),
			},
		},
	}
}

type Metadata struct {
	Name string      `json:"name" yaml:"name"`
	Lang models.Lang `json:"lang" yaml:"lang"`
}

type Beat struct {
	Name      string   `json:"name"      yaml:"name"`
	Key       string   `json:"key"       yaml:"key"`
	KeyPoints []string `json:"keyPoints" yaml:"keyPoints"`
	Purpose   string   `json:"purpose"   yaml:"purpose"`
	Scenes    Scenes   `json:"scenes"    yaml:"scenes"`
}

func (beat Beat) String() string {
	return fmt.Sprintf(
		"%s (%s)\nKey Points: %s\nPurpose: %s",
		beat.Name,
		beat.Scenes.String(),
		"\n\t- "+strings.Join(beat.KeyPoints, "\n\t- "),
		beat.Purpose,
	)
}

func (beat Beat) OutputSchema() any {
	return map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required":             []string{"key", "content", "title"},
		"properties": map[string]any{
			"key": map[string]any{
				"const": beat.Key,
			},
			"title": map[string]any{
				"type":        "string",
				"description": "A short title representing the beat.",
			},
			"content": map[string]any{
				"type": "string",
				"description": fmt.Sprintf(
					"A summary of the %s that make up the '%s' beat.\nKey Points: %s\nPurpose: %s",
					beat.Scenes.String(),
					beat.Name,
					"\n\t- "+strings.Join(beat.KeyPoints, "\n\t- "),
					beat.Purpose,
				),
			},
		},
	}
}

type Scenes struct {
	Exact *int `json:"exact,omitempty" yaml:"exact,omitempty"`
	Min   *int `json:"min,omitempty"   yaml:"min,omitempty"`
	Max   *int `json:"max,omitempty"   yaml:"max,omitempty"`
}

func (scenes Scenes) String() string {
	switch {
	case scenes.Exact != nil:
		if *scenes.Exact == 1 {
			return "exactly 1 scene"
		}

		return fmt.Sprintf("exactly %d scenes", *scenes.Exact)
	case scenes.Min != nil && scenes.Max != nil:
		return fmt.Sprintf("between %d and %d scenes", *scenes.Min, *scenes.Max)
	case scenes.Min != nil:
		return fmt.Sprintf("at least %d scenes", *scenes.Min)
	case scenes.Max != nil:
		return fmt.Sprintf("no more than %d scenes", *scenes.Max)
	default:
		return "any number of scenes"
	}
}
