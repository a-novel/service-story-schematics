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

	authModels "github.com/a-novel/service-authentication/models"
	authPkg "github.com/a-novel/service-authentication/pkg"

	"github.com/a-novel/service-story-schematics/internal/api"
	"github.com/a-novel/service-story-schematics/internal/api/codegen"
	apimocks "github.com/a-novel/service-story-schematics/internal/api/mocks"
	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
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

		params codegen.GetBeatsSheetParams

		selectBeatsSheetData *selectBeatsSheetData

		expect    codegen.GetBeatsSheetRes
		expectErr error
	}{
		{
			name: "Success",

			params: codegen.GetBeatsSheetParams{
				BeatsSheetID: codegen.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
			},

			selectBeatsSheetData: &selectBeatsSheetData{
				resp: &models.BeatsSheet{
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

			expect: &codegen.BeatsSheet{
				ID:          codegen.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				LoglineID:   codegen.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				StoryPlanID: codegen.StoryPlanID(uuid.MustParse("00000000-0000-0000-0100-000000000001")),
				Content: []codegen.Beat{
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
				Lang:      codegen.LangEn,
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "BeatsSheetNotFound",

			params: codegen.GetBeatsSheetParams{
				BeatsSheetID: codegen.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
			},

			selectBeatsSheetData: &selectBeatsSheetData{
				err: dao.ErrBeatsSheetNotFound,
			},

			expect: &codegen.NotFoundError{Error: dao.ErrBeatsSheetNotFound.Error()},
		},
		{
			name: "LoglineNotFound",

			params: codegen.GetBeatsSheetParams{
				BeatsSheetID: codegen.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
			},

			selectBeatsSheetData: &selectBeatsSheetData{
				err: dao.ErrLoglineNotFound,
			},

			expect: &codegen.NotFoundError{Error: dao.ErrLoglineNotFound.Error()},
		},
		{
			name: "Error",

			params: codegen.GetBeatsSheetParams{
				BeatsSheetID: codegen.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
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

			ctx := context.WithValue(t.Context(), authPkg.ClaimsContextKey{}, &authModels.AccessTokenClaims{
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
