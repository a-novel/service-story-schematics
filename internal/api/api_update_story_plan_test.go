package api_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/api"
	apimocks "github.com/a-novel/service-story-schematics/internal/api/mocks"
	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

func TestUpdateStoryPlan(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type updateStoryPlanData struct {
		resp *models.StoryPlan
		err  error
	}

	testCases := []struct {
		name string

		form *apimodels.UpdateStoryPlanForm

		updateStoryPlanData *updateStoryPlanData

		expect    apimodels.UpdateStoryPlanRes
		expectErr error
	}{
		{
			name: "Success",

			form: &apimodels.UpdateStoryPlanForm{
				Slug:        "slug",
				Name:        "name",
				Description: "description",
				Beats: []apimodels.BeatDefinition{
					{
						Name:      "Beat 1",
						Key:       "beat-1",
						KeyPoints: []string{"beat 1 - key point 1", "beat 1 - key point 2"},
						Purpose:   "purpose",
						MinScenes: 1,
						MaxScenes: 2,
					},
					{
						Name:      "Beat 2",
						Key:       "beat-2",
						KeyPoints: []string{"beat 2 - key point 1", "beat 2 - key point 2"},
						Purpose:   "purpose",
						MinScenes: 1,
						MaxScenes: 1,
					},
				},
				Lang: apimodels.LangEn,
			},

			updateStoryPlanData: &updateStoryPlanData{
				resp: &models.StoryPlan{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug:        "slug",
					Name:        "name",
					Description: "description",
					Beats: []models.BeatDefinition{
						{
							Name:      "Beat 1",
							Key:       "beat-1",
							KeyPoints: []string{"beat 1 - key point 1", "beat 1 - key point 2"},
							Purpose:   "purpose",
							MinScenes: 1,
							MaxScenes: 2,
						},
						{
							Name:      "Beat 2",
							Key:       "beat-2",
							KeyPoints: []string{"beat 2 - key point 1", "beat 2 - key point 2"},
							Purpose:   "purpose",
							MinScenes: 1,
							MaxScenes: 1,
						},
					},
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &apimodels.StoryPlan{
				ID:          apimodels.StoryPlanID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				Slug:        "slug",
				Name:        "name",
				Description: "description",
				Beats: []apimodels.BeatDefinition{
					{
						Name:      "Beat 1",
						Key:       "beat-1",
						KeyPoints: []string{"beat 1 - key point 1", "beat 1 - key point 2"},
						Purpose:   "purpose",
						MinScenes: 1,
						MaxScenes: 2,
					},
					{
						Name:      "Beat 2",
						Key:       "beat-2",
						KeyPoints: []string{"beat 2 - key point 1", "beat 2 - key point 2"},
						Purpose:   "purpose",
						MinScenes: 1,
						MaxScenes: 1,
					},
				},
				Lang:      apimodels.LangEn,
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "StoryPlanNotFound",

			form: &apimodels.UpdateStoryPlanForm{
				Slug:        "slug",
				Name:        "name",
				Description: "description",
				Beats: []apimodels.BeatDefinition{
					{
						Name:      "Beat 1",
						Key:       "beat-1",
						KeyPoints: []string{"beat 1 - key point 1", "beat 1 - key point 2"},
						Purpose:   "purpose",
						MinScenes: 1,
						MaxScenes: 2,
					},
					{
						Name:      "Beat 2",
						Key:       "beat-2",
						KeyPoints: []string{"beat 2 - key point 1", "beat 2 - key point 2"},
						Purpose:   "purpose",
						MinScenes: 1,
						MaxScenes: 1,
					},
				},
				Lang: apimodels.LangEn,
			},

			updateStoryPlanData: &updateStoryPlanData{
				err: dao.ErrStoryPlanNotFound,
			},

			expect: &apimodels.NotFoundError{Error: dao.ErrStoryPlanNotFound.Error()},
		},
		{
			name: "Error",

			form: &apimodels.UpdateStoryPlanForm{
				Slug:        "slug",
				Name:        "name",
				Description: "description",
				Beats: []apimodels.BeatDefinition{
					{
						Name:      "Beat 1",
						Key:       "beat-1",
						KeyPoints: []string{"beat 1 - key point 1", "beat 1 - key point 2"},
						Purpose:   "purpose",
						MinScenes: 1,
						MaxScenes: 2,
					},
					{
						Name:      "Beat 2",
						Key:       "beat-2",
						KeyPoints: []string{"beat 2 - key point 1", "beat 2 - key point 2"},
						Purpose:   "purpose",
						MinScenes: 1,
						MaxScenes: 1,
					},
				},
				Lang: apimodels.LangEn,
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

			source := apimocks.NewMockUpdateStoryPlanService(t)

			if testCase.updateStoryPlanData != nil {
				source.EXPECT().
					UpdateStoryPlan(mock.Anything, services.UpdateStoryPlanRequest{
						Slug:        models.Slug(testCase.form.GetSlug()),
						Name:        testCase.form.GetName(),
						Description: testCase.form.GetDescription(),
						Beats: lo.Map(
							testCase.form.GetBeats(),
							func(item apimodels.BeatDefinition, _ int) models.BeatDefinition {
								return models.BeatDefinition{
									Name:      item.GetName(),
									Key:       item.GetKey(),
									KeyPoints: item.GetKeyPoints(),
									Purpose:   item.GetPurpose(),
									MinScenes: item.GetMinScenes(),
									MaxScenes: item.GetMaxScenes(),
								}
							},
						),
						Lang: models.Lang(testCase.form.GetLang()),
					}).
					Return(testCase.updateStoryPlanData.resp, testCase.updateStoryPlanData.err)
			}

			handler := api.API{UpdateStoryPlanService: source}

			res, err := handler.UpdateStoryPlan(t.Context(), testCase.form)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)

			source.AssertExpectations(t)
		})
	}
}
