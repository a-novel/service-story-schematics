package api

import (
	"fmt"

	authapi "github.com/a-novel/service-authentication/api"

	"github.com/a-novel-kit/context"

	"github.com/a-novel/service-story-schematics/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type CreateLoglineService interface {
	CreateLogline(ctx context.Context, request services.CreateLoglineRequest) (*models.Logline, error)
}

func (api *API) CreateLogline(ctx context.Context, req *codegen.CreateLoglineForm) (codegen.CreateLoglineRes, error) {
	userID, err := authapi.RequireUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user ID: %w", err)
	}

	logline, err := api.CreateLoglineService.CreateLogline(ctx, services.CreateLoglineRequest{
		UserID:  userID,
		Slug:    models.Slug(req.GetSlug()),
		Name:    req.GetName(),
		Content: req.GetContent(),
	})
	if err != nil {
		return nil, fmt.Errorf("create logline: %w", err)
	}

	return &codegen.Logline{
		ID:        codegen.LoglineID(logline.ID),
		UserID:    codegen.UserID(logline.UserID),
		Slug:      codegen.Slug(logline.Slug),
		Name:      logline.Name,
		Content:   logline.Content,
		CreatedAt: logline.CreatedAt,
	}, nil
}
