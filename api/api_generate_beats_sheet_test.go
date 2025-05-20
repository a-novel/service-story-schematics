package api_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	authapi "github.com/a-novel/service-authentication/api"
	authmodels "github.com/a-novel/service-authentication/models"

	"github.com/a-novel/service-story-schematics/api"
	"github.com/a-novel/service-story-schematics/api/codegen"
	apimocks "github.com/a-novel/service-story-schematics/api/mocks"
	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
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

		form *codegen.GenerateBeatsSheetForm

		generateBeatsSheetData *generateBeatsSheetData

		expect    codegen.GenerateBeatsSheetRes
		expectErr error
	}{
		{
			name: "Success",

			form: &codegen.GenerateBeatsSheetForm{
				LoglineID:   codegen.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				StoryPlanID: codegen.StoryPlanID(uuid.New()),
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

			expect: &codegen.BeatsSheetIdea{
				Content: []codegen.Beat{
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
		},
		{
			name: "LoglineNotFound",

			form: &codegen.GenerateBeatsSheetForm{
				LoglineID:   codegen.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				StoryPlanID: codegen.StoryPlanID(uuid.New()),
			},

			generateBeatsSheetData: &generateBeatsSheetData{
				err: dao.ErrLoglineNotFound,
			},

			expect: &codegen.NotFoundError{Error: dao.ErrLoglineNotFound.Error()},
		},
		{
			name: "StoryPlanNotFound",

			form: &codegen.GenerateBeatsSheetForm{
				LoglineID:   codegen.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				StoryPlanID: codegen.StoryPlanID(uuid.New()),
			},

			generateBeatsSheetData: &generateBeatsSheetData{
				err: dao.ErrStoryPlanNotFound,
			},

			expect: &codegen.NotFoundError{Error: dao.ErrStoryPlanNotFound.Error()},
		},
		{
			name: "Error",

			form: &codegen.GenerateBeatsSheetForm{
				LoglineID:   codegen.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				StoryPlanID: codegen.StoryPlanID(uuid.New()),
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

			ctx := context.WithValue(t.Context(), authapi.ClaimsAPIKey{}, &authmodels.AccessTokenClaims{
				UserID: lo.ToPtr(uuid.MustParse("00000000-1000-0000-0000-000000000001")),
			})

			if testCase.generateBeatsSheetData != nil {
				source.EXPECT().
					GenerateBeatsSheet(ctx, services.GenerateBeatsSheetRequest{
						LoglineID:   uuid.UUID(testCase.form.GetLoglineID()),
						StoryPlanID: uuid.UUID(testCase.form.GetStoryPlanID()),
						UserID:      uuid.MustParse("00000000-1000-0000-0000-000000000001"),
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
