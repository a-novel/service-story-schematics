package daoai

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/samber/lo"

	"github.com/a-novel-kit/golm"

	"github.com/a-novel/service-story-schematics/config/prompts"
	"github.com/a-novel/service-story-schematics/models"
)

const expandBeatTemperature = 0.8

var ExpandBeatPrompts = struct {
	System *template.Template
	Input1 *template.Template
	Input2 *template.Template
}{
	System: template.Must(template.New("EN").Parse(prompts.Config.En.ExpandBeat.System)),
	Input1: template.Must(template.New("EN").Parse(prompts.Config.En.ExpandBeat.Input1)),
	Input2: template.Must(template.New("EN").Parse(prompts.Config.En.ExpandBeat.Input2)),
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

func (repository *ExpandBeatRepository) buildBeatsSheetResponse(request ExpandBeatRequest) golm.AssistantMessage {
	parts := make([]string, len(request.Beats))

	for i, beat := range request.Beats {
		parts[i] = fmt.Sprintf("%s\n%s", beat.Key, beat.Content)
	}

	return golm.NewAssistantMessage(strings.Join(parts, "\n\n"))
}

func (repository *ExpandBeatRepository) ExpandBeat(
	ctx context.Context, request ExpandBeatRequest,
) (*models.Beat, error) {
	if !lo.ContainsBy(request.Plan.Beats, func(item models.BeatDefinition) bool {
		return item.Key == request.TargetKey
	}) {
		return nil, NewErrExpandBeatRepository(ErrUnknownTargetKey)
	}

	storyPlanPrompt, err := StoryPlanToPrompt("EN", request.Plan)
	if err != nil {
		return nil, NewErrExpandBeatRepository(fmt.Errorf("parse story plan prompt: %w", err))
	}

	systemPrompt, err := golm.NewSystemMessageT(ExpandBeatPrompts.System, "EN", map[string]any{
		"StoryPlan": storyPlanPrompt,
		"Plan":      request.Plan,
	})
	if err != nil {
		return nil, NewErrExpandBeatRepository(fmt.Errorf("parse system message: %w", err))
	}

	userPrompt1, err := golm.NewUserMessageT(ExpandBeatPrompts.Input1, "EN", request)
	if err != nil {
		return nil, NewErrExpandBeatRepository(fmt.Errorf("parse user message 1: %w", err))
	}

	userPrompt2, err := golm.NewUserMessageT(ExpandBeatPrompts.Input2, "EN", request)
	if err != nil {
		return nil, NewErrExpandBeatRepository(fmt.Errorf("parse user message 2: %w", err))
	}

	chat := golm.Context(ctx)

	chat.SetSystem(systemPrompt)
	chat.PushHistory(
		userPrompt1,
		repository.buildBeatsSheetResponse(request),
	)

	params := golm.CompletionParams{
		Temperature: lo.ToPtr(expandBeatTemperature),
		User:        request.UserID,
	}

	var beat models.Beat

	if err = chat.CompletionJSON(ctx, userPrompt2, params, &beat); err != nil {
		return nil, NewErrExpandBeatRepository(err)
	}

	return &beat, nil
}

func NewExpandBeatRepository() *ExpandBeatRepository {
	return &ExpandBeatRepository{}
}
