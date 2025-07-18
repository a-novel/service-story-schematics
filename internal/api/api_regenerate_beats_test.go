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

func TestRegenerateBeats(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type regenerateBeatsData struct {
		resp []models.Beat
		err  error
	}

	testCases := []struct {
		name string

		form *apimodels.RegenerateBeatsForm

		regenerateBeatsData *regenerateBeatsData

		expect    apimodels.RegenerateBeatsRes
		expectErr error
	}{
		{
			name: "Success",

			form: &apimodels.RegenerateBeatsForm{
				BeatsSheetID: apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				RegenerateKeys: []string{
					"beat-1",
					"beat-2",
				},
			},

			regenerateBeatsData: &regenerateBeatsData{
				resp: []models.Beat{
					{
						Key:     "beat-1",
						Title:   "Regenerated Beat 1",
						Content: "Regenerated Content 1",
					},
					{
						Key:     "beat-2",
						Title:   "Regenerated Beat 2",
						Content: "Regenerated Content 2",
					},
				},
			},

			expect: &apimodels.Beats{
				{
					Key:     "beat-1",
					Title:   "Regenerated Beat 1",
					Content: "Regenerated Content 1",
				},
				{
					Key:     "beat-2",
					Title:   "Regenerated Beat 2",
					Content: "Regenerated Content 2",
				},
			},
		},
		{
			name: "BeatsSheetNotFound",

			form: &apimodels.RegenerateBeatsForm{
				BeatsSheetID: apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				RegenerateKeys: []string{
					"beat-1",
					"beat-2",
				},
			},

			regenerateBeatsData: &regenerateBeatsData{
				err: dao.ErrBeatsSheetNotFound,
			},

			expect: &apimodels.NotFoundError{Error: dao.ErrBeatsSheetNotFound.Error()},
		},
		{
			name: "LoglineNotFound",

			form: &apimodels.RegenerateBeatsForm{
				BeatsSheetID: apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				RegenerateKeys: []string{
					"beat-1",
					"beat-2",
				},
			},

			regenerateBeatsData: &regenerateBeatsData{
				err: dao.ErrLoglineNotFound,
			},

			expect: &apimodels.NotFoundError{Error: dao.ErrLoglineNotFound.Error()},
		},
		{
			name: "StoryPlanNotFound",

			form: &apimodels.RegenerateBeatsForm{
				BeatsSheetID: apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				RegenerateKeys: []string{
					"beat-1",
					"beat-2",
				},
			},

			regenerateBeatsData: &regenerateBeatsData{
				err: dao.ErrStoryPlanNotFound,
			},

			expect: &apimodels.NotFoundError{Error: dao.ErrStoryPlanNotFound.Error()},
		},
		{
			name: "Error",

			form: &apimodels.RegenerateBeatsForm{
				BeatsSheetID: apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				RegenerateKeys: []string{
					"beat-1",
					"beat-2",
				},
			},

			regenerateBeatsData: &regenerateBeatsData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			source := apimocks.NewMockRegenerateBeatsService(t)

			ctx := context.WithValue(t.Context(), authpkg.ClaimsContextKey{}, &authmodels.AccessTokenClaims{
				UserID: lo.ToPtr(uuid.MustParse("00000000-1000-0000-0000-000000000001")),
			})

			if testCase.regenerateBeatsData != nil {
				source.EXPECT().
					RegenerateBeats(mock.Anything, services.RegenerateBeatsRequest{
						BeatsSheetID:   uuid.UUID(testCase.form.GetBeatsSheetID()),
						UserID:         uuid.MustParse("00000000-1000-0000-0000-000000000001"),
						RegenerateKeys: testCase.form.GetRegenerateKeys(),
					}).
					Return(testCase.regenerateBeatsData.resp, testCase.regenerateBeatsData.err)
			}

			handler := api.API{RegenerateBeatsService: source}

			res, err := handler.RegenerateBeats(ctx, testCase.form)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)

			source.AssertExpectations(t)
		})
	}
}
