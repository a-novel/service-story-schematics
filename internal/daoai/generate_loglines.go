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

var ErrGenerateLoglinesRepository = errors.New("GenerateLoglinesRepository.GenerateLoglines")

func NewErrGenerateLoglinesRepository(err error) error {
	return errors.Join(err, ErrGenerateLoglinesRepository)
}

type GenerateLoglinesRequest struct {
	Count  int
	Theme  string
	UserID string
	Lang   models.Lang
}

type GenerateLoglinesRepository struct{}

func NewGenerateLoglinesRepository() *GenerateLoglinesRepository {
	return &GenerateLoglinesRepository{}
}

func (repository *GenerateLoglinesRepository) GenerateLoglines(
	ctx context.Context, request GenerateLoglinesRequest,
) ([]models.LoglineIdea, error) {
	span := sentry.StartSpan(ctx, "GenerateLoglinesRepository.GenerateLoglines")
	defer span.Finish()

	span.SetData("request.count", request.Count)
	span.SetData("request.theme", request.Theme)
	span.SetData("request.user_id", request.UserID)
	span.SetData("request.lang", request.Lang.String())

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
		span.SetData("prompt.error", err.Error())

		return nil, NewErrGenerateLoglinesRepository(fmt.Errorf("parse system message: %w", err))
	}

	chatCompletion, err := lib.OpenAIClient(span.Context()).
		Chat.Completions.
		New(span.Context(), openai.ChatCompletionNewParams{
			Model:       config.Groq.Model,
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
		span.SetData("chatCompletion.error", err.Error())

		return nil, NewErrGenerateLoglinesRepository(err)
	}

	var loglines struct {
		Loglines []models.LoglineIdea `json:"loglines"`
	}

	err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &loglines)
	if err != nil {
		span.SetData("unmarshal.error", err.Error())

		return nil, NewErrGenerateLoglinesRepository(err)
	}

	for i := range loglines.Loglines {
		loglines.Loglines[i].Lang = request.Lang
	}

	return loglines.Loglines, nil
}
