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

const expandLoglineTemperature = 0.8

var ExpandLoglinePrompts = struct {
	System *template.Template
}{
	System: RegisterTemplateLocales(prompts.ExpandLogline[models.LangEN].System, map[models.Lang]string{
		models.LangEN: prompts.ExpandLogline[models.LangEN].System,
		models.LangFR: prompts.ExpandLogline[models.LangFR].System,
	}),
}

var ExpandLoglineSchemas = map[models.Lang]any{
	models.LangEN: schemas.Logline[models.LangEN].Schema,
	models.LangFR: schemas.Logline[models.LangFR].Schema,
}

var ExpandLoglineDescriptions = map[models.Lang]string{
	models.LangEN: schemas.Logline[models.LangEN].Description,
	models.LangFR: schemas.Logline[models.LangFR].Description,
}

type ExpandLoglineRequest struct {
	Logline string
	UserID  string
	Lang    models.Lang
}

type ExpandLoglineRepository struct {
	config *config.OpenAI
}

func NewExpandLoglineRepository(config *config.OpenAI) *ExpandLoglineRepository {
	return &ExpandLoglineRepository{config: config}
}

func (repository *ExpandLoglineRepository) ExpandLogline(
	ctx context.Context, request ExpandLoglineRequest,
) (*models.LoglineIdea, error) {
	ctx, span := otel.Tracer().Start(ctx, "daoai.ExpandLogline")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.userID", request.UserID),
		attribute.String("request.lang", request.Lang.String()),
		attribute.String("request.logline", request.Logline),
	)

	systemPrompt := new(strings.Builder)

	err := ExpandLoglinePrompts.System.ExecuteTemplate(systemPrompt, request.Lang.String(), nil)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("execute system prompt: %w", err))
	}

	chatCompletion, err := repository.config.Client().
		Chat.Completions.
		New(ctx, openai.ChatCompletionNewParams{
			Model:       repository.config.Model,
			Temperature: param.NewOpt(expandLoglineTemperature),
			User:        param.NewOpt(request.UserID),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemPrompt.String()),
				openai.UserMessage(request.Logline),
			},
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
					JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
						Name:        "logline",
						Description: openai.String(ExpandLoglineDescriptions[request.Lang]),
						Schema:      ExpandLoglineSchemas[request.Lang],
						Strict:      openai.Bool(true),
					},
				},
			},
		})
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	var logline models.LoglineIdea

	err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &logline)
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	logline.Lang = request.Lang

	return otel.ReportSuccess(span, &logline), nil
}
