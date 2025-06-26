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

const generateBeatsSheetTemperature = 0.8

var GenerateBeatsSheetPrompts = struct {
	System *template.Template
}{
	System: RegisterTemplateLocales(prompts.Config.En.GenerateBeatsSheet, map[models.Lang]string{
		models.LangEN: prompts.Config.En.GenerateBeatsSheet,
		models.LangFR: prompts.Config.Fr.GenerateBeatsSheet,
	}),
}

var GenerateBeatsSheetSchemas = map[models.Lang]any{
	models.LangEN: schemas.Config.En.Beats.Schema,
	models.LangFR: schemas.Config.Fr.Beats.Schema,
}

var GenerateBeatsSheetDescriptions = map[models.Lang]string{
	models.LangEN: schemas.Config.En.Beats.Description,
	models.LangFR: schemas.Config.Fr.Beats.Description,
}

var ErrInvalidBeatSheet = errors.New("invalid beat sheet")

var ErrGenerateBeatsSheetRepository = errors.New("GenerateBeatsSheetRepository.GenerateBeatsSheet")

func NewErrGenerateBeatsSheetRepository(err error) error {
	return errors.Join(err, ErrGenerateBeatsSheetRepository)
}

type GenerateBeatsSheetRequest struct {
	Logline string
	Plan    models.StoryPlan
	UserID  string
	Lang    models.Lang
}

type GenerateBeatsSheetRepository struct{}

func (repository *GenerateBeatsSheetRepository) GenerateBeatsSheet(
	ctx context.Context, request GenerateBeatsSheetRequest,
) ([]models.Beat, error) {
	span := sentry.StartSpan(ctx, "GenerateBeatsSheetRepository.GenerateBeatsSheet")
	defer span.Finish()

	span.SetData("request.logline", request.Logline)
	span.SetData("request.storyPlan.id", request.Plan.ID)
	span.SetData("request.lang", request.Lang.String())
	span.SetData("request.userID", request.UserID)

	storyPlanPartialPrompt, err := StoryPlanToPrompt(request.Lang, request.Plan)
	if err != nil {
		span.SetData("storyPlan.toPrompt.error", err.Error())

		return nil, NewErrGenerateBeatsSheetRepository(fmt.Errorf("parse story plan prompt: %w", err))
	}

	systemPrompt := new(strings.Builder)

	err = GenerateBeatsSheetPrompts.System.ExecuteTemplate(systemPrompt, request.Lang.String(), map[string]any{
		"StoryPlan": storyPlanPartialPrompt,
		"Plan":      request.Plan,
	})
	if err != nil {
		span.SetData("prompt.error", err.Error())

		return nil, NewErrGenerateBeatsSheetRepository(fmt.Errorf("execute system prompt: %w", err))
	}

	chatCompletion, err := lib.OpenAIClient(span.Context()).
		Chat.Completions.
		New(span.Context(), openai.ChatCompletionNewParams{
			Model:       config.Groq.Model,
			Temperature: param.NewOpt(generateBeatsSheetTemperature),
			User:        param.NewOpt(request.UserID),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemPrompt.String()),
				openai.UserMessage(request.Logline),
			},
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
					JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
						Name:        "story_beats",
						Description: openai.String(GenerateBeatsSheetDescriptions[request.Lang]),
						Schema:      GenerateBeatsSheetSchemas[request.Lang],
						Strict:      openai.Bool(true),
					},
				},
			},
		})
	if err != nil {
		span.SetData("chatCompletion.error", err.Error())

		return nil, NewErrGenerateBeatsSheetRepository(err)
	}

	var beats struct {
		Beats []models.Beat `json:"beats"`
	}

	if err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &beats); err != nil {
		span.SetData("unmarshal.error", err.Error())

		return nil, NewErrGenerateBeatsSheetRepository(err)
	}

	if err = lib.CheckStoryPlan(beats.Beats, request.Plan.Beats); err != nil {
		span.SetData("checkStoryPlan.error", err.Error())

		return nil, NewErrGenerateBeatsSheetRepository(errors.Join(err, ErrInvalidBeatSheet))
	}

	return beats.Beats, nil
}

func NewGenerateBeatsSheetRepository() *GenerateBeatsSheetRepository {
	return &GenerateBeatsSheetRepository{}
}
