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
	"github.com/a-novel/service-story-schematics/config/prompts"
	"github.com/a-novel/service-story-schematics/config/schemas"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/models"
)

const generateLoglineTemperature = 1.0

var GenerateLoglinesPrompts = struct {
	Themed *template.Template
	Random *template.Template
}{
	Themed: RegisterTemplateLocales(prompts.Config.En.GenerateLogline.System.Themed, map[models.Lang]string{
		models.LangEN: prompts.Config.En.GenerateLogline.System.Themed,
		models.LangFR: prompts.Config.Fr.GenerateLogline.System.Themed,
	}),
	Random: RegisterTemplateLocales(prompts.Config.En.GenerateLogline.System.Random, map[models.Lang]string{
		models.LangEN: prompts.Config.En.GenerateLogline.System.Random,
		models.LangFR: prompts.Config.Fr.GenerateLogline.System.Random,
	}),
}

var GenerateLoglinesSchemas = map[models.Lang]any{
	models.LangEN: schemas.Config.En.Loglines.Schema,
	models.LangFR: schemas.Config.Fr.Loglines.Schema,
}

var GenerateLoglinesDescriptions = map[models.Lang]string{
	models.LangEN: schemas.Config.En.Loglines.Description,
	models.LangFR: schemas.Config.Fr.Loglines.Description,
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
