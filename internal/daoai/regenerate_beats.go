package daoai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/config/schemas"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"strings"
	"text/template"

	"github.com/samber/lo"

	"github.com/a-novel/service-story-schematics/config/prompts"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/models"
)

const regenerateBeatsTemperature = 0.6

var RegenerateBeatsPrompts = struct {
	System *template.Template
	Input1 *template.Template
	Input2 *template.Template
}{
	System: template.Must(template.New(string(models.LangEN)).Parse(prompts.Config.En.RegenerateBeats.System)),
	Input1: template.Must(template.New(string(models.LangEN)).Parse(prompts.Config.En.RegenerateBeats.Input1)),
	Input2: template.Must(template.New(string(models.LangEN)).Parse(prompts.Config.En.RegenerateBeats.Input2)),
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
	storyPlanPartialPrompt, err := StoryPlanToPrompt("EN", request.Plan)
	if err != nil {
		return nil, NewErrRegenerateBeatsRepository(fmt.Errorf("parse story plan prompt: %w", err))
	}

	systemPrompt := new(strings.Builder)

	err = RegenerateBeatsPrompts.System.ExecuteTemplate(systemPrompt, string(models.LangEN), map[string]any{
		"StoryPlan": storyPlanPartialPrompt,
		"Plan":      request.Plan,
	})
	if err != nil {
		return nil, NewErrRegenerateBeatsRepository(fmt.Errorf("parse system message: %w", err))
	}

	userPrompt1 := new(strings.Builder)

	err = RegenerateBeatsPrompts.Input1.ExecuteTemplate(userPrompt1, string(models.LangEN), map[string]any{
		"Logline": request.Logline,
	})
	if err != nil {
		return nil, NewErrRegenerateBeatsRepository(fmt.Errorf("parse user message: %w", err))
	}

	userPrompt2 := new(strings.Builder)

	err = RegenerateBeatsPrompts.Input2.ExecuteTemplate(userPrompt2, string(models.LangEN), map[string]any{
		"Beats": request.RegenerateKeys,
	})
	if err != nil {
		return nil, NewErrRegenerateBeatsRepository(fmt.Errorf("parse user message: %w", err))
	}

	chatCompletion, err := lib.OpenAIClient(ctx).Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       config.Groq.Model,
		Temperature: param.NewOpt(regenerateBeatsTemperature),
		User:        param.NewOpt(request.UserID),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt.String()),
			openai.UserMessage(userPrompt1.String()),
			openai.AssistantMessage(repository.extrudedBeatsSheet(request)),
			openai.UserMessage(userPrompt2.String()),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:        "beats",
					Description: openai.String(schemas.Config.En.Beats.Description),
					Schema:      schemas.Config.En.Beats.Schema,
					Strict:      openai.Bool(true),
				},
			},
		},
	})
	if err != nil {
		return nil, NewErrRegenerateBeatsRepository(err)
	}

	var beats struct {
		Beats []models.Beat `json:"beats"`
	}

	if err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &beats); err != nil {
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
