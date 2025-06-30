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
	"github.com/samber/lo"

	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/config/prompts"
	"github.com/a-novel/service-story-schematics/config/schemas"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/models"
)

const expandBeatTemperature = 0.8

var ExpandBeatPrompts = struct {
	System *template.Template
	Input1 *template.Template
	Input2 *template.Template
}{
	System: RegisterTemplateLocales(prompts.Config.En.ExpandBeat.System, map[models.Lang]string{
		models.LangEN: prompts.Config.En.ExpandBeat.System,
		models.LangFR: prompts.Config.Fr.ExpandBeat.System,
	}),
	Input1: RegisterTemplateLocales(prompts.Config.En.ExpandBeat.Input1, map[models.Lang]string{
		models.LangEN: prompts.Config.En.ExpandBeat.Input1,
		models.LangFR: prompts.Config.Fr.ExpandBeat.Input1,
	}),
	Input2: RegisterTemplateLocales(prompts.Config.En.ExpandBeat.Input2, map[models.Lang]string{
		models.LangEN: prompts.Config.En.ExpandBeat.Input2,
		models.LangFR: prompts.Config.Fr.ExpandBeat.Input2,
	}),
}

var ExpandBeatSchemas = map[models.Lang]any{
	models.LangEN: schemas.Config.En.Beat.Schema,
	models.LangFR: schemas.Config.Fr.Beat.Schema,
}

var ExpandBeatDescriptions = map[models.Lang]string{
	models.LangEN: schemas.Config.En.Beat.Description,
	models.LangFR: schemas.Config.Fr.Beat.Description,
}

var ErrExpandBeatRepository = errors.New(".ExpandBeat")

func NewErrExpandBeatRepository(err error) error {
	return errors.Join(err, ErrExpandBeatRepository)
}

var ErrUnknownTargetKey = errors.New("unknown target key")

type ExpandBeatRequest struct {
	Logline   string
	Beats     []models.Beat
	Plan      models.StoryPlan
	Lang      models.Lang
	TargetKey string
	UserID    string
}

type ExpandBeatRepository struct{}

func NewExpandBeatRepository() *ExpandBeatRepository {
	return &ExpandBeatRepository{}
}

func (repository *ExpandBeatRepository) ExpandBeat(
	ctx context.Context, request ExpandBeatRequest,
) (*models.Beat, error) {
	span := sentry.StartSpan(ctx, "ExpandBeatRepository.ExpandBeat")
	defer span.Finish()

	span.SetData("request.userID", request.UserID)
	span.SetData("request.storyPlan.id", request.Plan.ID.String())
	span.SetData("request.storyPlan.slug", request.Plan.Slug)
	span.SetData("request.storyPlan.name", request.Plan.Name)
	span.SetData("request.storyPlan.lang", request.Plan.Lang)
	span.SetData("request.targetKey", request.TargetKey)
	span.SetData("request.logline", request.Logline)

	if !lo.ContainsBy(request.Plan.Beats, func(item models.BeatDefinition) bool {
		return item.Key == request.TargetKey
	}) {
		span.SetData("error", "target key not found in story plan beats")

		return nil, NewErrExpandBeatRepository(ErrUnknownTargetKey)
	}

	storyPlanPartialPrompt, err := StoryPlanToPrompt(request.Lang, request.Plan)
	if err != nil {
		span.SetData("storyPlan.toPrompt.error", err.Error())

		return nil, NewErrExpandBeatRepository(fmt.Errorf("parse story plan prompt: %w", err))
	}

	systemPrompt := new(strings.Builder)

	err = ExpandBeatPrompts.System.ExecuteTemplate(systemPrompt, request.Lang.String(), map[string]any{
		"StoryPlan": storyPlanPartialPrompt,
		"Plan":      request.Plan,
	})
	if err != nil {
		span.SetData("systemPrompt.error", err.Error())

		return nil, NewErrExpandBeatRepository(fmt.Errorf("parse system message: %w", err))
	}

	userPrompt1 := new(strings.Builder)

	err = ExpandBeatPrompts.Input1.ExecuteTemplate(userPrompt1, request.Lang.String(), request)
	if err != nil {
		span.SetData("userPrompt1.error", err.Error())

		return nil, NewErrExpandBeatRepository(fmt.Errorf("parse user message 1: %w", err))
	}

	userPrompt2 := new(strings.Builder)

	err = ExpandBeatPrompts.Input2.ExecuteTemplate(userPrompt2, request.Lang.String(), request)
	if err != nil {
		span.SetData("userPrompt2.error", err.Error())

		return nil, NewErrExpandBeatRepository(fmt.Errorf("parse user message 2: %w", err))
	}

	chatCompletion, err := lib.OpenAIClient(span.Context()).
		Chat.Completions.
		New(span.Context(), openai.ChatCompletionNewParams{
			Model:       config.Groq.Model,
			Temperature: param.NewOpt(expandBeatTemperature),
			User:        param.NewOpt(request.UserID),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemPrompt.String()),
				openai.UserMessage(userPrompt1.String()),
				repository.buildBeatsSheetResponse(request),
				openai.UserMessage(userPrompt2.String()),
			},
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
					JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
						Name:        "storyBeat",
						Description: openai.String(ExpandBeatDescriptions[request.Lang]),
						Schema:      ExpandBeatSchemas[request.Lang],
						Strict:      openai.Bool(true),
					},
				},
			},
		})
	if err != nil {
		span.SetData("chatCompletion.error", err.Error())

		return nil, NewErrExpandBeatRepository(err)
	}

	var beat models.Beat

	err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &beat)
	if err != nil {
		span.SetData("unmarshal.error", err.Error())

		return nil, NewErrExpandBeatRepository(err)
	}

	return &beat, nil
}

func (repository *ExpandBeatRepository) buildBeatsSheetResponse(
	request ExpandBeatRequest,
) openai.ChatCompletionMessageParamUnion {
	parts := make([]string, len(request.Beats))

	for i, beat := range request.Beats {
		parts[i] = fmt.Sprintf("%s\n%s", beat.Key, beat.Content)
	}

	return openai.AssistantMessage(strings.Join(parts, "\n\n"))
}
