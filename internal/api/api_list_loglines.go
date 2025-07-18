package api

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/a-novel/golib/otel"
	authpkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

type ListLoglinesService interface {
	ListLoglines(ctx context.Context, request services.ListLoglinesRequest) ([]*models.LoglinePreview, error)
}

func (api *API) GetLoglines(ctx context.Context, params apimodels.GetLoglinesParams) (apimodels.GetLoglinesRes, error) {
	ctx, span := otel.Tracer().Start(ctx, "api.GetLoglines")
	defer span.End()

	userID, err := authpkg.RequireUserID(ctx)
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("get user ID: %w", err))
	}

	loglines, err := api.ListLoglinesService.ListLoglines(ctx, services.ListLoglinesRequest{
		UserID: userID,
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	})
	if err != nil {
		return nil, otel.ReportError(span, fmt.Errorf("list loglines: %w", err))
	}

	res := apimodels.GetLoglinesOKApplicationJSON(
		lo.Map(loglines, func(item *models.LoglinePreview, _ int) apimodels.LoglinePreview {
			return apimodels.LoglinePreview{
				Slug:      apimodels.Slug(item.Slug),
				Name:      item.Name,
				Content:   item.Content,
				Lang:      apimodels.Lang(item.Lang),
				CreatedAt: item.CreatedAt,
			}
		}),
	)

	return otel.ReportSuccess(span, &res), nil
}
