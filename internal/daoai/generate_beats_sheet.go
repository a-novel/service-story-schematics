package daoai

import (
	"errors"
	"fmt"
	"text/template"

	"github.com/samber/lo"

	"github.com/a-novel-kit/context"
	"github.com/a-novel-kit/golm"

	"github.com/a-novel/story-schematics/config/prompts"
	"github.com/a-novel/story-schematics/internal/lib"
	"github.com/a-novel/story-schematics/models"
)

const generateBeatsSheetTemperature = 0.8

var GenerateBeatsSheetPrompts = struct {
	System *template.Template
}{
	System: template.Must(template.New("EN").Parse(prompts.Config.En.GenerateBeatsSheet)),
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
}

type GenerateBeatsSheetRepository struct{}

func (repository *GenerateBeatsSheetRepository) GenerateBeatsSheet(
	ctx context.Context, request GenerateBeatsSheetRequest,
) ([]models.Beat, error) {
	chat := golm.Context(ctx)

	storyPlanPrompt, err := StoryPlanToPrompt("EN", request.Plan)
	if err != nil {
		return nil, NewErrGenerateBeatsSheetRepository(fmt.Errorf("parse story plan prompt: %w", err))
	}

	systemPrompt, err := golm.NewSystemMessageT(GenerateBeatsSheetPrompts.System, "EN", map[string]any{
		"StoryPlan": storyPlanPrompt,
		"Plan":      request.Plan,
	})
	if err != nil {
		return nil, NewErrGenerateBeatsSheetRepository(fmt.Errorf("parse system message: %w", err))
	}

	chat.SetSystem(systemPrompt)

	requestMessage := golm.NewUserMessage(request.Logline)

	params := golm.CompletionParams{
		Temperature: lo.ToPtr(generateBeatsSheetTemperature),
		User:        request.UserID,
	}

	var beats struct {
		Beats []models.Beat `json:"beats"`
	}

	if err = chat.CompletionJSON(ctx, requestMessage, params, &beats); err != nil {
		return nil, NewErrGenerateBeatsSheetRepository(err)
	}

	if err = lib.CheckStoryPlan(beats.Beats, request.Plan.Beats); err != nil {
		return nil, NewErrGenerateBeatsSheetRepository(errors.Join(err, ErrInvalidBeatSheet))
	}

	return beats.Beats, nil
}

func NewGenerateBeatsSheetRepository() *GenerateBeatsSheetRepository {
	return &GenerateBeatsSheetRepository{}
}
