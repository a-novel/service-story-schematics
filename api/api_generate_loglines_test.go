package api_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	authapi "github.com/a-novel/service-authentication/api"
	authmodels "github.com/a-novel/service-authentication/models"

	"github.com/a-novel-kit/context"

	"github.com/a-novel/service-story-schematics/api"
	"github.com/a-novel/service-story-schematics/api/codegen"
	apimocks "github.com/a-novel/service-story-schematics/api/mocks"
	"github.com/a-novel/service-story-schematics/internal/services"
	"github.com/a-novel/service-story-schematics/models"
)

func TestGenerateLoglines(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type generateLoglinesData struct {
		loglines []models.LoglineIdea
		err      error
	}

	testCases := []struct {
		name string

		form *codegen.GenerateLoglinesForm

		generateLoglinesData *generateLoglinesData

		expect    codegen.GenerateLoglinesRes
		expectErr error
	}{
		{
			name: "Success",

			form: &codegen.GenerateLoglinesForm{
				Count: 10,
				Theme: "theme",
			},

			generateLoglinesData: &generateLoglinesData{
				loglines: []models.LoglineIdea{
					{
						Name:    "Logline 1",
						Content: "Logline 1 content",
					},
					{
						Name:    "Logline 2",
						Content: "Logline 2 content",
					},
				},
			},

			expect: &codegen.GenerateLoglinesOKApplicationJSON{
				{
					Name:    "Logline 1",
					Content: "Logline 1 content",
				},
				{
					Name:    "Logline 2",
					Content: "Logline 2 content",
				},
			},
		},
		{
			name: "Error",

			form: &codegen.GenerateLoglinesForm{
				Count: 10,
				Theme: "theme",
			},

			generateLoglinesData: &generateLoglinesData{
				err: errFoo,
			},

			expectErr: errFoo,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			source := apimocks.NewMockGenerateLoglinesService(t)

			ctx := context.WithValue(t.Context(), authapi.ClaimsAPIKey{}, &authmodels.AccessTokenClaims{
				UserID: lo.ToPtr(uuid.MustParse("00000000-1000-0000-0000-000000000001")),
			})

			if testCase.generateLoglinesData != nil {
				source.EXPECT().
					GenerateLoglines(ctx, services.GenerateLoglinesRequest{
						Count:  testCase.form.GetCount(),
						Theme:  testCase.form.GetTheme(),
						UserID: uuid.MustParse("00000000-1000-0000-0000-000000000001"),
					}).
					Return(testCase.generateLoglinesData.loglines, testCase.generateLoglinesData.err)
			}

			handler := api.API{GenerateLoglinesService: source}

			res, err := handler.GenerateLoglines(ctx, testCase.form)
			require.ErrorIs(t, err, testCase.expectErr)
			require.Equal(t, testCase.expect, res)

			source.AssertExpectations(t)
		})
	}
}
