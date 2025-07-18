package api

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/a-novel/golib/otel"
	authpkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

type ListBeatsSheetsService interface {
	ListBeatsSheets(ctx context.Context, request services.ListBeatsSheetsRequest) ([]*models.BeatsSheetPreview, error)
}

func (api *API) GetBeatsSheets(
	ctx context.Context, params apimodels.GetBeatsSheetsParams,
) (apimodels.GetBeatsSheetsRes, error) {
	ctx, span := otel.Tracer().Start(ctx, "api.GetBeatsSheets")
	defer span.End()

	userID, err := authpkg.RequireUserID(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get user ID: %w", err))
	}

	beatsSheets, err := api.ListBeatsSheetsService.ListBeatsSheets(ctx, services.ListBeatsSheetsRequest{
		UserID:    userID,
		LoglineID: uuid.UUID(params.LoglineID),
		Limit:     params.Limit.Value,
		Offset:    params.Offset.Value,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("list beats sheets: %w", err))
	}

	res := apimodels.GetBeatsSheetsOKApplicationJSON(
		lo.Map(beatsSheets, func(item *models.BeatsSheetPreview, _ int) apimodels.BeatsSheetPreview {
			return apimodels.BeatsSheetPreview{
				ID:        apimodels.BeatsSheetID(item.ID),
				Lang:      apimodels.Lang(item.Lang),
				CreatedAt: item.CreatedAt,
			}
		}),
	)

	return otel.ReportSuccess(span, &res), nil
}
