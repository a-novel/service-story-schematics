package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	servicesmocks "github.com/a-novel/service-story-schematics/internal/services/mocks"
	"github.com/a-novel/service-story-schematics/models"
)

func TestCreateLogline(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type insertLoglineData struct {
		resp *dao.LoglineEntity
		err  error
	}

	type selectSlugIterationData struct {
		slug      models.Slug
		iteration int
		err       error
	}

	testCases := []struct {
		name string

		request services.CreateLoglineRequest

		insertLoglineData       *insertLoglineData
		selectSlugIterationData *selectSlugIterationData
		reinsertLoglineData     *insertLoglineData

		expect    *models.Logline
		expectErr error
	}{
		{
			name: "Success",

			request: services.CreateLoglineRequest{
				UserID:  uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:    "test-slug",
				Name:    "Test Logline",
				Content: "Once upon a time",
				Lang:    models.LangEN,
			},

			insertLoglineData: &insertLoglineData{
				resp: &dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Logline",
					Content:   "Once upon a time",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &models.Logline{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:      "test-slug",
				Name:      "Test Logline",
				Content:   "Once upon a time",
				Lang:      models.LangEN,
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "RetrySlug",

			request: services.CreateLoglineRequest{
				UserID:  uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:    "test-slug",
				Name:    "Test Logline",
				Content: "Once upon a time",
				Lang:    models.LangEN,
			},

			insertLoglineData: &insertLoglineData{
				err: dao.ErrLoglineAlreadyExists,
			},

			selectSlugIterationData: &selectSlugIterationData{
				slug:      "test-slug-2",
				iteration: 2,
			},

			reinsertLoglineData: &insertLoglineData{
				resp: &dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug-2",
					Name:      "Test Logline",
					Content:   "Once upon a time",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &models.Logline{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:      "test-slug-2",
				Name:      "Test Logline",
				Content:   "Once upon a time",
				Lang:      models.LangEN,
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "InsertError",

			request: services.CreateLoglineRequest{
				UserID:  uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:    "test-slug",
				Name:    "Test Logline",
				Content: "Once upon a time",
				Lang:    models.LangEN,
			},

			insertLoglineData: &insertLoglineData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "SlugIterationError",

			request: services.CreateLoglineRequest{
				UserID:  uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:    "test-slug",
				Name:    "Test Logline",
				Content: "Once upon a time",
				Lang:    models.LangEN,
			},

			insertLoglineData: &insertLoglineData{
				err: dao.ErrLoglineAlreadyExists,
			},

			selectSlugIterationData: &selectSlugIterationData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "RetryInsertError",

			request: services.CreateLoglineRequest{
				UserID:  uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:    "test-slug",
				Name:    "Test Logline",
				Content: "Once upon a time",
				Lang:    models.LangEN,
			},

			insertLoglineData: &insertLoglineData{
				err: dao.ErrLoglineAlreadyExists,
			},

			selectSlugIterationData: &selectSlugIterationData{
				slug:      "test-slug-2",
				iteration: 2,
			},

			reinsertLoglineData: &insertLoglineData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			source := servicesmocks.NewMockCreateLoglineSource(t)

			if testCase.insertLoglineData != nil {
				initialCall := source.EXPECT().
					InsertLogline(ctx, mock.MatchedBy(func(data dao.InsertLoglineData) bool {
						return assert.NotEqual(t, data.ID, uuid.Nil) &&
							assert.Equal(t, testCase.request.UserID, data.UserID) &&
							testCase.request.Slug == data.Slug &&
							assert.Equal(t, testCase.request.Name, data.Name) &&
							assert.Equal(t, testCase.request.Content, data.Content) &&
							assert.Equal(t, testCase.request.Lang, data.Lang) &&
							assert.WithinDuration(t, time.Now(), data.Now, time.Second)
					})).
					Return(testCase.insertLoglineData.resp, testCase.insertLoglineData.err).
					Once()

				if testCase.reinsertLoglineData != nil {
					source.EXPECT().
						InsertLogline(ctx, mock.MatchedBy(func(data dao.InsertLoglineData) bool {
							return assert.NotEqual(t, data.ID, uuid.Nil) &&
								assert.Equal(t, testCase.request.UserID, data.UserID) &&
								testCase.selectSlugIterationData.slug == data.Slug &&
								assert.Equal(t, testCase.request.Name, data.Name) &&
								assert.Equal(t, testCase.request.Content, data.Content) &&
								assert.Equal(t, testCase.request.Lang, data.Lang) &&
								assert.WithinDuration(t, time.Now(), data.Now, time.Second)
						})).
						Return(testCase.reinsertLoglineData.resp, testCase.reinsertLoglineData.err).
						NotBefore(initialCall)
				}
			}

			if testCase.selectSlugIterationData != nil {
				source.EXPECT().
					SelectSlugIteration(ctx, dao.SelectSlugIterationData{
						Slug:   testCase.request.Slug,
						Table:  "loglines",
						Filter: map[string][]any{"user_id = ?": {testCase.request.UserID}},
						Order:  []string{"created_at DESC"},
					}).
					Return(
						testCase.selectSlugIterationData.slug,
						testCase.selectSlugIterationData.iteration,
						testCase.selectSlugIterationData.err,
					)
			}

			service := services.NewCreateLoglineService(source)

			resp, err := service.CreateLogline(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			source.AssertExpectations(t)
		})
	}
}
