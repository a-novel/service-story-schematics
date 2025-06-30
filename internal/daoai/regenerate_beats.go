package daoai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/getsentry/sentry-go"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/samber/lo"

	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/config/prompts"
	"github.com/a-novel/service-story-schematics/config/schemas"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/models"
)

const regenerateBeatsTemperature = 0.6

var RegenerateBeatsPrompts = struct {
	System *template.Template
	Input1 *template.Template
	Input2 *template.Template
}{
	System: RegisterTemplateLocales(prompts.Config.En.RegenerateBeats.System, map[models.Lang]string{
		models.LangEN: prompts.Config.En.RegenerateBeats.System,
		models.LangFR: prompts.Config.Fr.RegenerateBeats.System,
	}),
	Input1: RegisterTemplateLocales(prompts.Config.En.RegenerateBeats.Input1, map[models.Lang]string{
		models.LangEN: prompts.Config.En.RegenerateBeats.Input1,
		models.LangFR: prompts.Config.Fr.RegenerateBeats.Input1,
	}),
	Input2: RegisterTemplateLocales(prompts.Config.En.RegenerateBeats.Input2, map[models.Lang]string{
		models.LangEN: prompts.Config.En.RegenerateBeats.Input2,
		models.LangFR: prompts.Config.Fr.RegenerateBeats.Input2,
	}),
}

var RegenerateBeatsSchemas = map[models.Lang]any{
	models.LangEN: schemas.Config.En.Beats.Schema,
	models.LangFR: schemas.Config.Fr.Beats.Schema,
}

var RegenerateBeatsDescriptions = map[models.Lang]string{
	models.LangEN: schemas.Config.En.Beats.Description,
	models.LangFR: schemas.Config.Fr.Beats.Description,
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
	Lang           models.Lang
}

type RegenerateBeatsRepository struct{}

func NewRegenerateBeatsRepository() *RegenerateBeatsRepository {
	return &RegenerateBeatsRepository{}
}

func (repository *RegenerateBeatsRepository) RegenerateBeats(
	ctx context.Context, request RegenerateBeatsRequest,
) ([]models.Beat, error) {
	span := sentry.StartSpan(ctx, "RegenerateBeatsRepository.RegenerateBeats")
	defer span.Finish()

	span.SetData("request.logline", request.Logline)
	span.SetData("request.userID", request.UserID)
	span.SetData("request.storyPlan.id", request.Plan.ID.String())
	span.SetData("request.regenerateKeys", request.RegenerateKeys)
	span.SetData("request.lang", request.Lang.String())

	storyPlanPartialPrompt, err := StoryPlanToPrompt(request.Lang, request.Plan)
	if err != nil {
		span.SetData("storyPlan.toPrompt.error", err.Error())

		return nil, NewErrRegenerateBeatsRepository(fmt.Errorf("parse story plan prompt: %w", err))
	}

	systemPrompt := new(strings.Builder)

	err = RegenerateBeatsPrompts.System.ExecuteTemplate(systemPrompt, request.Lang.String(), map[string]any{
		"StoryPlan": storyPlanPartialPrompt,
		"Plan":      request.Plan,
	})
	if err != nil {
		span.SetData("systemPrompt.error", err.Error())

		return nil, NewErrRegenerateBeatsRepository(fmt.Errorf("parse system message: %w", err))
	}

	userPrompt1 := new(strings.Builder)

	err = RegenerateBeatsPrompts.Input1.ExecuteTemplate(userPrompt1, request.Lang.String(), map[string]any{
		"Logline": request.Logline,
	})
	if err != nil {
		span.SetData("userPrompt1.error", err.Error())

		return nil, NewErrRegenerateBeatsRepository(fmt.Errorf("parse user message: %w", err))
	}

	userPrompt2 := new(strings.Builder)

	err = RegenerateBeatsPrompts.Input2.ExecuteTemplate(userPrompt2, request.Lang.String(), map[string]any{
		"Beats": request.RegenerateKeys,
	})
	if err != nil {
		span.SetData("userPrompt2.error", err.Error())

		return nil, NewErrRegenerateBeatsRepository(fmt.Errorf("parse user message: %w", err))
	}

	chatCompletion, err := lib.OpenAIClient(span.Context()).
		Chat.Completions.
		New(span.Context(), openai.ChatCompletionNewParams{
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
						Description: openai.String(RegenerateBeatsDescriptions[request.Lang]),
						Schema:      RegenerateBeatsSchemas[request.Lang],
						Strict:      openai.Bool(true),
					},
				},
			},
		})
	if err != nil {
		span.SetData("chatCompletion.error", err.Error())

		return nil, NewErrRegenerateBeatsRepository(err)
	}

	var beats struct {
		Beats []models.Beat `json:"beats"`
	}

	err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &beats)
	if err != nil {
		span.SetData("unmarshal.error", err.Error())

		return nil, NewErrRegenerateBeatsRepository(err)
	}

	err = lib.CheckStoryPlan(beats.Beats, repository.buildExpectedStoryPlan(request))
	if err != nil {
		span.SetData("checkStoryPlan.error", err.Error())

		return nil, NewErrRegenerateBeatsRepository(fmt.Errorf("check story plan: %w", err))
	}

	return repository.mergeSourceAndNewBeats(request, beats.Beats), nil
}

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
