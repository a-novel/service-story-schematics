package api

import (
	"context"
	"fmt"

	authapi "github.com/a-novel/service-authentication/api"

	"github.com/a-novel/service-story-schematics/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type ExpandLoglineService interface {
	ExpandLogline(ctx context.Context, request services.ExpandLoglineRequest) (*models.LoglineIdea, error)
}

func (api *API) ExpandLogline(ctx context.Context, req *codegen.LoglineIdea) (codegen.ExpandLoglineRes, error) {
	userID, err := authapi.RequireUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user ID: %w", err)
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
		return nil, fmt.Errorf("expand logline: %w", err)
	}

	return &codegen.LoglineIdea{
		Name:    logline.Name,
		Content: logline.Content,
		Lang:    codegen.Lang(logline.Lang),
	}, nil
}
