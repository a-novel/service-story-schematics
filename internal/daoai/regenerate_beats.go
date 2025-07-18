package daoai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/daoai/prompts"
	"github.com/a-novel/service-story-schematics/internal/daoai/schemas"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/config"
)

const regenerateBeatsTemperature = 0.6

var RegenerateBeatsPrompts = struct {
	System *template.Template
	Input1 *template.Template
	Input2 *template.Template
}{
	System: RegisterTemplateLocales(prompts.RegenerateBeats[models.LangEN].System, map[models.Lang]string{
		models.LangEN: prompts.RegenerateBeats[models.LangEN].System,
		models.LangFR: prompts.RegenerateBeats[models.LangFR].System,
	}),
	Input1: RegisterTemplateLocales(prompts.RegenerateBeats[models.LangEN].Input1, map[models.Lang]string{
		models.LangEN: prompts.RegenerateBeats[models.LangEN].Input1,
		models.LangFR: prompts.RegenerateBeats[models.LangFR].Input1,
	}),
	Input2: RegisterTemplateLocales(prompts.RegenerateBeats[models.LangEN].Input2, map[models.Lang]string{
		models.LangEN: prompts.RegenerateBeats[models.LangEN].Input2,
		models.LangFR: prompts.RegenerateBeats[models.LangFR].Input2,
	}),
}

var RegenerateBeatsSchemas = map[models.Lang]any{
	models.LangEN: schemas.Beats[models.LangEN].Schema,
	models.LangFR: schemas.Beats[models.LangFR].Schema,
}

var RegenerateBeatsDescriptions = map[models.Lang]string{
	models.LangEN: schemas.Beats[models.LangEN].Description,
	models.LangFR: schemas.Beats[models.LangFR].Description,
}

type RegenerateBeatsRequest struct {
	Logline        string
	Beats          []models.Beat
	Plan           models.StoryPlan
	RegenerateKeys []string
	UserID         string
	Lang           models.Lang
}

type RegenerateBeatsRepository struct {
	config *config.OpenAI
}

func NewRegenerateBeatsRepository(config *config.OpenAI) *RegenerateBeatsRepository {
	return &RegenerateBeatsRepository{config: config}
}

func (repository *RegenerateBeatsRepository) RegenerateBeats(
	ctx context.Context, request RegenerateBeatsRequest,
) ([]models.Beat, error) {
	ctx, span := otel.Tracer().Start(ctx, "daoai.RegenerateBeats")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.logline", request.Logline),
		attribute.String("request.userID", request.UserID),
		attribute.String("request.plan.id", request.Plan.ID.String()),
		attribute.StringSlice("request.regenerateKeys", request.RegenerateKeys),
		attribute.String("request.lang", request.Lang.String()),
	)

	storyPlanPartialPrompt, err := StoryPlanToPrompt(request.Lang, request.Plan)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("parse story plan prompt: %w", err))
	}

	systemPrompt := new(strings.Builder)

	err = RegenerateBeatsPrompts.System.ExecuteTemplate(systemPrompt, request.Lang.String(), map[string]any{
		"StoryPlan": storyPlanPartialPrompt,
		"Plan":      request.Plan,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("parse system message: %w", err))
	}

	userPrompt1 := new(strings.Builder)

	err = RegenerateBeatsPrompts.Input1.ExecuteTemplate(userPrompt1, request.Lang.String(), map[string]any{
		"Logline": request.Logline,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("parse user message: %w", err))
	}

	userPrompt2 := new(strings.Builder)

	err = RegenerateBeatsPrompts.Input2.ExecuteTemplate(userPrompt2, request.Lang.String(), map[string]any{
		"Beats": request.RegenerateKeys,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("parse user message: %w", err))
	}

	chatCompletion, err := repository.config.Client().
		Chat.Completions.
		New(ctx, openai.ChatCompletionNewParams{
			Model:       repository.config.Model,
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
		return nil, otel.ReportError(span, err)
	}

	var beats struct {
		Beats []models.Beat `json:"beats"`
	}

	err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &beats)
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	err = lib.CheckStoryPlan(beats.Beats, repository.buildExpectedStoryPlan(request))
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("check story plan: %w", err))
	}

	return otel.ReportSuccess(span, repository.mergeSourceAndNewBeats(request, beats.Beats)), nil
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
