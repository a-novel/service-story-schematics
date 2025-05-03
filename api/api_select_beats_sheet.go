package api

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"

	authapi "github.com/a-novel/service-authentication/api"

	"github.com/a-novel-kit/context"

	"github.com/a-novel/service-story-schematics/api/codegen"
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
	userID, err := authapi.RequireUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user ID: %w", err)
	}

	beatsSheet, err := api.SelectBeatsSheetService.SelectBeatsSheet(ctx, services.SelectBeatsSheetRequest{
		BeatsSheetID: uuid.UUID(params.BeatsSheetID),
		UserID:       userID,
	})

	switch {
	case errors.Is(err, dao.ErrBeatsSheetNotFound), errors.Is(err, dao.ErrLoglineNotFound):
		return &codegen.NotFoundError{Error: err.Error()}, nil
	case err != nil:
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
		CreatedAt: beatsSheet.CreatedAt,
	}, nil
}
