package api

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/samber/lo"

	authPkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/api/codegen"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

type ListLoglinesService interface {
	ListLoglines(ctx context.Context, request services.ListLoglinesRequest) ([]*models.LoglinePreview, error)
}

func (api *API) GetLoglines(ctx context.Context, params codegen.GetLoglinesParams) (codegen.GetLoglinesRes, error) {
	span := sentry.StartSpan(ctx, "API.GetLoglines")
	defer span.Finish()

	span.SetData("request.limit", params.Limit)
	span.SetData("request.offset", params.Offset)

	userID, err := authPkg.RequireUserID(ctx)
	if err != nil {
		span.SetData("request.userID.err", err.Error())

		return nil, fmt.Errorf("get user ID: %w", err)
	}

	span.SetData("request.userID", userID)

	loglines, err := api.ListLoglinesService.ListLoglines(span.Context(), services.ListLoglinesRequest{
		UserID: userID,
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	})
	if err != nil {
		span.SetData("service.err", err.Error())

		return nil, fmt.Errorf("list loglines: %w", err)
	}

	res := codegen.GetLoglinesOKApplicationJSON(
		lo.Map(loglines, func(item *models.LoglinePreview, _ int) codegen.LoglinePreview {
			return codegen.LoglinePreview{
				Slug:      codegen.Slug(item.Slug),
				Name:      item.Name,
				Content:   item.Content,
				Lang:      codegen.Lang(item.Lang),
				CreatedAt: item.CreatedAt,
			}
		}),
	)

	span.SetData("loglines.count", len(loglines))

	return &res, nil
}
