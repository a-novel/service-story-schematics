package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/codes"

	"github.com/a-novel/golib/otel"
	authpkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

type SelectBeatsSheetService interface {
	SelectBeatsSheet(ctx context.Context, request services.SelectBeatsSheetRequest) (*models.BeatsSheet, error)
}

func (api *API) GetBeatsSheet(
	ctx context.Context, params apimodels.GetBeatsSheetParams,
) (apimodels.GetBeatsSheetRes, error) {
	ctx, span := otel.Tracer().Start(ctx, "api.GetBeatsSheet")
	defer span.End()

	userID, err := authpkg.RequireUserID(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get user ID: %w", err))
	}

	beatsSheet, err := api.SelectBeatsSheetService.SelectBeatsSheet(ctx, services.SelectBeatsSheetRequest{
		BeatsSheetID: uuid.UUID(params.BeatsSheetID),
		UserID:       userID,
	})

	switch {
	case errors.Is(err, dao.ErrBeatsSheetNotFound), errors.Is(err, dao.ErrLoglineNotFound):
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return &apimodels.NotFoundError{Error: err.Error()}, nil
	case err != nil:
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return nil, fmt.Errorf("get beats sheet: %w", err)
	}

	return otel.ReportSuccess(span, &apimodels.BeatsSheet{
		ID:        apimodels.BeatsSheetID(beatsSheet.ID),
		LoglineID: apimodels.LoglineID(beatsSheet.LoglineID),
		Content: lo.Map(beatsSheet.Content, func(item models.Beat, _ int) apimodels.Beat {
			return apimodels.Beat{
				Key:     item.Key,
				Title:   item.Title,
				Content: item.Content,
			}
		}),
		Lang:      apimodels.Lang(beatsSheet.Lang),
		CreatedAt: beatsSheet.CreatedAt,
	}), nil
}
