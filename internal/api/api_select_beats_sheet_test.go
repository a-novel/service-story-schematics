package api_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	authmodels "github.com/a-novel/service-authentication/models"
	authpkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/api"
	apimocks "github.com/a-novel/service-story-schematics/internal/api/mocks"
	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

func TestSelectBeatsSheet(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type selectBeatsSheetData struct {
		resp *models.BeatsSheet
		err  error
	}

	testCases := []struct {
		name string

		params apimodels.GetBeatsSheetParams

		selectBeatsSheetData *selectBeatsSheetData

		expect    apimodels.GetBeatsSheetRes
		expectErr error
	}{
		{
			name: "Success",

			params: apimodels.GetBeatsSheetParams{
				BeatsSheetID: apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
			},

			selectBeatsSheetData: &selectBeatsSheetData{
				resp: &models.BeatsSheet{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					LoglineID: uuid.MustParse("00000000-0000-0000-1000-000000000001"),
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

			expect: &apimodels.BeatsSheet{
				ID:        apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				LoglineID: apimodels.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				Content: []apimodels.Beat{
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
				Lang:      apimodels.LangEn,
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "BeatsSheetNotFound",

			params: apimodels.GetBeatsSheetParams{
				BeatsSheetID: apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
			},

			selectBeatsSheetData: &selectBeatsSheetData{
				err: dao.ErrBeatsSheetNotFound,
			},

			expect: &apimodels.NotFoundError{Error: dao.ErrBeatsSheetNotFound.Error()},
		},
		{
			name: "LoglineNotFound",

			params: apimodels.GetBeatsSheetParams{
				BeatsSheetID: apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
			},

			selectBeatsSheetData: &selectBeatsSheetData{
				err: dao.ErrLoglineNotFound,
			},

			expect: &apimodels.NotFoundError{Error: dao.ErrLoglineNotFound.Error()},
		},
		{
			name: "Error",

			params: apimodels.GetBeatsSheetParams{
				BeatsSheetID: apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
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

			source := apimocks.NewMockSelectBeatsSheetService(t)

			ctx := context.WithValue(t.Context(), authpkg.ClaimsContextKey{}, &authmodels.AccessTokenClaims{
				UserID: lo.ToPtr(uuid.MustParse("00000000-1000-0000-0000-000000000001")),
			})

			if testCase.selectBeatsSheetData != nil {
				source.EXPECT().
					SelectBeatsSheet(mock.Anything, services.SelectBeatsSheetRequest{
						BeatsSheetID: uuid.UUID(testCase.params.BeatsSheetID),
						UserID:       uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					}).
					Return(testCase.selectBeatsSheetData.resp, testCase.selectBeatsSheetData.err)
			}

			handler := api.API{SelectBeatsSheetService: source}

			res, err := handler.GetBeatsSheet(ctx, testCase.params)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)

			source.AssertExpectations(t)
		})
	}
}
