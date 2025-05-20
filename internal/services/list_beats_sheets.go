package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrListBeatsSheetsService = errors.New("ListBeatsSheetsService.ListBeatsSheets")

func NewErrListBeatsSheetsService(err error) error {
	return errors.Join(err, ErrListBeatsSheetsService)
}

type ListBeatsSheetsSource interface {
	ListBeatsSheets(ctx context.Context, data dao.ListBeatsSheetsData) ([]*dao.BeatsSheetPreviewEntity, error)
	SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error)
}

type ListBeatsSheetsRequest struct {
	UserID    uuid.UUID
	LoglineID uuid.UUID
	Limit     int
	Offset    int
}

type ListBeatsSheetsService struct {
	source ListBeatsSheetsSource
}

func (service *ListBeatsSheetsService) ListBeatsSheets(
	ctx context.Context, request ListBeatsSheetsRequest,
) ([]*models.BeatsSheetPreview, error) {
	_, err := service.source.SelectLogline(ctx, dao.SelectLoglineData{
		ID:     request.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, NewErrListBeatsSheetsService(err)
	}

	data := dao.ListBeatsSheetsData{
		LoglineID: request.LoglineID,
		Limit:     request.Limit,
		Offset:    request.Offset,
	}

	resp, err := service.source.ListBeatsSheets(ctx, data)
	if err != nil {
		return nil, NewErrListBeatsSheetsService(err)
	}

	return lo.Map(resp, func(item *dao.BeatsSheetPreviewEntity, _ int) *models.BeatsSheetPreview {
		return &models.BeatsSheetPreview{
			ID:        item.ID,
			CreatedAt: item.CreatedAt,
		}
	}), nil
}

func NewListBeatsSheetsServiceSource(
	listBeatsSheetsDAO *dao.ListBeatsSheetsRepository,
	selectLoglineDAO *dao.SelectLoglineRepository,
) ListBeatsSheetsSource {
	return &struct {
		*dao.ListBeatsSheetsRepository
		*dao.SelectLoglineRepository
	}{
		ListBeatsSheetsRepository: listBeatsSheetsDAO,
		SelectLoglineRepository:   selectLoglineDAO,
	}
}

func NewListBeatsSheetsService(source ListBeatsSheetsSource) *ListBeatsSheetsService {
	return &ListBeatsSheetsService{source: source}
}
