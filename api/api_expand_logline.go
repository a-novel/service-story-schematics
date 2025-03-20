package api

import (
	"fmt"

	authapi "github.com/a-novel/authentication/api"

	"github.com/a-novel-kit/context"

	"github.com/a-novel/story-schematics/api/codegen"
	"github.com/a-novel/story-schematics/internal/services"
	"github.com/a-novel/story-schematics/models"
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
		},
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("expand logline: %w", err)
	}

	return &codegen.LoglineIdea{
		Name:    logline.Name,
		Content: logline.Content,
	}, nil
}
