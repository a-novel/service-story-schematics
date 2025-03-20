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

type ListLoglinesService interface {
	ListLoglines(ctx context.Context, request services.ListLoglinesRequest) ([]*models.LoglinePreview, error)
}

func (api *API) GetLoglines(ctx context.Context, params codegen.GetLoglinesParams) (codegen.GetLoglinesRes, error) {
	userID, err := authapi.RequireUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user ID: %w", err)
	}

	loglines, err := api.ListLoglinesService.ListLoglines(ctx, services.ListLoglinesRequest{
		UserID: userID,
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("list loglines: %w", err)
	}

	res := codegen.GetLoglinesOKApplicationJSON(
		lo.Map(loglines, func(item *models.LoglinePreview, _ int) codegen.LoglinePreview {
			return codegen.LoglinePreview{
				Slug:      codegen.Slug(item.Slug),
				Name:      item.Name,
				Content:   item.Content,
				CreatedAt: item.CreatedAt,
			}
		}),
	)

	return &res, nil
}
