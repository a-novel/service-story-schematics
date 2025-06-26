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

type GenerateBeatsSheetService interface {
	GenerateBeatsSheet(ctx context.Context, request services.GenerateBeatsSheetRequest) ([]models.Beat, error)
}

func (api *API) GenerateBeatsSheet(
	ctx context.Context, req *codegen.GenerateBeatsSheetForm,
) (codegen.GenerateBeatsSheetRes, error) {
	span := sentry.StartSpan(ctx, "API.GenerateBeatsSheet")
	defer span.Finish()

	span.SetData("request.loglineID", req.GetLoglineID())
	span.SetData("request.storyPlanID", req.GetStoryPlanID())
	span.SetData("request.lang", req.GetLang())

	userID, err := authapi.RequireUserID(ctx)
	if err != nil {
		span.SetData("request.userID.err", err.Error())

		return nil, fmt.Errorf("get user ID: %w", err)
	}

	span.SetData("request.userID", userID)

	beatsSheet, err := api.GenerateBeatsSheetService.GenerateBeatsSheet(
		span.Context(),
		services.GenerateBeatsSheetRequest{
			LoglineID:   uuid.UUID(req.GetLoglineID()),
			StoryPlanID: uuid.UUID(req.GetStoryPlanID()),
			UserID:      userID,
			Lang:        models.Lang(req.GetLang()),
		},
	)

	switch {
	case errors.Is(err, dao.ErrLoglineNotFound), errors.Is(err, dao.ErrStoryPlanNotFound):
		span.SetData("service.err", err.Error())

		return &codegen.NotFoundError{Error: err.Error()}, nil
	case err != nil:
		span.SetData("service.err", err.Error())

		return nil, fmt.Errorf("generate beats sheet: %w", err)
	}

	return &codegen.BeatsSheetIdea{
		Content: lo.Map(beatsSheet, func(item models.Beat, _ int) codegen.Beat {
			return codegen.Beat{
				Key:     item.Key,
				Title:   item.Title,
				Content: item.Content,
			}
		}),
		Lang: req.GetLang(),
	}, nil
}
