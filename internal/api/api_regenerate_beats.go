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

type RegenerateBeatsService interface {
	RegenerateBeats(ctx context.Context, request services.RegenerateBeatsRequest) ([]models.Beat, error)
}

func (api *API) RegenerateBeats(
	ctx context.Context, req *apimodels.RegenerateBeatsForm,
) (apimodels.RegenerateBeatsRes, error) {
	ctx, span := otel.Tracer().Start(ctx, "api.RegenerateBeats")
	defer span.End()

	userID, err := authpkg.RequireUserID(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get user ID: %w", err))
	}

	beats, err := api.RegenerateBeatsService.RegenerateBeats(ctx, services.RegenerateBeatsRequest{
		BeatsSheetID:   uuid.UUID(req.GetBeatsSheetID()),
		UserID:         userID,
		RegenerateKeys: req.GetRegenerateKeys(),
	})

	switch {
	case errors.Is(err, dao.ErrBeatsSheetNotFound),
		errors.Is(err, dao.ErrLoglineNotFound),
		errors.Is(err, services.ErrStoryPlanNotFound):
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return &apimodels.NotFoundError{Error: err.Error()}, nil
	case err != nil:
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return nil, fmt.Errorf("regenerate beats: %w", err)
	}

	var res apimodels.Beats = lo.Map(beats, func(item models.Beat, _ int) apimodels.Beat {
		return apimodels.Beat{
			Key:     item.Key,
			Title:   item.Title,
			Content: item.Content,
		}
	})

	return otel.ReportSuccess(span, &res), nil
}
