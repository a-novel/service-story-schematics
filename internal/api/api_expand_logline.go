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

type ExpandLoglineService interface {
	ExpandLogline(ctx context.Context, request services.ExpandLoglineRequest) (*models.LoglineIdea, error)
}

func (api *API) ExpandLogline(ctx context.Context, req *apimodels.LoglineIdea) (apimodels.ExpandLoglineRes, error) {
	ctx, span := otel.Tracer().Start(ctx, "api.ExpandLogline")
	defer span.End()

	userID, err := authpkg.RequireUserID(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get user ID: %w", err))
	}

	logline, err := api.ExpandLoglineService.ExpandLogline(ctx, services.ExpandLoglineRequest{
		Logline: models.LoglineIdea{
			Name:    req.GetName(),
			Content: req.GetContent(),
			Lang:    models.Lang(req.GetLang()),
		},
		UserID: userID,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("expand logline: %w", err))
	}

	return otel.ReportSuccess(span, &apimodels.LoglineIdea{
		Name:    logline.Name,
		Content: logline.Content,
		Lang:    apimodels.Lang(logline.Lang),
	}), nil
}
