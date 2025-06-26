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

	authapi "github.com/a-novel/service-authentication/api"
	authmodels "github.com/a-novel/service-authentication/models"

	"github.com/a-novel/service-story-schematics/api"
	"github.com/a-novel/service-story-schematics/api/codegen"
	apimocks "github.com/a-novel/service-story-schematics/api/mocks"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
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

		params codegen.GetBeatsSheetsParams

		listBeatsSheetsData *listBeatsSheetsData

		expect    codegen.GetBeatsSheetsRes
		expectErr error
	}{
		{
			name: "Success",

			params: codegen.GetBeatsSheetsParams{
				LoglineID: codegen.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				Limit:     codegen.OptInt{Value: 10, Set: true},
				Offset:    codegen.OptInt{Value: 2, Set: true},
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

			expect: &codegen.GetBeatsSheetsOKApplicationJSON{
				{
					ID:        codegen.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					Lang:      codegen.LangEn,
					CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        codegen.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
					Lang:      codegen.LangEn,
					CreatedAt: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "Error",

			params: codegen.GetBeatsSheetsParams{
				LoglineID: codegen.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				Limit:     codegen.OptInt{Value: 10, Set: true},
				Offset:    codegen.OptInt{Value: 2, Set: true},
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

			ctx := context.WithValue(t.Context(), authapi.ClaimsAPIKey{}, &authmodels.AccessTokenClaims{
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
