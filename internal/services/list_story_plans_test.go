package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/story-schematics/internal/dao"
	"github.com/a-novel/story-schematics/internal/services"
	servicesmocks "github.com/a-novel/story-schematics/internal/services/mocks"
	"github.com/a-novel/story-schematics/models"
)

func TestListStoryPlans(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type listStoryPlansData struct {
		resp []*dao.StoryPlanPreviewEntity
		err  error
	}

	testCases := []struct {
		name string

		request services.ListStoryPlansRequest

		listStoryPlansData *listStoryPlansData

		expect    []*models.StoryPlanPreview
		expectErr error
	}{
		{
			name: "Success",

			request: services.ListStoryPlansRequest{
				Limit:  10,
				Offset: 20,
			},

			listStoryPlansData: &listStoryPlansData{
				resp: []*dao.StoryPlanPreviewEntity{
					{
						ID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						Slug: "test-slug-2",

						Name:        "Test Name 2",
						Description: "Test Description 2, a lot going on here.",

						CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:   uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						Slug: "test-slug-3",

						Name:        "Test Name 3",
						Description: "Test Description 3, a lot going on here.",

						CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Slug: "test-slug-1",

						Name:        "Test Name",
						Description: "Test Description, a lot going on here.",

						CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},

			expect: []*models.StoryPlanPreview{
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Slug: "test-slug-2",

					Name:        "Test Name 2",
					Description: "Test Description 2, a lot going on here.",

					CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					Slug: "test-slug-3",

					Name:        "Test Name 3",
					Description: "Test Description 3, a lot going on here.",

					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug: "test-slug-1",

					Name:        "Test Name",
					Description: "Test Description, a lot going on here.",

					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "Error",

			request: services.ListStoryPlansRequest{
				Limit:  10,
				Offset: 20,
			},

			listStoryPlansData: &listStoryPlansData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			source := servicesmocks.NewMockListStoryPlansSource(t)

			if testCase.listStoryPlansData != nil {
				source.EXPECT().
					ListStoryPlans(ctx, dao.ListStoryPlansData{
						Limit:  testCase.request.Limit,
						Offset: testCase.request.Offset,
					}).
					Return(testCase.listStoryPlansData.resp, testCase.listStoryPlansData.err)
			}

			service := services.NewListStoryPlansService(source)

			resp, err := service.ListStoryPlans(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			source.AssertExpectations(t)
		})
	}
}
