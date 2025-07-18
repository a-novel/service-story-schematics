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
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

func TestListBeatsSheets(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type listBeatsSheetsData struct {
		resp []*models.BeatsSheetPreview
		err  error
	}

	testCases := []struct {
		name string

		params apimodels.GetBeatsSheetsParams

		listBeatsSheetsData *listBeatsSheetsData

		expect    apimodels.GetBeatsSheetsRes
		expectErr error
	}{
		{
			name: "Success",

			params: apimodels.GetBeatsSheetsParams{
				LoglineID: apimodels.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				Limit:     apimodels.OptInt{Value: 10, Set: true},
				Offset:    apimodels.OptInt{Value: 2, Set: true},
			},

			listBeatsSheetsData: &listBeatsSheetsData{
				resp: []*models.BeatsSheetPreview{
					{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Lang:      models.LangEN,
						CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
						Lang:      models.LangEN,
						CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
					},
				},
			},

			expect: &apimodels.GetBeatsSheetsOKApplicationJSON{
				{
					ID:        apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					Lang:      apimodels.LangEn,
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
					Lang:      apimodels.LangEn,
					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "Error",

			params: apimodels.GetBeatsSheetsParams{
				LoglineID: apimodels.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				Limit:     apimodels.OptInt{Value: 10, Set: true},
				Offset:    apimodels.OptInt{Value: 2, Set: true},
			},

			listBeatsSheetsData: &listBeatsSheetsData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			source := apimocks.NewMockListBeatsSheetsService(t)

			ctx := context.WithValue(t.Context(), authpkg.ClaimsContextKey{}, &authmodels.AccessTokenClaims{
				UserID: lo.ToPtr(uuid.MustParse("00000000-1000-0000-0000-000000000001")),
			})

			if testCase.listBeatsSheetsData != nil {
				source.EXPECT().
					ListBeatsSheets(mock.Anything, services.ListBeatsSheetsRequest{
						UserID:    uuid.MustParse("00000000-1000-0000-0000-000000000001"),
						LoglineID: uuid.UUID(testCase.params.LoglineID),
						Limit:     testCase.params.Limit.Value,
						Offset:    testCase.params.Offset.Value,
					}).
					Return(testCase.listBeatsSheetsData.resp, testCase.listBeatsSheetsData.err)
			}

			handler := api.API{ListBeatsSheetsService: source}

			res, err := handler.GetBeatsSheets(ctx, testCase.params)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)

			source.AssertExpectations(t)
		})
	}
}
