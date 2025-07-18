package api

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/a-novel/golib/otel"
	authpkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

type GenerateLoglinesService interface {
	GenerateLoglines(ctx context.Context, request services.GenerateLoglinesRequest) ([]models.LoglineIdea, error)
}

func (api *API) GenerateLoglines(
	ctx context.Context, req *apimodels.GenerateLoglinesForm,
) (apimodels.GenerateLoglinesRes, error) {
	ctx, span := otel.Tracer().Start(ctx, "api.GenerateLoglines")
	defer span.End()

	userID, err := authpkg.RequireUserID(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get user ID: %w", err))
	}

	loglines, err := api.GenerateLoglinesService.GenerateLoglines(ctx, services.GenerateLoglinesRequest{
		Count:  req.GetCount(),
		Theme:  req.GetTheme(),
		UserID: userID,
		Lang:   models.Lang(req.GetLang()),
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("generate loglines: %w", err))
	}

	res := apimodels.GenerateLoglinesOKApplicationJSON(
		lo.Map(loglines, func(item models.LoglineIdea, _ int) apimodels.LoglineIdea {
			return apimodels.LoglineIdea{
				Name:    item.Name,
				Content: item.Content,
				Lang:    apimodels.Lang(item.Lang),
			}
		}),
	)

	return otel.ReportSuccess(span, &res), nil
}
