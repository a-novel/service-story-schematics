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
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
	"github.com/a-novel/service-story-schematics/models/api"
)

func TestExpandLogline(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type expandLoglineData struct {
		resp *models.LoglineIdea
		err  error
	}

	testCases := []struct {
		name string

		form *apimodels.LoglineIdea

		expandLoglineData *expandLoglineData

		expect    apimodels.ExpandLoglineRes
		expectErr error
	}{
		{
			name: "Success",

			form: &apimodels.LoglineIdea{
				Name:    "Logline 1",
				Content: "Logline 1 content",
				Lang:    apimodels.LangEn,
			},

			expandLoglineData: &expandLoglineData{
				resp: &models.LoglineIdea{
					Name:    "Logline 1 expanded",
					Content: "Logline 1 content expanded",
					Lang:    models.LangEN,
				},
			},

			expect: &apimodels.LoglineIdea{
				Name:    "Logline 1 expanded",
				Content: "Logline 1 content expanded",
				Lang:    apimodels.LangEn,
			},
		},
		{
			name: "Error",

			form: &apimodels.LoglineIdea{
				Name:    "Logline 1",
				Content: "Logline 1 content",
				Lang:    apimodels.LangEn,
			},

			expandLoglineData: &expandLoglineData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			source := apimocks.NewMockExpandLoglineService(t)

			ctx := context.WithValue(t.Context(), authpkg.ClaimsContextKey{}, &authmodels.AccessTokenClaims{
				UserID: lo.ToPtr(uuid.MustParse("00000000-1000-0000-0000-000000000001")),
			})

			if testCase.expandLoglineData != nil {
				source.EXPECT().
					ExpandLogline(mock.Anything, services.ExpandLoglineRequest{
						Logline: models.LoglineIdea{
							Name:    testCase.form.GetName(),
							Content: testCase.form.GetContent(),
							Lang:    models.Lang(testCase.form.GetLang()),
						},
						UserID: uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					}).
					Return(testCase.expandLoglineData.resp, testCase.expandLoglineData.err)
			}

			handler := api.API{ExpandLoglineService: source}

			res, err := handler.ExpandLogline(ctx, testCase.form)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)

			source.AssertExpectations(t)
		})
	}
}
