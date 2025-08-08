package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/internal/services"
	servicesmocks "github.com/a-novel/service-story-schematics/internal/services/mocks"
	"github.com/a-novel/service-story-schematics/models"
	storyplanmodel "github.com/a-novel/service-story-schematics/models/story_plan"
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
		resp *storyplanmodel.Plan
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
				LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Lang:      models.LangEN,
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
				resp: &storyplanmodel.Plan{
					Metadata: storyplanmodel.Metadata{
						Name: "Test Story Plan",
						Lang: models.LangEN,
					},
					Beats: []storyplanmodel.Beat{
						{
							Name: "Beat 1",
							Key:  "beat-1",
							KeyPoints: []string{
								"Key Point 1",
								"Key Point 2",
							},
							Purpose: "Purpose 1",
						},
					},
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
				LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Lang:      models.LangEN,
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
				resp: &storyplanmodel.Plan{
					Metadata: storyplanmodel.Metadata{
						Name: "Test Story Plan",
						Lang: models.LangEN,
					},
					Beats: []storyplanmodel.Beat{
						{
							Name: "Beat 1",
							Key:  "beat-1",
							KeyPoints: []string{
								"Key Point 1",
								"Key Point 2",
							},
							Purpose: "Purpose 1",
						},
					},
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
				LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Lang:      models.LangEN,
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
				LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				UserID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Lang:      models.LangEN,
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
						Plan:    testCase.selectStoryPlanData.resp,
						Lang:    testCase.request.Lang,
						UserID:  testCase.request.UserID.String(),
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
					SelectStoryPlan(
						mock.Anything,
						services.SelectStoryPlanRequest{Lang: testCase.request.Lang},
					).
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
