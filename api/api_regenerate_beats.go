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

type RegenerateBeatsService interface {
	RegenerateBeats(ctx context.Context, request services.RegenerateBeatsRequest) ([]models.Beat, error)
}

func (api *API) RegenerateBeats(
	ctx context.Context, req *codegen.RegenerateBeatsForm,
) (codegen.RegenerateBeatsRes, error) {
	userID, err := authapi.RequireUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user ID: %w", err)
	}

	beats, err := api.RegenerateBeatsService.RegenerateBeats(ctx, services.RegenerateBeatsRequest{
		BeatsSheetID:   uuid.UUID(req.GetBeatsSheetID()),
		UserID:         userID,
		RegenerateKeys: req.GetRegenerateKeys(),
	})

	switch {
	case errors.Is(err, dao.ErrBeatsSheetNotFound),
		errors.Is(err, dao.ErrLoglineNotFound),
		errors.Is(err, dao.ErrStoryPlanNotFound):
		return &codegen.NotFoundError{Error: err.Error()}, nil
	case err != nil:
		return nil, fmt.Errorf("regenerate beats: %w", err)
	}

	return &codegen.BeatsSheet{
		Content: lo.Map(beats, func(item models.Beat, _ int) codegen.Beat {
			return codegen.Beat{
				Key:     item.Key,
				Title:   item.Title,
				Content: item.Content,
			}
		}),
	}, nil
}
