package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"

	"github.com/a-novel/golib/otel"
	authpkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

type ExpandBeatService interface {
	ExpandBeat(ctx context.Context, request services.ExpandBeatRequest) (*models.Beat, error)
}

func (api *API) ExpandBeat(ctx context.Context, req *apimodels.ExpandBeatForm) (apimodels.ExpandBeatRes, error) {
	ctx, span := otel.Tracer().Start(ctx, "api.ExpandBeat")
	defer span.End()

	userID, err := authpkg.RequireUserID(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get user ID: %w", err))
	}

	beat, err := api.ExpandBeatService.ExpandBeat(ctx, services.ExpandBeatRequest{
		BeatsSheetID: uuid.UUID(req.GetBeatsSheetID()),
		TargetKey:    req.GetTargetKey(),
		UserID:       userID,
	})

	switch {
	case errors.Is(err, dao.ErrBeatsSheetNotFound), errors.Is(err, services.ErrStoryPlanNotFound):
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return &apimodels.NotFoundError{Error: err.Error()}, nil
	case errors.Is(err, daoai.ErrUnknownTargetKey):
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return &apimodels.UnprocessableEntityError{Error: err.Error()}, nil
	case err != nil:
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return nil, fmt.Errorf("expand beat: %w", err)
	}

	return otel.ReportSuccess(span, &apimodels.Beat{
		Key:     beat.Key,
		Title:   beat.Title,
		Content: beat.Content,
	}), nil
}
