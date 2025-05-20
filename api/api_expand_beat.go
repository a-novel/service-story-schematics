package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	authapi "github.com/a-novel/service-authentication/api"

	"github.com/a-novel/service-story-schematics/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type ExpandBeatService interface {
	ExpandBeat(ctx context.Context, request services.ExpandBeatRequest) (*models.Beat, error)
}

func (api *API) ExpandBeat(ctx context.Context, req *codegen.ExpandBeatForm) (codegen.ExpandBeatRes, error) {
	userID, err := authapi.RequireUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user ID: %w", err)
	}

	beat, err := api.ExpandBeatService.ExpandBeat(ctx, services.ExpandBeatRequest{
		BeatsSheetID: uuid.UUID(req.GetBeatsSheetID()),
		TargetKey:    req.GetTargetKey(),
		UserID:       userID,
	})

	switch {
	case errors.Is(err, dao.ErrBeatsSheetNotFound), errors.Is(err, dao.ErrStoryPlanNotFound):
		return &codegen.NotFoundError{Error: err.Error()}, nil
	case errors.Is(err, daoai.ErrUnknownTargetKey):
		return &codegen.UnprocessableEntityError{Error: err.Error()}, nil
	case err != nil:
		return nil, fmt.Errorf("expand beat: %w", err)
	}

	return &codegen.Beat{
		Key:     beat.Key,
		Title:   beat.Title,
		Content: beat.Content,
	}, nil
}
