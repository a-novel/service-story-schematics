package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/codes"

	"github.com/a-novel/golib/otel"
	authpkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
	storyplanmodel "github.com/a-novel/service-story-schematics/models/story_plan"
)

type CreateBeatsSheetService interface {
	CreateBeatsSheet(ctx context.Context, request services.CreateBeatsSheetRequest) (*models.BeatsSheet, error)
}

func (api *API) CreateBeatsSheet(
	ctx context.Context, req *apimodels.CreateBeatsSheetForm,
) (apimodels.CreateBeatsSheetRes, error) {
	ctx, span := otel.Tracer().Start(ctx, "api.CreateBeatsSheet")
	defer span.End()

	userID, err := authpkg.RequireUserID(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get user ID: %w", err))
	}

	beatsSheet, err := api.CreateBeatsSheetService.CreateBeatsSheet(ctx, services.CreateBeatsSheetRequest{
		LoglineID: uuid.UUID(req.GetLoglineID()),
		UserID:    userID,
		Content: lo.Map(req.GetContent(), func(item apimodels.Beat, _ int) models.Beat {
			return models.Beat{
				Key:     item.GetKey(),
				Title:   item.GetTitle(),
				Content: item.GetContent(),
			}
		}),
		Lang: models.Lang(req.GetLang()),
	})

	switch {
	case errors.Is(err, dao.ErrLoglineNotFound), errors.Is(err, services.ErrStoryPlanNotFound):
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return &apimodels.NotFoundError{Error: err.Error()}, nil
	case errors.Is(err, storyplanmodel.ErrInvalidPlan):
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return &apimodels.UnprocessableEntityError{Error: err.Error()}, nil
	case err != nil:
		span.RecordError(err)
		span.SetStatus(codes.Error, "")

		return nil, fmt.Errorf("create beats sheet: %w", err)
	}

	return otel.ReportSuccess(span, &apimodels.BeatsSheet{
		ID:        apimodels.BeatsSheetID(beatsSheet.ID),
		LoglineID: apimodels.LoglineID(beatsSheet.LoglineID),
		Content: lo.Map(beatsSheet.Content, func(item models.Beat, _ int) apimodels.Beat {
			return apimodels.Beat{
				Key:     item.Key,
				Title:   item.Title,
				Content: item.Content,
			}
		}),
		Lang:      apimodels.Lang(beatsSheet.Lang),
		CreatedAt: beatsSheet.CreatedAt,
	}), nil
}
