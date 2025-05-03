package daoai

import (
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/samber/lo"

	"github.com/a-novel-kit/context"
	"github.com/a-novel-kit/golm"

	"github.com/a-novel/service-story-schematics/config/prompts"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/models"
)

const regenerateBeatsTemperature = 0.6

var RegenerateBeatsPrompts = struct {
	System *template.Template
	Input  *template.Template
}{
	System: template.Must(template.New("EN").Parse(prompts.Config.En.RegenerateBeats.System)),
	Input:  template.Must(template.New("EN").Parse(prompts.Config.En.RegenerateBeats.Input)),
}

var ErrRegenerateBeatsRepository = errors.New("RegenerateBeatsRepository.RegenerateBeats")

func NewErrRegenerateBeatsRepository(err error) error {
	return errors.Join(err, ErrRegenerateBeatsRepository)
}

type RegenerateBeatsRequest struct {
	Logline        string
	Beats          []models.Beat
	Plan           models.StoryPlan
	RegenerateKeys []string
	UserID         string
}

type RegenerateBeatsRepository struct{}

func (repository *RegenerateBeatsRepository) extrudedBeatsSheet(request RegenerateBeatsRequest) string {
	parts := make([]string, 0, len(request.Beats))

	for _, beat := range request.Beats {
		if lo.Contains(request.RegenerateKeys, beat.Key) {
			continue
		}

		parts = append(parts, fmt.Sprintf("%s\n%s", beat.Key, beat.Content))
	}

	return strings.Join(parts, "\n\n")
}

func (repository *RegenerateBeatsRepository) buildExpectedStoryPlan(
	request RegenerateBeatsRequest,
) []models.BeatDefinition {
	beats := make([]models.BeatDefinition, 0, len(request.Plan.Beats))

	for _, beat := range request.Plan.Beats {
		if !lo.Contains(request.RegenerateKeys, beat.Key) {
			continue
		}

		beats = append(beats, beat)
	}

	return beats
}

func (repository *RegenerateBeatsRepository) mergeSourceAndNewBeats(
	request RegenerateBeatsRequest, newBeats []models.Beat,
) []models.Beat {
	var output []models.Beat

	for _, beat := range request.Beats {
		if lo.Contains(request.RegenerateKeys, beat.Key) {
			for _, newBeat := range newBeats {
				if newBeat.Key == beat.Key {
					output = append(output, newBeat)

					break
				}
			}
		} else {
			output = append(output, beat)
		}
	}

	return output
}

func (repository *RegenerateBeatsRepository) RegenerateBeats(
	ctx context.Context, request RegenerateBeatsRequest,
) ([]models.Beat, error) {
	storyPlanPrompt, err := StoryPlanToPrompt("EN", request.Plan)
	if err != nil {
		return nil, NewErrRegenerateBeatsRepository(fmt.Errorf("parse story plan prompt: %w", err))
	}

	systemPrompt, err := golm.NewSystemMessageT(RegenerateBeatsPrompts.System, "EN", map[string]interface{}{
		"StoryPlan": storyPlanPrompt,
		"Plan":      request.Plan,
	})
	if err != nil {
		return nil, NewErrRegenerateBeatsRepository(fmt.Errorf("parse system message: %w", err))
	}

	userPrompt, err := golm.NewUserMessageT(RegenerateBeatsPrompts.Input, "EN", map[string]any{
		"Logline":    request.Logline,
		"BeatsSheet": repository.extrudedBeatsSheet(request),
	})
	if err != nil {
		return nil, NewErrRegenerateBeatsRepository(fmt.Errorf("parse user message: %w", err))
	}

	chat := golm.Context(ctx)
	chat.SetSystem(systemPrompt)

	params := golm.CompletionParams{
		Temperature: lo.ToPtr(regenerateBeatsTemperature),
		User:        request.UserID,
	}

	var beats struct {
		Beats []models.Beat `json:"beats"`
	}

	if err = chat.CompletionJSON(ctx, userPrompt, params, &beats); err != nil {
		return nil, NewErrRegenerateBeatsRepository(err)
	}

	if err = lib.CheckStoryPlan(beats.Beats, repository.buildExpectedStoryPlan(request)); err != nil {
		return nil, NewErrRegenerateBeatsRepository(fmt.Errorf("check story plan: %w", err))
	}

	return repository.mergeSourceAndNewBeats(request, beats.Beats), nil
}

func NewRegenerateBeatsRepository() *RegenerateBeatsRepository {
	return &RegenerateBeatsRepository{}
}
