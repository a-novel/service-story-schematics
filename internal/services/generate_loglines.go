package services

import (
	"context"
	"errors"
	"github.com/getsentry/sentry-go"

	"github.com/google/uuid"

	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrGenerateLoglinesService = errors.New("GenerateLoglinesService.GenerateLoglines")

func NewErrGenerateLoglinesService(err error) error {
	return errors.Join(err, ErrGenerateLoglinesService)
}

type GenerateLoglinesSource interface {
	GenerateLoglines(ctx context.Context, request daoai.GenerateLoglinesRequest) ([]models.LoglineIdea, error)
}

type GenerateLoglinesRequest struct {
	Count  int
	Theme  string
	UserID uuid.UUID
	Lang   models.Lang
}

type GenerateLoglinesService struct {
	source GenerateLoglinesSource
}

func (service *GenerateLoglinesService) GenerateLoglines(
	ctx context.Context, request GenerateLoglinesRequest,
) ([]models.LoglineIdea, error) {
	span := sentry.StartSpan(ctx, "GenerateLoglinesService.GenerateLoglines")
	defer span.Finish()

	span.SetData("request.count", request.Count)
	span.SetData("request.theme", request.Theme)
	span.SetData("request.userID", request.UserID)
	span.SetData("request.lang", request.Lang)

	resp, err := service.source.GenerateLoglines(span.Context(), daoai.GenerateLoglinesRequest{
		Count:  request.Count,
		Theme:  request.Theme,
		UserID: request.UserID.String(),
		Lang:   request.Lang,
	})
	if err != nil {
		span.SetData("daoai.generateLoglines.error", err.Error())

		return nil, NewErrGenerateLoglinesService(err)
	}

	return resp, nil
}

func NewGenerateLoglinesService(source GenerateLoglinesSource) *GenerateLoglinesService {
	return &GenerateLoglinesService{source: source}
}
