package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/a-novel-kit/context"

	"github.com/a-novel/story-schematics/internal/dao"
	"github.com/a-novel/story-schematics/models"
)

var ErrSelectLoglineService = errors.New("SelectLoglineService.SelectLogline")

func NewErrSelectLoglineService(err error) error {
	return errors.Join(err, ErrSelectLoglineService)
}

type SelectLoglineSource interface {
	SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error)
	SelectLoglineBySlug(ctx context.Context, data dao.SelectLoglineBySlugData) (*dao.LoglineEntity, error)
}

type SelectLoglineRequest struct {
	UserID uuid.UUID
	Slug   *models.Slug
	ID     *uuid.UUID
}

type SelectLoglineService struct {
	source SelectLoglineSource
}

func (service *SelectLoglineService) SelectLogline(
	ctx context.Context, request SelectLoglineRequest,
) (*models.Logline, error) {
	if request.Slug != nil {
		data, err := service.source.SelectLoglineBySlug(ctx, dao.SelectLoglineBySlugData{
			Slug:   lo.FromPtr(request.Slug),
			UserID: request.UserID,
		})
		if err != nil {
			return nil, NewErrSelectLoglineService(err)
		}

		return &models.Logline{
			ID:        data.ID,
			UserID:    data.UserID,
			Slug:      data.Slug,
			Name:      data.Name,
			Content:   data.Content,
			CreatedAt: data.CreatedAt,
		}, nil
	}

	data, err := service.source.SelectLogline(ctx, dao.SelectLoglineData{
		ID:     lo.FromPtr(request.ID),
		UserID: request.UserID,
	})
	if err != nil {
		return nil, NewErrSelectLoglineService(err)
	}

	return &models.Logline{
		ID:        data.ID,
		UserID:    data.UserID,
		Slug:      data.Slug,
		Name:      data.Name,
		Content:   data.Content,
		CreatedAt: data.CreatedAt,
	}, nil
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

func NewSelectLoglineService(source SelectLoglineSource) *SelectLoglineService {
	return &SelectLoglineService{source: source}
}
