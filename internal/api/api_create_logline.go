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

type CreateLoglineService interface {
	CreateLogline(ctx context.Context, request services.CreateLoglineRequest) (*models.Logline, error)
}

func (api *API) CreateLogline(ctx context.Context, req *codegen.CreateLoglineForm) (codegen.CreateLoglineRes, error) {
	span := sentry.StartSpan(ctx, "API.CreateLogline")
	defer span.Finish()

	span.SetData("request.slug", req.GetSlug())
	span.SetData("request.name", req.GetName())
	span.SetData("request.lang", req.GetLang())

	userID, err := authPkg.RequireUserID(ctx)
	if err != nil {
		span.SetData("request.userID.err", err.Error())

		return nil, fmt.Errorf("get user ID: %w", err)
	}

	span.SetData("request.userID", userID)

	logline, err := api.CreateLoglineService.CreateLogline(span.Context(), services.CreateLoglineRequest{
		UserID:  userID,
		Slug:    models.Slug(req.GetSlug()),
		Name:    req.GetName(),
		Content: req.GetContent(),
		Lang:    models.Lang(req.GetLang()),
	})
	if err != nil {
		span.SetData("service.err", err.Error())

		return nil, fmt.Errorf("create logline: %w", err)
	}

	span.SetData("logline.id", logline.ID.String())

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
