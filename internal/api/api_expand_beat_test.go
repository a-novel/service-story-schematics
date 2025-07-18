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
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

func TestExpandBeat(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type expandBeatData struct {
		resp *models.Beat
		err  error
	}

	testCases := []struct {
		name string

		form *apimodels.ExpandBeatForm

		expandBeatData *expandBeatData

		expect    apimodels.ExpandBeatRes
		expectErr error
	}{
		{
			name: "Success",

			form: &apimodels.ExpandBeatForm{
				BeatsSheetID: apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				TargetKey:    "beat-1",
			},

			expandBeatData: &expandBeatData{
				resp: &models.Beat{
					Key:     "beat-1",
					Title:   "Beat 1 expanded",
					Content: "Beat 1 content expanded",
				},
			},

			expect: &apimodels.Beat{
				Key:     "beat-1",
				Title:   "Beat 1 expanded",
				Content: "Beat 1 content expanded",
			},
		},
		{
			name: "BeatsSheetNotFound",

			form: &apimodels.ExpandBeatForm{
				BeatsSheetID: apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				TargetKey:    "beat-1",
			},

			expandBeatData: &expandBeatData{
				err: dao.ErrBeatsSheetNotFound,
			},

			expect: &apimodels.NotFoundError{Error: dao.ErrBeatsSheetNotFound.Error()},
		},
		{
			name: "StoryPlanNotFound",

			form: &apimodels.ExpandBeatForm{
				BeatsSheetID: apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				TargetKey:    "beat-1",
			},

			expandBeatData: &expandBeatData{
				err: dao.ErrStoryPlanNotFound,
			},

			expect: &apimodels.NotFoundError{Error: dao.ErrStoryPlanNotFound.Error()},
		},
		{
			name: "UnknownTargetKey",

			form: &apimodels.ExpandBeatForm{
				BeatsSheetID: apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				TargetKey:    "beat-1",
			},

			expandBeatData: &expandBeatData{
				err: daoai.ErrUnknownTargetKey,
			},

			expect: &apimodels.UnprocessableEntityError{Error: daoai.ErrUnknownTargetKey.Error()},
		},
		{
			name: "Error",

			form: &apimodels.ExpandBeatForm{
				BeatsSheetID: apimodels.BeatsSheetID(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
				TargetKey:    "beat-1",
			},

			expandBeatData: &expandBeatData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			source := apimocks.NewMockExpandBeatService(t)

			ctx := context.WithValue(t.Context(), authpkg.ClaimsContextKey{}, &authmodels.AccessTokenClaims{
				UserID: lo.ToPtr(uuid.MustParse("00000000-1000-0000-0000-000000000001")),
			})

			if testCase.expandBeatData != nil {
				source.EXPECT().
					ExpandBeat(mock.Anything, services.ExpandBeatRequest{
						BeatsSheetID: uuid.UUID(testCase.form.GetBeatsSheetID()),
						TargetKey:    testCase.form.GetTargetKey(),
						UserID:       uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					}).
					Return(testCase.expandBeatData.resp, testCase.expandBeatData.err)
			}

			handler := api.API{ExpandBeatService: source}

			res, err := handler.ExpandBeat(ctx, testCase.form)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)

			source.AssertExpectations(t)
		})
	}
}
