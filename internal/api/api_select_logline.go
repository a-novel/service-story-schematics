package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/a-novel/golib/otel"
	authpkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

type SelectLoglineService interface {
	SelectLogline(ctx context.Context, request services.SelectLoglineRequest) (*models.Logline, error)
}

func (api *API) GetLogline(ctx context.Context, params apimodels.GetLoglineParams) (apimodels.GetLoglineRes, error) {
	ctx, span := otel.Tracer().Start(ctx, "api.GetLogline")
	defer span.End()

	userID, err := authpkg.RequireUserID(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get user ID: %w", err))
	}

	logline, err := api.SelectLoglineService.SelectLogline(ctx, services.SelectLoglineRequest{
		UserID: userID,
		Slug:   lo.Ternary(params.Slug.IsSet(), lo.ToPtr(models.Slug(params.Slug.Value)), nil),
		ID:     lo.Ternary(params.ID.IsSet(), lo.ToPtr(uuid.UUID(params.ID.Value)), nil),
	})

	switch {
	case errors.Is(err, dao.ErrLoglineNotFound):
		_ = otel.ReportError(span, err)

		return &apimodels.NotFoundError{Error: err.Error()}, nil
	case err != nil:
		_ = otel.ReportError(span, err)

		return nil, fmt.Errorf("get logline: %w", err)
	}

	return otel.ReportSuccess(span, &apimodels.Logline{
		ID:        apimodels.LoglineID(logline.ID),
		UserID:    apimodels.UserID(logline.UserID),
		Slug:      apimodels.Slug(logline.Slug),
		Name:      logline.Name,
		Content:   logline.Content,
		Lang:      apimodels.Lang(logline.Lang),
		CreatedAt: logline.CreatedAt,
	}), nil
}
