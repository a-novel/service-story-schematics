package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/story-schematics/internal/dao"
	"github.com/a-novel/story-schematics/internal/daoai"
	"github.com/a-novel/story-schematics/internal/services"
	servicesmocks "github.com/a-novel/story-schematics/internal/services/mocks"
	"github.com/a-novel/story-schematics/models"
)

func TestRegenerateBeats(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type regenerateBeatsData struct {
		resp []models.Beat
		err  error
	}

	type selectBeatsSheetData struct {
		resp *dao.BeatsSheetEntity
		err  error
	}

	type selectLoglineData struct {
		resp *dao.LoglineEntity
		err  error
	}

	type selectStoryPlanData struct {
		resp *dao.StoryPlanEntity
		err  error
	}

	testCases := []struct {
		name string

		request services.RegenerateBeatsRequest

		selectBeatsSheetData *selectBeatsSheetData
		selectLoglineData    *selectLoglineData
		selectStoryPlanData  *selectStoryPlanData
		regenerateBeatsData  *regenerateBeatsData

		expect    []models.Beat
		expectErr error
	}{
		{
			name: "Success",

			request: services.RegenerateBeatsRequest{
				BeatsSheetID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				UserID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				RegenerateKeys: []string{
					"beat-1",
					"beat-2",
				},
			},

			selectBeatsSheetData: &selectBeatsSheetData{
				resp: &dao.BeatsSheetEntity{
					ID:          uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					LoglineID:   uuid.MustParse("00000000-0000-1000-0000-000000000001"),
					StoryPlanID: uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					Content: []models.Beat{
						{
							Key:     "beat-1",
							Title:   "Generated Beat 1",
							Content: "Generated Content 1",
						},
						{
							Key:     "beat-2",
							Title:   "Generated Beat 2",
							Content: "Generated Content 2",
						},
					},
					CreatedAt: time.Now(),
				},
			},

			selectLoglineData: &selectLoglineData{
				resp: &dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-1000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug:      "logline-1",
					Name:      "Logline 1",
					Content:   "Content 1",
					CreatedAt: time.Now(),
				},
			},

			selectStoryPlanData: &selectStoryPlanData{
				resp: &dao.StoryPlanEntity{
					ID:          uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					Slug:        "story-plan-1",
					Name:        "Story Plan 1",
					Description: "Description 1",
					Beats: []models.BeatDefinition{
						{
							Name: "Beat 1",
							Key:  "beat-1",
							KeyPoints: []string{
								"Key Point 1",
								"Key Point 2",
							},
							Purpose:   "Purpose 1",
							MinScenes: 1,
						},
					},
					CreatedAt: time.Now(),
				},
			},

			regenerateBeatsData: &regenerateBeatsData{
				resp: []models.Beat{
					{
						Key:     "beat-1",
						Title:   "Regenerated Beat 1",
						Content: "Regenerated Content 1",
					},
					{
						Key:     "beat-2",
						Title:   "Regenerated Beat 2",
						Content: "Regenerated Content 2",
					},
				},
			},

			expect: []models.Beat{
				{
					Key:     "beat-1",
					Title:   "Regenerated Beat 1",
					Content: "Regenerated Content 1",
				},
				{
					Key:     "beat-2",
					Title:   "Regenerated Beat 2",
					Content: "Regenerated Content 2",
				},
			},
		},
		{
			name: "RegenerateBeats/Error",

			request: services.RegenerateBeatsRequest{
				BeatsSheetID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				UserID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				RegenerateKeys: []string{
					"beat-1",
					"beat-2",
				},
			},

			selectBeatsSheetData: &selectBeatsSheetData{
				resp: &dao.BeatsSheetEntity{
					ID:          uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					LoglineID:   uuid.MustParse("00000000-0000-1000-0000-000000000001"),
					StoryPlanID: uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					Content: []models.Beat{
						{
							Key:     "beat-1",
							Title:   "Generated Beat 1",
							Content: "Generated Content 1",
						},
						{
							Key:     "beat-2",
							Title:   "Generated Beat 2",
							Content: "Generated Content 2",
						},
					},
					CreatedAt: time.Now(),
				},
			},

			selectLoglineData: &selectLoglineData{
				resp: &dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-1000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug:      "logline-1",
					Name:      "Logline 1",
					Content:   "Content 1",
					CreatedAt: time.Now(),
				},
			},

			selectStoryPlanData: &selectStoryPlanData{
				resp: &dao.StoryPlanEntity{
					ID:          uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					Slug:        "story-plan-1",
					Name:        "Story Plan 1",
					Description: "Description 1",
					Beats: []models.BeatDefinition{
						{
							Name: "Beat 1",
							Key:  "beat-1",
							KeyPoints: []string{
								"Key Point 1",
								"Key Point 2",
							},
							Purpose:   "Purpose 1",
							MinScenes: 1,
						},
					},
					CreatedAt: time.Now(),
				},
			},

			regenerateBeatsData: &regenerateBeatsData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "SelectStoryPlan/Error",

			request: services.RegenerateBeatsRequest{
				BeatsSheetID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				UserID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				RegenerateKeys: []string{
					"beat-1",
					"beat-2",
				},
			},

			selectBeatsSheetData: &selectBeatsSheetData{
				resp: &dao.BeatsSheetEntity{
					ID:          uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					LoglineID:   uuid.MustParse("00000000-0000-1000-0000-000000000001"),
					StoryPlanID: uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					Content: []models.Beat{
						{
							Key:     "beat-1",
							Title:   "Generated Beat 1",
							Content: "Generated Content 1",
						},
						{
							Key:     "beat-2",
							Title:   "Generated Beat 2",
							Content: "Generated Content 2",
						},
					},
					CreatedAt: time.Now(),
				},
			},

			selectLoglineData: &selectLoglineData{
				resp: &dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-1000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Slug:      "logline-1",
					Name:      "Logline 1",
					Content:   "Content 1",
					CreatedAt: time.Now(),
				},
			},

			selectStoryPlanData: &selectStoryPlanData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "SelectLogline/Error",

			request: services.RegenerateBeatsRequest{
				BeatsSheetID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				UserID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				RegenerateKeys: []string{
					"beat-1",
					"beat-2",
				},
			},

			selectBeatsSheetData: &selectBeatsSheetData{
				resp: &dao.BeatsSheetEntity{
					ID:          uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					LoglineID:   uuid.MustParse("00000000-0000-1000-0000-000000000001"),
					StoryPlanID: uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					Content: []models.Beat{
						{
							Key:     "beat-1",
							Title:   "Generated Beat 1",
							Content: "Generated Content 1",
						},
						{
							Key:     "beat-2",
							Title:   "Generated Beat 2",
							Content: "Generated Content 2",
						},
					},
					CreatedAt: time.Now(),
				},
			},

			selectLoglineData: &selectLoglineData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "SelectBeatsSheet/Error",

			request: services.RegenerateBeatsRequest{
				BeatsSheetID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				UserID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				RegenerateKeys: []string{
					"beat-1",
					"beat-2",
				},
			},

			selectBeatsSheetData: &selectBeatsSheetData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			source := servicesmocks.NewMockRegenerateBeatsSource(t)

			if testCase.selectBeatsSheetData != nil {
				source.EXPECT().
					SelectBeatsSheet(ctx, testCase.request.BeatsSheetID).
					Return(testCase.selectBeatsSheetData.resp, testCase.selectBeatsSheetData.err)
			}

			if testCase.selectLoglineData != nil {
				source.EXPECT().
					SelectLogline(ctx, dao.SelectLoglineData{
						ID:     testCase.selectBeatsSheetData.resp.LoglineID,
						UserID: testCase.request.UserID,
					}).
					Return(testCase.selectLoglineData.resp, testCase.selectLoglineData.err)
			}

			if testCase.selectStoryPlanData != nil {
				source.EXPECT().
					SelectStoryPlan(ctx, testCase.selectBeatsSheetData.resp.StoryPlanID).
					Return(testCase.selectStoryPlanData.resp, testCase.selectStoryPlanData.err)
			}

			if testCase.regenerateBeatsData != nil {
				source.EXPECT().
					RegenerateBeats(ctx, daoai.RegenerateBeatsRequest{
						Logline: testCase.selectLoglineData.resp.Name + "\n\n" + testCase.selectLoglineData.resp.Content,
						Plan: models.StoryPlan{
							ID:          testCase.selectStoryPlanData.resp.ID,
							Slug:        testCase.selectStoryPlanData.resp.Slug,
							Name:        testCase.selectStoryPlanData.resp.Name,
							Description: testCase.selectStoryPlanData.resp.Description,
							Beats:       testCase.selectStoryPlanData.resp.Beats,
							CreatedAt:   testCase.selectStoryPlanData.resp.CreatedAt,
						},
						Beats:          testCase.selectBeatsSheetData.resp.Content,
						RegenerateKeys: testCase.request.RegenerateKeys,
					}).
					Return(testCase.regenerateBeatsData.resp, testCase.regenerateBeatsData.err)
			}

			service := services.NewRegenerateBeatsService(source)

			resp, err := service.RegenerateBeats(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			source.AssertExpectations(t)
		})
	}
}
