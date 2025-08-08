package api_test

import (
	"context"
	"errors"
	"testing"

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

func TestGenerateBeatsSheet(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type generateBeatsSheetData struct {
		resp []models.Beat
		err  error
	}

	testCases := []struct {
		name string

		form *apimodels.GenerateBeatsSheetForm

		generateBeatsSheetData *generateBeatsSheetData

		expect    apimodels.GenerateBeatsSheetRes
		expectErr error
	}{
		{
			name: "Success",

			form: &apimodels.GenerateBeatsSheetForm{
				LoglineID: apimodels.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				Lang:      apimodels.LangEn,
			},

			generateBeatsSheetData: &generateBeatsSheetData{
				resp: []models.Beat{
					{
						Key:     "beat-1",
						Title:   "Beat 1",
						Content: "Beat 1 content",
					},
					{
						Key:     "beat-2",
						Title:   "Beat 2",
						Content: "Beat 2 content",
					},
				},
			},

			expect: &apimodels.BeatsSheetIdea{
				Content: []apimodels.Beat{
					{
						Key:     "beat-1",
						Title:   "Beat 1",
						Content: "Beat 1 content",
					},
					{
						Key:     "beat-2",
						Title:   "Beat 2",
						Content: "Beat 2 content",
					},
				},
				Lang: apimodels.LangEn,
			},
		},
		{
			name: "LoglineNotFound",

			form: &apimodels.GenerateBeatsSheetForm{
				LoglineID: apimodels.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				Lang:      apimodels.LangEn,
			},

			generateBeatsSheetData: &generateBeatsSheetData{
				err: dao.ErrLoglineNotFound,
			},

			expect: &apimodels.NotFoundError{Error: dao.ErrLoglineNotFound.Error()},
		},
		{
			name: "StoryPlanNotFound",

			form: &apimodels.GenerateBeatsSheetForm{
				LoglineID: apimodels.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				Lang:      apimodels.LangEn,
			},

			generateBeatsSheetData: &generateBeatsSheetData{
				err: services.ErrStoryPlanNotFound,
			},

			expect: &apimodels.NotFoundError{Error: services.ErrStoryPlanNotFound.Error()},
		},
		{
			name: "Error",

			form: &apimodels.GenerateBeatsSheetForm{
				LoglineID: apimodels.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				Lang:      apimodels.LangEn,
			},

			generateBeatsSheetData: &generateBeatsSheetData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			source := apimocks.NewMockGenerateBeatsSheetService(t)

			ctx := context.WithValue(t.Context(), authpkg.ClaimsContextKey{}, &authmodels.AccessTokenClaims{
				UserID: lo.ToPtr(uuid.MustParse("00000000-1000-0000-0000-000000000001")),
			})

			if testCase.generateBeatsSheetData != nil {
				source.EXPECT().
					GenerateBeatsSheet(mock.Anything, services.GenerateBeatsSheetRequest{
						LoglineID: uuid.UUID(testCase.form.GetLoglineID()),
						UserID:    uuid.MustParse("00000000-1000-0000-0000-000000000001"),
						Lang:      models.Lang(testCase.form.GetLang()),
					}).
					Return(testCase.generateBeatsSheetData.resp, testCase.generateBeatsSheetData.err)
			}

			handler := api.API{GenerateBeatsSheetService: source}

			res, err := handler.GenerateBeatsSheet(ctx, testCase.form)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)

			source.AssertExpectations(t)
		})
	}
}
