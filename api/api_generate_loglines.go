package api

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"

	"github.com/samber/lo"

	authapi "github.com/a-novel/service-authentication/api"

	"github.com/a-novel/service-story-schematics/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type GenerateLoglinesService interface {
	GenerateLoglines(ctx context.Context, request services.GenerateLoglinesRequest) ([]models.LoglineIdea, error)
}

func (api *API) GenerateLoglines(
	ctx context.Context, req *codegen.GenerateLoglinesForm,
) (codegen.GenerateLoglinesRes, error) {
	span := sentry.StartSpan(ctx, "API.GenerateLoglines")
	defer span.Finish()

	span.SetData("request.count", req.GetCount())
	span.SetData("request.theme", req.GetTheme())
	span.SetData("request.lang", req.GetLang())

	userID, err := authapi.RequireUserID(ctx)
	if err != nil {
		span.SetData("request.userID.err", err.Error())

		return nil, fmt.Errorf("get user ID: %w", err)
	}

	span.SetData("request.userID", userID)

	loglines, err := api.GenerateLoglinesService.GenerateLoglines(span.Context(), services.GenerateLoglinesRequest{
		Count:  req.GetCount(),
		Theme:  req.GetTheme(),
		UserID: userID,
		Lang:   models.Lang(req.GetLang()),
	})
	if err != nil {
		span.SetData("service.err", err.Error())

		return nil, fmt.Errorf("generate loglines: %w", err)
	}

	res := codegen.GenerateLoglinesOKApplicationJSON(
		lo.Map(loglines, func(item models.LoglineIdea, _ int) codegen.LoglineIdea {
			return codegen.LoglineIdea{
				Name:    item.Name,
				Content: item.Content,
				Lang:    codegen.Lang(item.Lang),
			}
		}),
	)

	return &res, nil
}
