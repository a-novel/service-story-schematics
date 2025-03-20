package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/story-schematics/internal/dao"
	"github.com/a-novel/story-schematics/internal/services"
	servicesmocks "github.com/a-novel/story-schematics/internal/services/mocks"
	"github.com/a-novel/story-schematics/models"
)

func TestCreateStoryPlan(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type insertStoryPlanData struct {
		resp *dao.StoryPlanEntity
		err  error
	}

	type selectSlugIterationData struct {
		slug      models.Slug
		iteration int
		err       error
	}

	testCases := []struct {
		name string

		request services.CreateStoryPlanRequest

		insertStoryPlanData     *insertStoryPlanData
		selectSlugIterationData *selectSlugIterationData
		reinsertStoryPlanData   *insertStoryPlanData

		expect    *models.StoryPlan
		expectErr error
	}{
		{
			name: "Success",

			request: services.CreateStoryPlanRequest{
				Slug:        "plan-slug",
				Name:        "Plan",
				Description: "Plan plan",
				Beats: []models.BeatDefinition{
					{
						Name: "beat-1",
						Key:  "beat-1-key",
						KeyPoints: []string{
							"key-point-1",
							"key-point-2",
						},
						Purpose: "beat 1 purpose",
					},
					{
						Name: "beat-2",
						Key:  "beat-2-key",
						KeyPoints: []string{
							"key-point-3",
							"key-point-4",
						},
						Purpose: "beat 2 purpose",
					},
				},
			},

			insertStoryPlanData: &insertStoryPlanData{
				resp: &dao.StoryPlanEntity{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug:        "plan-slug",
					Name:        "Plan",
					Description: "Plan plan",
					Beats: []models.BeatDefinition{
						{
							Name: "beat-1",
							Key:  "beat-1-key",
							KeyPoints: []string{
								"key-point-1",
								"key-point-2",
							},
							Purpose: "beat 1 purpose",
						},
						{
							Name: "beat-2",
							Key:  "beat-2-key",
							KeyPoints: []string{
								"key-point-3",
								"key-point-4",
							},
							Purpose: "beat 2 purpose",
						},
					},
					CreatedAt: time.Date(2019, time.August, 26, 1, 2, 0, 0, time.UTC),
				},
			},

			expect: &models.StoryPlan{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Slug:        "plan-slug",
				Name:        "Plan",
				Description: "Plan plan",
				Beats: []models.BeatDefinition{
					{
						Name: "beat-1",
						Key:  "beat-1-key",
						KeyPoints: []string{
							"key-point-1",
							"key-point-2",
						},
						Purpose: "beat 1 purpose",
					},
					{
						Name: "beat-2",
						Key:  "beat-2-key",
						KeyPoints: []string{
							"key-point-3",
							"key-point-4",
						},
						Purpose: "beat 2 purpose",
					},
				},
				CreatedAt: time.Date(2019, time.August, 26, 1, 2, 0, 0, time.UTC),
			},
		},
		{
			name: "RetrySlug",

			request: services.CreateStoryPlanRequest{
				Slug:        "plan-slug",
				Name:        "Plan",
				Description: "Plan plan",
				Beats: []models.BeatDefinition{
					{
						Name: "beat-1",
						Key:  "beat-1-key",
						KeyPoints: []string{
							"key-point-1",
							"key-point-2",
						},
						Purpose: "beat 1 purpose",
					},
					{
						Name: "beat-2",
						Key:  "beat-2-key",
						KeyPoints: []string{
							"key-point-3",
							"key-point-4",
						},
						Purpose: "beat 2 purpose",
					},
				},
			},

			insertStoryPlanData: &insertStoryPlanData{
				err: dao.ErrStoryPlanAlreadyExists,
			},

			selectSlugIterationData: &selectSlugIterationData{
				slug:      "plan-slug-2",
				iteration: 2,
			},

			reinsertStoryPlanData: &insertStoryPlanData{
				resp: &dao.StoryPlanEntity{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug:        "plan-slug-2",
					Name:        "Plan",
					Description: "Plan plan",
					Beats: []models.BeatDefinition{
						{
							Name: "beat-1",
							Key:  "beat-1-key",
							KeyPoints: []string{
								"key-point-1",
								"key-point-2",
							},
							Purpose: "beat 1 purpose",
						},
						{
							Name: "beat-2",
							Key:  "beat-2-key",
							KeyPoints: []string{
								"key-point-3",
								"key-point-4",
							},
							Purpose: "beat 2 purpose",
						},
					},
					CreatedAt: time.Date(2019, time.August, 26, 1, 2, 0, 0, time.UTC),
				},
			},

			expect: &models.StoryPlan{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Slug:        "plan-slug-2",
				Name:        "Plan",
				Description: "Plan plan",
				Beats: []models.BeatDefinition{
					{
						Name: "beat-1",
						Key:  "beat-1-key",
						KeyPoints: []string{
							"key-point-1",
							"key-point-2",
						},
						Purpose: "beat 1 purpose",
					},
					{
						Name: "beat-2",
						Key:  "beat-2-key",
						KeyPoints: []string{
							"key-point-3",
							"key-point-4",
						},
						Purpose: "beat 2 purpose",
					},
				},
				CreatedAt: time.Date(2019, time.August, 26, 1, 2, 0, 0, time.UTC),
			},
		},
		{
			name: "InsertError",

			request: services.CreateStoryPlanRequest{
				Slug:        "plan-slug",
				Name:        "Plan",
				Description: "Plan plan",
				Beats: []models.BeatDefinition{
					{
						Name: "beat-1",
						Key:  "beat-1-key",
						KeyPoints: []string{
							"key-point-1",
							"key-point-2",
						},
						Purpose: "beat 1 purpose",
					},
					{
						Name: "beat-2",
						Key:  "beat-2-key",
						KeyPoints: []string{
							"key-point-3",
							"key-point-4",
						},
						Purpose: "beat 2 purpose",
					},
				},
			},

			insertStoryPlanData: &insertStoryPlanData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "SlugIterationError",

			request: services.CreateStoryPlanRequest{
				Slug:        "plan-slug",
				Name:        "Plan",
				Description: "Plan plan",
				Beats: []models.BeatDefinition{
					{
						Name: "beat-1",
						Key:  "beat-1-key",
						KeyPoints: []string{
							"key-point-1",
							"key-point-2",
						},
						Purpose: "beat 1 purpose",
					},
					{
						Name: "beat-2",
						Key:  "beat-2-key",
						KeyPoints: []string{
							"key-point-3",
							"key-point-4",
						},
						Purpose: "beat 2 purpose",
					},
				},
			},

			insertStoryPlanData: &insertStoryPlanData{
				err: dao.ErrStoryPlanAlreadyExists,
			},

			selectSlugIterationData: &selectSlugIterationData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "RetryInsertError",

			request: services.CreateStoryPlanRequest{
				Slug:        "plan-slug",
				Name:        "Plan",
				Description: "Plan plan",
				Beats: []models.BeatDefinition{
					{
						Name: "beat-1",
						Key:  "beat-1-key",
						KeyPoints: []string{
							"key-point-1",
							"key-point-2",
						},
						Purpose: "beat 1 purpose",
					},
					{
						Name: "beat-2",
						Key:  "beat-2-key",
						KeyPoints: []string{
							"key-point-3",
							"key-point-4",
						},
						Purpose: "beat 2 purpose",
					},
				},
			},

			insertStoryPlanData: &insertStoryPlanData{
				err: dao.ErrStoryPlanAlreadyExists,
			},

			selectSlugIterationData: &selectSlugIterationData{
				slug:      "plan-slug-2",
				iteration: 2,
			},

			reinsertStoryPlanData: &insertStoryPlanData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			source := servicesmocks.NewMockCreateStoryPlanSource(t)

			if testCase.insertStoryPlanData != nil {
				initialCall := source.EXPECT().
					InsertStoryPlan(ctx, mock.MatchedBy(func(data dao.InsertStoryPlanData) bool {
						return assert.NotEqual(t, data.Plan.ID, uuid.Nil) &&
							testCase.request.Slug == data.Plan.Slug &&
							assert.Equal(t, testCase.request.Name, data.Plan.Name) &&
							assert.Equal(t, testCase.request.Description, data.Plan.Description) &&
							assert.Equal(t, testCase.request.Beats, data.Plan.Beats) &&
							assert.WithinDuration(t, time.Now(), data.Plan.CreatedAt, time.Second)
					})).
					Return(testCase.insertStoryPlanData.resp, testCase.insertStoryPlanData.err).
					Once()

				if testCase.reinsertStoryPlanData != nil {
					source.EXPECT().
						InsertStoryPlan(ctx, mock.MatchedBy(func(data dao.InsertStoryPlanData) bool {
							return assert.NotEqual(t, data.Plan.ID, uuid.Nil) &&
								testCase.selectSlugIterationData.slug == data.Plan.Slug &&
								assert.Equal(t, testCase.request.Name, data.Plan.Name) &&
								assert.Equal(t, testCase.request.Description, data.Plan.Description) &&
								assert.Equal(t, testCase.request.Beats, data.Plan.Beats) &&
								assert.WithinDuration(t, time.Now(), data.Plan.CreatedAt, time.Second)
						})).
						Return(testCase.reinsertStoryPlanData.resp, testCase.reinsertStoryPlanData.err).
						NotBefore(initialCall)
				}
			}

			if testCase.selectSlugIterationData != nil {
				source.EXPECT().
					SelectSlugIteration(ctx, dao.SelectSlugIterationData{
						Slug:  testCase.request.Slug,
						Table: "story_plans",
						Order: []string{"created_at DESC"},
					}).
					Return(
						testCase.selectSlugIterationData.slug,
						testCase.selectSlugIterationData.iteration,
						testCase.selectSlugIterationData.err,
					)
			}

			service := services.NewCreateStoryPlanService(source)

			resp, err := service.CreateStoryPlan(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			source.AssertExpectations(t)
		})
	}
}
