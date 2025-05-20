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

const expandLoglineTemperature = 0.8

var ExpandLoglinePrompts = struct {
	System *template.Template
}{
	System: template.Must(template.New("EN").Parse(prompts.Config.En.ExpandLogline)),
}

var ErrExpandLoglineRepository = errors.New("ExpandLoglineRepository.ExpandLogline")

func NewErrExpandLoglineRepository(err error) error {
	return errors.Join(err, ErrExpandLoglineRepository)
}

type ExpandLoglineRequest struct {
	Logline string
	UserID  string
}

type ExpandLoglineRepository struct{}

func (repository *ExpandLoglineRepository) ExpandLogline(
	ctx context.Context, request ExpandLoglineRequest,
) (*models.LoglineIdea, error) {
	chat := golm.Context(ctx)

	systemMessage, err := golm.NewSystemMessageT(ExpandLoglinePrompts.System, "EN", nil)
	if err != nil {
		return nil, NewErrExpandLoglineRepository(fmt.Errorf("parse system message: %w", err))
	}

	chat.SetSystem(systemMessage)

	requestMessage := golm.NewUserMessage(request.Logline)

	params := golm.CompletionParams{
		Temperature: lo.ToPtr(expandLoglineTemperature),
		User:        request.UserID,
	}

	var logline models.LoglineIdea

	if err = chat.CompletionJSON(ctx, requestMessage, params, &logline); err != nil {
		return nil, NewErrExpandLoglineRepository(err)
	}

	return &logline, nil
}

func NewExpandLoglineRepository() *ExpandLoglineRepository {
	return &ExpandLoglineRepository{}
}
