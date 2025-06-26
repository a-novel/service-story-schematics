package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/samber/lo"

	authapi "github.com/a-novel/service-authentication/api"

	"github.com/a-novel/service-story-schematics/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type RegenerateBeatsService interface {
	RegenerateBeats(ctx context.Context, request services.RegenerateBeatsRequest) ([]models.Beat, error)
}

func (api *API) RegenerateBeats(
	ctx context.Context, req *codegen.RegenerateBeatsForm,
) (codegen.RegenerateBeatsRes, error) {
	span := sentry.StartSpan(ctx, "API.RegenerateBeats")
	defer span.Finish()

	span.SetData("request.beatsSheetID", req.GetBeatsSheetID())
	span.SetData("request.regenerateKeys", req.GetRegenerateKeys())

	userID, err := authapi.RequireUserID(ctx)
	if err != nil {
		span.SetData("request.userID.err", err.Error())

		return nil, fmt.Errorf("get user ID: %w", err)
	}

	span.SetData("request.userID", userID)

	beats, err := api.RegenerateBeatsService.RegenerateBeats(span.Context(), services.RegenerateBeatsRequest{
		BeatsSheetID:   uuid.UUID(req.GetBeatsSheetID()),
		UserID:         userID,
		RegenerateKeys: req.GetRegenerateKeys(),
	})

	switch {
	case errors.Is(err, dao.ErrBeatsSheetNotFound),
		errors.Is(err, dao.ErrLoglineNotFound),
		errors.Is(err, dao.ErrStoryPlanNotFound):
		span.SetData("service.err", err.Error())

		return &codegen.NotFoundError{Error: err.Error()}, nil
	case err != nil:
		span.SetData("service.err", err.Error())

		return nil, fmt.Errorf("regenerate beats: %w", err)
	}

	var res codegen.Beats = lo.Map(beats, func(item models.Beat, _ int) codegen.Beat {
		return codegen.Beat{
			Key:     item.Key,
			Title:   item.Title,
			Content: item.Content,
		}
	})

	return &res, nil
}
