package daoai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/packages/param"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"

	"github.com/a-novel/golib/otel"

	"github.com/a-novel/service-story-schematics/internal/daoai/prompts"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/config"
	storyplanmodel "github.com/a-novel/service-story-schematics/models/story_plan"
)

var ExpandBeatPrompts = struct {
	System *template.Template
	Input1 *template.Template
	Input2 *template.Template
}{
	System: template.Must(template.New("").Parse(prompts.ExpandBeat.System)),
	Input1: template.Must(template.New("").Parse(prompts.ExpandBeat.Input1)),
	Input2: template.Must(template.New("").Parse(prompts.ExpandBeat.Input2)),
}

var ErrUnknownTargetKey = errors.New("unknown target key")

type ExpandBeatRequest struct {
	Logline   string
	Beats     []models.Beat
	Plan      *storyplanmodel.Plan
	Lang      models.Lang
	TargetKey string
	UserID    string
}

type ExpandBeatRepository struct {
	config *config.OpenAI
}

func NewExpandBeatRepository(config *config.OpenAI) *ExpandBeatRepository {
	return &ExpandBeatRepository{config: config}
}

func (repository *ExpandBeatRepository) ExpandBeat(
	ctx context.Context, request ExpandBeatRequest,
) (*models.Beat, error) {
	ctx, span := otel.Tracer().Start(ctx, "daoai.ExpandBeat")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.userID", request.UserID),
		attribute.String("request.Lang", request.Lang.String()),
		attribute.String("request.targetKey", request.TargetKey),
		attribute.String("request.logline", request.Logline),
	)

	targetBeat, err := request.Plan.GetBeat(request.TargetKey)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get target beat: %w", err))
	}

	systemPrompt := new(strings.Builder)

	err = ExpandBeatPrompts.System.Execute(systemPrompt, map[string]any{
		"PlanName": request.Plan.Metadata.Name,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("parse system message: %w", err))
	}

	userPrompt1 := new(strings.Builder)

	err = ExpandBeatPrompts.Input1.Execute(userPrompt1, request)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("parse user message 1: %w", err))
	}

	userPrompt2 := new(strings.Builder)

	err = ExpandBeatPrompts.Input2.Execute(userPrompt2, request)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("parse user message 2: %w", err))
	}

	chatCompletion, err := repository.config.Client().
		Chat.Completions.
		New(ctx, openai.ChatCompletionNewParams{
			Model: repository.config.Model,
			User:  param.NewOpt(request.UserID),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemPrompt.String()),
				openai.UserMessage(userPrompt1.String()),
				repository.buildBeatsSheetResponse(request),
				openai.UserMessage(ForceNextAnswerLocale(request.Lang, userPrompt2.String())),
			},
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
					JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
						Name: "storyBeat",
						Description: openai.String(
							fmt.Sprintf("The expanded version of the '%s' beat.", request.TargetKey),
						),
						Schema: targetBeat.OutputSchema(),
						Strict: openai.Bool(true),
					},
				},
			},
		})
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	var beat models.Beat

	err = json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &beat)
	if err != nil {
		return nil, otel.ReportError(span, err)
	}

	return otel.ReportSuccess(span, &beat), nil
}

func (repository *ExpandBeatRepository) buildBeatsSheetResponse(
	request ExpandBeatRequest,
) openai.ChatCompletionMessageParamUnion {
	return openai.AssistantMessage(strings.Join(lo.Map(request.Beats, func(item models.Beat, _ int) string {
		return item.String()
	}), "\n\n"))
}
