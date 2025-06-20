package api

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"

	authapi "github.com/a-novel/service-authentication/api"

	"github.com/a-novel/service-story-schematics/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type ListBeatsSheetsService interface {
	ListBeatsSheets(ctx context.Context, request services.ListBeatsSheetsRequest) ([]*models.BeatsSheetPreview, error)
}

func (api *API) GetBeatsSheets(
	ctx context.Context, params codegen.GetBeatsSheetsParams,
) (codegen.GetBeatsSheetsRes, error) {
	userID, err := authapi.RequireUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user ID: %w", err)
	}

	beatsSheets, err := api.ListBeatsSheetsService.ListBeatsSheets(ctx, services.ListBeatsSheetsRequest{
		UserID:    userID,
		LoglineID: uuid.UUID(params.LoglineID),
		Limit:     params.Limit.Value,
		Offset:    params.Offset.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("list beats sheets: %w", err)
	}

	res := codegen.GetBeatsSheetsOKApplicationJSON(
		lo.Map(beatsSheets, func(item *models.BeatsSheetPreview, _ int) codegen.BeatsSheetPreview {
			return codegen.BeatsSheetPreview{
				ID:        codegen.BeatsSheetID(item.ID),
				Lang:      codegen.Lang(item.Lang),
				CreatedAt: item.CreatedAt,
			}
		}),
	)

	return &res, nil
}
