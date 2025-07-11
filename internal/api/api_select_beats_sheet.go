package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/samber/lo"

	authPkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type SelectBeatsSheetService interface {
	SelectBeatsSheet(ctx context.Context, request services.SelectBeatsSheetRequest) (*models.BeatsSheet, error)
}

func (api *API) GetBeatsSheet(
	ctx context.Context, params codegen.GetBeatsSheetParams,
) (codegen.GetBeatsSheetRes, error) {
	span := sentry.StartSpan(ctx, "API.GetBeatsSheet")
	defer span.Finish()

	span.SetData("request.beatsSheetID", params.BeatsSheetID)

	userID, err := authPkg.RequireUserID(ctx)
	if err != nil {
		span.SetData("request.userID.err", err.Error())

		return nil, fmt.Errorf("get user ID: %w", err)
	}

	span.SetData("request.userID", userID)

	beatsSheet, err := api.SelectBeatsSheetService.SelectBeatsSheet(span.Context(), services.SelectBeatsSheetRequest{
		BeatsSheetID: uuid.UUID(params.BeatsSheetID),
		UserID:       userID,
	})

	switch {
	case errors.Is(err, dao.ErrBeatsSheetNotFound), errors.Is(err, dao.ErrLoglineNotFound):
		span.SetData("service.err", err.Error())

		return &codegen.NotFoundError{Error: err.Error()}, nil
	case err != nil:
		span.SetData("service.err", err.Error())

		return nil, fmt.Errorf("get beats sheet: %w", err)
	}

	return &codegen.BeatsSheet{
		ID:          codegen.BeatsSheetID(beatsSheet.ID),
		LoglineID:   codegen.LoglineID(beatsSheet.LoglineID),
		StoryPlanID: codegen.StoryPlanID(beatsSheet.StoryPlanID),
		Content: lo.Map(beatsSheet.Content, func(item models.Beat, _ int) codegen.Beat {
			return codegen.Beat{
				Key:     item.Key,
				Title:   item.Title,
				Content: item.Content,
			}
		}),
		Lang:      codegen.Lang(beatsSheet.Lang),
		CreatedAt: beatsSheet.CreatedAt,
	}, nil
}
