package daoai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/config/schemas"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"strings"
	"text/template"

	"github.com/a-novel/service-story-schematics/config/prompts"
	"github.com/a-novel/service-story-schematics/models"
)

const generateLoglineTemperature = 1.0

var GenerateLoglinesPrompts = struct {
	Themed *template.Template
	Random *template.Template
}{
	Themed: template.Must(template.New(string(models.LangEN)).Parse(prompts.Config.En.GenerateLogline.System.Themed)),
	Random: template.Must(template.New(string(models.LangEN)).Parse(prompts.Config.En.GenerateLogline.System.Random)),
}

var ErrGenerateLoglinesRepository = errors.New("GenerateLoglinesRepository.GenerateLoglines")

func NewErrGenerateLoglinesRepository(err error) error {
	return errors.Join(err, ErrGenerateLoglinesRepository)
}

type GenerateLoglinesRequest struct {
	Count  int
	Theme  string
	UserID string
}

type GenerateLoglinesRepository struct{}

func (repository *GenerateLoglinesRepository) GenerateLoglines(
	ctx context.Context, request GenerateLoglinesRequest,
) ([]models.LoglineIdea, error) {
	var (
		err      error
		messages []openai.ChatCompletionMessageParamUnion
	)

	if request.Theme != "" {
		systemPrompt := new(strings.Builder)
		err = GenerateLoglinesPrompts.Themed.ExecuteTemplate(systemPrompt, string(models.LangEN), request)

		messages = []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt.String()),
			openai.UserMessage(request.Theme),
		}
	} else {
		userPrompt := new(strings.Builder)
		err = GenerateLoglinesPrompts.Random.ExecuteTemplate(userPrompt, string(models.LangEN), request)

		messages = []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(userPrompt.String()),
		}
	}

	if err != nil {
		return nil, NewErrGenerateLoglinesRepository(fmt.Errorf("parse system message: %w", err))
	}

	chatCompletion, err := lib.OpenAIClient(ctx).Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:       config.Groq.Model,
		Temperature: param.NewOpt(generateLoglineTemperature),
		User:        param.NewOpt(request.UserID),
		Messages:    messages,
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:        "logline",
					Description: openai.String(schemas.Config.En.Loglines.Description),
					Schema:      schemas.Config.En.Loglines.Schema,
					Strict:      openai.Bool(true),
				},
			},
		},
	})
	if err != nil {
		return nil, NewErrGenerateLoglinesRepository(err)
	}

	var loglines struct {
		Loglines []models.LoglineIdea `json:"loglines"`
	}

	if err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &loglines); err != nil {
		return nil, NewErrGenerateLoglinesRepository(err)
	}

	return loglines.Loglines, nil
}

func NewGenerateLoglinesRepository() *GenerateLoglinesRepository {
	return &GenerateLoglinesRepository{}
}
