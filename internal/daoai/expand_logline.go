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

	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/internal/daoai/prompts"
	"github.com/a-novel/service-story-schematics/internal/daoai/schemas"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/models"
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

var ErrExpandLoglineRepository = errors.New("ExpandLoglineRepository.ExpandLogline")

func NewErrExpandLoglineRepository(err error) error {
	return errors.Join(err, ErrExpandLoglineRepository)
}

type ExpandLoglineRequest struct {
	Logline string
	UserID  string
	Lang    models.Lang
}

type ExpandLoglineRepository struct{}

func NewExpandLoglineRepository() *ExpandLoglineRepository {
	return &ExpandLoglineRepository{}
}

func (repository *ExpandLoglineRepository) ExpandLogline(
	ctx context.Context, request ExpandLoglineRequest,
) (*models.LoglineIdea, error) {
	span := sentry.StartSpan(ctx, "ExpandLoglineRepository.ExpandLogline")
	defer span.Finish()

	span.SetData("request.logline", request.Logline)
	span.SetData("request.userID", request.UserID)
	span.SetData("request.lang", request.Lang.String())

	systemPrompt := new(strings.Builder)

	err := ExpandLoglinePrompts.System.ExecuteTemplate(systemPrompt, request.Lang.String(), nil)
	if err != nil {
		span.SetData("prompt.error", err.Error())

		return nil, NewErrExpandLoglineRepository(fmt.Errorf("execute system prompt: %w", err))
	}

	chatCompletion, err := lib.OpenAIClient(span.Context()).
		Chat.Completions.
		New(span.Context(), openai.ChatCompletionNewParams{
			Model:       config.Groq.Model,
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
		span.SetData("chatCompletion.error", err.Error())

		return nil, NewErrExpandLoglineRepository(err)
	}

	var logline models.LoglineIdea

	err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &logline)
	if err != nil {
		span.SetData("unmarshal.error", err.Error())

		return nil, NewErrExpandLoglineRepository(err)
	}

	logline.Lang = request.Lang

	return &logline, nil
}
