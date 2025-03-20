package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/a-novel/story-schematics/internal/dao"
	"github.com/a-novel/story-schematics/internal/services"
	servicesmocks "github.com/a-novel/story-schematics/internal/services/mocks"
	"github.com/a-novel/story-schematics/models"
)

func TestListLoglines(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type listLoglinesData struct {
		resp []*dao.LoglinePreviewEntity
		err  error
	}

	testCases := []struct {
		name string

		request services.ListLoglinesRequest

		listLoglinesData *listLoglinesData

		expect    []*models.LoglinePreview
		expectErr error
	}{
		{
			name: "Success",

			request: services.ListLoglinesRequest{
				UserID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Limit:  10,
				Offset: 20,
			},

			listLoglinesData: &listLoglinesData{
				resp: []*dao.LoglinePreviewEntity{
					{
						Slug:      "test-slug",
						Name:      "Test Name",
						Content:   "Lorem ipsum dolor sit amet",
						CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
					},
					{
						Slug:      "test-slug-3",
						Name:      "Test Name 3",
						Content:   "Lorem ipsum dolor sit amet 3",
						CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					},
					{
						Slug:      "test-slug-2",
						Name:      "Test Name 2",
						Content:   "Lorem ipsum dolor sit amet 2",
						CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},

			expect: []*models.LoglinePreview{
				{
					Slug:      "test-slug",
					Name:      "Test Name",
					Content:   "Lorem ipsum dolor sit amet",
					CreatedAt: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					Slug:      "test-slug-3",
					Name:      "Test Name 3",
					Content:   "Lorem ipsum dolor sit amet 3",
					CreatedAt: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					Slug:      "test-slug-2",
					Name:      "Test Name 2",
					Content:   "Lorem ipsum dolor sit amet 2",
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "Error",

			request: services.ListLoglinesRequest{
				UserID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
				Limit:  10,
				Offset: 20,
			},

			listLoglinesData: &listLoglinesData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			source := servicesmocks.NewMockListLoglinesSource(t)

			if testCase.listLoglinesData != nil {
				source.EXPECT().
					ListLoglines(ctx, dao.ListLoglinesData{
						UserID: testCase.request.UserID,
						Limit:  testCase.request.Limit,
						Offset: testCase.request.Offset,
					}).
					Return(testCase.listLoglinesData.resp, testCase.listLoglinesData.err)
			}

			service := services.NewListLoglinesService(source)

			resp, err := service.ListLoglines(ctx, testCase.request)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, resp)

			source.AssertExpectations(t)
		})
	}
}
