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
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type CreateBeatsSheetService interface {
	CreateBeatsSheet(ctx context.Context, request services.CreateBeatsSheetRequest) (*models.BeatsSheet, error)
}

func (api *API) CreateBeatsSheet(
	ctx context.Context, req *codegen.CreateBeatsSheetForm,
) (codegen.CreateBeatsSheetRes, error) {
	span := sentry.StartSpan(ctx, "API.CreateBeatsSheet")
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

	beatsSheet, err := api.CreateBeatsSheetService.CreateBeatsSheet(span.Context(), services.CreateBeatsSheetRequest{
		LoglineID:   uuid.UUID(req.GetLoglineID()),
		UserID:      userID,
		StoryPlanID: uuid.UUID(req.GetStoryPlanID()),
		Content: lo.Map(req.GetContent(), func(item codegen.Beat, _ int) models.Beat {
			return models.Beat{
				Key:     item.GetKey(),
				Title:   item.GetTitle(),
				Content: item.GetContent(),
			}
		}),
		Lang: models.Lang(req.GetLang()),
	})

	switch {
	case errors.Is(err, dao.ErrLoglineNotFound), errors.Is(err, dao.ErrStoryPlanNotFound):
		span.SetData("service.err", err.Error())

		return &codegen.NotFoundError{Error: err.Error()}, nil
	case errors.Is(err, lib.ErrInvalidStoryPlan):
		span.SetData("service.err", err.Error())

		return &codegen.UnprocessableEntityError{Error: err.Error()}, nil
	case err != nil:
		span.SetData("service.err", err.Error())

		return nil, fmt.Errorf("create beats sheet: %w", err)
	}

	span.SetData("service.beatsSheet.id", beatsSheet.ID)

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
