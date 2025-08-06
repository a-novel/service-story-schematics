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
	Themed: template.Must(template.New("").Parse(prompts.GenerateLoglines.System.Themed)),
	Random: template.Must(template.New("").Parse(prompts.GenerateLoglines.System.Random)),
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
		err = GenerateLoglinesPrompts.Themed.Execute(systemPrompt, request)

		messages = []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(ForceNextAnswerLocale(request.Lang, systemPrompt.String())),
			openai.UserMessage(request.Theme),
		}
	} else {
		userPrompt := new(strings.Builder)
		err = GenerateLoglinesPrompts.Random.Execute(userPrompt, request)

		messages = []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(ForceNextAnswerLocale(request.Lang, "")),
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
						Description: openai.String(schemas.Loglines.Description),
						Schema:      schemas.Loglines.Schema,
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
