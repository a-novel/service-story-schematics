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
	"github.com/a-novel/story-schematics/internal/lib"
	"github.com/a-novel/story-schematics/internal/services"
	servicesmocks "github.com/a-novel/story-schematics/internal/services/mocks"
	"github.com/a-novel/story-schematics/models"
)

func TestCreateBeatsSheet(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type insertBeatsSheetData struct {
		resp *dao.BeatsSheetEntity
		err  error
	}

	type selectStoryPlanData struct {
		resp *dao.StoryPlanEntity
		err  error
	}

	type selectLoglineData struct {
		resp *dao.LoglineEntity
		err  error
	}

	testCases := []struct {
		name string

		request services.CreateBeatsSheetRequest

		selectLoglineData    *selectLoglineData
		selectStoryPlanData  *selectStoryPlanData
		insertBeatsSheetData *insertBeatsSheetData

		expect    *models.BeatsSheet
		expectErr error
	}{
		{
			name: "Success",

			request: services.CreateBeatsSheetRequest{
				UserID:      uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				StoryPlanID: uuid.MustParse("00000000-0000-1000-0000-000000000001"),
				LoglineID:   uuid.MustParse("00000000-1000-0000-0000-000000000001"),
				Content: []models.Beat{
					{
						Key:     "test-beat",
						Title:   "Test Beat",
						Content: "Test Beat Content",
					},
					{
						Key:     "test-beat-2",
						Title:   "Test Beat 2",
						Content: "Test Beat Content 2",
					},
				},
			},

			selectLoglineData: &selectLoglineData{
				resp: &dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			selectStoryPlanData: &selectStoryPlanData{
				resp: &dao.StoryPlanEntity{
					ID:          uuid.MustParse("00000000-0000-1000-0000-000000000001"),
					Slug:        "test-slug",
					Name:        "Test Name",
					Description: "Lorem ipsum dolor sit amet",
					Beats: []models.BeatDefinition{
						{
							Name: "Test Beat",
							Key:  "test-beat",
							KeyPoints: []string{
								"Test Key Point",
							},
							Purpose: "Test Purpose",
						},
						{
							Name: "Test Beat 2",
							Key:  "test-beat-2",
							KeyPoints: []string{
								"Test Key Point 2",
							},
							Purpose: "Test Purpose 2",
						},
					},
				},
			},

			insertBeatsSheetData: &insertBeatsSheetData{
				resp: &dao.BeatsSheetEntity{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					LoglineID:   uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					StoryPlanID: uuid.MustParse("00000000-0000-1000-0000-000000000001"),
					Content: []models.Beat{
						{
							Key:     "test-beat",
							Title:   "Test Beat",
							Content: "Test Beat Content",
						},
						{
							Key:     "test-beat-2",
							Title:   "Test Beat 2",
							Content: "Test Beat Content 2",
						},
					},
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &models.BeatsSheet{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				LoglineID:   uuid.MustParse("00000000-1000-0000-0000-000000000001"),
				StoryPlanID: uuid.MustParse("00000000-0000-1000-0000-000000000001"),
				Content: []models.Beat{
					{
						Key:     "test-beat",
						Title:   "Test Beat",
						Content: "Test Beat Content",
					},
					{
						Key:     "test-beat-2",
						Title:   "Test Beat 2",
						Content: "Test Beat Content 2",
					},
				},
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Insert/Error",

			request: services.CreateBeatsSheetRequest{
				UserID:      uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				StoryPlanID: uuid.MustParse("00000000-0000-1000-0000-000000000001"),
				LoglineID:   uuid.MustParse("00000000-1000-0000-0000-000000000001"),
				Content: []models.Beat{
					{
						Key:     "test-beat",
						Title:   "Test Beat",
						Content: "Test Beat Content",
					},
					{
						Key:     "test-beat-2",
						Title:   "Test Beat 2",
						Content: "Test Beat Content 2",
					},
				},
			},

			selectLoglineData: &selectLoglineData{
				resp: &dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			selectStoryPlanData: &selectStoryPlanData{
				resp: &dao.StoryPlanEntity{
					ID:          uuid.MustParse("00000000-0000-1000-0000-000000000001"),
					Slug:        "test-slug",
					Name:        "Test Name",
					Description: "Lorem ipsum dolor sit amet",
					Beats: []models.BeatDefinition{
						{
							Name: "Test Beat",
							Key:  "test-beat",
							KeyPoints: []string{
								"Test Key Point",
							},
							Purpose: "Test Purpose",
						},
						{
							Name: "Test Beat 2",
							Key:  "test-beat-2",
							KeyPoints: []string{
								"Test Key Point 2",
							},
							Purpose: "Test Purpose 2",
						},
					},
				},
			},

			insertBeatsSheetData: &insertBeatsSheetData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "SelectStoryPlan/Error",

			request: services.CreateBeatsSheetRequest{
				UserID:      uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				StoryPlanID: uuid.MustParse("00000000-0000-1000-0000-000000000001"),
				LoglineID:   uuid.MustParse("00000000-1000-0000-0000-000000000001"),
				Content: []models.Beat{
					{
						Key:     "test-beat",
						Title:   "Test Beat",
						Content: "Test Beat Content",
					},
					{
						Key:     "test-beat-2",
						Title:   "Test Beat 2",
						Content: "Test Beat Content 2",
					},
				},
			},

			selectLoglineData: &selectLoglineData{
				resp: &dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			selectStoryPlanData: &selectStoryPlanData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "CheckLogline/Error",

			request: services.CreateBeatsSheetRequest{
				UserID:      uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				StoryPlanID: uuid.MustParse("00000000-0000-1000-0000-000000000001"),
				LoglineID:   uuid.MustParse("00000000-1000-0000-0000-000000000001"),
				Content: []models.Beat{
					{
						Key:     "test-beat",
						Title:   "Test Beat",
						Content: "Test Beat Content",
					},
					{
						Key:     "test-beat-2",
						Title:   "Test Beat 2",
						Content: "Test Beat Content 2",
					},
				},
			},

			selectLoglineData: &selectLoglineData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "StoryPlanMismatch",

			request: services.CreateBeatsSheetRequest{
				UserID:      uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				StoryPlanID: uuid.MustParse("00000000-0000-1000-0000-000000000001"),
				LoglineID:   uuid.MustParse("00000000-1000-0000-0000-000000000001"),
				Content: []models.Beat{
					{
						Key:     "test-beat",
						Title:   "Test Beat",
						Content: "Test Beat Content",
					},
					{
						Key:     "test-beat-2",
						Title:   "Test Beat 2",
						Content: "Test Beat Content 2",
					},
				},
			},

			selectLoglineData: &selectLoglineData{
				resp: &dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			selectStoryPlanData: &selectStoryPlanData{
				resp: &dao.StoryPlanEntity{
					ID:          uuid.MustParse("00000000-0000-1000-0000-000000000001"),
					Slug:        "test-slug",
					Name:        "Test Name",
					Description: "Lorem ipsum dolor sit amet",
					Beats: []models.BeatDefinition{
						{
							Name: "Test Beat",
							Key:  "test-beat",
							KeyPoints: []string{
								"Test Key Point",
							},
							Purpose: "Test Purpose",
						},
					},
				},
			},

			expectErr: lib.ErrInvalidStoryPlan,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			source := servicesmocks.NewMockCreateBeatsSheetSource(t)

			if testCase.selectLoglineData != nil {
				source.EXPECT().
					SelectLogline(ctx, dao.SelectLoglineData{
						ID:     testCase.request.LoglineID,
						UserID: testCase.request.UserID,
					}).
					Return(testCase.selectLoglineData.resp, testCase.selectLoglineData.err)
			}

			if testCase.selectStoryPlanData != nil {
				source.EXPECT().
					SelectStoryPlan(ctx, testCase.request.StoryPlanID).
					Return(testCase.selectStoryPlanData.resp, testCase.selectStoryPlanData.err)
			}

			if testCase.insertBeatsSheetData != nil {
				source.EXPECT().
					InsertBeatsSheet(ctx, mock.MatchedBy(func(data dao.InsertBeatsSheetData) bool {
						return assert.NotEqual(t, data.Sheet.ID, uuid.Nil) &&
							assert.Equal(t, testCase.request.LoglineID, data.Sheet.LoglineID) &&
							assert.Equal(t, testCase.request.StoryPlanID, data.Sheet.StoryPlanID) &&
							assert.Equal(t, testCase.request.Content, data.Sheet.Content) &&
							assert.WithinDuration(t, time.Now(), data.Sheet.CreatedAt, time.Second)
					})).
					Return(testCase.insertBeatsSheetData.resp, testCase.insertBeatsSheetData.err)
			}

			service := services.NewCreateBeatsSheetService(source)

			resp, err := service.CreateBeatsSheet(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			source.AssertExpectations(t)
		})
	}
}
