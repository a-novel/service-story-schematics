package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/samber/lo"

	authapi "github.com/a-novel/service-authentication/api"

	"github.com/a-novel/service-story-schematics/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type SelectLoglineService interface {
	SelectLogline(ctx context.Context, request services.SelectLoglineRequest) (*models.Logline, error)
}

func (api *API) GetLogline(ctx context.Context, params codegen.GetLoglineParams) (codegen.GetLoglineRes, error) {
	span := sentry.StartSpan(ctx, "API.GetLogline")
	defer span.Finish()

	span.SetData("request.slug", params.Slug)
	span.SetData("request.id", params.ID)

	userID, err := authapi.RequireUserID(ctx)
	if err != nil {
		span.SetData("request.userID.err", err.Error())

		return nil, fmt.Errorf("get user ID: %w", err)
	}

	span.SetData("request.userID", userID)

	logline, err := api.SelectLoglineService.SelectLogline(span.Context(), services.SelectLoglineRequest{
		UserID: userID,
		Slug:   lo.Ternary(params.Slug.IsSet(), lo.ToPtr(models.Slug(params.Slug.Value)), nil),
		ID:     lo.Ternary(params.ID.IsSet(), lo.ToPtr(uuid.UUID(params.ID.Value)), nil),
	})

	switch {
	case errors.Is(err, dao.ErrLoglineNotFound):
		span.SetData("service.err", err.Error())

		return &codegen.NotFoundError{Error: err.Error()}, nil
	case err != nil:
		span.SetData("service.err", err.Error())

		return nil, fmt.Errorf("get logline: %w", err)
	}

	return &codegen.Logline{
		ID:        codegen.LoglineID(logline.ID),
		UserID:    codegen.UserID(logline.UserID),
		Slug:      codegen.Slug(logline.Slug),
		Name:      logline.Name,
		Content:   logline.Content,
		Lang:      codegen.Lang(logline.Lang),
		CreatedAt: logline.CreatedAt,
	}, nil
}
