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

type GenerateBeatsSheetService interface {
	GenerateBeatsSheet(ctx context.Context, request services.GenerateBeatsSheetRequest) ([]models.Beat, error)
}

func (api *API) GenerateBeatsSheet(
	ctx context.Context, req *apimodels.GenerateBeatsSheetForm,
) (apimodels.GenerateBeatsSheetRes, error) {
	ctx, span := otel.Tracer().Start(ctx, "api.GenerateBeatsSheet")
	defer span.End()

	userID, err := authpkg.RequireUserID(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get user ID: %w", err))
	}

	beatsSheet, err := api.GenerateBeatsSheetService.GenerateBeatsSheet(
		ctx,
		services.GenerateBeatsSheetRequest{
			LoglineID: uuid.UUID(req.GetLoglineID()),
			UserID:    userID,
			Lang:      models.Lang(req.GetLang()),
		},
	)

	switch {
	case errors.Is(err, dao.ErrLoglineNotFound), errors.Is(err, services.ErrStoryPlanNotFound):
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return &apimodels.NotFoundError{Error: err.Error()}, nil
	case err != nil:
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return nil, fmt.Errorf("generate beats sheet: %w", err)
	}

	return otel.ReportSuccess(span, &apimodels.BeatsSheetIdea{
		Content: lo.Map(beatsSheet, func(item models.Beat, _ int) apimodels.Beat {
			return apimodels.Beat{
				Key:     item.Key,
				Title:   item.Title,
				Content: item.Content,
			}
		}),
		Lang: req.GetLang(),
	}), nil
}
