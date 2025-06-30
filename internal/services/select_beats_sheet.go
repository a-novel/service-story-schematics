package services

import (
	"context"
	"errors"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrSelectBeatsSheetService = errors.New("SelectBeatsSheetService.SelectBeatsSheet")

func NewErrSelectBeatsSheetService(err error) error {
	return errors.Join(err, ErrSelectBeatsSheetService)
}

type SelectBeatsSheetSource interface {
	SelectBeatsSheet(ctx context.Context, data uuid.UUID) (*dao.BeatsSheetEntity, error)
	SelectLogline(ctx context.Context, data dao.SelectLoglineData) (*dao.LoglineEntity, error)
}

func NewSelectBeatsSheetServiceSource(
	selectBeatsSheetDAO *dao.SelectBeatsSheetRepository,
	selectLoglineDAO *dao.SelectLoglineRepository,
) SelectBeatsSheetSource {
	return &struct {
		*dao.SelectBeatsSheetRepository
		*dao.SelectLoglineRepository
	}{
		SelectBeatsSheetRepository: selectBeatsSheetDAO,
		SelectLoglineRepository:    selectLoglineDAO,
	}
}

type SelectBeatsSheetRequest struct {
	BeatsSheetID uuid.UUID
	UserID       uuid.UUID
}

type SelectBeatsSheetService struct {
	source SelectBeatsSheetSource
}

func NewSelectBeatsSheetService(source SelectBeatsSheetSource) *SelectBeatsSheetService {
	return &SelectBeatsSheetService{source: source}
}

func (service *SelectBeatsSheetService) SelectBeatsSheet(
	ctx context.Context, request SelectBeatsSheetRequest,
) (*models.BeatsSheet, error) {
	span := sentry.StartSpan(ctx, "SelectBeatsSheetService.SelectBeatsSheet")
	defer span.Finish()

	span.SetData("request.beatsSheetID", request.BeatsSheetID)
	span.SetData("request.userID", request.UserID)

	data, err := service.source.SelectBeatsSheet(span.Context(), request.BeatsSheetID)
	if err != nil {
		span.SetData("dao.selectBeatsSheet.err", err.Error())

		return nil, NewErrSelectBeatsSheetService(err)
	}

	// Make sure the selected beats sheet is linked to a logline that belongs to the user.
	_, err = service.source.SelectLogline(span.Context(), dao.SelectLoglineData{
		ID:     data.LoglineID,
		UserID: request.UserID,
	})
	if err != nil {
		span.SetData("dao.selectLogline.err", err.Error())

		return nil, NewErrSelectBeatsSheetService(err)
	}

	return &models.BeatsSheet{
		ID:          data.ID,
		LoglineID:   data.LoglineID,
		StoryPlanID: data.StoryPlanID,
		Content:     data.Content,
		Lang:        data.Lang,
		CreatedAt:   data.CreatedAt,
	}, nil
}
