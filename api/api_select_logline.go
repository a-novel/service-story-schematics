package api

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"

	authapi "github.com/a-novel/authentication/api"

	"github.com/a-novel-kit/context"

	"github.com/a-novel/story-schematics/api/codegen"
	"github.com/a-novel/story-schematics/internal/dao"
	"github.com/a-novel/story-schematics/internal/services"
	"github.com/a-novel/story-schematics/models"
)

type SelectLoglineService interface {
	SelectLogline(ctx context.Context, request services.SelectLoglineRequest) (*models.Logline, error)
}

func (api *API) GetLogline(ctx context.Context, params codegen.GetLoglineParams) (codegen.GetLoglineRes, error) {
	userID, err := authapi.RequireUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user ID: %w", err)
	}

	logline, err := api.SelectLoglineService.SelectLogline(ctx, services.SelectLoglineRequest{
		UserID: userID,
		Slug:   lo.Ternary(params.Slug.IsSet(), lo.ToPtr(models.Slug(params.Slug.Value)), nil),
		ID:     lo.Ternary(params.ID.IsSet(), lo.ToPtr(uuid.UUID(params.ID.Value)), nil),
	})

	switch {
	case errors.Is(err, dao.ErrLoglineNotFound):
		return &codegen.NotFoundError{Error: err.Error()}, nil
	case err != nil:
		return nil, fmt.Errorf("get logline: %w", err)
	}

	return &codegen.Logline{
		ID:        codegen.LoglineID(logline.ID),
		UserID:    codegen.UserID(logline.UserID),
		Slug:      codegen.Slug(logline.Slug),
		Name:      logline.Name,
		Content:   logline.Content,
		CreatedAt: logline.CreatedAt,
	}, nil
}
