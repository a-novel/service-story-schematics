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

func TestUpdateStoryPlan(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type updateStoryPlanData struct {
		resp *dao.StoryPlanEntity
		err  error
	}

	testCases := []struct {
		name string

		request services.UpdateStoryPlanRequest

		updateStoryPlanData *updateStoryPlanData

		expect    *models.StoryPlan
		expectErr error
	}{
		{
			name: "Success",

			request: services.UpdateStoryPlanRequest{
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

			updateStoryPlanData: &updateStoryPlanData{
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
			name: "Error",

			request: services.UpdateStoryPlanRequest{
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

			updateStoryPlanData: &updateStoryPlanData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			source := servicesmocks.NewMockUpdateStoryPlanSource(t)

			if testCase.updateStoryPlanData != nil {
				source.EXPECT().
					UpdateStoryPlan(ctx, mock.MatchedBy(func(data dao.UpdateStoryPlanData) bool {
						return assert.NotEqual(t, data.Plan.ID, uuid.Nil) &&
							assert.Equal(t, testCase.request.Slug, data.Plan.Slug) &&
							assert.Equal(t, testCase.request.Name, data.Plan.Name) &&
							assert.Equal(t, testCase.request.Description, data.Plan.Description) &&
							assert.Equal(t, testCase.request.Beats, data.Plan.Beats) &&
							assert.WithinDuration(t, time.Now(), data.Plan.CreatedAt, time.Second)
					})).
					Return(testCase.updateStoryPlanData.resp, testCase.updateStoryPlanData.err)
			}

			service := services.NewUpdateStoryPlanService(source)

			resp, err := service.UpdateStoryPlan(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			source.AssertExpectations(t)
		})
	}
}
