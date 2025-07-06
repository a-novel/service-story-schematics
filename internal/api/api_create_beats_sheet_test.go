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
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

func TestCreateBeatsSheet(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type createBeatsSheetData struct {
		resp *models.BeatsSheet
		err  error
	}

	testCases := []struct {
		name string

		form *codegen.CreateBeatsSheetForm

		createBeatsSheetData *createBeatsSheetData

		expect    codegen.CreateBeatsSheetRes
		expectErr error
	}{
		{
			name: "Success",

			form: &codegen.CreateBeatsSheetForm{
				LoglineID:   codegen.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				StoryPlanID: codegen.StoryPlanID(uuid.MustParse("00000000-0000-1000-0000-000000000001")),
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
				Lang: codegen.LangEn,
			},

			createBeatsSheetData: &createBeatsSheetData{
				resp: &models.BeatsSheet{
					ID:          uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					LoglineID:   uuid.MustParse("00000000-0000-0000-1000-000000000001"),
					StoryPlanID: uuid.MustParse("00000000-0000-1000-0000-000000000001"),
					Content: []models.Beat{
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
					Lang:      models.LangEN,
					CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},

			expect: &codegen.BeatsSheet{
				ID:          codegen.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				LoglineID:   codegen.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				StoryPlanID: codegen.StoryPlanID(uuid.MustParse("00000000-0000-1000-0000-000000000001")),
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
				Lang:      codegen.LangEn,
				CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Error/LoglineNotFound",

			form: &codegen.CreateBeatsSheetForm{
				LoglineID:   codegen.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				StoryPlanID: codegen.StoryPlanID(uuid.MustParse("00000000-0000-1000-0000-000000000001")),
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
				Lang: codegen.LangEn,
			},

			createBeatsSheetData: &createBeatsSheetData{
				err: dao.ErrLoglineNotFound,
			},

			expect: &codegen.NotFoundError{Error: dao.ErrLoglineNotFound.Error()},
		},
		{
			name: "Error/StoryPlanNotFound",

			form: &codegen.CreateBeatsSheetForm{
				LoglineID:   codegen.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				StoryPlanID: codegen.StoryPlanID(uuid.MustParse("00000000-0000-1000-0000-000000000001")),
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
				Lang: codegen.LangEn,
			},

			createBeatsSheetData: &createBeatsSheetData{
				err: dao.ErrStoryPlanNotFound,
			},

			expect: &codegen.NotFoundError{Error: dao.ErrStoryPlanNotFound.Error()},
		},
		{
			name: "Error/InvalidStoryPlan",

			form: &codegen.CreateBeatsSheetForm{
				LoglineID:   codegen.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				StoryPlanID: codegen.StoryPlanID(uuid.MustParse("00000000-0000-1000-0000-000000000001")),
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
				Lang: codegen.LangEn,
			},

			createBeatsSheetData: &createBeatsSheetData{
				err: lib.ErrInvalidStoryPlan,
			},

			expect: &codegen.UnprocessableEntityError{Error: lib.ErrInvalidStoryPlan.Error()},
		},
		{
			name: "Error/CreateBeatsSheet",

			form: &codegen.CreateBeatsSheetForm{
				LoglineID:   codegen.LoglineID(uuid.MustParse("00000000-0000-0000-1000-000000000001")),
				StoryPlanID: codegen.StoryPlanID(uuid.MustParse("00000000-0000-1000-0000-000000000001")),
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
				Lang: codegen.LangEn,
			},

			createBeatsSheetData: &createBeatsSheetData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			source := apimocks.NewMockCreateBeatsSheetService(t)

			ctx := context.WithValue(t.Context(), authPkg.ClaimsContextKey{}, &authModels.AccessTokenClaims{
				UserID: lo.ToPtr(uuid.MustParse("00000000-1000-0000-0000-000000000001")),
			})

			if testCase.createBeatsSheetData != nil {
				source.EXPECT().
					CreateBeatsSheet(mock.Anything, services.CreateBeatsSheetRequest{
						LoglineID:   uuid.UUID(testCase.form.GetLoglineID()),
						UserID:      uuid.MustParse("00000000-1000-0000-0000-000000000001"),
						StoryPlanID: uuid.UUID(testCase.form.GetStoryPlanID()),
						Lang:        models.Lang(testCase.form.GetLang()),
						Content: lo.Map(testCase.form.GetContent(), func(item codegen.Beat, _ int) models.Beat {
							return models.Beat{
								Key:     item.GetKey(),
								Title:   item.GetTitle(),
								Content: item.GetContent(),
							}
						}),
					}).
					Return(testCase.createBeatsSheetData.resp, testCase.createBeatsSheetData.err)
			}

			handler := api.API{CreateBeatsSheetService: source}

			res, err := handler.CreateBeatsSheet(ctx, testCase.form)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)

			source.AssertExpectations(t)
		})
	}
}
