package api

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
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
	span := sentry.StartSpan(ctx, "API.GetBeatsSheets")
	defer span.Finish()

	span.SetData("request.loglineID", params.LoglineID)
	span.SetData("request.limit", params.Limit)
	span.SetData("request.offset", params.Offset)

	userID, err := authapi.RequireUserID(ctx)
	if err != nil {
		span.SetData("request.userID.err", err.Error())

		return nil, fmt.Errorf("get user ID: %w", err)
	}

	span.SetData("request.userID", userID)

	beatsSheets, err := api.ListBeatsSheetsService.ListBeatsSheets(span.Context(), services.ListBeatsSheetsRequest{
		UserID:    userID,
		LoglineID: uuid.UUID(params.LoglineID),
		Limit:     params.Limit.Value,
		Offset:    params.Offset.Value,
	})
	if err != nil {
		span.SetData("service.err", err.Error())

		return nil, fmt.Errorf("list beats sheets: %w", err)
	}

	span.SetData("beatsSheets.count", len(beatsSheets))

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
