package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"

	authPkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type ExpandBeatService interface {
	ExpandBeat(ctx context.Context, request services.ExpandBeatRequest) (*models.Beat, error)
}

func (api *API) ExpandBeat(ctx context.Context, req *codegen.ExpandBeatForm) (codegen.ExpandBeatRes, error) {
	span := sentry.StartSpan(ctx, "API.ExpandBeat")
	defer span.Finish()

	span.SetData("request.beatsSheetID", req.GetBeatsSheetID())
	span.SetData("request.targetKey", req.GetTargetKey())

	userID, err := authPkg.RequireUserID(ctx)
	if err != nil {
		span.SetData("request.userID.err", err.Error())

		return nil, fmt.Errorf("get user ID: %w", err)
	}

	span.SetData("request.userID", userID)

	beat, err := api.ExpandBeatService.ExpandBeat(span.Context(), services.ExpandBeatRequest{
		BeatsSheetID: uuid.UUID(req.GetBeatsSheetID()),
		TargetKey:    req.GetTargetKey(),
		UserID:       userID,
	})

	switch {
	case errors.Is(err, dao.ErrBeatsSheetNotFound), errors.Is(err, dao.ErrStoryPlanNotFound):
		span.SetData("service.err", err.Error())

		return &codegen.NotFoundError{Error: err.Error()}, nil
	case errors.Is(err, daoai.ErrUnknownTargetKey):
		span.SetData("service.err", err.Error())

		return &codegen.UnprocessableEntityError{Error: err.Error()}, nil
	case err != nil:
		span.SetData("service.err", err.Error())

		return nil, fmt.Errorf("expand beat: %w", err)
	}

	return &codegen.Beat{
		Key:     beat.Key,
		Title:   beat.Title,
		Content: beat.Content,
	}, nil
}
