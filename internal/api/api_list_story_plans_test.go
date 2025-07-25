package api_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/api"
	apimocks "github.com/a-novel/service-story-schematics/internal/api/mocks"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

func TestListStoryPlans(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type listStoryPlansData struct {
		resp []*models.StoryPlanPreview
		err  error
	}

	testCases := []struct {
		name string

		params apimodels.GetStoryPlansParams

		listStoryPlansData *listStoryPlansData

		expect    apimodels.GetStoryPlansRes
		expectErr error
	}{
		{
			name: "Success",

			params: apimodels.GetStoryPlansParams{
				Limit:  apimodels.OptInt{Value: 10, Set: true},
				Offset: apimodels.OptInt{Value: 2, Set: true},
			},

			listStoryPlansData: &listStoryPlansData{
				resp: []*models.StoryPlanPreview{
					{
						ID:   uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						Slug: "test-slug-2",

						Name:        "Test Name 2",
						Description: "Test Description 2, a lot going on here.",
						Lang:        models.LangEN,

						CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:   uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						Slug: "test-slug-3",

						Name:        "Test Name 3",
						Description: "Test Description 3, a lot going on here.",
						Lang:        models.LangEN,

						CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
					},
				},
			},

			expect: &apimodels.GetStoryPlansOKApplicationJSON{
				{
					ID:   apimodels.StoryPlanID(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
					Slug: "test-slug-2",

					Name:        "Test Name 2",
					Description: "Test Description 2, a lot going on here.",
					Lang:        apimodels.LangEn,

					CreatedAt: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:   apimodels.StoryPlanID(uuid.MustParse("00000000-0000-0000-0000-000000000003")),
					Slug: "test-slug-3",

					Name:        "Test Name 3",
					Description: "Test Description 3, a lot going on here.",
					Lang:        apimodels.LangEn,

					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "Error",

			params: apimodels.GetStoryPlansParams{
				Limit:  apimodels.OptInt{Value: 10, Set: true},
				Offset: apimodels.OptInt{Value: 2, Set: true},
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

			source := apimocks.NewMockListStoryPlansService(t)

			ctx := t.Context()

			if testCase.listStoryPlansData != nil {
				source.EXPECT().
					ListStoryPlans(mock.Anything, services.ListStoryPlansRequest{
						Limit:  testCase.params.Limit.Value,
						Offset: testCase.params.Offset.Value,
					}).
					Return(testCase.listStoryPlansData.resp, testCase.listStoryPlansData.err)
			}

			handler := api.API{ListStoryPlansService: source}

			res, err := handler.GetStoryPlans(ctx, testCase.params)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)

			source.AssertExpectations(t)
		})
	}
}
