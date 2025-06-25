package services_test

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	servicesmocks "github.com/a-novel/service-story-schematics/internal/services/mocks"
	"github.com/a-novel/service-story-schematics/models"
)

func TestSelectLogline(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type selectLoglineData struct {
		resp *dao.LoglineEntity
		err  error
	}

	testCases := []struct {
		name string

		request services.SelectLoglineRequest

		selectLoglineData       *selectLoglineData
		selectLoglineBySlugData *selectLoglineData

		expect    *models.Logline
		expectErr error
	}{
		{
			name: "Success/ID",

			request: services.SelectLoglineRequest{
				ID:     lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				UserID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
			},

			selectLoglineData: &selectLoglineData{
				resp: &dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &models.Logline{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:      "test-slug",
				Name:      "Test Name 2",
				Content:   "Lorem ipsum dolor sit amet 2",
				Lang:      models.LangEN,
				CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Success/Slug",

			request: services.SelectLoglineRequest{
				Slug:   lo.ToPtr(models.Slug("test-slug")),
				UserID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
			},

			selectLoglineBySlugData: &selectLoglineData{
				resp: &dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &models.Logline{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Slug:      "test-slug",
				Name:      "Test Name 2",
				Content:   "Lorem ipsum dolor sit amet 2",
				Lang:      models.LangEN,
				CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Error/ID",

			request: services.SelectLoglineRequest{
				ID:     lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				UserID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
			},

			selectLoglineData: &selectLoglineData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "Error/Slug",

			request: services.SelectLoglineRequest{
				Slug:   lo.ToPtr(models.Slug("test-slug")),
				UserID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
			},

			selectLoglineBySlugData: &selectLoglineData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			source := servicesmocks.NewMockSelectLoglineSource(t)

			if testCase.selectLoglineData != nil {
				source.EXPECT().
					SelectLogline(mock.Anything, dao.SelectLoglineData{
						ID:     lo.FromPtr(testCase.request.ID),
						UserID: testCase.request.UserID,
					}).
					Return(testCase.selectLoglineData.resp, testCase.selectLoglineData.err)
			}

			if testCase.selectLoglineBySlugData != nil {
				source.EXPECT().
					SelectLoglineBySlug(mock.Anything, dao.SelectLoglineBySlugData{
						Slug:   lo.FromPtr(testCase.request.Slug),
						UserID: testCase.request.UserID,
					}).
					Return(testCase.selectLoglineBySlugData.resp, testCase.selectLoglineBySlugData.err)
			}

			service := services.NewSelectLoglineService(source)

			resp, err := service.SelectLogline(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			source.AssertExpectations(t)
		})
	}
}
