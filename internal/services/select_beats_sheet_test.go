package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	servicesmocks "github.com/a-novel/service-story-schematics/internal/services/mocks"
	"github.com/a-novel/service-story-schematics/models"
)

func TestSelectBeatsSheet(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type selectBeatsSheetData struct {
		resp *dao.BeatsSheetEntity
		err  error
	}

	type selectLoglineData struct {
		resp *dao.LoglineEntity
		err  error
	}

	testCases := []struct {
		name string

		request services.SelectBeatsSheetRequest

		selectBeatsSheetData *selectBeatsSheetData
		selectLoglineData    *selectLoglineData

		expect    *models.BeatsSheet
		expectErr error
	}{
		{
			name: "Success",

			request: services.SelectBeatsSheetRequest{
				BeatsSheetID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:       uuid.MustParse("00000000-0000-0000-1000-000000000001"),
			},

			selectBeatsSheetData: &selectBeatsSheetData{
				resp: &dao.BeatsSheetEntity{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					LoglineID:   uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					StoryPlanID: uuid.MustParse("00000000-0000-0000-0100-000000000001"),
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
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			selectLoglineData: &selectLoglineData{
				resp: &dao.LoglineEntity{
					ID:        uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					Slug:      "test-slug",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &models.BeatsSheet{
				ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				LoglineID:   uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				StoryPlanID: uuid.MustParse("00000000-0000-0000-0100-000000000001"),
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
				Lang:      models.LangEN,
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "SelectError",

			request: services.SelectBeatsSheetRequest{
				BeatsSheetID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:       uuid.MustParse("00000000-0000-0000-1000-000000000001"),
			},

			selectBeatsSheetData: &selectBeatsSheetData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "CheckLoglineError",

			request: services.SelectBeatsSheetRequest{
				BeatsSheetID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				UserID:       uuid.MustParse("00000000-0000-0000-1000-000000000001"),
			},

			selectBeatsSheetData: &selectBeatsSheetData{
				resp: &dao.BeatsSheetEntity{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					LoglineID:   uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					StoryPlanID: uuid.MustParse("00000000-0000-0000-0100-000000000001"),
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
					Lang:      models.LangEN,
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
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

			source := servicesmocks.NewMockSelectBeatsSheetSource(t)

			if testCase.selectLoglineData != nil {
				source.EXPECT().
					SelectLogline(ctx, dao.SelectLoglineData{
						ID:     testCase.selectBeatsSheetData.resp.LoglineID,
						UserID: testCase.request.UserID,
					}).
					Return(testCase.selectLoglineData.resp, testCase.selectLoglineData.err)
			}

			if testCase.selectBeatsSheetData != nil {
				source.EXPECT().
					SelectBeatsSheet(ctx, testCase.request.BeatsSheetID).
					Return(testCase.selectBeatsSheetData.resp, testCase.selectBeatsSheetData.err)
			}

			service := services.NewSelectBeatsSheetService(source)

			resp, err := service.SelectBeatsSheet(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			source.AssertExpectations(t)
		})
	}
}
