package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/a-novel-kit/context"

	"github.com/a-novel/story-schematics/internal/dao"
	"github.com/a-novel/story-schematics/models"
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
}

type CreateLoglineService struct {
	source CreateLoglineSource
}

func (service *CreateLoglineService) CreateLogline(
	ctx context.Context, request CreateLoglineRequest,
) (*models.Logline, error) {
	data := dao.InsertLoglineData{
		ID:      uuid.New(),
		UserID:  request.UserID,
		Slug:    request.Slug,
		Name:    request.Name,
		Content: request.Content,
		Now:     time.Now(),
	}

	resp, err := service.source.InsertLogline(ctx, data)

	// If slug is taken, try to modify it by appending a version number.
	if errors.Is(err, dao.ErrLoglineAlreadyExists) {
		data.Slug, _, err = service.source.SelectSlugIteration(ctx, dao.SelectSlugIterationData{
			Slug:  data.Slug,
			Table: "loglines",
			Filter: map[string][]any{
				"user_id = ?": {data.UserID},
			},
			Order: []string{"created_at DESC"},
		})
		if err != nil {
			return nil, NewErrCreateLoglineService(err)
		}

		resp, err = service.source.InsertLogline(ctx, data)
	}

	if err != nil {
		return nil, NewErrCreateLoglineService(err)
	}

	return &models.Logline{
		ID:        resp.ID,
		UserID:    resp.UserID,
		Slug:      resp.Slug,
		Name:      resp.Name,
		Content:   resp.Content,
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
