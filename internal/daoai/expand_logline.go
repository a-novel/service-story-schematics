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

const expandLoglineTemperature = 0.8

var ExpandLoglinePrompts = struct {
	System *template.Template
}{
	System: RegisterTemplateLocales(prompts.Config.En.ExpandLogline, map[models.Lang]string{
		models.LangEN: prompts.Config.En.ExpandLogline,
		models.LangFR: prompts.Config.Fr.ExpandLogline,
	}),
}

var ExpandLoglineSchemas = map[models.Lang]any{
	models.LangEN: schemas.Config.En.Logline.Schema,
	models.LangFR: schemas.Config.Fr.Logline.Schema,
}

var ExpandLoglineDescriptions = map[models.Lang]string{
	models.LangEN: schemas.Config.En.Logline.Description,
	models.LangFR: schemas.Config.Fr.Logline.Description,
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

func (repository *ExpandLoglineRepository) ExpandLogline(
	ctx context.Context, request ExpandLoglineRequest,
) (*models.LoglineIdea, error) {
	systemPrompt := new(strings.Builder)
	if err := ExpandLoglinePrompts.System.ExecuteTemplate(systemPrompt, request.Lang.String(), nil); err != nil {
		return nil, NewErrExpandLoglineRepository(fmt.Errorf("execute system prompt: %w", err))
	}

	chatCompletion, err := lib.OpenAIClient(ctx).Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
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
		return nil, NewErrExpandLoglineRepository(err)
	}

	var logline models.LoglineIdea

	if err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &logline); err != nil {
		return nil, NewErrExpandLoglineRepository(err)
	}

	logline.Lang = request.Lang

	return &logline, nil
}

func NewExpandLoglineRepository() *ExpandLoglineRepository {
	return &ExpandLoglineRepository{}
}
