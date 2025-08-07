package daoai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/packages/param"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/daoai/prompts"
	"github.com/a-novel/service-story-schematics/internal/daoai/schemas"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/config"
)

var ExpandLoglinePrompts = struct {
	System *template.Template
}{
	System: template.Must(template.New("").Parse(prompts.ExpandLogline.System)),
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

	err := ExpandLoglinePrompts.System.Execute(systemPrompt, nil)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("execute system prompt: %w", err))
	}

	chatCompletion, err := repository.config.Client().
		Chat.Completions.
		New(ctx, openai.ChatCompletionNewParams{
			Model: repository.config.Model,
			User:  param.NewOpt(request.UserID),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(ForceNextAnswerLocale(request.Lang, systemPrompt.String())),
				openai.UserMessage(request.Logline),
			},
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
					JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
						Name:        "logline",
						Description: openai.String(schemas.Logline.Description),
						Schema:      schemas.Logline.Schema,
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
