package services

import (
	"context"
	"errors"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrSelectLoglineService = errors.New("SelectLoglineService.SelectLogline")

func NewErrSelectLoglineService(err error) error {
	return errors.Join(err, ErrSelectLoglineService)
}

type SelectLoglineSource interface {
	SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error)
	SelectLoglineBySlug(ctx context.Context, data dao.SelectLoglineBySlugData) (*dao.LoglineEntity, error)
}

func NewSelectLoglineServiceSource(
	selectDAO *dao.SelectLoglineRepository,
	selectBySlugRepository *dao.SelectLoglineBySlugRepository,
) SelectLoglineSource {
	return &struct {
		*dao.SelectLoglineRepository
		*dao.SelectLoglineBySlugRepository
	}{
		SelectLoglineRepository:       selectDAO,
		SelectLoglineBySlugRepository: selectBySlugRepository,
	}
}

type SelectLoglineRequest struct {
	UserID uuid.UUID
	Slug   *models.Slug
	ID     *uuid.UUID
}

type SelectLoglineService struct {
	source SelectLoglineSource
}

func NewSelectLoglineService(source SelectLoglineSource) *SelectLoglineService {
	return &SelectLoglineService{source: source}
}

func (service *SelectLoglineService) SelectLogline(
	ctx context.Context, request SelectLoglineRequest,
) (*models.Logline, error) {
	span := sentry.StartSpan(ctx, "SelectLoglineService.SelectLogline")
	defer span.Finish()

	span.SetData("request.userID", request.UserID)
	span.SetData("request.slug", request.Slug)
	span.SetData("request.id", request.ID)

	if request.Slug != nil {
		data, err := service.source.SelectLoglineBySlug(span.Context(), dao.SelectLoglineBySlugData{
			Slug:   lo.FromPtr(request.Slug),
			UserID: request.UserID,
		})
		if err != nil {
			span.SetData("dao.selectLoglineBySlug.err", err.Error())

			return nil, NewErrSelectLoglineService(err)
		}

		return &models.Logline{
			ID:        data.ID,
			UserID:    data.UserID,
			Slug:      data.Slug,
			Name:      data.Name,
			Content:   data.Content,
			Lang:      data.Lang,
			CreatedAt: data.CreatedAt,
		}, nil
	}

	data, err := service.source.SelectLogline(span.Context(), dao.SelectLoglineData{
		ID:     lo.FromPtr(request.ID),
		UserID: request.UserID,
	})
	if err != nil {
		span.SetData("dao.selectLogline.err", err.Error())

		return nil, NewErrSelectLoglineService(err)
	}

	return &models.Logline{
		ID:        data.ID,
		UserID:    data.UserID,
		Slug:      data.Slug,
		Name:      data.Name,
		Content:   data.Content,
		Lang:      data.Lang,
		CreatedAt: data.CreatedAt,
	}, nil
}
