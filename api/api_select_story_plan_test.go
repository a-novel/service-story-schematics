package api_test

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/api"
	"github.com/a-novel/service-story-schematics/api/codegen"
	apimocks "github.com/a-novel/service-story-schematics/api/mocks"
	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

func TestSelectStoryPlan(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type selectStoryPlanData struct {
		resp *models.StoryPlan
		err  error
	}

	testCases := []struct {
		name string

		params codegen.GetStoryPlanParams

		selectStoryPlanData *selectStoryPlanData

		expect    codegen.GetStoryPlanRes
		expectErr error
	}{
		{
			name: "Success/ID",

			params: codegen.GetStoryPlanParams{
				ID: codegen.OptStoryPlanID{
					Value: codegen.StoryPlanID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					Set:   true,
				},
				Slug: codegen.OptSlug{},
			},

			selectStoryPlanData: &selectStoryPlanData{
				resp: &models.StoryPlan{
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
					Lang: models.LangEN,

					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &codegen.StoryPlan{
				ID:   codegen.StoryPlanID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				Slug: "test-slug",

				Name:        "Test Name",
				Description: "Test Description, a lot going on here.",

				Beats: []codegen.BeatDefinition{
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
				Lang: codegen.LangEn,

				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Success/Slug",

			params: codegen.GetStoryPlanParams{
				ID: codegen.OptStoryPlanID{},
				Slug: codegen.OptSlug{
					Value: "test-slug",
					Set:   true,
				},
			},

			selectStoryPlanData: &selectStoryPlanData{
				resp: &models.StoryPlan{
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
					Lang: models.LangEN,

					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &codegen.StoryPlan{
				ID:   codegen.StoryPlanID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				Slug: "test-slug",

				Name:        "Test Name",
				Description: "Test Description, a lot going on here.",

				Beats: []codegen.BeatDefinition{
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
				Lang: codegen.LangEn,

				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "LoglineNotFound",

			params: codegen.GetStoryPlanParams{
				ID: codegen.OptStoryPlanID{
					Value: codegen.StoryPlanID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					Set:   true,
				},
				Slug: codegen.OptSlug{},
			},

			selectStoryPlanData: &selectStoryPlanData{
				err: dao.ErrStoryPlanNotFound,
			},

			expect: &codegen.NotFoundError{Error: dao.ErrStoryPlanNotFound.Error()},
		},
		{
			name: "Error",

			params: codegen.GetStoryPlanParams{
				ID: codegen.OptStoryPlanID{
					Value: codegen.StoryPlanID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					Set:   true,
				},
				Slug: codegen.OptSlug{},
			},

			selectStoryPlanData: &selectStoryPlanData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			source := apimocks.NewMockSelectStoryPlanService(t)

			ctx := t.Context()

			if testCase.selectStoryPlanData != nil {
				source.EXPECT().
					SelectStoryPlan(mock.Anything, services.SelectStoryPlanRequest{
						Slug: lo.Ternary(
							testCase.params.Slug.IsSet(),
							lo.ToPtr(models.Slug(testCase.params.Slug.Value)),
							nil,
						),
						ID: lo.Ternary(
							testCase.params.ID.IsSet(),
							lo.ToPtr(uuid.UUID(testCase.params.ID.Value)),
							nil,
						),
					}).
					Return(testCase.selectStoryPlanData.resp, testCase.selectStoryPlanData.err)
			}

			handler := api.API{SelectStoryPlanService: source}

			res, err := handler.GetStoryPlan(ctx, testCase.params)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)

			source.AssertExpectations(t)
		})
	}
}
