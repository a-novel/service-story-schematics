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
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
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

		form *codegen.LoglineIdea

		expandLoglineData *expandLoglineData

		expect    codegen.ExpandLoglineRes
		expectErr error
	}{
		{
			name: "Success",

			form: &codegen.LoglineIdea{
				Name:    "Logline 1",
				Content: "Logline 1 content",
			},

			expandLoglineData: &expandLoglineData{
				resp: &models.LoglineIdea{
					Name:    "Logline 1 expanded",
					Content: "Logline 1 content expanded",
				},
			},

			expect: &codegen.LoglineIdea{
				Name:    "Logline 1 expanded",
				Content: "Logline 1 content expanded",
			},
		},
		{
			name: "Error",

			form: &codegen.LoglineIdea{
				Name:    "Logline 1",
				Content: "Logline 1 content",
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

			ctx := context.WithValue(t.Context(), authapi.ClaimsAPIKey{}, &authmodels.AccessTokenClaims{
				UserID: lo.ToPtr(uuid.MustParse("00000000-1000-0000-0000-000000000001")),
			})

			if testCase.expandLoglineData != nil {
				source.EXPECT().
					ExpandLogline(ctx, services.ExpandLoglineRequest{
						Logline: models.LoglineIdea{
							Name:    testCase.form.GetName(),
							Content: testCase.form.GetContent(),
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
