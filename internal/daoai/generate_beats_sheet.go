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
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/daoai/prompts"
	"github.com/a-novel/service-story-schematics/internal/daoai/schemas"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/config"
)

const generateBeatsSheetTemperature = 0.8

var GenerateBeatsSheetPrompts = struct {
	System *template.Template
}{
	System: template.Must(template.New("").Parse(prompts.GenerateBeatsSheet.System)),
}

var ErrInvalidBeatSheet = errors.New("invalid beat sheet")

type GenerateBeatsSheetRequest struct {
	Logline string
	Plan    models.StoryPlan
	UserID  string
	Lang    models.Lang
}

type GenerateBeatsSheetRepository struct {
	config *config.OpenAI
}

func NewGenerateBeatsSheetRepository(config *config.OpenAI) *GenerateBeatsSheetRepository {
	return &GenerateBeatsSheetRepository{config: config}
}

func (repository *GenerateBeatsSheetRepository) GenerateBeatsSheet(
	ctx context.Context, request GenerateBeatsSheetRequest,
) ([]models.Beat, error) {
	ctx, span := otel.Tracer().Start(ctx, "daoai.GenerateBeatsSheet")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.logline", request.Logline),
		attribute.String("request.plan.id", request.Plan.ID.String()),
		attribute.String("request.lang", request.Lang.String()),
		attribute.String("request.userID", request.UserID),
	)

	storyPlanPartialPrompt, err := StoryPlanToPrompt(request.Plan)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("parse story plan prompt: %w", err))
	}

	systemPrompt := new(strings.Builder)

	err = GenerateBeatsSheetPrompts.System.Execute(systemPrompt, map[string]any{
		"StoryPlan": storyPlanPartialPrompt,
		"Plan":      request.Plan,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("execute system prompt: %w", err))
	}

	chatCompletion, err := repository.config.Client().
		Chat.Completions.
		New(ctx, openai.ChatCompletionNewParams{
			Model:       repository.config.Model,
			Temperature: param.NewOpt(generateBeatsSheetTemperature),
			User:        param.NewOpt(request.UserID),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(ForceNextAnswerLocale(request.Lang, systemPrompt.String())),
				openai.UserMessage(request.Logline),
			},
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
					JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
						Name:        "story_beats",
						Description: openai.String(schemas.Beats.Description),
						Schema:      schemas.Beats.Schema,
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

	err = lib.CheckStoryPlan(beats.Beats, request.Plan.Beats)
	if err != nil {
		return nil, otel.ReportError(span, errors.Join(err, ErrInvalidBeatSheet))
	}

	return otel.ReportSuccess(span, beats.Beats), nil
}
