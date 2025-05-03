package services

import (
	"errors"

	"github.com/google/uuid"

	"github.com/a-novel-kit/context"

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
}

type GenerateLoglinesService struct {
	source GenerateLoglinesSource
}

func (service *GenerateLoglinesService) GenerateLoglines(
	ctx context.Context, request GenerateLoglinesRequest,
) ([]models.LoglineIdea, error) {
	resp, err := service.source.GenerateLoglines(ctx, daoai.GenerateLoglinesRequest{
		Count:  request.Count,
		Theme:  request.Theme,
		UserID: request.UserID.String(),
	})
	if err != nil {
		return nil, NewErrGenerateLoglinesService(err)
	}

	return resp, nil
}

func NewGenerateLoglinesService(source GenerateLoglinesSource) *GenerateLoglinesService {
	return &GenerateLoglinesService{source: source}
}
