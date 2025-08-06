package daoai

import (
	"context"
	"encoding/json"
	"errors"
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
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/config"
)

const expandBeatTemperature = 0.8

var ExpandBeatPrompts = struct {
	System *template.Template
	Input1 *template.Template
	Input2 *template.Template
}{
	System: template.Must(template.New("").Parse(prompts.ExpandBeat.System)),
	Input1: template.Must(template.New("").Parse(prompts.ExpandBeat.Input1)),
	Input2: template.Must(template.New("").Parse(prompts.ExpandBeat.Input2)),
}

var ErrUnknownTargetKey = errors.New("unknown target key")

type ExpandBeatRequest struct {
	Logline   string
	Beats     []models.Beat
	Plan      models.StoryPlan
	Lang      models.Lang
	TargetKey string
	UserID    string
}

type ExpandBeatRepository struct {
	config *config.OpenAI
}

func NewExpandBeatRepository(config *config.OpenAI) *ExpandBeatRepository {
	return &ExpandBeatRepository{config: config}
}

func (repository *ExpandBeatRepository) ExpandBeat(
	ctx context.Context, request ExpandBeatRequest,
) (*models.Beat, error) {
	ctx, span := otel.Tracer().Start(ctx, "daoai.ExpandBeat")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.userID", request.UserID),
		attribute.String("request.storyPlan.id", request.Plan.ID.String()),
		attribute.String("request.storyPlan.slug", request.Plan.Slug.String()),
		attribute.String("request.storyPlan.name", request.Plan.Name),
		attribute.String("request.storyPlan.lang", request.Plan.Lang.String()),
		attribute.String("request.targetKey", request.TargetKey),
		attribute.String("request.logline", request.Logline),
	)

	if !lo.ContainsBy(request.Plan.Beats, func(item models.BeatDefinition) bool {
		return item.Key == request.TargetKey
	}) {
		return nil, otel.ReportError(span, ErrUnknownTargetKey)
	}

	storyPlanPartialPrompt, err := StoryPlanToPrompt(request.Plan)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("parse story plan prompt: %w", err))
	}

	systemPrompt := new(strings.Builder)

	err = ExpandBeatPrompts.System.Execute(systemPrompt, map[string]any{
		"StoryPlan": storyPlanPartialPrompt,
		"Plan":      request.Plan,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("parse system message: %w", err))
	}

	userPrompt1 := new(strings.Builder)

	err = ExpandBeatPrompts.Input1.Execute(userPrompt1, request)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("parse user message 1: %w", err))
	}

	userPrompt2 := new(strings.Builder)

	err = ExpandBeatPrompts.Input2.Execute(userPrompt2, request)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("parse user message 2: %w", err))
	}

	chatCompletion, err := repository.config.Client().
		Chat.Completions.
		New(ctx, openai.ChatCompletionNewParams{
			Model:       repository.config.Model,
			Temperature: param.NewOpt(expandBeatTemperature),
			User:        param.NewOpt(request.UserID),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemPrompt.String()),
				openai.UserMessage(userPrompt1.String()),
				repository.buildBeatsSheetResponse(request),
				openai.UserMessage(ForceNextAnswerLocale(request.Lang, userPrompt2.String())),
			},
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
					JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
						Name:        "storyBeat",
						Description: openai.String(schemas.Beat.Description),
						Schema:      schemas.Beat.Schema,
						Strict:      openai.Bool(true),
					},
				},
			},
		})
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	var beat models.Beat

	err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &beat)
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	return otel.ReportSuccess(span, &beat), nil
}

func (repository *ExpandBeatRepository) buildBeatsSheetResponse(
	request ExpandBeatRequest,
) openai.ChatCompletionMessageParamUnion {
	parts := make([]string, len(request.Beats))

	for i, beat := range request.Beats {
		parts[i] = fmt.Sprintf("%s\n%s", beat.Key, beat.Content)
	}

	return openai.AssistantMessage(strings.Join(parts, "\n\n"))
}
