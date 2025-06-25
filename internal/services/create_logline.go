package services

import (
	"context"
	"errors"
	"github.com/getsentry/sentry-go"
	"time"

	"github.com/google/uuid"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrCreateLoglineService = errors.New("CreateLoglineService.CreateLogline")

func NewErrCreateLoglineService(err error) error {
	return errors.Join(err, ErrCreateLoglineService)
}

type CreateLoglineSource interface {
	InsertLogline(ctx context.Context, data dao.InsertLoglineData) (*dao.LoglineEntity, error)
	SelectSlugIteration(ctx context.Context, data dao.SelectSlugIterationData) (models.Slug, int, error)
}

type CreateLoglineRequest struct {
	UserID  uuid.UUID
	Slug    models.Slug
	Name    string
	Content string
	Lang    models.Lang
}

type CreateLoglineService struct {
	source CreateLoglineSource
}

func (service *CreateLoglineService) CreateLogline(
	ctx context.Context, request CreateLoglineRequest,
) (*models.Logline, error) {
	span := sentry.StartSpan(ctx, "CreateLoglineService.CreateLogline")
	defer span.Finish()

	span.SetData("request.userID", request.UserID)
	span.SetData("request.slug", request.Slug)
	span.SetData("request.name", request.Name)
	span.SetData("request.lang", request.Lang)

	data := dao.InsertLoglineData{
		ID:      uuid.New(),
		UserID:  request.UserID,
		Slug:    request.Slug,
		Name:    request.Name,
		Content: request.Content,
		Lang:    request.Lang,
		Now:     time.Now(),
	}

	resp, err := service.source.InsertLogline(span.Context(), data)

	// If slug is taken, try to modify it by appending a version number.
	if errors.Is(err, dao.ErrLoglineAlreadyExists) {
		span.SetData("dao.insertLogline.slug.taken", err.Error())

		data.Slug, _, err = service.source.SelectSlugIteration(span.Context(), dao.SelectSlugIterationData{
			Slug:  data.Slug,
			Table: "loglines",
			Filter: map[string][]any{
				"user_id = ?": {data.UserID},
			},
			Order: []string{"created_at DESC"},
		})
		if err != nil {
			span.SetData("dao.selectSlugIteration.err", err.Error())

			return nil, NewErrCreateLoglineService(err)
		}

		resp, err = service.source.InsertLogline(span.Context(), data)
	}

	if err != nil {
		span.SetData("dao.insertLogline.err", err.Error())

		return nil, NewErrCreateLoglineService(err)
	}

	span.SetData("dao.insertLogline.id", resp.ID)

	return &models.Logline{
		ID:        resp.ID,
		UserID:    resp.UserID,
		Slug:      resp.Slug,
		Name:      resp.Name,
		Content:   resp.Content,
		Lang:      resp.Lang,
		CreatedAt: resp.CreatedAt,
	}, nil
}

func NewCreateLoglineServiceSource(
	insertLoglineDAO *dao.InsertLoglineRepository,
	selectSlugIterationDAO *dao.SelectSlugIterationRepository,
) CreateLoglineSource {
	return &struct {
		*dao.InsertLoglineRepository
		*dao.SelectSlugIterationRepository
	}{
		InsertLoglineRepository:       insertLoglineDAO,
		SelectSlugIterationRepository: selectSlugIterationDAO,
	}
}

func NewCreateLoglineService(source CreateLoglineSource) *CreateLoglineService {
	return &CreateLoglineService{source: source}
}
