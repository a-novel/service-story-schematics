package services

import (
	"context"
	"errors"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"

	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrExpandLoglineService = errors.New("ExpandLoglineService.ExpandLogline")

func NewErrExpandLoglineService(err error) error {
	return errors.Join(err, ErrExpandLoglineService)
}

type ExpandLoglineSource interface {
	ExpandLogline(ctx context.Context, request daoai.ExpandLoglineRequest) (*models.LoglineIdea, error)
}

type ExpandLoglineRequest struct {
	Logline models.LoglineIdea
	UserID  uuid.UUID
}

type ExpandLoglineService struct {
	source ExpandLoglineSource
}

func (service *ExpandLoglineService) ExpandLogline(
	ctx context.Context, request ExpandLoglineRequest,
) (*models.LoglineIdea, error) {
	span := sentry.StartSpan(ctx, "ExpandLoglineService.ExpandLogline")
	defer span.Finish()

	span.SetData("request.logline.name", request.Logline.Name)
	span.SetData("request.logline.lang", request.Logline.Lang)
	span.SetData("request.userID", request.UserID)

	resp, err := service.source.ExpandLogline(span.Context(), daoai.ExpandLoglineRequest{
		Logline: request.Logline.Name + "\n\n" + request.Logline.Content,
		UserID:  request.UserID.String(),
		Lang:    request.Logline.Lang,
	})
	if err != nil {
		span.SetData("service.error", err.Error())

		return nil, NewErrExpandLoglineService(err)
	}

	return resp, nil
}

func NewExpandLoglineService(source ExpandLoglineSource) *ExpandLoglineService {
	return &ExpandLoglineService{source: source}
}
