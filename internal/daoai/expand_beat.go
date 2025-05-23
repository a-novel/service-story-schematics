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

	"github.com/samber/lo"

	"github.com/a-novel/service-story-schematics/config/prompts"
	"github.com/a-novel/service-story-schematics/models"
)

const expandBeatTemperature = 0.8

var ExpandBeatPrompts = struct {
	System *template.Template
	Input1 *template.Template
	Input2 *template.Template
}{
	System: template.Must(template.New(string(models.LangEN)).Parse(prompts.Config.En.ExpandBeat.System)),
	Input1: template.Must(template.New(string(models.LangEN)).Parse(prompts.Config.En.ExpandBeat.Input1)),
	Input2: template.Must(template.New(string(models.LangEN)).Parse(prompts.Config.En.ExpandBeat.Input2)),
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
	TargetKey string
	UserID    string
}

type ExpandBeatRepository struct{}

func (repository *ExpandBeatRepository) buildBeatsSheetResponse(
	request ExpandBeatRequest,
) openai.ChatCompletionMessageParamUnion {
	parts := make([]string, len(request.Beats))

	for i, beat := range request.Beats {
		parts[i] = fmt.Sprintf("%s\n%s", beat.Key, beat.Content)
	}

	return openai.AssistantMessage(strings.Join(parts, "\n\n"))
}

func (repository *ExpandBeatRepository) ExpandBeat(
	ctx context.Context, request ExpandBeatRequest,
) (*models.Beat, error) {
	if !lo.ContainsBy(request.Plan.Beats, func(item models.BeatDefinition) bool {
		return item.Key == request.TargetKey
	}) {
		return nil, NewErrExpandBeatRepository(ErrUnknownTargetKey)
	}

	storyPlanPartialPrompt, err := StoryPlanToPrompt("EN", request.Plan)
	if err != nil {
		return nil, NewErrExpandBeatRepository(fmt.Errorf("parse story plan prompt: %w", err))
	}

	systemPrompt := new(strings.Builder)

	err = ExpandBeatPrompts.System.ExecuteTemplate(systemPrompt, string(models.LangEN), map[string]any{
		"StoryPlan": storyPlanPartialPrompt,
		"Plan":      request.Plan,
	})
	if err != nil {
		return nil, NewErrExpandBeatRepository(fmt.Errorf("parse system message: %w", err))
	}

	userPrompt1 := new(strings.Builder)

	err = ExpandBeatPrompts.Input1.ExecuteTemplate(userPrompt1, string(models.LangEN), request)
	if err != nil {
		return nil, NewErrExpandBeatRepository(fmt.Errorf("parse user message 1: %w", err))
	}

	userPrompt2 := new(strings.Builder)

	err = ExpandBeatPrompts.Input2.ExecuteTemplate(userPrompt2, string(models.LangEN), request)
	if err != nil {
		return nil, NewErrExpandBeatRepository(fmt.Errorf("parse user message 2: %w", err))
	}

	chatCompletion, err := lib.OpenAIClient(ctx).Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
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
					Description: openai.String(schemas.Config.En.Beat.Description),
					Schema:      schemas.Config.En.Beat.Schema,
					Strict:      openai.Bool(true),
				},
			},
		},
	})
	if err != nil {
		return nil, NewErrExpandBeatRepository(err)
	}

	var beat models.Beat

	if err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &beat); err != nil {
		return nil, NewErrExpandBeatRepository(err)
	}

	return &beat, nil
}

func NewExpandBeatRepository() *ExpandBeatRepository {
	return &ExpandBeatRepository{}
}
