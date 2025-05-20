package daoai

import (
	"context"
	"errors"
	"fmt"
	"text/template"

	"github.com/samber/lo"

	"github.com/a-novel-kit/golm"

	"github.com/a-novel/service-story-schematics/config/prompts"
	"github.com/a-novel/service-story-schematics/models"
)

const generateLoglineTemperature = 1.0

var GenerateLoglinesPrompts = struct {
	Themed *template.Template
	Random *template.Template
}{
	Themed: template.Must(template.New("EN").Parse(prompts.Config.En.GenerateLogline.System.Themed)),
	Random: template.Must(template.New("EN").Parse(prompts.Config.En.GenerateLogline.System.Random)),
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
	chat := golm.Context(ctx)

	var (
		err            error
		systemMessage  *golm.SystemMessage
		requestMessage golm.UserMessage
	)

	if request.Theme != "" {
		systemMessage, err = golm.NewSystemMessageT(GenerateLoglinesPrompts.Themed, "EN", request)
		requestMessage = golm.NewUserMessage(request.Theme)
	} else {
		requestMessage, err = golm.NewUserMessageT(GenerateLoglinesPrompts.Random, "EN", request)
	}

	if err != nil {
		return nil, NewErrGenerateLoglinesRepository(fmt.Errorf("parse system message: %w", err))
	}

	chat.SetSystem(systemMessage)

	params := golm.CompletionParams{
		Temperature: lo.ToPtr(generateLoglineTemperature),
		User:        request.UserID,
	}

	var loglines struct {
		Loglines []models.LoglineIdea `json:"loglines"`
	}

	if err = chat.CompletionJSON(ctx, requestMessage, params, &loglines); err != nil {
		return nil, NewErrGenerateLoglinesRepository(err)
	}

	return loglines.Loglines, nil
}

func NewGenerateLoglinesRepository() *GenerateLoglinesRepository {
	return &GenerateLoglinesRepository{}
}
