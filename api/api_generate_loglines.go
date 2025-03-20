package api

import (
	"fmt"

	"github.com/samber/lo"

	authapi "github.com/a-novel/authentication/api"

	"github.com/a-novel-kit/context"

	"github.com/a-novel/story-schematics/api/codegen"
	"github.com/a-novel/story-schematics/internal/services"
	"github.com/a-novel/story-schematics/models"
)

type GenerateLoglinesService interface {
	GenerateLoglines(ctx context.Context, request services.GenerateLoglinesRequest) ([]models.LoglineIdea, error)
}

func (api *API) GenerateLoglines(
	ctx context.Context, req *codegen.GenerateLoglinesForm,
) (codegen.GenerateLoglinesRes, error) {
	userID, err := authapi.RequireUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user ID: %w", err)
	}

	loglines, err := api.GenerateLoglinesService.GenerateLoglines(ctx, services.GenerateLoglinesRequest{
		Count:  req.GetCount(),
		Theme:  req.GetTheme(),
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("generate loglines: %w", err)
	}

	res := codegen.GenerateLoglinesOKApplicationJSON(
		lo.Map(loglines, func(item models.LoglineIdea, _ int) codegen.LoglineIdea {
			return codegen.LoglineIdea{
				Name:    item.Name,
				Content: item.Content,
			}
		}),
	)

	return &res, nil
}
