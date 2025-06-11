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

func TestListBeatsSheets(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type listBeatsSheetsData struct {
		resp []*dao.BeatsSheetPreviewEntity
		err  error
	}

	type selectLoglineData struct {
		resp *dao.LoglineEntity
		err  error
	}

	testCases := []struct {
		name string

		request services.ListBeatsSheetsRequest

		listBeatsSheetsData *listBeatsSheetsData
		selectLoglineData   *selectLoglineData

		expect    []*models.BeatsSheetPreview
		expectErr error
	}{
		{
			name: "Success",

			request: services.ListBeatsSheetsRequest{
				UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				LoglineID: uuid.MustParse("00000000-0000-1000-0000-000000000001"),
				Limit:     10,
				Offset:    20,
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

			listBeatsSheetsData: &listBeatsSheetsData{
				resp: []*dao.BeatsSheetPreviewEntity{
					{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Lang:      models.LangEN,
						CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						Lang:      models.LangEN,
						CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					},
				},
			},

			expect: []*models.BeatsSheetPreview{
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					Lang:      models.LangEN,
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "ListError",

			request: services.ListBeatsSheetsRequest{
				UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				LoglineID: uuid.MustParse("00000000-0000-1000-0000-000000000001"),
				Limit:     10,
				Offset:    20,
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

			listBeatsSheetsData: &listBeatsSheetsData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
		{
			name: "CheckLoglineError",

			request: services.ListBeatsSheetsRequest{
				UserID:    uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				LoglineID: uuid.MustParse("00000000-0000-1000-0000-000000000001"),
				Limit:     10,
				Offset:    20,
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

			source := servicesmocks.NewMockListBeatsSheetsSource(t)

			if testCase.listBeatsSheetsData != nil {
				source.EXPECT().
					ListBeatsSheets(ctx, dao.ListBeatsSheetsData{
						LoglineID: testCase.request.LoglineID,
						Limit:     testCase.request.Limit,
						Offset:    testCase.request.Offset,
					}).
					Return(testCase.listBeatsSheetsData.resp, testCase.listBeatsSheetsData.err)
			}

			if testCase.selectLoglineData != nil {
				source.EXPECT().
					SelectLogline(ctx, dao.SelectLoglineData{
						ID:     testCase.request.LoglineID,
						UserID: testCase.request.UserID,
					}).
					Return(testCase.selectLoglineData.resp, testCase.selectLoglineData.err)
			}

			service := services.NewListBeatsSheetsService(source)

			resp, err := service.ListBeatsSheets(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			source.AssertExpectations(t)
		})
	}
}
