package api

import (
	"context"
	"errors"
	"fmt"

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
	userID, err := authapi.RequireUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user ID: %w", err)
	}

	beatsSheet, err := api.CreateBeatsSheetService.CreateBeatsSheet(ctx, services.CreateBeatsSheetRequest{
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
		return &codegen.NotFoundError{Error: err.Error()}, nil
	case errors.Is(err, lib.ErrInvalidStoryPlan):
		return &codegen.UnprocessableEntityError{Error: err.Error()}, nil
	case err != nil:
		return nil, fmt.Errorf("create beats sheet: %w", err)
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
