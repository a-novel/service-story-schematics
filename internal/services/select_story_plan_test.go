package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/story-schematics/internal/dao"
	"github.com/a-novel/story-schematics/internal/services"
	servicesmocks "github.com/a-novel/story-schematics/internal/services/mocks"
	"github.com/a-novel/story-schematics/models"
)

func TestSelectStoryPlan(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type selectStoryPlanData struct {
		resp *dao.StoryPlanEntity
		err  error
	}

	testCases := []struct {
		name string

		request services.SelectStoryPlanRequest

		selectStoryPlanData       *selectStoryPlanData
		selectStoryPlanBySlugData *selectStoryPlanData

		expect    *models.StoryPlan
		expectErr error
	}{
		{
			name: "Success/ID",

			request: services.SelectStoryPlanRequest{
				ID: lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
			},

			selectStoryPlanData: &selectStoryPlanData{
				resp: &dao.StoryPlanEntity{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug: "test-slug",

					Name:        "Test Name",
					Description: "Test Description, a lot going on here.",

					Beats: []models.BeatDefinition{
						{
							Name:      "Test Beat",
							Key:       "test-beat",
							KeyPoints: []string{"The key point of the beat, in a single sentence."},
							Purpose:   "The purpose of the beat, in a single sentence.",
						},
						{
							Name:      "Test Beat 2",
							Key:       "test-beat-2",
							KeyPoints: []string{"The key point of the second beat, in a single sentence."},
							Purpose:   "The purpose of the plot second point, in a single sentence.",
						},
					},

					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &models.StoryPlan{
				ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Slug: "test-slug",

				Name:        "Test Name",
				Description: "Test Description, a lot going on here.",

				Beats: []models.BeatDefinition{
					{
						Name:      "Test Beat",
						Key:       "test-beat",
						KeyPoints: []string{"The key point of the beat, in a single sentence."},
						Purpose:   "The purpose of the beat, in a single sentence.",
					},
					{
						Name:      "Test Beat 2",
						Key:       "test-beat-2",
						KeyPoints: []string{"The key point of the second beat, in a single sentence."},
						Purpose:   "The purpose of the plot second point, in a single sentence.",
					},
				},

				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Success/Slug",

			request: services.SelectStoryPlanRequest{
				Slug: lo.ToPtr(models.Slug("test-slug")),
			},

			selectStoryPlanBySlugData: &selectStoryPlanData{
				resp: &dao.StoryPlanEntity{
					ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug: "test-slug",

					Name:        "Test Name",
					Description: "Test Description, a lot going on here.",

					Beats: []models.BeatDefinition{
						{
							Name:      "Test Beat",
							Key:       "test-beat",
							KeyPoints: []string{"The key point of the beat, in a single sentence."},
							Purpose:   "The purpose of the beat, in a single sentence.",
						},
						{
							Name:      "Test Beat 2",
							Key:       "test-beat-2",
							KeyPoints: []string{"The key point of the second beat, in a single sentence."},
							Purpose:   "The purpose of the plot second point, in a single sentence.",
						},
					},

					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &models.StoryPlan{
				ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Slug: "test-slug",

				Name:        "Test Name",
				Description: "Test Description, a lot going on here.",

				Beats: []models.BeatDefinition{
					{
						Name:      "Test Beat",
						Key:       "test-beat",
						KeyPoints: []string{"The key point of the beat, in a single sentence."},
						Purpose:   "The purpose of the beat, in a single sentence.",
					},
					{
						Name:      "Test Beat 2",
						Key:       "test-beat-2",
						KeyPoints: []string{"The key point of the second beat, in a single sentence."},
						Purpose:   "The purpose of the plot second point, in a single sentence.",
					},
				},

				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Error/ID",

			request: services.SelectStoryPlanRequest{
				ID: lo.ToPtr(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
			},

			selectStoryPlanData: &selectStoryPlanData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "Error/Slug",

			request: services.SelectStoryPlanRequest{
				Slug: lo.ToPtr(models.Slug("test-slug")),
			},

			selectStoryPlanBySlugData: &selectStoryPlanData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			source := servicesmocks.NewMockSelectStoryPlanSource(t)

			if testCase.selectStoryPlanData != nil {
				source.EXPECT().
					SelectStoryPlan(ctx, lo.FromPtr(testCase.request.ID)).
					Return(testCase.selectStoryPlanData.resp, testCase.selectStoryPlanData.err)
			}

			if testCase.selectStoryPlanBySlugData != nil {
				source.EXPECT().
					SelectStoryPlanBySlug(ctx, lo.FromPtr(testCase.request.Slug)).
					Return(testCase.selectStoryPlanBySlugData.resp, testCase.selectStoryPlanBySlugData.err)
			}

			service := services.NewSelectStoryPlanService(source)

			resp, err := service.SelectStoryPlan(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			source.AssertExpectations(t)
		})
	}
}
