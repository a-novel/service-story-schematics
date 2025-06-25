package services_test

import (
	"errors"
	"github.com/stretchr/testify/mock"
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

func TestGenerateBeatsSheet(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type generateBeatsSheetData struct {
		resp []models.Beat
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

		request services.GenerateBeatsSheetRequest

		selectLoglineData      *selectLoglineData
		selectStoryPlanData    *selectStoryPlanData
		generateBeatsSheetData *generateBeatsSheetData

		expect    []models.Beat
		expectErr error
	}{
		{
			name: "Success",

			request: services.GenerateBeatsSheetRequest{
				LoglineID:   uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				StoryPlanID: uuid.MustParse("00000000-0000-1000-0000-000000000001"),
				UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Lang:        models.LangEN,
			},

			selectLoglineData: &selectLoglineData{
				resp: &dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-1000-000000000001"),
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
					ID:          uuid.MustParse("00000000-0000-1000-0000-000000000001"),
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

			generateBeatsSheetData: &generateBeatsSheetData{
				resp: []models.Beat{
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
			},

			expect: []models.Beat{
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
		},
		{
			name: "GenerateBeatsSheet/Error",

			request: services.GenerateBeatsSheetRequest{
				LoglineID:   uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				StoryPlanID: uuid.MustParse("00000000-0000-1000-0000-000000000001"),
				UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Lang:        models.LangEN,
			},

			selectLoglineData: &selectLoglineData{
				resp: &dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-1000-000000000001"),
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
					ID:          uuid.MustParse("00000000-0000-1000-0000-000000000001"),
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

			generateBeatsSheetData: &generateBeatsSheetData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "SelectStoryPlan/Error",

			request: services.GenerateBeatsSheetRequest{
				LoglineID:   uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				StoryPlanID: uuid.MustParse("00000000-0000-1000-0000-000000000001"),
				UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Lang:        models.LangEN,
			},

			selectLoglineData: &selectLoglineData{
				resp: &dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-0000-0000-1000-000000000001"),
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

			request: services.GenerateBeatsSheetRequest{
				LoglineID:   uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				StoryPlanID: uuid.MustParse("00000000-0000-1000-0000-000000000001"),
				UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Lang:        models.LangEN,
			},

			selectLoglineData: &selectLoglineData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			source := servicesmocks.NewMockGenerateBeatsSheetSource(t)

			if testCase.generateBeatsSheetData != nil {
				source.EXPECT().
					GenerateBeatsSheet(mock.Anything, daoai.GenerateBeatsSheetRequest{
						Logline: testCase.selectLoglineData.resp.Name + "\n\n" + testCase.selectLoglineData.resp.Content,
						Plan: models.StoryPlan{
							ID:          testCase.selectStoryPlanData.resp.ID,
							Slug:        testCase.selectStoryPlanData.resp.Slug,
							Name:        testCase.selectStoryPlanData.resp.Name,
							Description: testCase.selectStoryPlanData.resp.Description,
							Beats:       testCase.selectStoryPlanData.resp.Beats,
							Lang:        testCase.selectStoryPlanData.resp.Lang,
							CreatedAt:   testCase.selectStoryPlanData.resp.CreatedAt,
						},
						Lang:   testCase.request.Lang,
						UserID: testCase.request.UserID.String(),
					}).
					Return(testCase.generateBeatsSheetData.resp, testCase.generateBeatsSheetData.err)
			}

			if testCase.selectLoglineData != nil {
				source.EXPECT().
					SelectLogline(mock.Anything, dao.SelectLoglineData{
						ID:     testCase.request.LoglineID,
						UserID: testCase.request.UserID,
					}).
					Return(testCase.selectLoglineData.resp, testCase.selectLoglineData.err)
			}

			if testCase.selectStoryPlanData != nil {
				source.EXPECT().
					SelectStoryPlan(mock.Anything, testCase.request.StoryPlanID).
					Return(testCase.selectStoryPlanData.resp, testCase.selectStoryPlanData.err)
			}

			service := services.NewGenerateBeatsSheetService(source)

			resp, err := service.GenerateBeatsSheet(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			source.AssertExpectations(t)
		})
	}
}
