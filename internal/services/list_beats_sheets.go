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

var ErrListBeatsSheetsService = errors.New("ListBeatsSheetsService.ListBeatsSheets")

func NewErrListBeatsSheetsService(err error) error {
	return errors.Join(err, ErrListBeatsSheetsService)
}

type ListBeatsSheetsSource interface {
	ListBeatsSheets(ctx context.Context, data dao.ListBeatsSheetsData) ([]*dao.BeatsSheetPreviewEntity, error)
	SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error)
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

type ListBeatsSheetsRequest struct {
	UserID    uuid.UUID
	LoglineID uuid.UUID
	Limit     int
	Offset    int
}

type ListBeatsSheetsService struct {
	source ListBeatsSheetsSource
}

func NewListBeatsSheetsService(source ListBeatsSheetsSource) *ListBeatsSheetsService {
	return &ListBeatsSheetsService{source: source}
}

func (service *ListBeatsSheetsService) ListBeatsSheets(
	ctx context.Context, request ListBeatsSheetsRequest,
) ([]*models.BeatsSheetPreview, error) {
	span := sentry.StartSpan(ctx, "ListBeatsSheetsService.ListBeatsSheets")
	defer span.Finish()

	span.SetData("request.userID", request.UserID)
	span.SetData("request.loglineID", request.LoglineID)
	span.SetData("request.limit", request.Limit)
	span.SetData("request.offset", request.Offset)

	_, err := service.source.SelectLogline(span.Context(), dao.SelectLoglineData{
		ID:     request.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		span.SetData("dao.selectLogline.err", err.Error())

		return nil, NewErrListBeatsSheetsService(err)
	}

	data := dao.ListBeatsSheetsData{
		LoglineID: request.LoglineID,
		Limit:     request.Limit,
		Offset:    request.Offset,
	}

	resp, err := service.source.ListBeatsSheets(span.Context(), data)
	if err != nil {
		span.SetData("dao.listBeatsSheets.err", err.Error())

		return nil, NewErrListBeatsSheetsService(err)
	}

	span.SetData("dao.listBeatsSheets.count", len(resp))

	return lo.Map(resp, func(item *dao.BeatsSheetPreviewEntity, _ int) *models.BeatsSheetPreview {
		return &models.BeatsSheetPreview{
			ID:        item.ID,
			Lang:      item.Lang,
			CreatedAt: item.CreatedAt,
		}
	}), nil
}
