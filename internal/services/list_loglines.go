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

var ErrListLoglinesService = errors.New("ListLoglinesService.ListLoglines")

func NewErrListLoglinesService(err error) error {
	return errors.Join(err, ErrListLoglinesService)
}

type ListLoglinesSource interface {
	ListLoglines(ctx context.Context, data dao.ListLoglinesData) ([]*dao.LoglinePreviewEntity, error)
}

type ListLoglinesRequest struct {
	UserID uuid.UUID
	Limit  int
	Offset int
}

type ListLoglinesService struct {
	source ListLoglinesSource
}

func (service *ListLoglinesService) ListLoglines(
	ctx context.Context, request ListLoglinesRequest,
) ([]*models.LoglinePreview, error) {
	span := sentry.StartSpan(ctx, "ListLoglinesService.ListLoglines")
	defer span.Finish()

	span.SetData("request.userID", request.UserID)
	span.SetData("request.limit", request.Limit)
	span.SetData("request.offset", request.Offset)

	resp, err := service.source.ListLoglines(span.Context(), dao.ListLoglinesData{
		UserID: request.UserID,
		Limit:  request.Limit,
		Offset: request.Offset,
	})
	if err != nil {
		span.SetData("dao.listLoglines.err", err.Error())

		return nil, NewErrListLoglinesService(err)
	}

	span.SetData("dao.listLoglines.count", len(resp))

	return lo.Map(resp, func(item *dao.LoglinePreviewEntity, _ int) *models.LoglinePreview {
		return &models.LoglinePreview{
			Slug:      item.Slug,
			Name:      item.Name,
			Content:   item.Content,
			Lang:      item.Lang,
			CreatedAt: item.CreatedAt,
		}
	}), nil
}

func NewListLoglinesService(source ListLoglinesSource) *ListLoglinesService {
	return &ListLoglinesService{source: source}
}
