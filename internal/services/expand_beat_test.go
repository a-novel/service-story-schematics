package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/internal/services"
	servicesmocks "github.com/a-novel/service-story-schematics/internal/services/mocks"
	"github.com/a-novel/service-story-schematics/models"
)

func TestExpandBeat(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type expandBeatData struct {
		resp *models.Beat
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

		request services.ExpandBeatRequest

		selectBeatsSheetData *selectBeatsSheetData
		selectLoglineData    *selectLoglineData
		selectStoryPlanData  *selectStoryPlanData
		expandBeatData       *expandBeatData

		expect    *models.Beat
		expectErr error
	}{
		{
			name: "Success",

			request: services.ExpandBeatRequest{
				BeatsSheetID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				TargetKey:    "test",
				UserID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
					Lang:      models.LangEN,
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
					Lang:      models.LangEN,
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
					Lang:      models.LangEN,
					CreatedAt: time.Now(),
				},
			},

			expandBeatData: &expandBeatData{
				resp: &models.Beat{
					Key:     "beat-1",
					Title:   "Generated Beat 1 (expanded)",
					Content: "Generated Content 1 (expanded)",
				},
			},

			expect: &models.Beat{
				Key:     "beat-1",
				Title:   "Generated Beat 1 (expanded)",
				Content: "Generated Content 1 (expanded)",
			},
		},
		{
			name: "ExpandBeat/Error",

			request: services.ExpandBeatRequest{
				BeatsSheetID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				TargetKey:    "test",
				UserID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
					Lang:      models.LangEN,
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
					Lang:      models.LangEN,
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
					Lang:      models.LangEN,
					CreatedAt: time.Now(),
				},
			},

			expandBeatData: &expandBeatData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "SelectStoryPlan/Error",

			request: services.ExpandBeatRequest{
				BeatsSheetID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				TargetKey:    "test",
				UserID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
					Lang:      models.LangEN,
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
					Lang:      models.LangEN,
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

			request: services.ExpandBeatRequest{
				BeatsSheetID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				TargetKey:    "test",
				UserID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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
					Lang:      models.LangEN,
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

			request: services.ExpandBeatRequest{
				BeatsSheetID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				TargetKey:    "test",
				UserID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
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

			source := servicesmocks.NewMockExpandBeatSource(t)

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

			if testCase.expandBeatData != nil {
				source.EXPECT().
					ExpandBeat(ctx, daoai.ExpandBeatRequest{
						Logline: testCase.selectLoglineData.resp.Name + "\n\n" + testCase.selectLoglineData.resp.Content,
						Beats:   testCase.selectBeatsSheetData.resp.Content,
						Plan: models.StoryPlan{
							ID:          testCase.selectStoryPlanData.resp.ID,
							Slug:        testCase.selectStoryPlanData.resp.Slug,
							Name:        testCase.selectStoryPlanData.resp.Name,
							Description: testCase.selectStoryPlanData.resp.Description,
							Beats:       testCase.selectStoryPlanData.resp.Beats,
							Lang:        testCase.selectBeatsSheetData.resp.Lang,
							CreatedAt:   testCase.selectStoryPlanData.resp.CreatedAt,
						},
						Lang:      testCase.selectBeatsSheetData.resp.Lang,
						TargetKey: testCase.request.TargetKey,
						UserID:    testCase.request.UserID.String(),
					}).
					Return(testCase.expandBeatData.resp, testCase.expandBeatData.err)
			}

			service := services.NewExpandBeatService(source)

			resp, err := service.ExpandBeat(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			source.AssertExpectations(t)
		})
	}
}
