package daoai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/daoai/prompts"
	"github.com/a-novel/service-story-schematics/internal/daoai/schemas"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/config"
)

const generateLoglineTemperature = 1.0

var GenerateLoglinesPrompts = struct {
	Themed *template.Template
	Random *template.Template
}{
	Themed: RegisterTemplateLocales(prompts.GenerateLoglines[models.LangEN].System.Themed, map[models.Lang]string{
		models.LangEN: prompts.GenerateLoglines[models.LangEN].System.Themed,
		models.LangFR: prompts.GenerateLoglines[models.LangFR].System.Themed,
	}),
	Random: RegisterTemplateLocales(prompts.GenerateLoglines[models.LangEN].System.Random, map[models.Lang]string{
		models.LangEN: prompts.GenerateLoglines[models.LangEN].System.Random,
		models.LangFR: prompts.GenerateLoglines[models.LangFR].System.Random,
	}),
}

var GenerateLoglinesSchemas = map[models.Lang]any{
	models.LangEN: schemas.Loglines[models.LangEN].Schema,
	models.LangFR: schemas.Loglines[models.LangFR].Schema,
}

var GenerateLoglinesDescriptions = map[models.Lang]string{
	models.LangEN: schemas.Loglines[models.LangEN].Description,
	models.LangFR: schemas.Loglines[models.LangFR].Description,
}

type GenerateLoglinesRequest struct {
	Count  int
	Theme  string
	UserID string
	Lang   models.Lang
}

type GenerateLoglinesRepository struct {
	config *config.OpenAI
}

func NewGenerateLoglinesRepository(config *config.OpenAI) *GenerateLoglinesRepository {
	return &GenerateLoglinesRepository{config: config}
}

func (repository *GenerateLoglinesRepository) GenerateLoglines(
	ctx context.Context, request GenerateLoglinesRequest,
) ([]models.LoglineIdea, error) {
	ctx, span := otel.Tracer().Start(ctx, "daoai.GenerateLoglines")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.user_id", request.UserID),
		attribute.String("request.lang", request.Lang.String()),
		attribute.Int("request.count", request.Count),
		attribute.String("request.theme", request.Theme),
	)

	var (
		err      error
		messages []openai.ChatCompletionMessageParamUnion
	)

	if request.Theme != "" {
		systemPrompt := new(strings.Builder)
		err = GenerateLoglinesPrompts.Themed.ExecuteTemplate(systemPrompt, request.Lang.String(), request)

		messages = []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt.String()),
			openai.UserMessage(request.Theme),
		}
	} else {
		userPrompt := new(strings.Builder)
		err = GenerateLoglinesPrompts.Random.ExecuteTemplate(userPrompt, request.Lang.String(), request)

		messages = []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(userPrompt.String()),
		}
	}

	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("parse system message: %w", err))
	}

	chatCompletion, err := repository.config.Client().
		Chat.Completions.
		New(ctx, openai.ChatCompletionNewParams{
			Model:       repository.config.Model,
			Temperature: param.NewOpt(generateLoglineTemperature),
			User:        param.NewOpt(request.UserID),
			Messages:    messages,
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
					JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
						Name:        "logline",
						Description: openai.String(GenerateLoglinesDescriptions[request.Lang]),
						Schema:      GenerateLoglinesSchemas[request.Lang],
						Strict:      openai.Bool(true),
					},
				},
			},
		})
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	var loglines struct {
		Loglines []models.LoglineIdea `json:"loglines"`
	}

	err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &loglines)
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	for i := range loglines.Loglines {
		loglines.Loglines[i].Lang = request.Lang
	}

	return otel.ReportSuccess(span, loglines.Loglines), nil
}
