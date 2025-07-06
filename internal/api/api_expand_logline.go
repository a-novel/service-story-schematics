package api

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"

	authPkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type ExpandLoglineService interface {
	ExpandLogline(ctx context.Context, request services.ExpandLoglineRequest) (*models.LoglineIdea, error)
}

func (api *API) ExpandLogline(ctx context.Context, req *codegen.LoglineIdea) (codegen.ExpandLoglineRes, error) {
	span := sentry.StartSpan(ctx, "API.ExpandLogline")
	defer span.Finish()

	span.SetData("request.name", req.GetName())
	span.SetData("request.lang", req.GetLang())

	userID, err := authPkg.RequireUserID(ctx)
	if err != nil {
		span.SetData("request.userID.err", err.Error())

		return nil, fmt.Errorf("get user ID: %w", err)
	}

	span.SetData("request.userID", userID)

	logline, err := api.ExpandLoglineService.ExpandLogline(span.Context(), services.ExpandLoglineRequest{
		Logline: models.LoglineIdea{
			Name:    req.GetName(),
			Content: req.GetContent(),
			Lang:    models.Lang(req.GetLang()),
		},
		UserID: userID,
	})
	if err != nil {
		span.SetData("service.err", err.Error())

		return nil, fmt.Errorf("expand logline: %w", err)
	}

	return &codegen.LoglineIdea{
		Name:    logline.Name,
		Content: logline.Content,
		Lang:    codegen.Lang(logline.Lang),
	}, nil
}
