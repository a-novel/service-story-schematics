package api

import (
	"context"
	"fmt"

	"github.com/a-novel/golib/otel"
	authpkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

type CreateLoglineService interface {
	CreateLogline(ctx context.Context, request services.CreateLoglineRequest) (*models.Logline, error)
}

func (api *API) CreateLogline(
	ctx context.Context, req *apimodels.CreateLoglineForm,
) (apimodels.CreateLoglineRes, error) {
	ctx, span := otel.Tracer().Start(ctx, "api.CreateLogline")
	defer span.End()

	userID, err := authpkg.RequireUserID(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get user ID: %w", err))
	}

	logline, err := api.CreateLoglineService.CreateLogline(ctx, services.CreateLoglineRequest{
		UserID:  userID,
		Slug:    models.Slug(req.GetSlug()),
		Name:    req.GetName(),
		Content: req.GetContent(),
		Lang:    models.Lang(req.GetLang()),
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("create logline: %w", err))
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
